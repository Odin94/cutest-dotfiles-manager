package state

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const hashesFilename = "hashes.toml"

func HashesPath(root string) string {
	return filepath.Join(root, ".cdm", hashesFilename)
}

func LoadHashes(root string) (map[string]string, error) {
	path := HashesPath(root)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	var out map[string]string
	if err := toml.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = make(map[string]string)
	}
	return out, nil
}

func SaveHashes(root string, hashes map[string]string) error {
	dir := filepath.Join(root, ".cdm")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := HashesPath(root)
	data, err := toml.Marshal(hashes)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
