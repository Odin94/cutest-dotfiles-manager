package addcmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Odin94/cutest-dotfiles-manager/internal/config"
	"github.com/pelletier/go-toml/v2"
)

func Run(root, sourcePath, targetDest string) error {
	absSrc, err := resolveForRead(sourcePath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(absSrc)
	if err != nil {
		return err
	}
	destRel := filepath.ToSlash(filepath.Base(absSrc))
	if targetDest != "" && !filepath.IsAbs(targetDest) && !strings.HasPrefix(targetDest, "$") {
		destRel = filepath.ToSlash(filepath.Clean(targetDest))
	}
	destAbs := filepath.Join(root, filepath.Clean(destRel))
	destDir := filepath.Dir(destAbs)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(destAbs, data, 0644); err != nil {
		return err
	}
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	if cfg.Mappings == nil {
		cfg.Mappings = make(map[string]string)
	}
	cfg.Mappings[destRel] = sourcePath
	return writeConfig(root, cfg)
}

func resolveForRead(path string) (string, error) {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			return home, nil
		}
		path = filepath.Join(home, path[2:])
	}
	return filepath.Abs(path)
}

func writeConfig(root string, cfg *config.Config) error {
	path := filepath.Join(root, config.ConfigFilename)
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
