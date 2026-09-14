package rbac

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

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
	db         *gorm.DB
	enforcer   *casbin.Enforcer
	watcher    persist.Watcher
	reloadMu   sync.Mutex
	lastReload time.Time
}

func NewRBACService(cfg CasbinConfig) (*RBACService, error) {
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

	type cbSetter interface {
		SetUpdateCallback(func(string)) error
	}
	if setter, ok := w.(cbSetter); ok {
		_ = setter.SetUpdateCallback(func(msg string) {
			svc.reloadMu.Lock()
			since := time.Since(svc.lastReload)
			svc.reloadMu.Unlock()
			if since < 500*time.Millisecond {
				return
			}

			if err := svc.loadPoliciesFromDB(); err != nil {
				slog.Error("failed to reload policies from DB on watcher callback", "err", err)
			}
		})
	}

	e.SetWatcher(w)
	svc.watcher = w

	return svc, nil
}

func (s *RBACService) loadPoliciesFromDB() error {
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

	s.reloadMu.Lock()
	s.lastReload = time.Now().UTC()
	s.reloadMu.Unlock()

	return nil
}

func NewRBACServiceFromEnv() (*RBACService, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("load RBAC config from env: %w", err)
	}
	return NewRBACService(cfg)
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
	s.enforcer.RemoveFilteredPolicy(0, roleName)
	if s.watcher != nil {
		type remFiltered interface {
			UpdateForRemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error
		}
		if rf, ok := s.watcher.(remFiltered); ok {
			_ = rf.UpdateForRemoveFilteredPolicy("p", "p", 0, roleName)
		}
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

	s.enforcer.AddPolicy(roleName, resourceName, actionName)
	if s.watcher != nil {
		type addUpd interface {
			UpdateForAddPolicy(sec, ptype string, params ...string) error
		}
		if au, ok := s.watcher.(addUpd); ok {
			_ = au.UpdateForAddPolicy("p", "p", roleName, resourceName, actionName)
		}
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

	s.enforcer.RemovePolicy(roleName, resourceName, actionName)
	if s.watcher != nil {
		type remUpd interface {
			UpdateForRemovePolicy(sec, ptype string, params ...string) error
		}
		if ru, ok := s.watcher.(remUpd); ok {
			_ = ru.UpdateForRemovePolicy("p", "p", roleName, resourceName, actionName)
		}
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

	s.enforcer.AddGroupingPolicy(userID.String(), roleName)
	if s.watcher != nil {
		type addUpd interface {
			UpdateForAddPolicy(sec, ptype string, params ...string) error
		}
		if au, ok := s.watcher.(addUpd); ok {
			_ = au.UpdateForAddPolicy("g", "g", userID.String(), roleName)
		}
	}

	return nil
}

func (s *RBACService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	if err := s.db.WithContext(ctx).Delete(&UserRole{},
		"user_id = ? AND role_name = ?", userID, roleName).Error; err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	s.enforcer.RemoveGroupingPolicy(userID.String(), roleName)
	if s.watcher != nil {
		type remUpd interface {
			UpdateForRemovePolicy(sec, ptype string, params ...string) error
		}
		if ru, ok := s.watcher.(remUpd); ok {
			_ = ru.UpdateForRemovePolicy("g", "g", userID.String(), roleName)
		}
	}

	return nil
}
