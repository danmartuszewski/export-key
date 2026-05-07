---
name: export-key
description: Load API keys and secrets from 1Password (or .env files) into the current shell using the `ek` command. Use whenever the user asks you to run a command that needs a secret, token, or API key, or when a command fails because an environment variable is unset. Avoids exposing secret values in chat.
---

# export-key (`ek`)

`ek` exports a secret as an environment variable in the current shell, sourced from 1Password or a configured `.env` file. The user installs and configures it once; you only invoke it.

## When to use

Trigger this skill when:

- The user asks you to run a command that needs a secret, token, or API key.
- A command fails with a missing-env-var or unauthorized error (e.g. `OPENAI_API_KEY not set`, `401 Unauthorized`).
- The user says something like "use my X key" or "load the Y token".

Use `ek` **before** asking the user to paste a value. Never request a secret value in chat.

## Commands

```bash
# Export a single key by env-var name
# Matches an exact 1Password item title or the BASE-* pattern
ek OPENAI_API_KEY

# Export every key tagged with a project
ek .myproject

# Open the interactive picker (when the name is unknown)
ek

# List available keys without exporting anything
export-key list
```

## Verifying without leaking

Never `echo` or `printenv` a secret directly. Confirm presence only:

```bash
[ -n "$OPENAI_API_KEY" ] && echo "OPENAI_API_KEY is set" || echo "missing"
```

## Same-session rule

`ek` modifies the **current** shell. Run dependent commands in the same shell invocation as `ek`, not in a fresh subshell:

```bash
# Correct — same shell
ek OPENAI_API_KEY && curl -H "Authorization: Bearer $OPENAI_API_KEY" https://api.example.com/v1/whoami
```

If your tooling spawns a fresh shell per command, chain commands with `&&` or run them as a single shell invocation.

## If `ek` is not found

Shell integration isn't loaded. Tell the user — do not modify their shell config without permission. Setup instructions: <https://github.com/danmartuszewski/export-key>.

## If a key is not found

`ek` exits non-zero and prints a message to stderr. Don't guess at alternative names — tell the user the key wasn't found and ask whether they want to add it to their backend.
