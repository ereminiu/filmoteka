package config

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"sslmode"`
}

const (
	testPath = "./config/test.yaml"
	prodPath = "./config/prod.yaml"
)

func getPath() string {
	if os.Getenv("MODE") == "prod" {
		return prodPath
	}
	return testPath
}

func LoadConfigs() (*Config, error) {
	path := getPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) GetApiAdress() string {
	return "localhost:" + cfg.Port
}
