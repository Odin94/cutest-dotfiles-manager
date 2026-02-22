package cdm

import (
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/Odin94/cutest-dotfiles-manager/internal/undo"
	"github.com/spf13/cobra"
)

func undoConflictApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undo-conflict-apply",
		Short: "Restore target files from .cdm/temp/ from the latest apply",
		RunE:  runUndoConflictApply,
	}
}

func runUndoConflictApply(cmd *cobra.Command, args []string) error {
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found")
		return errNoConfig
	}
	return undo.Run(root)
}
