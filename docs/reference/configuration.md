---
description: Chief configuration reference. CLI flags and settings for customizing Chief's behavior.
---

# Configuration

Chief is designed to work with zero configuration. All state lives in `.chief/` and settings are passed via CLI flags.

## CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--max-iterations <n>`, `-n` | Loop iteration limit | `10` |
| `--no-sound` | Disable completion sound | `false` |
| `--verbose` | Show raw Claude output in log | `false` |
| `--merge` | Auto-merge progress on conversion conflicts | `false` |
| `--force` | Auto-overwrite on conversion conflicts | `false` |

## Claude Code Configuration

Chief invokes Claude Code under the hood. Claude Code has its own configuration:

```bash
# Authentication
claude login

# Model selection (if you have access)
claude config set model claude-3-opus-20240229
```

See [Claude Code documentation](https://github.com/anthropics/claude-code) for details.

## Permission Handling

By default, Claude Code asks for permission before executing bash commands, writing files, and making network requests. Chief automatically disables these prompts when invoking Claude to enable autonomous operation.

::: warning
Chief runs Claude with full permissions to modify your codebase. Only run Chief on PRDs you trust.
:::

## No Global Config

Intentionally, Chief has no global configuration file. This ensures:

1. **Portability**: Project works the same on any machine
2. **Reproducibility**: No hidden state affecting behavior
3. **Simplicity**: One place to look for all settings
