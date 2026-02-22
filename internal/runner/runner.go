package runner

import (
	"os"
	"path/filepath"

	"github.com/cutest-tools/cutest-dotfiles-manager/internal/config"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/ui"
)

const maxTraverseLevels = 5

type Options struct {
	TraverseUpPrompt bool
}

func GetConfigRoot(opts Options) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	if config.ExistsInDir(cwd) {
		return cwd, true
	}
	root, found := config.FindConfigDir(cwd, false, 0)
	if found {
		return root, true
	}
	if !opts.TraverseUpPrompt {
		return "", false
	}
	ok, err := ui.ConfirmTraverseUp()
	if err != nil || !ok {
		return "", false
	}
	root, found = config.FindConfigDir(cwd, true, maxTraverseLevels)
	return root, found
}

func AbsPath(root, rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(root, rel)
}
