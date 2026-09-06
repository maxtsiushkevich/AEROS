package rbac

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type CasbinConfig struct {
	Host       string `yaml:"casbin_db_host"`
	Port       int    `yaml:"casbin_db_port"`
	User       string `yaml:"casbin_db_user"`
	Password   string `yaml:"casbin_db_pass"`
	DbName     string `yaml:"casbin_db_name"`
	ConfigPath string `yaml:"casbin_config_path"`
}

func LoadConfig(filePath string) (CasbinConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return CasbinConfig{}, fmt.Errorf("open RBAC config %q: %w", filePath, err)
	}
	defer file.Close()

	config := CasbinConfig{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return CasbinConfig{}, fmt.Errorf("decode RBAC config %q: %w", filePath, err)
	}

	if config.Host == "" || config.User == "" || config.Password == "" || config.DbName == "" || config.ConfigPath == "" {
		return CasbinConfig{}, fmt.Errorf("RBAC config file %q is incomplete", filePath)
	}
	if config.Port <= 0 {
		return CasbinConfig{}, fmt.Errorf("RBAC config file %q has invalid casbin_db_port", filePath)
	}

	return config, nil
}

func LoadConfigFromEnv() (CasbinConfig, error) {
	filePath := strings.TrimSpace(os.Getenv("RBAC_CONFIG_PATH"))
	if filePath == "" {
		return CasbinConfig{}, fmt.Errorf("RBAC_CONFIG_PATH is required")
	}

	return LoadConfig(filePath)
}
