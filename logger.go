package lazytest

import (
	"github.com/k0kubun/go-ansi"
)

func Log(text string) {
	ansi.Println(text)
}

func Render(report chan Report) {
	for {
		r := <- report
		Log(r.Message)
	}
}
