package lazytest

import (
	"fmt"
	"go/parser"
	"go/token"
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
	if !w.isIncluded(path, false) {
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

	eventChannel <- Mod{
		FilePath: e.Name,
		Package:  packageName(e.Name),
	}
	// TODO: remove old watches on delete, add new watches on create, do both on rename
}

func (w *fileWatcher) isIncluded(path string, isFile bool) bool {
	include := len(w.extensions) == 0

	if !isFile {
		include = true
	} else {
		for _, e := range w.extensions {
			if filepath.Ext(path) == e {
				include = true
			}
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

	if w.isIncluded(path, true) {
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

func packageName(path string) string {
	fset := token.NewFileSet()
	// parse the go source file, but only the package clause
	astFile, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
	if err != nil {
		log(err.Error())
		return ""
	}
	if astFile.Name == nil {
		log("no name")
		return ""
	}
	pkg := filepath.Dir(path)
	lastSlash := strings.LastIndex(pkg, string(filepath.Separator)) + 1
	pkg = pkg[0:lastSlash]
	pkg = pkg + astFile.Name.Name

	gopath := os.Getenv("GOPATH")
	gosrc := gopath + string(filepath.Separator) + "src" + string(filepath.Separator)

	pkg = strings.TrimPrefix(pkg, gosrc)

	return pkg
}
