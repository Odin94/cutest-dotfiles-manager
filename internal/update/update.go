package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Odin94/cutest-dotfiles-manager/internal/ui"
)

const repo = "Odin94/cutest-dotfiles-manager"

type githubRelease struct {
	TagName string         `json:"tag_name"`
	Assets  []githubAsset  `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func Run() error {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github API: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rel githubRelease
	if err := json.Unmarshal(body, &rel); err != nil {
		return err
	}
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	suffix := ""
	if goos == "windows" {
		suffix = ".exe"
	}
	want := fmt.Sprintf("-%s-%s%s", goos, goarch, suffix)
	var downloadURL string
	for _, a := range rel.Assets {
		if strings.HasSuffix(a.Name, want) || strings.Contains(a.Name, want) {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no asset found for %s/%s (looking for name containing %q)", goos, goarch, want)
	}
	resp2, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("download: %s", resp2.Status)
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(exe)
	newName := "cdm.new"
	if goos == "windows" {
		newName = "cdm.exe.new"
	}
	dest := filepath.Join(dir, newName)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp2.Body)
	_ = out.Close()
	if err != nil {
		_ = os.Remove(dest)
		return err
	}
	if goos != "windows" {
		_ = os.Chmod(dest, 0755)
	}
	ui.PrintSuccess(fmt.Sprintf("Downloaded %s to %s", rel.TagName, dest))
	fmt.Println()
	fmt.Println("To finish updating:")
	fmt.Printf("  1. Close any terminals using cdm.\n")
	fmt.Printf("  2. Replace your current cdm binary with the new one:\n")
	fmt.Printf("     %s\n", dest)
	fmt.Printf("     -> %s\n", exe)
	return nil
}
