package config

import (
	"errors"
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

func LoadConfigs(mode string) (*Config, error) {
	var path string
	if mode == "test" {
		path = testPath
	} else if mode == "prod" {
		path = prodPath
	} else {
		return nil, errors.New("unknown mode")
	}
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
