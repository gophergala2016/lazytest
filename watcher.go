package lazytest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type fileWatcher struct {
	extensions []string
	exclude    []string
	watcher    *fsnotify.Watcher
}

type Mod struct {
	Package  string
	FilePath string
	Function string
	Line     int
}

func (w *fileWatcher) handleDir(path string) error {
	if !w.isIncluded(path) {
		return filepath.SkipDir
	}

	if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
		return filepath.SkipDir
	}

	return w.watcher.Add(path)
}

func (w *fileWatcher) handleEvent(e fsnotify.Event, eventChannel chan Mod) {
	if e.Op|fsnotify.Rename == e.Op || e.Op|fsnotify.Chmod == e.Op {
		return
	}

	eventChannel <- Mod{FilePath: e.Name}
	// TODO: remove old watches on delete, add new watches on create, do both on rename
}

func (w *fileWatcher) isIncluded(path string) bool {
	include := len(w.extensions) == 0

	for _, e := range w.extensions {
		if filepath.Ext(path) == e {
			include = true
		}
	}

	for _, e := range w.exclude {
		if filepath.HasPrefix(path, e) {
		   return false
		}
	}

	return include
}

func (w *fileWatcher) listenForEvents(eventChannel chan Mod) {
	for {
		select {
		case e := <-w.watcher.Events:
			w.handleEvent(e, eventChannel)

		case err := <-w.watcher.Errors:
			log(fmt.Sprintf("Watcher error %v", err))
		}
	}
}

func (w *fileWatcher) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return w.handleDir(path)
	}

	if w.isIncluded(path) {
		return w.watcher.Add(path)
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

	w := &fileWatcher{
		extensions: extensions,
		exclude:    exclude,
		watcher:    watcher,
	}

	events := make(chan Mod, 50)
	go w.listenForEvents(events)
	return events, filepath.Walk(absolutePath, w.walk)
}
