package rbac

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/uuid"
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
}

func NewRBACService(cfg CasbinConfig) (*RBACService, error) {
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.DbName == "" || cfg.ConfigPath == "" {
		return nil, fmt.Errorf("RBAC config is incomplete")
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)

	a, err := gormadapter.NewAdapter("postgres", connString, true)
	if err != nil {
		return nil, fmt.Errorf("create Casbin adapter: %w", err)
	}

	e, err := casbin.NewEnforcer(cfg.ConfigPath, a)
	if err != nil {
		return nil, fmt.Errorf("create Casbin enforcer: %w", err)
	}
	if err := e.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load Casbin policy: %w", err)
	}

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

	return &RBACService{db: db, enforcer: e}, nil
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

	// Update Casbin policies
	s.enforcer.RemoveFilteredPolicy(0, roleName)
	s.enforcer.SavePolicy()

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
	s.enforcer.SavePolicy()

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
	s.enforcer.SavePolicy()

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
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) RemoveRoleFromUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	if err := s.db.WithContext(ctx).Delete(&UserRole{},
		"user_id = ? AND role_name = ?", userID, roleName).Error; err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	s.enforcer.RemoveGroupingPolicy(userID.String(), roleName)
	s.enforcer.SavePolicy()

	return nil
}
