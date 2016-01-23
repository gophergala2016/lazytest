package lazytest

type TestResult uint8

const (
	TestFailed TestResult = iota
	TestPassed
	TestErrored
)

type TestReport struct {
	Name    string
	Result  TestResult
	Message string
}

type Report []TestReport

var rep chan Report

func Runner(batch chan Batch) chan Report {
	rep = make(chan Report, 50)

	go testRunner(batch)

	return rep
}

func testRunner(batch chan Batch) {

}
