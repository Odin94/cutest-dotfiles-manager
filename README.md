# cutest-dotfiles-manager (cdm)

A dotfiles manager that keeps a mapping of repo files to config destinations and copies them with `cdm apply`. Built with Go, Cobra, and Fang for a pretty CLI and shell completion.

## Install

- **From a release**: [Releases](https://github.com/Odin94/cutest-dotfiles-manager/releases) → download the binary for your OS/arch and put it on your `PATH`.
- **From source**: `go install github.com/Odin94/cutest-dotfiles-manager/cmd/cdm@latest` (ensure `$GOPATH/bin` or `$GOBIN` is on your `PATH`).
- **Build locally**: `go build -o cdm ./cmd/cdm`.

See [docs/INSTALL.md](docs/INSTALL.md) for details and how to create releases.

## Quick start

```bash
cd your-dotfiles-repo
cdm init
cdm add ~/.bashrc
# Edit .cdm.toml if needed (e.g. set variables in .local.cdm.toml)
cdm apply
```

## Config

- **`.cdm.toml`** (committed): `variables`, `[mappings]`, optional `[mappings.windows]` / `[mappings.macos]` / `[mappings.linux]`, and `[scripts]` with `pre_apply` / `post_apply`.
- **`.local.cdm.toml`** (gitignored): variable values (e.g. `HOME = "/Users/me"`). Resolved with higher priority than env.

Example `.cdm.toml` with inline comments:

```toml
# Variable names to substitute in mapping destinations. Values come from
# .local.cdm.toml first, then from the environment.
variables = ["HOME", "XDG_CONFIG_HOME"]

# Default mappings: source path in repo (key) -> destination path (value).
# Use $VAR for substitution. Paths are relative to the repo root.
[mappings]
".bashrc" = "$HOME/.bashrc"
".config/nvim/init.lua" = "$XDG_CONFIG_HOME/nvim/init.lua"

# Optional: OS-specific mappings (merged with [mappings] for the current OS).
# Only the section matching your OS (windows / darwin=macos / linux) is used.
[mappings.windows]
"windows-only.conf" = "$HOME/AppData/Roaming/app/settings.conf"

[mappings.macos]
# [mappings.macos] entries here

[mappings.linux]
# [mappings.linux] entries here

# Scripts run around apply: paths relative to repo root. Missing or failed
# scripts are reported; apply continues.
[scripts]
pre_apply = [".scripts/pre.sh"]
post_apply = [".scripts/post.sh"]
```

## Commands

| Command | Description |
|--------|-------------|
| `cdm init [repo]` | Create `.cdm/` and config; optionally clone `repo` then init. |
| `cdm apply` | Copy mapped files to targets. Use `-d` / `--dry-run` to preview. |
| `cdm diff` | Show diff between source and target for each mapping. |
| `cdm status` | List source files changed since last apply. |
| `cdm add <path> [target-path]` | Copy file into repo and add mapping (relative source path). |
| `cdm watch` | Watch source files and auto-apply (Ctrl+C to stop). |
| `cdm undo-conflict-apply` | Restore targets from `.cdm/temp/` from the last apply. |
| `cdm update` | Download latest release from GitHub; save as `cdm.new` and print replace instructions. |
| `cdm help` / `cdm -h` | List commands and usage. |
| `cdm completion bash \| zsh \| fish` | Generate shell completions. |

## Shell completion

```bash
cdm completion bash > /path/to/cdm-completion.bash   # or cdm completion fish/zsh/...
# source it or install under your shell's completion dir
```

## State

- **`.cdm/`**: `hashes.toml`, `log.txt`, `temp/` (backups), `last_backup.toml` (for undo).

If you run cdm in a directory without `.cdm.toml`, it will ask whether to traverse up (up to 5 levels) to find one.
