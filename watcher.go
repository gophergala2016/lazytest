package lazytest

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

var (
	watcher    *fsnotify.Watcher
	events     chan Mod
	root       string
	extensions []string
	exclude    []string

	/*fileWatches   = []*fsnotify.Watcher{}
	folderWatches = []*fsnotify.Watcher{}
	lock          sync.Mutex*/
)

type Mod struct {
	Package  string
	FilePath string
	Function string
	Line     int
}

func Watch(r string, ext []string, excl []string) (chan Mod, error) {
	var err error
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		return nil, err
	}
	root = r
	exclude = excl
	extensions = ext

	go handleEvents()

	// get absolute path
	absP, err := filepath.Abs(r)
	if err != nil {
		return nil, err
	}

	// walk folder tree starting at the path
	filepath.Walk(absP, func(path string, info os.FileInfo, err error) error {
		log.Println(path)
		if info.IsDir() {
			if !isIncluded(path) {
				return filepath.SkipDir
			}
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}
			if err := watcher.Add(path); err != nil {
				panic(err) // this needs better error handling
			}
		} else if isIncluded(path) {
			if err := watcher.Add(path); err != nil {
				panic(err) // this needs better error handling
			}
		}
		return err
	})

	return events, nil
}

func handleEvents() {
	for {
		select {
		case e := <-watcher.Events:
			// don't trigger events on file being renamed or chmod changed
			if (e.Op|fsnotify.Write == e.Op) || (e.Op|fsnotify.Create == e.Op) || (e.Op|fsnotify.Remove == e.Op) {
				events <- Mod{
					FilePath: e.Name,
				}
			}

		case err := <-watcher.Errors:
			log.Printf("Watcher error %v", err)
			return
		}
	}
}

func isIncluded(path string) bool {
	include := false

	// if no extensions were provided match all
	if len(extensions) == 0 {
		include = true
	} else {
		ext := filepath.Ext(path)
		for _, e := range extensions {
			if ext == e {
				include = true
			}
		}
	}

	for _, e := range exclude {
		if filepath.HasPrefix(path, e) {
			include = false
		}
	}

	return include
}
