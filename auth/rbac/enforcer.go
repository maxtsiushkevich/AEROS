package rbac

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type AuthorizationService interface {
	CanAccess(string, string, string) bool
}

type CasbinService struct {
	config   *CasbinConfig
	Enforcer *casbin.Enforcer
}

func (cs *CasbinService) CanAccess(id string, path string, action string) bool {
	return true
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
		Enforcer: e,
	}
}
