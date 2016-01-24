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

func (w *fileWalker) handleDir(path string) error {
	if !w.isIncluded(path) {
		return filepath.SkipDir
	}

	if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
		return filepath.SkipDir
	}

	return w.watcher.Add(path)
}

func handleEvent(e fsnotify.Event, eventChannel chan Mod) {
	if e.Op|fsnotify.Rename == e.Op || e.Op|fsnotify.Chmod == e.Op {
		return
	}

	eventChannel <- Mod{FilePath: e.Name}
	// TODO: remove old watches on delete, add new watches on create, do both on rename
}

func (w *fileWalker) isIncluded(path string) bool {
	include := len(w.extensions) == 0

	if !include {
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

func (w *fileWalker) listenForEvents(eventChannel chan Mod) {
	for {
		select {
		case e := <-w.watcher.Events:
			handleEvent(e, eventChannel)

		case err := <-w.watcher.Errors:
			log(fmt.Sprintf("Watcher error %v", err))
		}
	}
}

func (w *fileWalker) walkFunction(path string, info os.FileInfo,
	err error) error {

	if info.IsDir() {
		return w.handleDir(path)
	}

	if w.isIncluded(path) {
		return w.watcher.Add(path)
	}

	return nil
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
	go walker.listenForEvents(events)
	return events, filepath.Walk(absolutePath, walker.walkFunction)
}
