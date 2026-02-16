# export-key

Export API keys from 1Password (or dotenv files) as environment variables with a single command.

```
$ ek
```

An interactive TUI picker lets you fuzzy-search, navigate with arrow keys, or type a number to instantly export any key.

## How it works

A Go binary outputs `export VAR="value"` to stdout. A thin shell function `eval`s it, setting the variable in your current shell. The TUI renders on stderr so it doesn't interfere with the export. Same proven pattern as direnv, fnm, and rbenv.

## Install

### Homebrew

```bash
brew install danmartuszewski/tap/export-key
```

### From source

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

## Usage

```bash
# Interactive picker (fuzzy search, arrow keys, or type a number)
ek

# Export by number (from the list)
ek 3

# Export by name (fuzzy match)
ek openai

# Export all keys for a project
ek .myapp
ek -p myapp  # long form

# List available keys
export-key list
```

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

## License

MIT
