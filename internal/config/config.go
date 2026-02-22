package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	ConfigFilename     = ".cdm.toml"
	LocalConfigFilename = ".local.cdm.toml"
)

type Config struct {
	Variables []string         `toml:"variables"`
	Mappings  map[string]string `toml:"mappings"`
	Scripts   ScriptsConfig    `toml:"scripts"`
	Windows   map[string]string `toml:"mappings.windows,omitempty"`
	Macos     map[string]string `toml:"mappings.macos,omitempty"`
	Linux     map[string]string `toml:"mappings.linux,omitempty"`
}

type ScriptsConfig struct {
	PreApply  []string `toml:"pre_apply"`
	PostApply []string `toml:"post_apply"`
}

type LocalConfig struct {
	Values map[string]string `toml:",inline"`
}

func Load(root string) (*Config, error) {
	path := filepath.Join(root, ConfigFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Mappings == nil {
		cfg.Mappings = make(map[string]string)
	}
	return &cfg, nil
}

func LoadLocal(root string) (*LocalConfig, error) {
	path := filepath.Join(root, LocalConfigFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &LocalConfig{Values: make(map[string]string)}, nil
		}
		return nil, err
	}
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for k, v := range raw {
		if s, ok := v.(string); ok {
			values[k] = s
		}
	}
	return &LocalConfig{Values: values}, nil
}

func FindConfigDir(start string, traverseUp bool, maxLevels int) (string, bool) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for i := 0; i <= maxLevels; i++ {
		p := filepath.Join(dir, ConfigFilename)
		if _, err := os.Stat(p); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
		if !traverseUp && i == 0 {
			return "", false
		}
	}
	return "", false
}

func ExistsInDir(dir string) bool {
	path := filepath.Join(dir, ConfigFilename)
	_, err := os.Stat(path)
	return err == nil
}
