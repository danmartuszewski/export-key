package backend

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Dotenv struct {
	Paths []string
	// entries caches parsed KEY=VALUE pairs from all files.
	entries map[string]string
	// order preserves insertion order for stable numbering.
	order []string
}

func NewDotenv(paths []string) *Dotenv {
	return &Dotenv{Paths: paths}
}

func (d *Dotenv) ListItems() ([]string, error) {
	if err := d.load(); err != nil {
		return nil, err
	}

	titles := make([]string, len(d.order))
	copy(titles, d.order)
	return titles, nil
}

func (d *Dotenv) GetSecret(title string) (string, error) {
	if err := d.load(); err != nil {
		return "", err
	}

	val, ok := d.entries[title]
	if !ok {
		return "", fmt.Errorf("key %q not found in dotenv files", title)
	}
	return val, nil
}

func (d *Dotenv) load() error {
	if d.entries != nil {
		return nil
	}

	d.entries = make(map[string]string)

	for _, path := range d.Paths {
		expanded := expandPath(path)
		if err := d.loadFile(expanded); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to read %s: %w", path, err)
		}
	}

	if len(d.entries) == 0 {
		return fmt.Errorf("no keys found in dotenv files: %v", d.Paths)
	}

	return nil
}

func (d *Dotenv) loadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = unquote(val)

		if _, exists := d.entries[key]; !exists {
			d.order = append(d.order, key)
		}
		d.entries[key] = val
	}

	return scanner.Err()
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return home + path[1:]
		}
	}
	return path
}
