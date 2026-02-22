package diff

import (
	"fmt"
	"os"

	"github.com/aymanbagabas/go-udiff"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/config"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/runner"
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
	for srcRel, targetPath := range resolved {
		srcPath := runner.AbsPath(root, srcRel)
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			ui.PrintWarn(srcRel + ": " + err.Error())
			continue
		}
		targetData, err := os.ReadFile(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				ui.PrintWarn(targetPath + " (target does not exist)")
			} else {
				ui.PrintWarn(targetPath + ": " + err.Error())
			}
			continue
		}
		d := udiff.Unified(srcRel, targetPath, string(targetData), string(srcData))
		if d != "" {
			fmt.Println("--- " + srcRel + " -> " + targetPath)
			fmt.Println(d)
		}
	}
	return nil
}
