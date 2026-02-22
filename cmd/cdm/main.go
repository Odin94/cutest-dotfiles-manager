package cdm

import (
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var errNoConfig = errors.New("no config")

func Run() {
	root := &cobra.Command{
		Use:   "cdm",
		Short: "Cutest dotfiles manager – apply, diff, and manage dotfiles from a repo",
		Long:  "cdm keeps a mapping of repo files to config destinations and copies them with 'cdm apply'. Use -h for commands.",
	}
	root.AddCommand(applyCmd())
	root.AddCommand(diffCmd())
	root.AddCommand(statusCmd())
	root.AddCommand(initCmd())
	root.AddCommand(addCmd())
	root.AddCommand(watchCmd())
	root.AddCommand(undoConflictApplyCmd())
	root.AddCommand(updateCmd())

	if err := fang.Execute(context.Background(), root); err != nil {
		os.Exit(1)
	}
}
