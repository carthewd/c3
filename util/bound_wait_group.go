package util

import "sync"

// BoundWaitGroup adds a semaphore to wait groups by way of a buffered channel
type BoundWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

// NewBoundWaitGroup returns a new wait group with channel semaphore
func NewBoundWaitGroup(cap int) BoundWaitGroup {
	return BoundWaitGroup{ch: make(chan struct{}, cap)}
}

// Add updates the wait group and semaphore channel
func (bwg *BoundWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

// Done
func (bwg *BoundWaitGroup) Done() {
	bwg.Add(-1)
}

// Wait on wait group
func (bwg *BoundWaitGroup) Wait() {
	bwg.wg.Wait()
}
