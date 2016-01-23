package lazytest

type Report struct {
	Message string
}

func Runner(batch chan Batch) chan Report {
	report := make(chan Report, 50)
	return report
}
