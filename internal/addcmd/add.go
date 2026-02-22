package addcmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Odin94/cutest-dotfiles-manager/internal/config"
	"github.com/pelletier/go-toml/v2"
)

func Run(root, sourcePath, targetDest string) error {
	absSrc, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(absSrc)
	if err != nil {
		return err
	}
	baseName := filepath.Base(absSrc)
	destRel := baseName
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
	destPath := targetDest
	if destPath == "" {
		destPath = "$HOME/." + baseName
	}
	cfg.Mappings[destRel] = destPath
	return writeConfig(root, cfg)
}

func writeConfig(root string, cfg *config.Config) error {
	path := filepath.Join(root, config.ConfigFilename)
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
