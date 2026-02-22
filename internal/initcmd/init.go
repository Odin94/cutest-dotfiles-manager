package initcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cdm/internal/config"
	"cdm/internal/ui"
)

func Run(root string, cloneRepo string) error {
	if cloneRepo != "" {
		if err := cloneAndChdir(&root, cloneRepo); err != nil {
			return err
		}
	}
	cdmDir := filepath.Join(root, ".cdm")
	if err := os.MkdirAll(filepath.Join(cdmDir, "temp"), 0755); err != nil {
		return err
	}
	configPath := filepath.Join(root, config.ConfigFilename)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		content := "# cdm config\n# variables = [\"HOME\"]\n# [mappings]\n# \".bashrc\" = \"$HOME/.bashrc\"\n"
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			return err
		}
		ui.PrintSuccess("created " + configPath)
	}
	gitignore := filepath.Join(root, ".gitignore")
	const cdmIgnore = "\n.cdm/\n.local.cdm.toml\n"
	existing, _ := os.ReadFile(gitignore)
	if len(existing) > 0 {
		existing = append(existing, '\n')
	}
	if !strings.Contains(string(existing), ".cdm/") {
		if err := os.WriteFile(gitignore, append(existing, []byte(cdmIgnore)...), 0644); err != nil {
			return err
		}
		ui.PrintSuccess("updated .gitignore")
	} else if !strings.Contains(string(existing), ".local.cdm.toml") {
		if err := os.WriteFile(gitignore, append(existing, []byte(".local.cdm.toml\n")...), 0644); err != nil {
			return err
		}
		ui.PrintSuccess("updated .gitignore")
	}
	ui.PrintSuccess("init done at " + root)
	return nil
}

func cloneAndChdir(root *string, repo string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	base := filepath.Base(strings.TrimSuffix(strings.TrimSuffix(repo, ".git"), "/"))
	cloneDir := filepath.Join(cwd, base)
	if err := runGitClone(repo, cloneDir); err != nil {
		return fmt.Errorf("clone: %w", err)
	}
	*root = cloneDir
	return nil
}

func runGitClone(repo, dir string) error {
	cmd := exec.Command("git", "clone", repo, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
