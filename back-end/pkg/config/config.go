package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB    *Database `yaml:"database,omitempty"`
	JWT   *JWT      `yaml:"JWT,omitempty"`
	REDIS *Redis    `yaml:"redis,omitempty"`
}

// Load reads the config file and returns a Config struct
func Load(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	cfg := new(Config)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config file: %w", err)
	}
	return cfg, nil
}

type Database struct {
	Host     string `yaml:"host,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Port     string `yaml:"port,omitempty"`
}

// JWT configuration
type JWT struct {
	Secret     string `yaml:"secret,omitempty"`
	Expiration int    `yaml:"expiration,omitempty"`
}

// Redis redis struct
type Redis struct {
	Addr   string `yaml:"host,omitempty"`
	Pass   string `yaml:"password,omitempty"`
	Expire int    `yaml:"expiration,omitempty"`
}