# cutest-dotfiles-manager (cdm)

A dotfiles manager that keeps a mapping of repo files to config destinations and copies them with `cdm apply`. Built with Go, Cobra, and Fang for a pretty CLI and shell completion.

## Install

```bash
go install github.com/cutest-tools/cutest-dotfiles-manager/cmd/cdm@latest
```

Or build from source: `go build -o cdm ./cmd/cdm`

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
- **`.local.cdm.toml`** (gitignored): variable values (e.g. `HOME = "/Users/me"`). Resolved after env.

## Commands

| Command | Description |
|--------|-------------|
| `cdm apply` | Copy mapped files to targets. Use `-d` / `--dry-run` to preview. |
| `cdm diff` | Show diff between source and target for each mapping. |
| `cdm status` | List source files changed since last apply. |
| `cdm init [repo]` | Create `.cdm/` and config; optionally clone `repo` then init. |
| `cdm add <path> [target-path]` | Copy file into repo and add mapping (relative source path). |
| `cdm watch` | Watch source files and auto-apply (Ctrl+C to stop). |
| `cdm undo-conflict-apply` | Restore targets from `.cdm/temp/` from the last apply. |
| `cdm help` / `cdm -h` | List commands and usage. |
| `cdm completion bash \| zsh \| fish` | Generate shell completions. |

## Shell completion

```bash
cdm completion bash > /path/to/cdm-completion.bash
# source it or install under your shell's completion dir
```

## State

- **`.cdm/`**: `hashes.toml`, `log.txt`, `temp/` (backups), `last_backup.toml` (for undo).

If you run cdm in a directory without `.cdm.toml`, it will ask whether to traverse up (up to 5 levels) to find one.
