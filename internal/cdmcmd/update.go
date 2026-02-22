package cdmcmd

import (
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/Odin94/cutest-dotfiles-manager/internal/update"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Download the latest release from GitHub and save as cdm.new (replace current binary to finish)",
		RunE:  runUpdate,
	}
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if err := update.Run(); err != nil {
		ui.PrintError(err.Error())
		return err
	}
	return nil
}
