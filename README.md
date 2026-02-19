<p align="center">
  <img src="assets/icon.png" height="128">
</p>

<h1 align="center">export-key</h1>

<p align="center">
  Stop hunting for API keys. Just <code>ek</code> and they're in your shell.
</p>

## Why export-key?

```bash
# Before: copy-pasting from 1Password, managing .env files
op item get "OPENAI_API_KEY" --vault CLI --fields credential | pbcopy
export OPENAI_API_KEY="sk-..."

# After
ek openai
```

```bash
ek                           # interactive fuzzy picker
ek 3                         # export by number
ek openai                    # fuzzy match by name
ek .myapp                    # export all keys for a project
export-key list              # list available keys
```

## How it works

A Go binary outputs `export VAR="value"` to stdout. A thin shell function `eval`s it, setting the variable in your current shell. The TUI renders on stderr so it doesn't interfere with the export. Same pattern used by direnv, fnm, and rbenv.

## Install

### Homebrew (macOS/Linux)

```bash
brew install danmartuszewski/tap/export-key
```

### Go

```bash
go install github.com/danmartuszewski/export-key/cmd/export-key@latest
```

### Shell integration

Add to your `~/.zshrc` (or `~/.bashrc`, `~/.config/fish/config.fish`):

```bash
eval "$(export-key init zsh)"
```

Or if you prefer a minimal one-liner instead:

```bash
ek() { eval "$(export-key select "$@")"; }
```

Both define the `ek` function that wraps `export-key select`. The `init` variant stays in sync with binary updates automatically.

## Key naming convention

Item titles in 1Password follow this pattern:

```
ENV_VAR_NAME           -> env var: ENV_VAR_NAME, no project
ENV_VAR_NAME-project   -> env var: ENV_VAR_NAME, project: project
```

Everything before the first `-` becomes the environment variable name. Everything after is the project tag.

### Project grouping

Keys sharing the same suffix belong to the same project:

```
OPENAI_API_KEY-myapp   -> project: myapp
STRIPE_KEY-myapp       -> project: myapp
AWS_ACCESS_KEY         -> no project
```

`ek .myapp` (or `ek -p myapp`) exports all keys tagged with that project at once.

## Security

Common ways secrets get leaked and how export-key handles them:

| Risk | Traditional approach | export-key |
|---|---|---|
| **Secrets in dotfiles** | `export KEY="sk-..."` hardcoded in `.bashrc` or `.zshrc` | Keys are fetched on demand, never written to shell config |
| **Secrets in shell history** | `export KEY="sk-..."` saved in `.zsh_history` | Only `ek openai` appears in history, never the secret value |
| **Secrets in clipboard** | Copy-paste from 1Password leaves secrets in clipboard history | Secrets go directly from the vault to your shell, no clipboard involved |
| **Secrets committed to git** | `.env` files accidentally checked in | With 1Password backend, there are no local files to commit |
| **Secrets shared across sessions** | `.env` loaded at shell startup persists across all terminals | Keys live only in the shell session where you exported them |
| **Secrets read by AI agents** | Coding agents (Claude Code, Codex, Warp AI) can read `.env` files, shell history, and dotfiles in your project | No local secret files to read - keys exist only as in-memory env vars in your shell session |

Even with the dotenv backend, export-key is safer than `source .env` - keys are loaded selectively (only what you pick) and on demand (not at every shell startup).

### AI coding agents and your secrets

Coding agents like Claude Code, Codex, and Warp AI read your shell configuration files (`.bashrc`, `.zshrc`, `.profile`) as part of their normal operation. If your API keys are in those files, every agent session has direct access to them.

Even with the dotenv backend, storing secrets in a dedicated `.env` file that only export-key references is much better than putting them in `.zshrc`. Your `.env` can live outside the project directory (e.g. `~/.secrets.env`) where agents have no reason to look. No agent knows the file exists or where to find it. You can also block agents from reading it through their configuration (e.g. `.claude/settings.json` for Claude Code or `.codexrc` for Codex).

Agents will always read `.zshrc` because it's a standard file in a known location. A custom `.env` path that only export-key references through its own config is invisible to them by default.

## Backends

### 1Password (default)

Uses the `op` CLI. Make sure you have [1Password CLI](https://developer.1password.com/docs/cli) installed and signed in.

### Dotenv

For users without 1Password. Reads plaintext `.env` files.

> **Note:** Dotenv files store secrets in plaintext. Use 1Password for better security.

## Configuration

Config file at `~/.config/export-key/config.yaml`:

```yaml
backend: 1password
onepassword:
  vault: CLI
  field: credential
dotenv:
  paths: [.env, ~/.secrets.env]
```

### Environment variable overrides

| Variable | Description |
|---|---|
| `EK_BACKEND` | Backend to use (`1password`, `dotenv`) |
| `EK_VAULT` | 1Password vault name |
| `EK_FIELD` | 1Password field name |
| `EK_DOTENV_PATHS` | Colon-separated dotenv file paths |

## Project Structure

```
export-key/
├── cmd/export-key/    # Main entry point
├── internal/
│   ├── backend/       # Backend interface & implementations (1Password, dotenv)
│   ├── cmd/           # CLI commands (cobra)
│   ├── config/        # Configuration loading
│   ├── keyitem/       # Key naming parser (ENV_VAR-project)
│   ├── shell/         # Shell init script generation
│   └── tui/           # Interactive picker (bubbletea)
├── Makefile
└── README.md
```

## License

MIT License - see [LICENSE](LICENSE) for details.
