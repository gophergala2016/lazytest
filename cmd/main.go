package main

import (
	"flag"
	"os"
	"strings"

	"github.com/gophergala2016/lazytest"
	"github.com/mattn/go-colorable"
)

var (
	flags      = flag.NewFlagSet("lazytest", flag.ExitOnError)
	root       = flags.String("root", ".", "watch root")
	exclude    = flags.String("exclude", "/vendor/", "exclude paths")
	extensions = flags.String("extensions", "go,tpl,html", "file extensions to watch") // comma separated list of watched extensions
)

func init() {
	flags.Parse(os.Args[1:])
}

func main() {
	exts := strings.Split(*extensions, ",")
	excl := strings.Split(*exclude, ",")

	events, err := lazytest.Watch(*root, exts, excl)
	if err != nil {
		log.Println(err)
		return
	}
	testBatch := lazytest.MatchTests(events)
	report := lazytest.Runner(testBatch)
	lazytest.Render(report)
}
