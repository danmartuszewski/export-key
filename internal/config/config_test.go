package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Backend != "1password" {
		t.Errorf("Backend = %q, want %q", cfg.Backend, "1password")
	}
	if cfg.OnePassword.Vault != "CLI" {
		t.Errorf("Vault = %q, want %q", cfg.OnePassword.Vault, "CLI")
	}
	if cfg.OnePassword.Field != "credential" {
		t.Errorf("Field = %q, want %q", cfg.OnePassword.Field, "credential")
	}
	if len(cfg.Dotenv.Paths) != 1 || cfg.Dotenv.Paths[0] != ".env" {
		t.Errorf("Dotenv.Paths = %v, want [.env]", cfg.Dotenv.Paths)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	cfg := DefaultConfig()

	t.Run("EK_BACKEND", func(t *testing.T) {
		t.Setenv("EK_BACKEND", "dotenv")
		result := applyEnvOverrides(cfg)
		if result.Backend != "dotenv" {
			t.Errorf("Backend = %q, want %q", result.Backend, "dotenv")
		}
	})

	t.Run("EK_VAULT", func(t *testing.T) {
		t.Setenv("EK_VAULT", "Personal")
		result := applyEnvOverrides(cfg)
		if result.OnePassword.Vault != "Personal" {
			t.Errorf("Vault = %q, want %q", result.OnePassword.Vault, "Personal")
		}
	})

	t.Run("EK_FIELD", func(t *testing.T) {
		t.Setenv("EK_FIELD", "password")
		result := applyEnvOverrides(cfg)
		if result.OnePassword.Field != "password" {
			t.Errorf("Field = %q, want %q", result.OnePassword.Field, "password")
		}
	})

	t.Run("EK_DOTENV_PATHS", func(t *testing.T) {
		t.Setenv("EK_DOTENV_PATHS", "/a/.env:/b/.env")
		result := applyEnvOverrides(cfg)
		if len(result.Dotenv.Paths) != 2 {
			t.Fatalf("Paths len = %d, want 2", len(result.Dotenv.Paths))
		}
		if result.Dotenv.Paths[0] != "/a/.env" || result.Dotenv.Paths[1] != "/b/.env" {
			t.Errorf("Paths = %v", result.Dotenv.Paths)
		}
	})

	t.Run("unset vars use defaults", func(t *testing.T) {
		result := applyEnvOverrides(cfg)
		if result.Backend != "1password" {
			t.Errorf("Backend = %q, want default", result.Backend)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".config", "export-key")
	os.MkdirAll(cfgDir, 0755)

	cfgContent := `backend: dotenv
onepassword:
  vault: Personal
  field: password
dotenv:
  paths: [/secrets/.env]
`
	os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(cfgContent), 0644)

	// Point HOME to temp dir so Load() finds our config
	t.Setenv("HOME", dir)
	// Clear any env overrides
	t.Setenv("EK_BACKEND", "")
	t.Setenv("EK_VAULT", "")
	t.Setenv("EK_FIELD", "")
	t.Setenv("EK_DOTENV_PATHS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Backend != "dotenv" {
		t.Errorf("Backend = %q, want %q", cfg.Backend, "dotenv")
	}
	if cfg.OnePassword.Vault != "Personal" {
		t.Errorf("Vault = %q, want %q", cfg.OnePassword.Vault, "Personal")
	}
	if cfg.OnePassword.Field != "password" {
		t.Errorf("Field = %q, want %q", cfg.OnePassword.Field, "password")
	}
	if len(cfg.Dotenv.Paths) != 1 || cfg.Dotenv.Paths[0] != "/secrets/.env" {
		t.Errorf("Paths = %v", cfg.Dotenv.Paths)
	}
}

func TestLoadEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".config", "export-key")
	os.MkdirAll(cfgDir, 0755)

	// Config file says 1password
	os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("backend: 1password\n"), 0644)

	t.Setenv("HOME", dir)
	// Env var overrides to dotenv
	t.Setenv("EK_BACKEND", "dotenv")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Backend != "dotenv" {
		t.Errorf("Backend = %q, want %q (env should override file)", cfg.Backend, "dotenv")
	}
}

func TestLoadNoConfigFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("EK_BACKEND", "")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	// Should return defaults
	if cfg.Backend != "1password" {
		t.Errorf("Backend = %q, want default %q", cfg.Backend, "1password")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".config", "export-key")
	os.MkdirAll(cfgDir, 0755)

	os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte("{{invalid yaml"), 0644)

	t.Setenv("HOME", dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
