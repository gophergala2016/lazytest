package lazytest

import (
	log "github.com/Sirupsen/logrus"
)

func Render(report chan Report) {
	for {
		r := <- report
		log.Info(r.Message)
	}
	close(report)
}
