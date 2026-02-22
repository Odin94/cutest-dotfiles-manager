package cdmcmd

import (
	"os"

	"github.com/Odin94/cutest-dotfiles-manager/internal/initcmd"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [repo]",
		Short: "Create .cdm dir and config; optionally clone repo and init there",
		RunE:  runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repo := ""
	if len(args) > 0 {
		repo = args[0]
	}
	return initcmd.Run(cwd, repo)
}
