package rbac

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type AuthorizationService interface {
	AddUserToRole(id string, role Role)
	IsAuthenticated(sub string, obj string, act Action) (bool, error)
	AddPermission(role Role, obj string, act Action)
}

type CasbinService struct {
	config   *CasbinConfig
	enforcer *casbin.Enforcer
	logger   *slog.Logger
}

func (cs *CasbinService) AddUserToRole(id string, role Role) {
	cs.enforcer.AddGroupingPolicy(id, string(role))
}

func (cs *CasbinService) IsAuthenticated(sub string, obj string, act Action) (bool, error) {
	ok, reason, err := cs.enforcer.EnforceEx(sub, obj, string(act))
	if err != nil {
		cs.logger.Error("Error occurred while enforcing policy", "err", err)
		return ok, err
	}

	cs.logger.Debug("Enforce result", "sub", sub, "obj", obj, "act", act, "ok", ok, "reason", reason)
	return ok, nil
}

func (cs *CasbinService) AddPermission(role Role, obj string, act Action) {
	_, _ = cs.enforcer.AddPolicy(string(role), obj, string(act))
}

func NewCasbinService(configPath *string, logger *slog.Logger) *CasbinService {
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

	return &CasbinService{
		logger:   logger,
		enforcer: e,
	}
}
