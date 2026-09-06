package rbac

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var rbacService = NewRBACService()

type AuthorizationService interface {
	// AddUserToRbac(id uuid.UUID)
	IsAuthenticated(sub string, obj string, act string) (bool, error)
	CreateRole(name, description string) (*Role, error)
	DeleteRole(roleName string) error
	CreateAction(name string) (*Action, error)
	CreateResource(name, description string) (*Resource, error)
	CreatePermission(resourceName, actionName string) (*Permission, error)
	GrantPermissionToRole(roleName, resourceName, actionName string) error
	RevokePermissionFromRole(roleName, resourceName, actionName string) error
	AssignRoleToUser(userID uuid.UUID, roleName string) error
	RemoveRoleFromUser(userID uuid.UUID, roleName string) error
}

type RBACService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
}

func AddUserToRbac(id uuid.UUID) {
	rbacService.AssignRoleToUser(id, "user")
}

func NewRBACService() *RBACService {
	cfg, err := LoadConfig("../rbac/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName)

	a, err := gormadapter.NewAdapter("postgres", connString, true)
	if err != nil {
		log.Fatal(err)
	}

	e, err := casbin.NewEnforcer(cfg.ConfigPath, a)
	if err != nil {
		log.Fatal("Failed to create enforcer:", err)
	}

	if err := e.LoadPolicy(); err != nil {
		log.Fatal("Failed to load policy:", err)
	}

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&Role{},
		&Action{},
		&Resource{},
		&Permission{},
		&RolePermission{},
		&UserRole{},
	); err != nil {
		log.Fatal(err)
	}

	return &RBACService{
		db:       db,
		enforcer: e,
	}
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

func (s *RBACService) CreateRole(name, description string) (*Role, error) {
	role := &Role{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (s *RBACService) DeleteRole(roleName string) error {
	// Delete RolePermission
	if err := s.db.Delete(&RolePermission{}, "role_name = ?", roleName).Error; err != nil {
		return fmt.Errorf("failed to delete role permissions: %w", err)
	}

	// Delete Role
	if err := s.db.Delete(&Role{}, "name = ?", roleName).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	// Update Casbin policies
	s.enforcer.RemoveFilteredPolicy(0, roleName)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) CreateAction(name string) (*Action, error) {
	action := &Action{
		Name: name,
	}

	if err := s.db.Create(action).Error; err != nil {
		return nil, fmt.Errorf("failed to create action: %w", err)
	}

	return action, nil
}

func (s *RBACService) CreateResource(name, description string) (*Resource, error) {
	resource := &Resource{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(resource).Error; err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (s *RBACService) CreatePermission(resourceName, actionName string) (*Permission, error) {
	var resource Resource
	if err := s.db.First(&resource, "name = ?", resourceName).Error; err != nil {
		return nil, fmt.Errorf("resource not found: %w", err)
	}

	var action Action
	if err := s.db.First(&action, "name = ?", actionName).Error; err != nil {
		return nil, fmt.Errorf("action not found: %w", err)
	}

	perm := &Permission{
		ResourceName: resourceName,
		ActionName:   actionName,
	}

	if err := s.db.Create(perm).Error; err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return perm, nil
}

func (s *RBACService) GrantPermissionToRole(roleName, resourceName, actionName string) error {
	var role Role
	if err := s.db.First(&role, "name = ?", roleName).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	var perm Permission
	if err := s.db.First(&perm, "resource_name = ? AND action_name = ?", resourceName, actionName).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	rp := &RolePermission{
		RoleName:     roleName,
		ResourceName: resourceName,
		ActionName:   actionName,
	}

	if err := s.db.Create(rp).Error; err != nil {
		return fmt.Errorf("failed to grant permission to role: %w", err)
	}

	s.enforcer.AddPolicy(roleName, resourceName, actionName)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) RevokePermissionFromRole(roleName, resourceName, actionName string) error {
	var perm Permission
	if err := s.db.First(&perm, "resource_name = ? AND action_name = ?", resourceName, actionName).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	if err := s.db.Delete(&RolePermission{},
		"role_name = ? AND resource_name = ? AND action_name = ?",
		roleName, resourceName, actionName).Error; err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	s.enforcer.RemovePolicy(roleName, resourceName, actionName)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) AssignRoleToUser(userID uuid.UUID, roleName string) error {
	var role Role
	if err := s.db.First(&role, "name = ?", roleName).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	ur := &UserRole{
		UserID:   userID,
		RoleName: roleName,
	}

	if err := s.db.Create(ur).Error; err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	s.enforcer.AddGroupingPolicy(userID.String(), roleName)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) RemoveRoleFromUser(userID uuid.UUID, roleName string) error {
	if err := s.db.Delete(&UserRole{},
		"user_id = ? AND role_name = ?", userID, roleName).Error; err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	s.enforcer.RemoveGroupingPolicy(userID.String(), roleName)
	s.enforcer.SavePolicy()

	return nil
}
