package lazytest

import (
	"github.com/k0kubun/go-ansi"
)

func log(text string) {
	ansi.Println(text)
}

func Render(report chan Report) {
	for {
		r := <-report
		for _, test := range r {
			log(test.Message)
		}
	}
}
