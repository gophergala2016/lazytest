package main

import (
	"flag"
	"os"

	"github.com/gophergala2016/lazytest"
)

var (
	flags      = flag.NewFlagSet("lazytest", flag.ExitOnError)
	root       = flags.String("root", ".", "watch root")
	exclude    = flags.String("exclude", "./vendor/*", "exclude paths")
	extensions = flags.String("extensions", "go,tpl,html", "file extensions to watch") // comma separated list of watched extensions
	exclude    = flags.String("exclude", "", "exclude from watch")                     // modify to accept a slice
)

func main() {
	flags.Parse(os.Args[1:])

	// add error handling to all 4
	events := lazytest.Watch(root, exclude, extensions)
	testBatch := lazytest.MatchTests(events)
	report := lazytest.Runner(testBatch)
	lazytest.Render(report)
}
