package state

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const logFilename = "log.txt"

func LogPath(root string) string {
	return filepath.Join(root, ".cdm", logFilename)
}

func AppendLog(root, line string) error {
	dir := filepath.Join(root, ".cdm")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := LogPath(root)
	ts := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("%s %s\n", ts, line)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(entry)
	return err
}
