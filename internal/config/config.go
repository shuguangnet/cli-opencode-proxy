package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig              `yaml:"server"`
	Upstream   UpstreamConfig            `yaml:"upstream"`
	Models     map[string]string         `yaml:"models"`
	Accounts   map[string]AccountConfig  `yaml:"accounts"`
	Keys       map[string]KeyConfig      `yaml:"keys"`
	Mapping    MappingConfig             `yaml:"mapping"`
	RateLimit  RateLimitConfig           `yaml:"rate_limit"`
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type UpstreamConfig struct {
	Binary  string        `yaml:"binary"`
	Attach  string        `yaml:"attach"`
	Timeout time.Duration `yaml:"timeout"`
}

type AccountConfig struct {
	AuthMode  string `yaml:"auth_mode"`
	Token     string `yaml:"token"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type KeyConfig struct {
	Account       string   `yaml:"account"`
	AllowedModels []string `yaml:"allowed_models"`
}

type MappingConfig struct {
	Temperature TemperatureConfig `yaml:"temperature"`
}

type TemperatureConfig struct {
	TargetMin float64 `yaml:"target_min"`
	TargetMax float64 `yaml:"target_max"`
}

type RateLimitConfig struct {
	Enabled bool `yaml:"enabled"`
	RPM     int  `yaml:"rpm"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	expanded := os.Expand(string(data), func(key string) string {
		return os.Getenv(key)
	})

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.ReadTimeout == 0 {
		c.Server.ReadTimeout = 15 * time.Second
	}
	if c.Upstream.Binary == "" {
		c.Upstream.Binary = "opencode"
	}
	if c.Upstream.Timeout == 0 {
		c.Upstream.Timeout = 120 * time.Second
	}
	if c.Mapping.Temperature.TargetMax == 0 && c.Mapping.Temperature.TargetMin == 0 {
		c.Mapping.Temperature.TargetMax = 1
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Upstream.Binary) == "" {
		return fmt.Errorf("upstream.binary is required")
	}
	if len(c.Models) == 0 {
		return fmt.Errorf("at least one model mapping is required")
	}
	if len(c.Keys) == 0 {
		return fmt.Errorf("at least one gateway key is required")
	}
	for key, keyCfg := range c.Keys {
		if keyCfg.Account == "" {
			return fmt.Errorf("key %s missing account", key)
		}
		if _, ok := c.Accounts[keyCfg.Account]; !ok {
			return fmt.Errorf("key %s references unknown account %s", key, keyCfg.Account)
		}
	}
	return nil
}
