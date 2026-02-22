# Installing cdm

## Option 1: From a GitHub release (recommended)

1. Open [Releases](https://github.com/Odin94/cutest-dotfiles-manager/releases) and download the archive or binary for your OS and architecture (e.g. `cdm-v1.0.0-windows-amd64.exe`).
2. Put the binary in a directory that is on your `PATH`:
   - **Windows**: e.g. `C:\Program Files\cdm\` or a folder you added to `Path` in system environment variables.
   - **macOS / Linux**: e.g. `/usr/local/bin` or `~/.local/bin` (then run `chmod +x cdm`).

## Option 2: Go install (from source)

Ensure Go is installed and `$GOPATH/bin` (or `$GOBIN`) is on your `PATH`, then:

```bash
go install github.com/Odin94/cutest-dotfiles-manager/cmd/cdm@latest
```

The binary will be named `cdm` (or `cdm.exe` on Windows) and placed in `$GOPATH/bin` (or `$GOBIN`).

## Option 3: Build from source

```bash
git clone https://github.com/Odin94/cutest-dotfiles-manager.git
cd cutest-dotfiles-manager
go build -o cdm ./cmd/cdm
# Move 'cdm' (or cdm.exe on Windows) to a directory on your PATH
```

## Updating

If you installed from a release or from source, you can update to the latest release with:

```bash
cdm update
```

This downloads the latest release from GitHub and saves the new binary next to your current `cdm` executable (e.g. `cdm.new` or `cdm.exe.new`). Replace your current binary with it when convenient (e.g. after closing any terminal using `cdm`).

---

# Creating GitHub releases

## Manual release

1. Tag a version: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. On GitHub: **Releases** → **Draft a new release**, choose the tag, add notes, upload built binaries (see asset names below).

## Automated releases (GitHub Actions)

The repo includes a workflow that runs when you push a tag like `v*` (e.g. `v1.0.0`):

1. **Tag and push:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The workflow (`.github/workflows/release.yml`) will:
   - Build binaries for Windows (amd64, arm64), Linux (amd64, arm64), and macOS (amd64, arm64)
   - Create a GitHub release for that tag and attach the binaries as assets

3. Asset names follow: `cdm-<version>-<os>-<arch>` with `.exe` for Windows (e.g. `cdm-v1.0.0-windows-amd64.exe`). The `cdm update` command looks for these names to pick the right file for your platform.

## One-time setup for automated releases

- **Write access**: Pushing a tag that matches `v*` triggers the workflow. You need write access to the repo to push tags.
- **Permissions**: The workflow uses `GITHUB_TOKEN` to create the release and upload assets; no extra secrets are required for public repos.
