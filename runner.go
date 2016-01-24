package lazytest

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type testQueue struct {
	tests []Batch
}

const (
	RunnerIdle int32 = iota
	RunnerBusy
)

type TestStatus int8

const (
	StatusPending TestStatus = iota
	StatusSkipped
	StatusFailed
	StatusPanicked
	StatusPassed
)

type TestReport struct {
	Name    string
	Package string
	Status  TestStatus
	Message string
}

var (
	runnerDone   chan struct{}
	runnerStatus int32
	mux          sync.Mutex
	queue        *testQueue = &testQueue{}
	rep          chan Report
)

type Report []TestReport

func Runner(batch chan Batch) chan Report {
	rep = make(chan Report, 50)
	go queueTests(batch, rep)
	return rep
}

func (t *testQueue) run() {
	packageTests := make(map[string][]string)
	for _, test := range t.tests {
		packageTests[test.Package] = append(packageTests[test.Package], regexp.QuoteMeta(test.TestName))
	}
	for pkg, tests := range packageTests {
		testRegexp := fmt.Sprintf("'(%s)'", strings.Join(tests, "|"))

		cmd := exec.Command("go", "test", pkg, "-run", testRegexp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log(err.Error())
		}
		log(string(out))
	}
	atomic.StoreInt32(&runnerStatus, RunnerIdle)
	runnerDone <- struct{}{}
}

func queueTests(batch chan Batch, rep chan Report) {
	var delay *time.Timer
	for {
		select {
		case b := <-batch:
			mux.Lock()
			if delay == nil {
				delay = time.NewTimer(time.Second * 2)
			}
			queue.tests = append(queue.tests, b)
			mux.Unlock()

		case <-delay.C:
			mux.Lock()
			if atomic.CompareAndSwapInt32(&runnerStatus, RunnerIdle, RunnerBusy) {
				delay = nil
				go queue.run()
				queue = &testQueue{}
			}
			mux.Unlock()

		case <-runnerDone:
			mux.Lock()
			if delay == nil && len(queue.tests) > 0 {
				atomic.StoreInt32(&runnerStatus, RunnerBusy)
				go queue.run()
				queue = &testQueue{}
			}
			mux.Unlock()
		}
	}
}
