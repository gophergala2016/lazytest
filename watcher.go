package lazytest

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-fsnotify/fsnotify"
)

var (
	fileWatches   = []*fsnotify.Watcher{}
	folderWatches = []*fsnotify.Watcher{}
	lock          sync.Mutex
)

type Mod struct {
	Package  string
	FilePath string
	Function string
	Line     int
}

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) && !ev.IsAttrib() {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watch() {
	for _, p := range settings.WatchPaths {
		p, _ = filepath.Abs(p)
		filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if isExcluded(path) {
					return filepath.SkipDir
				}
				if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}
				watchFolder(path)
			}
			return err
		})
	}
}

func Watch(include []FileMatch, exclude []FileMatch) chan Mod {
	mods := make(chan Mod, 50)
	return mods
}
