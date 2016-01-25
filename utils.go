package lazytest

/*
 * Batch is a struct holding information about a batch of unit tests.
 */
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

/*
 * MatchTests launches a go routine to match file change events to unit tests.
 */
func MatchTests(events chan Mod) chan Batch {
	batch := make(chan Batch, 50)
	go match(events, batch)
	return batch
}
