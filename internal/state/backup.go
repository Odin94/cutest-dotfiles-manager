package state

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

const (
	tempDirName         = "temp"
	lastBackupFilename  = "last_backup.toml"
	backupTimeFormat    = "20060102_150405"
)

func TempDir(root string) string {
	return filepath.Join(root, ".cdm", tempDirName)
}

func LastBackupPath(root string) string {
	return filepath.Join(root, ".cdm", lastBackupFilename)
}

func BackupPath(root, targetPath string) string {
	dir := filepath.Join(root, ".cdm", tempDirName)
	base := filepath.Base(targetPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	if name == "" {
		name = base
		ext = ""
	}
	ts := time.Now().Format(backupTimeFormat)
	return filepath.Join(dir, name+"_"+ts+ext)
}

func EnsureTempDir(root string) error {
	return os.MkdirAll(TempDir(root), 0755)
}

func WriteLastBackupManifest(root string, tempToTarget map[string]string) error {
	dir := filepath.Join(root, ".cdm")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(tempToTarget)
	if err != nil {
		return err
	}
	return os.WriteFile(LastBackupPath(root), data, 0644)
}

func LoadLastBackupManifest(root string) (map[string]string, error) {
	path := LastBackupPath(root)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out map[string]string
	if err := toml.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func ClearLastBackupManifest(root string) error {
	return os.Remove(LastBackupPath(root))
}
