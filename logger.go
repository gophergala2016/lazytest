package lazytest

import (
	"github.com/k0kubun/go-ansi"
)

func Render(report chan Report) {
	for {
		r := <- report
		ansi.Println(r.Message)
	}
}
