package rbac

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type AuthorizationService interface {
	AddUserToRole(id string, role Role)
	IsAuthenticated(sub string, obj string, act string) (bool, error)
}

type CasbinService struct {
	config   *CasbinConfig
	enforcer *casbin.Enforcer
}

func (cs *CasbinService) AddUserToRole(id string, role Role) {
	cs.enforcer.AddGroupingPolicy(id, string(role))
}

func (cs *CasbinService) IsAuthenticated(sub string, obj string, act string) (bool, error) {
	ok, err := cs.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

func NewCasbinService(configPath string) *CasbinService {
	cfg, err := LoadConfig(configPath)
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
		enforcer: e,
	}
}
