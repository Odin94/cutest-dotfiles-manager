package cdmcmd

import (
	"fmt"

	"github.com/Odin94/cutest-dotfiles-manager/internal/apply"
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Copy mapped dotfiles from repo to their target paths",
		RunE:  runApply,
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Only show what would be done; do not write files or state")
	return cmd
}

func runApply(cmd *cobra.Command, _ []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found (run from a dotfiles repo or confirm traverse up)")
		return fmt.Errorf("no config")
	}
	result, err := apply.Run(root, dryRun)
	if err != nil {
		ui.PrintError(err.Error())
		return err
	}
	for _, w := range result.Warnings {
		ui.PrintWarn(w)
	}
	if len(result.Errors) > 0 {
		ui.PrintError("Summary of errors:")
		for _, e := range result.Errors {
			if e.Source != "" || e.Target != "" {
				ui.PrintError(fmt.Sprintf("  %s -> %s: %v", e.Source, e.Target, e.Err))
			} else {
				ui.PrintError(fmt.Sprintf("  %v", e.Err))
			}
		}
		return fmt.Errorf("%d error(s)", len(result.Errors))
	}
	if !dryRun && result.Applied > 0 {
		ui.PrintSuccess(fmt.Sprintf("Applied %d file(s)", result.Applied))
	}
	if result.Backups > 0 && !dryRun {
		ui.PrintWarn(fmt.Sprintf("%d target(s) were backed up to .cdm/temp/", result.Backups))
	}
	return nil
}
