package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Backend     string            `yaml:"backend"`
	OnePassword OnePasswordConfig `yaml:"onepassword"`
	Dotenv      DotenvConfig      `yaml:"dotenv"`
}

type OnePasswordConfig struct {
	Vault string `yaml:"vault"`
	Field string `yaml:"field"`
}

type DotenvConfig struct {
	Paths []string `yaml:"paths"`
}

func DefaultConfig() Config {
	return Config{
		Backend: "1password",
		OnePassword: OnePasswordConfig{
			Vault: "CLI",
			Field: "credential",
		},
		Dotenv: DotenvConfig{
			Paths: []string{".env"},
		},
	}
}

// Load reads config from ~/.config/export-key/config.yaml with env var overrides.
func Load() (Config, error) {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not determine home directory: %v\n", err)
		return applyEnvOverrides(cfg), nil
	}

	cfgPath := filepath.Join(home, ".config", "export-key", "config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return applyEnvOverrides(cfg), nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return applyEnvOverrides(cfg), nil
}

func applyEnvOverrides(cfg Config) Config {
	if v := os.Getenv("EK_BACKEND"); v != "" {
		cfg.Backend = v
	}
	if v := os.Getenv("EK_VAULT"); v != "" {
		cfg.OnePassword.Vault = v
	}
	if v := os.Getenv("EK_FIELD"); v != "" {
		cfg.OnePassword.Field = v
	}
	if v := os.Getenv("EK_DOTENV_PATHS"); v != "" {
		cfg.Dotenv.Paths = strings.Split(v, ":")
	}
	return cfg
}
