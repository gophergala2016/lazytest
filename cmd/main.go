package main

import (
	"flag"
	"os"

	"github.com/gophergala2016/lazytest"
	"github.com/mattn/go-colorable"
	log "github.com/Sirupsen/logrus"
)

var (
	flags      = flag.NewFlagSet("lazytest", flag.ExitOnError)
	include    = flags.String("include", "./*", "watch path")                            // modify to accept a slice
	extensions = flags.String("extensions", "go,tpl,html", "file extensions to watch") // comma separated list of watched extensions
	exclude    = flags.String("exclude", "", "exclude from watch")                     // modify to accept a slice
)

func init() {
	flags.Parse(os.Args[1:])
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	events    := lazytest.Watch(nil, nil)
	testBatch := lazytest.MatchTests(events)
	report    := lazytest.Runner(testBatch)
	lazytest.Render(report)
}
