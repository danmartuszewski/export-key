package backend

import (
	"fmt"

	"github.com/danmartuszewski/export-key/internal/config"
)

// New creates a Backend from the given config.
func New(cfg config.Config) (Backend, error) {
	switch cfg.Backend {
	case "1password", "onepassword", "op":
		return NewOnePassword(cfg.OnePassword.Vault, cfg.OnePassword.Field), nil
	case "dotenv", "env":
		return NewDotenv(cfg.Dotenv.Paths), nil
	default:
		return nil, fmt.Errorf("unknown backend: %q (supported: 1password, dotenv)", cfg.Backend)
	}
}
