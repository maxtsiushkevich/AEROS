package rbac

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthorizationService interface {
	AddUserToRbac(id uuid.UUID)
	IsAuthenticated(sub string, obj string, act string) (bool, error)
	CreateRole(name, description string) (*Role, error)
	DeleteRole(roleID uint) error
	CreateAction(name string) (*Action, error)
	CreateResource(name, description string) (*Resource, error)
	CreatePermission(resourceID, actionID uint) (*Permission, error)
	GrantPermissionToRole(roleID, permissionID uint) error
	RevokePermissionFromRole(roleID, permissionID uint) error
	AssignRoleToUser(userID uuid.UUID, roleID uint) error
	RemoveRoleFromUser(userID uuid.UUID, roleID uint) error
}

type RBACService struct {
	db       *gorm.DB
	enforcer *casbin.Enforcer
	logger   *slog.Logger
}

func NewRBACService(configPath *string, logger *slog.Logger) *RBACService {
	cfg, err := LoadConfig(*configPath)
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

	// if err := db.AutoMigrate(
	// 	&Role{},
	// 	&Action{},
	// 	&Resource{},
	// 	&Permission{},
	// 	&RolePermission{},
	// 	&UserRole{},
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	return &RBACService{
		db:       db,
		enforcer: e,
		logger:   logger,
	}
}

func (cs *RBACService) IsAuthenticated(sub string, obj string, act string) (bool, error) {
	ok, reason, err := cs.enforcer.EnforceEx(sub, obj, act)
	if err != nil {
		cs.logger.Error("Error occurred while enforcing policy", "err", err)
		return ok, err
	}

	cs.logger.Debug("Enforce result", "sub", sub, "obj", obj, "act", act, "ok", ok, "reason", reason)
	return ok, nil
}

func (cs *RBACService) AddUserToRbac(id uuid.UUID) {
	cs.enforcer.AddGroupingPolicy(id.String(), "User")
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

func (s *RBACService) DeleteRole(roleID uint) error {
	// Delete role permissions
	if err := s.db.Delete(&RolePermission{}, "role_id = ?", roleID).Error; err != nil {
		return err
	}

	// Delete role
	if err := s.db.Delete(&Role{}, "id = ?", roleID).Error; err != nil {
		return err
	}

	// Update Casbin policies
	s.enforcer.RemoveFilteredPolicy(0, strconv.FormatUint(uint64(roleID), 10))
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) CreateAction(name string) (*Action, error) {
	action := &Action{
		Name: name,
	}

	if err := s.db.Create(action).Error; err != nil {
		return nil, err
	}

	return action, nil
}

func (s *RBACService) CreateResource(name, description string) (*Resource, error) {
	resource := &Resource{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(resource).Error; err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *RBACService) CreatePermission(resourceID, actionID uint) (*Permission, error) {
	perm := &Permission{
		ResourceID: resourceID,
		ActionID:   actionID,
	}

	if err := s.db.Create(perm).Error; err != nil {
		return nil, err
	}

	return perm, nil
}

func (s *RBACService) GrantPermissionToRole(roleID, permissionID uint) error {
	var perm Permission
	if err := s.db.First(&perm, "id = ?", permissionID).Error; err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	rp := &RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	if err := s.db.Create(rp).Error; err != nil {
		return err
	}

	// Add to Casbin: sub, obj, act
	s.enforcer.AddPolicy(
		strconv.FormatUint(uint64(roleID), 10),
		strconv.FormatUint(uint64(perm.ResourceID), 10),
		strconv.FormatUint(uint64(perm.ActionID), 10),
	)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) RevokePermissionFromRole(roleID, permissionID uint) error {
	var perm Permission
	if err := s.db.First(&perm, "id = ?", permissionID).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&RolePermission{},
		"role_id = ? AND permission_id = ?", roleID, permissionID).Error; err != nil {
		return err
	}

	// Delete from Casbin
	s.enforcer.RemovePolicy(
		strconv.FormatUint(uint64(roleID), 10),
		strconv.FormatUint(uint64(perm.ResourceID), 10),
		strconv.FormatUint(uint64(perm.ActionID), 10),
	)
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) AssignRoleToUser(userID uuid.UUID, roleID uint) error {
	ur := &UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := s.db.Create(ur).Error; err != nil {
		return err
	}

	// user_id has role_id
	s.enforcer.AddGroupingPolicy(userID.String(), strconv.FormatUint(uint64(roleID), 10))
	s.enforcer.SavePolicy()

	return nil
}

func (s *RBACService) RemoveRoleFromUser(userID uuid.UUID, roleID uint) error {
	if err := s.db.Delete(&UserRole{},
		"user_id = ? AND role_id = ?", userID, roleID).Error; err != nil {
		return err
	}

	s.enforcer.RemoveGroupingPolicy(userID.String(), strconv.FormatUint(uint64(roleID), 10))
	s.enforcer.SavePolicy()

	return nil
}
