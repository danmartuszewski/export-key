package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const listCacheTTL = 30 * time.Second

type OnePassword struct {
	Vault string
	Field string
}

type opItem struct {
	Title string `json:"title"`
}

func NewOnePassword(vault, field string) *OnePassword {
	return &OnePassword{Vault: vault, Field: field}
}

func (op *OnePassword) ListItems() ([]string, error) {
	// Try cache first
	if titles, ok := op.readCache(); ok {
		return titles, nil
	}

	args := []string{"item", "list", "--format=json"}
	if op.Vault != "" {
		args = append(args, "--vault", op.Vault)
	}

	out, err := exec.Command("op", args...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("op item list failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("op not found; install 1Password CLI: https://developer.1password.com/docs/cli")
	}

	var items []opItem
	if err := json.Unmarshal(out, &items); err != nil {
		return nil, fmt.Errorf("failed to parse op output: %w", err)
	}

	titles := make([]string, len(items))
	for i, item := range items {
		titles[i] = item.Title
	}

	// Write cache (best effort)
	op.writeCache(titles)

	return titles, nil
}

func (op *OnePassword) cachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	vault := op.Vault
	if vault == "" {
		vault = "_default"
	}
	return filepath.Join(dir, "export-key", "items-"+vault+".json")
}

func (op *OnePassword) readCache() ([]string, bool) {
	path := op.cachePath()
	info, err := os.Stat(path)
	if err != nil || time.Since(info.ModTime()) > listCacheTTL {
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var titles []string
	if err := json.Unmarshal(data, &titles); err != nil {
		return nil, false
	}

	return titles, true
}

func (op *OnePassword) writeCache(titles []string) {
	path := op.cachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}

	data, err := json.Marshal(titles)
	if err != nil {
		return
	}

	_ = os.WriteFile(path, data, 0600)
}

func (op *OnePassword) GetSecret(title string) (string, error) {
	args := []string{"item", "get", title, "--fields", op.Field, "--reveal"}
	if op.Vault != "" {
		args = append(args, "--vault", op.Vault)
	}

	out, err := exec.Command("op", args...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("op item get failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("failed to get secret for %q: %w", title, err)
	}

	return strings.TrimSpace(string(out)), nil
}
