package undo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cutest-tools/cutest-dotfiles-manager/internal/state"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/ui"
)

func Run(root string) error {
	manifest, err := state.LoadLastBackupManifest(root)
	if err != nil {
		return err
	}
	if len(manifest) == 0 {
		ui.PrintWarn("no backup manifest found (nothing to undo)")
		return nil
	}
	for tempPath, targetPath := range manifest {
		data, err := os.ReadFile(tempPath)
		if err != nil {
			ui.PrintWarn(tempPath + ": " + err.Error())
			continue
		}
		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			ui.PrintWarn(targetPath + ": " + err.Error())
			continue
		}
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			ui.PrintWarn(targetPath + ": " + err.Error())
			continue
		}
		ui.PrintSuccess("restored " + targetPath)
		_ = state.AppendLog(root, fmt.Sprintf("undo: %s -> %s", tempPath, targetPath))
	}
	_ = state.ClearLastBackupManifest(root)
	return nil
}
