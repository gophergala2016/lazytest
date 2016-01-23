package main

import (
	"flag"
	"os"

	"github.com/gophergala2016/lazytest"
)

var (
	flags      = flag.NewFlagSet("lazytest", flag.ExitOnError)
	include    = flags.String("include", ".", "watch path")                            // modify to accept a slice
	extensions = flags.String("extensions", "go,tpl,html", "file extensions to watch") // comma separated list of watched extensions
	exclude    = flags.String("exclude", "", "exclude from watch")                     // modify to accept a slice
)

func main() {
	flags.Parse(os.Args[1:])

	// add error handling to all 4
	events := lazytest.Watch(nil, nil)
	testBatch := lazytest.MatchTests(events)
	report := lazytest.Runner(testBatch)
	lazytest.Render(report)
}
