package cdmcmd

import (
	"os"
	"os/signal"

	"github.com/Odin94/cutest-dotfiles-manager/internal/apply"
	"github.com/Odin94/cutest-dotfiles-manager/internal/runner"
	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
	"github.com/Odin94/cutest-dotfiles-manager/internal/watch"
	"github.com/spf13/cobra"
)

func watchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Watch source files and auto-apply on change",
		RunE:  runWatch,
	}
}

func runWatch(cmd *cobra.Command, args []string) error {
	root, ok := runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})
	if !ok {
		ui.PrintError("no .cdm.toml found")
		return errNoConfig
	}
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer stop()
	applyFn := func() (*apply.ApplyResult, error) {
		return apply.Run(root, false)
	}
	return watch.Run(ctx, root, applyFn)
}
