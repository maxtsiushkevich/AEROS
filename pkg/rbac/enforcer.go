package rbac

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/persist"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/google/uuid"
	rds "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthorizationService interface {
	IsAuthenticated(sub string, obj string, act string) (bool, error)
	CreateRole(ctx context.Context, name, description string) (*Role, error)
	DeleteRole(ctx context.Context, roleName string) error
	CreateAction(ctx context.Context, name string) (*Action, error)
	CreateResource(ctx context.Context, name, description string) (*Resource, error)
	CreatePermission(ctx context.Context, resourceName, actionName string) (*Permission, error)
	GrantPermissionToRole(ctx context.Context, roleName, resourceName, actionName string) (*RolePermission, error)
	RevokePermissionFromRole(ctx context.Context, roleName, resourceName, actionName string) error
	AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error
	RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error
}

type RBACService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
	watcher  persist.WatcherEx
	reloadMu sync.Mutex
}

func newRBACService(cfg CasbinConfig) (*RBACService, error) {
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.DbName == "" || cfg.ConfigPath == "" || cfg.RedisAddr == "" || cfg.RedisChannel == "" {
		return nil, fmt.Errorf("RBAC config is incomplete: missing required DB or Redis fields")
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)

	// DB is the source of truth
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect Postgres: %w", err)
	}

	if err := db.AutoMigrate(
		&Role{},
		&Action{},
		&Resource{},
		&Permission{},
		&RolePermission{},
		&UserRole{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate RBAC tables: %w", err)
	}

	// Create an in-memory Casbin enforcer
	e, err := casbin.NewEnforcer(cfg.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("create Casbin enforcer: %w", err)
	}

	svc := &RBACService{db: db, enforcer: e}

	// Load policies
	if err := svc.loadPoliciesFromDB(); err != nil {
		return nil, fmt.Errorf("load policies from DB: %w", err)
	}

	watchAddr := cfg.RedisAddr
	if strings.Contains(watchAddr, "://") {
		if u, err := url.Parse(watchAddr); err == nil {
			if u.Host != "" {
				watchAddr = u.Host
			}
		}
	} else if strings.Contains(watchAddr, "@") {
		// strip potential userinfo
		if parts := strings.SplitN(watchAddr, "@", 2); len(parts) == 2 {
			watchAddr = parts[1]
		}
	}

	opts := rediswatcher.WatcherOptions{}
	if cfg.RedisPass != "" {
		opts.Options = rds.Options{Password: cfg.RedisPass}
	}
	if cfg.RedisChannel != "" {
		opts.Channel = cfg.RedisChannel
	}
	opts.IgnoreSelf = true
	opts.LocalID = uuid.NewString()

	w, err := rediswatcher.NewWatcher(watchAddr, opts)
	if err != nil {
		return nil, fmt.Errorf("create redis watcher: %w", err)
	}
	watcher, ok := w.(persist.WatcherEx)
	if !ok {
		return nil, fmt.Errorf("redis watcher does not implement persist.WatcherEx")
	}

	if err := w.SetUpdateCallback(func(msg string) {
		if err := svc.loadPoliciesFromDB(); err != nil {
			slog.Error("failed to reload policies from DB on watcher callback", "err", err)
		}
	}); err != nil {
		return nil, fmt.Errorf("configure redis watcher callback: %w", err)
	}

	if err := e.SetWatcher(w); err != nil {
		return nil, fmt.Errorf("set Casbin watcher: %w", err)
	}
	e.EnableAutoNotifyWatcher(false)
	svc.watcher = watcher

	return svc, nil
}

func (s *RBACService) loadPoliciesFromDB() error {
	s.reloadMu.Lock()
	defer s.reloadMu.Unlock()

	s.enforcer.ClearPolicy()

	var rps []RolePermission
	if err := s.db.Find(&rps).Error; err != nil {
		return fmt.Errorf("fetch role permissions: %w", err)
	}
	var loadedPolicies []string
	for _, rp := range rps {
		s.enforcer.AddPolicy(rp.RoleName, rp.ResourceName, rp.ActionName)
		loadedPolicies = append(loadedPolicies, fmt.Sprintf("%s,%s,%s", rp.RoleName, rp.ResourceName, rp.ActionName))
	}

	var urs []UserRole
	if err := s.db.Find(&urs).Error; err != nil {
		return fmt.Errorf("fetch user roles: %w", err)
	}
	var loadedGroupings []string
	for _, ur := range urs {
		s.enforcer.AddGroupingPolicy(ur.UserID.String(), ur.RoleName)
		loadedGroupings = append(loadedGroupings, fmt.Sprintf("%s->%s", ur.UserID.String(), ur.RoleName))
	}

	slog.Info("Loaded policies from DB", "count_policies", len(loadedPolicies), "policies", loadedPolicies, "count_groupings", len(loadedGroupings), "groupings", loadedGroupings)

	return nil
}

func NewRBACService() (*RBACService, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("load RBAC config from env: %w", err)
	}
	return newRBACService(cfg)
}

func (cs *RBACService) Close() error {
	database, err := cs.db.DB()
	if err != nil {
		return err
	}
	err = database.Close()
	if err != nil {
		return err
	}

	cs.watcher.Close()
	return nil

}

func (cs *RBACService) IsAuthenticated(sub string, obj string, act string) (bool, error) {
	ok, reason, err := cs.enforcer.EnforceEx(sub, obj, act)
	if err != nil {
		slog.Error("Error occurred while enforcing policy", "err", err)
		return ok, err
	}

	slog.Debug("Enforce result", "sub", sub, "obj", obj, "act", act, "ok", ok, "reason", reason)
	return ok, nil
}

func (s *RBACService) CreateRole(ctx context.Context, name, description string) (*Role, error) {
	var existingRole Role
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&existingRole).Error; err == nil {
		return &existingRole, RoleExistsError
	}

	role := &Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (s *RBACService) DeleteRole(ctx context.Context, roleName string) error {
	var role Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", roleName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return RoleNotFound
		}
		return fmt.Errorf("failed to find role: %w", err)
	}

	// Delete RolePermission
	if err := s.db.WithContext(ctx).Delete(&RolePermission{}, "role_name = ?", roleName).Error; err != nil {
		return fmt.Errorf("failed to delete role permissions: %w", err)
	}

	// Delete Role
	if err := s.db.WithContext(ctx).Delete(&Role{}, "name = ?", roleName).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Update Casbin cache and notify other instances via watcher
	if _, err := s.enforcer.RemoveFilteredPolicy(0, roleName); err != nil {
		return fmt.Errorf("remove role policies from enforcer: %w", err)
	}
	if err := s.watcher.UpdateForRemoveFilteredPolicy("p", "p", 0, roleName); err != nil {
		return fmt.Errorf("notify role policy removal: %w", err)
	}

	return nil
}

func (s *RBACService) CreateAction(ctx context.Context, name string) (*Action, error) {
	var existingAction Action
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&existingAction).Error; err == nil {
		return &existingAction, ActionExistsError
	}

	action := &Action{
		Name: name,
	}

	if err := s.db.WithContext(ctx).Create(action).Error; err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return action, nil
}

func (s *RBACService) CreateResource(ctx context.Context, name, description string) (*Resource, error) {
	var existingResource Resource
	if err := s.db.WithContext(ctx).Where("name = ?", name).First(&existingResource).Error; err == nil {
		return &existingResource, ResourceExistsError
	}

	resource := &Resource{
		Name:        name,
		Description: description,
	}

	if err := s.db.WithContext(ctx).Create(resource).Error; err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (s *RBACService) CreatePermission(ctx context.Context, resourceName, actionName string) (*Permission, error) {
	var resource *Resource
	if err := s.db.WithContext(ctx).First(&resource, "name = ?", resourceName).Error; err != nil {
		return nil, ResourceNotFound
	}

	var action *Action
	if err := s.db.WithContext(ctx).First(&action, "name = ?", actionName).Error; err != nil {
		return nil, ActionNotFound
	}

	var existingPermission Permission
	if err := s.db.WithContext(ctx).Where("resource_name = ? AND action_name = ?", resourceName, actionName).First(&existingPermission).Error; err == nil {
		existingPermission.Resource = resource
		existingPermission.Action = action
		return &existingPermission, PermissionExistsError
	}

	perm := &Permission{
		ResourceName: resourceName,
		ActionName:   actionName,
		Resource:     resource,
		Action:       action,
	}

	if err := s.db.WithContext(ctx).Create(perm).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return perm, nil
}

func (s *RBACService) GrantPermissionToRole(ctx context.Context, roleName, resourceName, actionName string) (*RolePermission, error) {
	var role *Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", roleName).Error; err != nil {
		return nil, RoleNotFound
	}

	var perm *Permission
	if err := s.db.WithContext(ctx).First(&perm, "resource_name = ? AND action_name = ?", resourceName, actionName).Error; err != nil {
		return nil, PermissionNotFound
	}

	var existingRolePermission RolePermission
	if err := s.db.WithContext(ctx).Where("role_name = ? AND resource_name = ? AND action_name = ?", roleName, resourceName, actionName).First(&existingRolePermission).Error; err == nil {
		existingRolePermission.RoleName = roleName
		existingRolePermission.ResourceName = resourceName
		existingRolePermission.ActionName = actionName
		existingRolePermission.Role = role
		existingRolePermission.Permission = perm
		return &existingRolePermission, RolePermissionExistsError
	}

	rp := &RolePermission{
		RoleName:     roleName,
		ResourceName: resourceName,
		ActionName:   actionName,
		Role:         role,
		Permission:   perm,
	}

	if err := s.db.WithContext(ctx).Create(rp).Error; err != nil {
		return nil, fmt.Errorf("failed to grant permission to role: %w", err)
	}

	if _, err := s.enforcer.AddPolicy(roleName, resourceName, actionName); err != nil {
		return nil, fmt.Errorf("add policy to enforcer: %w", err)
	}
	if err := s.watcher.UpdateForAddPolicy("p", "p", roleName, resourceName, actionName); err != nil {
		return nil, fmt.Errorf("notify policy addition: %w", err)
	}

	return rp, nil
}

func (s *RBACService) RevokePermissionFromRole(ctx context.Context, roleName, resourceName, actionName string) error {
	var perm Permission
	if err := s.db.WithContext(ctx).First(&perm, "resource_name = ? AND action_name = ?", resourceName, actionName).Error; err != nil {
		return PermissionNotFound
	}

	var role Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", roleName).Error; err != nil {
		return RoleNotFound
	}

	if err := s.db.WithContext(ctx).Delete(&RolePermission{},
		"role_name = ? AND resource_name = ? AND action_name = ?",
		roleName, resourceName, actionName).Error; err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	if _, err := s.enforcer.RemovePolicy(roleName, resourceName, actionName); err != nil {
		return fmt.Errorf("remove policy from enforcer: %w", err)
	}
	if err := s.watcher.UpdateForRemovePolicy("p", "p", roleName, resourceName, actionName); err != nil {
		return fmt.Errorf("notify policy removal: %w", err)
	}

	return nil
}

func (s *RBACService) AssignRoleToUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	var role Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", roleName).Error; err != nil {
		return RoleNotFound
	}

	ur := &UserRole{
		UserID:   userID,
		RoleName: roleName,
	}

	if err := s.db.WithContext(ctx).Create(ur).Error; err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	if _, err := s.enforcer.AddGroupingPolicy(userID.String(), roleName); err != nil {
		return fmt.Errorf("add grouping policy to enforcer: %w", err)
	}
	if err := s.watcher.UpdateForAddPolicy("g", "g", userID.String(), roleName); err != nil {
		return fmt.Errorf("notify grouping policy addition: %w", err)
	}

	return nil
}

func (s *RBACService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	if err := s.db.WithContext(ctx).Delete(&UserRole{},
		"user_id = ? AND role_name = ?", userID, roleName).Error; err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	if _, err := s.enforcer.RemoveGroupingPolicy(userID.String(), roleName); err != nil {
		return fmt.Errorf("remove grouping policy from enforcer: %w", err)
	}
	if err := s.watcher.UpdateForRemovePolicy("g", "g", userID.String(), roleName); err != nil {
		return fmt.Errorf("notify grouping policy removal: %w", err)
	}

	return nil
}
