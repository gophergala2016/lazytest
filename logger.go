package lazytest

import (
	"github.com/k0kubun/go-ansi"
)

func log(text string) {
	ansi.Println(text)
}

/*
 * Render listens on a provided channel and logs incoming messages
 */
func Render(report chan Report) {
	for {
		r := <-report
		for _, test := range r {
			log(test.Message)
		}
	}
}
