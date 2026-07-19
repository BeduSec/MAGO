// Copyright (c) BeduSec. All rights reserved.
package config

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Store       StoreConfig       `yaml:"store"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
	WAF         WAFConfig         `yaml:"waf"`
	AdminToken  string            `yaml:"admin_token"`
	LogLevel    string            `yaml:"log_level"`
	LogJSON     bool              `yaml:"log_json"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type StoreConfig struct {
	Type     string `yaml:"type"`
	RedisURL string `yaml:"redis_url"`
}

type RateLimiterConfig struct {
	DefaultRate     float64 `yaml:"default_rate"`
	DefaultBurst    int     `yaml:"default_burst"`
	CleanupInterval int     `yaml:"cleanup_interval_sec"`
}

type WAFConfig struct {
	RulesFile string `yaml:"rules_file"`
	DryRun    bool   `yaml:"dry_run"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Store: StoreConfig{
			Type: "memory",
		},
		RateLimiter: RateLimiterConfig{
			DefaultRate:     100.0,
			DefaultBurst:    200,
			CleanupInterval: 60,
		},
		WAF: WAFConfig{
			RulesFile: "rules.json",
			DryRun:    false,
		},
		LogLevel: "info",
		LogJSON:  true,
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("MAGO_SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("MAGO_SERVER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := os.Getenv("MAGO_STORE_TYPE"); v != "" {
		cfg.Store.Type = v
	}
	if v := os.Getenv("MAGO_REDIS_URL"); v != "" {
		cfg.Store.RedisURL = v
	}
	if v := os.Getenv("MAGO_RATE_DEFAULT"); v != "" {
		if r, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.RateLimiter.DefaultRate = r
		}
	}
	if v := os.Getenv("MAGO_RATE_BURST"); v != "" {
		if b, err := strconv.Atoi(v); err == nil {
			cfg.RateLimiter.DefaultBurst = b
		}
	}
	if v := os.Getenv("MAGO_WAF_RULES_FILE"); v != "" {
		cfg.WAF.RulesFile = v
	}
	if v := os.Getenv("MAGO_WAF_DRY_RUN"); v != "" {
		cfg.WAF.DryRun = strings.ToLower(v) == "true"
	}
	if v := os.Getenv("MAGO_ADMIN_TOKEN"); v != "" {
		cfg.AdminToken = v
	}
	if v := os.Getenv("MAGO_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("MAGO_LOG_JSON"); v != "" {
		cfg.LogJSON = strings.ToLower(v) == "true"
	}
}