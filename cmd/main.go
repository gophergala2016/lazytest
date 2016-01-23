package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gophergala2016/lazytest"
)

var (
	flags      = flag.NewFlagSet("lazytest", flag.ExitOnError)
	root       = flags.String("root", ".", "watch root")
	exclude    = flags.String("exclude", "./vendor/*", "exclude paths")
	extensions = flags.String("extensions", "go,tpl,html", "file extensions to watch") // comma separated list of watched extensions
)

func init() {
	flags.Parse(os.Args[1:])
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	events := lazytest.Watch(nil, nil)
	testBatch := lazytest.MatchTests(events)
	report := lazytest.Runner(testBatch)
	lazytest.Render(report)
}
