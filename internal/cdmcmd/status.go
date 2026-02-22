package cdmcmd

import (
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/status"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show which source files have been edited since last apply",
		RunE:  runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found")
		return errNoConfig
	}
	return status.Run(root)
}
