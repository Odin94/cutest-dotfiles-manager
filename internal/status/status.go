package status

import (
	"fmt"

	"github.com/cutest-tools/cutest-dotfiles-manager/internal/config"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/runner"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/state"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/ui"
)

func Run(root string) error {
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	local, _ := config.LoadLocal(root)
	resolved, missingVars := cfg.ResolvedMappings(local)
	if len(missingVars) > 0 {
		missing := cfg.MissingVariables(local)
		ui.PrintMissingVarsWarning(missing, missingVars)
	}
	hashes, err := state.LoadHashes(root)
	if err != nil {
		return err
	}
	for srcRel, targetPath := range resolved {
		srcPath := runner.AbsPath(root, srcRel)
		srcHash, err := state.FileHash(srcPath)
		if err != nil {
			ui.PrintWarn(srcRel + ": " + err.Error())
			continue
		}
		stored, ok := hashes[targetPath]
		if !ok {
			fmt.Println(ui.AnsiYellow + "? " + srcRel + " (not yet applied)" + ui.AnsiReset)
			continue
		}
		if srcHash != stored {
			fmt.Println(ui.AnsiYellow + "M " + srcRel + ui.AnsiReset)
		}
	}
	return nil
}
