package cdmcmd

import (
	"github.com/Odin94/cutest-dotfiles-manager/internal/diff"
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/spf13/cobra"
)

func diffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Show diff between source and target files for each mapping",
		RunE:  runDiff,
	}
}

func runDiff(cmd *cobra.Command, args []string) error {
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found")
		return errNoConfig
	}
	return diff.Run(root)
}
