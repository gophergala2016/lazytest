package lazytest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Mod struct {
	Package  string
	FilePath string
	Function string
	Line     int
}

type fileWalker struct {
	extensions []string
	exclude    []string
	watcher    *fsnotify.Watcher
}

func (w *fileWalker) walkFunction(path string, info os.FileInfo,
	err error) error {

	if info.IsDir() {
		if !w.isIncluded(path) {
			return filepath.SkipDir
		}
		if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}
		if err := w.watcher.Add(path); err != nil {
			panic(err) // TODO: better error handling
		}
		return err
	}

	if w.isIncluded(path) {
		if err := w.watcher.Add(path); err != nil {
			panic(err) // TODO: better error handling
		}
	}

	return err
}

func Watch(root string, extensions []string, exclude []string) (chan Mod,
	error) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absolutePath, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	walker := &fileWalker{
		extensions: extensions,
		exclude:    exclude,
		watcher:    watcher,
	}

	events := make(chan Mod, 50)
	go walker.handleEvents(events)
	return events, filepath.Walk(absolutePath, walker.walkFunction)
}

func (w *fileWalker) handleEvents(events chan Mod) {
	for {
		select {
		case e := <-w.watcher.Events:
			// don't trigger events on file being renamed or chmod changed
			if (e.Op|fsnotify.Write == e.Op) || (e.Op|fsnotify.Create == e.Op) || (e.Op|fsnotify.Remove == e.Op) {
				events <- Mod{
					FilePath: e.Name,
				}
			}
			// TODO: remove old watches on delete, add new watches on create, do both on rename

		case err := <-w.watcher.Errors:
			log(fmt.Sprintf("Watcher error %v", err))
		}
	}
}

func (w *fileWalker) isIncluded(path string) bool {
	include := false

	// if no extensions were provided match all
	if len(w.extensions) == 0 {
		include = true
	} else {
		ext := filepath.Ext(path)
		for _, e := range w.extensions {
			if ext == e {
				include = true
			}
		}
	}

	for _, e := range w.exclude {
		if filepath.HasPrefix(path, e) {
			include = false
		}
	}

	return include
}
