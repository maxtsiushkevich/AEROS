package rbac

import (
	"os"

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
		return CasbinConfig{}, err
	}
	defer file.Close()

	config := CasbinConfig{}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return CasbinConfig{}, err
	}

	return config, nil
}
