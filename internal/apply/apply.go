package apply

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cutest-tools/cutest-dotfiles-manager/internal/config"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/runner"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/state"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/ui"
)

type ApplyResult struct {
	Applied  int
	Skipped  int
	Backups  int
	Errors   []ApplyError
	Warnings []string
}

type ApplyError struct {
	Source string
	Target string
	Err    error
}

func Run(root string, dryRun bool) (*ApplyResult, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return nil, err
	}
	local, err := config.LoadLocal(root)
	if err != nil {
		return nil, err
	}
	resolved, missingVars := cfg.ResolvedMappings(local)
	if len(missingVars) > 0 {
		missing := cfg.MissingVariables(local)
		ui.PrintMissingVarsWarning(missing, missingVars)
	}

	if dryRun {
		ui.PrintDryRunBanner()
	}

	hashes, err := state.LoadHashes(root)
	if err != nil {
		return nil, err
	}

	result := &ApplyResult{}

	if !dryRun {
		for _, script := range cfg.Scripts.PreApply {
			scriptPath := runner.AbsPath(root, script)
			if _, err := os.Stat(scriptPath); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("pre_apply script missing: %s", scriptPath))
				_ = state.AppendLog(root, fmt.Sprintf("pre_apply script missing: %s", scriptPath))
				continue
			}
			// Run script (simplified: exec; would use exec.Command)
			_ = runScript(scriptPath, root)
		}
	}

	backupManifest := make(map[string]string)

	for srcRel, targetPath := range resolved {
		srcPath := runner.AbsPath(root, srcRel)
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			result.Errors = append(result.Errors, ApplyError{Source: srcRel, Target: targetPath, Err: err})
			continue
		}
		sourceHash := state.ContentHash(srcData)

		targetDir := filepath.Dir(targetPath)
		if _, err := os.Stat(targetDir); err != nil {
			if os.IsNotExist(err) {
				if dryRun {
					result.Warnings = append(result.Warnings, fmt.Sprintf("would create directory: %s", targetDir))
				} else {
					if mkErr := os.MkdirAll(targetDir, 0755); mkErr != nil {
						result.Errors = append(result.Errors, ApplyError{Source: srcRel, Target: targetPath, Err: mkErr})
						continue
					}
					result.Warnings = append(result.Warnings, fmt.Sprintf("created directory: %s", targetDir))
					_ = state.AppendLog(root, fmt.Sprintf("created directory: %s", targetDir))
				}
			} else {
				result.Errors = append(result.Errors, ApplyError{Source: srcRel, Target: targetPath, Err: err})
				continue
			}
		}

		existingHash, hadHash := hashes[targetPath]
		targetExists := false
		var currentTargetHash string
		if st, err := os.Stat(targetPath); err == nil && st.Mode().IsRegular() {
			targetExists = true
			currentTargetHash, _ = state.FileHash(targetPath)
		}

		if targetExists && hadHash && currentTargetHash != existingHash {
			if dryRun {
				result.Warnings = append(result.Warnings, fmt.Sprintf("would backup (target changed): %s -> .cdm/temp/...", targetPath))
				result.Backups++
			} else {
				_ = state.EnsureTempDir(root)
				backupPath := state.BackupPath(root, targetPath)
				targetData, _ := os.ReadFile(targetPath)
				if writeErr := os.WriteFile(backupPath, targetData, 0644); writeErr != nil {
					result.Errors = append(result.Errors, ApplyError{Source: srcRel, Target: targetPath, Err: writeErr})
					continue
				}
				backupManifest[backupPath] = targetPath
				result.Backups++
				ui.PrintWarn(fmt.Sprintf("target changed; backed up to %s", backupPath))
				_ = state.AppendLog(root, fmt.Sprintf("backup: %s -> %s", targetPath, backupPath))
			}
		}

		if dryRun {
			result.Warnings = append(result.Warnings, fmt.Sprintf("would copy: %s -> %s", srcRel, targetPath))
			result.Applied++
			continue
		}

		if writeErr := os.WriteFile(targetPath, srcData, 0644); writeErr != nil {
			result.Errors = append(result.Errors, ApplyError{Source: srcRel, Target: targetPath, Err: writeErr})
			continue
		}
		hashes[targetPath] = sourceHash
		result.Applied++
		_ = state.AppendLog(root, fmt.Sprintf("apply: %s -> %s", srcRel, targetPath))
	}

	if !dryRun {
		if err := state.SaveHashes(root, hashes); err != nil {
			result.Errors = append(result.Errors, ApplyError{Err: err})
		}
		if len(backupManifest) > 0 {
			_ = state.WriteLastBackupManifest(root, backupManifest)
		}
		for _, script := range cfg.Scripts.PostApply {
			scriptPath := runner.AbsPath(root, script)
			if _, err := os.Stat(scriptPath); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("post_apply script missing: %s", scriptPath))
				continue
			}
			_ = runScript(scriptPath, root)
		}
	}

	return result, nil
}

func runScript(scriptPath, workDir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", scriptPath)
	} else {
		cmd = exec.Command("sh", scriptPath)
	}
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
