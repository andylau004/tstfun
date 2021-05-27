package utils

import "sync"

type Barrier struct {
	curCnt int
	maxCnt int
	cond   *sync.Cond
}

func NewBarrier(maxCnt int) *Barrier {
	mutex := new(sync.Mutex)
	cond := sync.NewCond(mutex)
	return &Barrier{curCnt: maxCnt, maxCnt: maxCnt, cond: cond}
}

func (barrier *Barrier) BarrierWait() {
	barrier.cond.L.Lock()
	if barrier.curCnt--; barrier.curCnt > 0 {
		barrier.cond.Wait()
	} else {
		barrier.cond.Broadcast()
		barrier.curCnt = barrier.maxCnt
	}
	barrier.cond.L.Unlock()
}
