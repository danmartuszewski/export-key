package backend

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

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
	return titles, nil
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
