package util

import "sync"

type BoundWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

func NewBoundWaitGroup(cap int) BoundWaitGroup {
	return BoundWaitGroup{ch: make(chan struct{}, cap)}
}

func (bwg *BoundWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

func (bwg *BoundWaitGroup) Done() {
	bwg.Add(-1)
}

func (bwg *BoundWaitGroup) Wait() {
	bwg.wg.Wait()
}
