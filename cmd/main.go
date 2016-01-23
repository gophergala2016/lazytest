package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gophergala2016/lazytest"
)

var flags struct {
	root       string
	exclude    string
	extensions string
}

func init() {
	flag.StringVar(&flags.root, "root", ".", "watch root")
	flag.StringVar(&flags.exclude, "exclude", "/vendor/", "exclude paths")
	flag.StringVar(&flags.extensions, "extensions", "go,tpl,html",
		"file extensions to watch")
	flag.Parse()
}

func main() {
	testBatch := lazytest.MatchTests(watch())
	report := lazytest.Runner(testBatch)
	lazytest.Render(report)
}

func watch() chan lazytest.Mod {
	exclude := strings.Split(flags.exclude, ",")
	extensions := strings.Split(flags.extensions, ",")

	events, err := lazytest.Watch(flags.root, extensions, exclude)
	if err != nil {
		log.Fatal(err)
	}

	return events
}
