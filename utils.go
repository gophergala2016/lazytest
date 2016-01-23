package lazytest

type Batch struct {
	Package  string
	TestName string
}

func match(events chan Mod, batch chan Batch) {
	for {
		event := <-events
		batch <- Batch{Package: event.Package}
	}
}

func MatchTests(events chan Mod) chan Batch {
	batch := make(chan Batch, 50)
	go match(events, batch)
	return batch
}
