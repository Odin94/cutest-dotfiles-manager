package watch

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/cutest-tools/cutest-dotfiles-manager/internal/apply"
	"github.com/cutest-tools/cutest-dotfiles-manager/internal/config"
	"github.com/fsnotify/fsnotify"
)

const debounceMs = 400

func Run(ctx context.Context, root string, applyFn func() (*apply.ApplyResult, error)) error {
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	local, _ := config.LoadLocal(root)
	resolved, _ := cfg.ResolvedMappings(local)
	var paths []string
	for srcRel := range resolved {
		paths = append(paths, filepath.Join(root, srcRel))
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	dirs := make(map[string]struct{})
	for _, p := range paths {
		dir := filepath.Dir(p)
		if _, ok := dirs[dir]; !ok {
			dirs[dir] = struct{}{}
			if err := watcher.Add(dir); err != nil {
				continue
			}
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var mu sync.Mutex
	var timer *time.Timer
	scheduleApply := func() {
		mu.Lock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(debounceMs*time.Millisecond, func() {
			mu.Lock()
			timer = nil
			mu.Unlock()
			applyFn()
		})
		mu.Unlock()
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				scheduleApply()
			case <-watcher.Errors:
			}
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}
