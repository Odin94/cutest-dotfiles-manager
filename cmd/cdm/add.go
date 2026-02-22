package cdm

import (
	"github.com/Odin94/cutest-dotfiles-manager/internal/addcmd"
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <path> [target-path]",
		Short: "Copy file into repo and add mapping in .cdm.toml",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runAdd,
	}
}

func runAdd(cmd *cobra.Command, args []string) error {
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found (run cdm init first)")
		return errNoConfig
	}
	sourcePath := args[0]
	targetPath := ""
	if len(args) > 1 {
		targetPath = args[1]
	}
	if err := addcmd.Run(root, sourcePath, targetPath); err != nil {
		ui.PrintError(err.Error())
		return err
	}
	ui.PrintSuccess("added " + sourcePath)
	return nil
}
