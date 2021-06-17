package main

import (
	"sync/atomic"
)

type MyAtomic interface {
	IncreaseAllElse()
	IncreaseA()
	IncreaseB()
}

type NoPad struct {
	a uint64
	b uint64
	c uint64
}

func (this *NoPad) IncreaseAllElse() {
	atomic.AddUint64(&this.a, 1)
	atomic.AddUint64(&this.b, 1)
	atomic.AddUint64(&this.c, 1)
}
func (this *NoPad) IncreaseA() {
	atomic.AddUint64(&this.a, 1)
}

func (this *NoPad) IncreaseB() {
	atomic.AddUint64(&this.b, 1)
}

type Pad struct {
	a   uint64
	_p1 [8]uint64
	b   uint64
	_p2 [8]uint64
	c   uint64
	_p3 [8]uint64
}

func (this *Pad) IncreaseAllElse() {
	atomic.AddUint64(&this.a, 1)
	atomic.AddUint64(&this.b, 1)
	atomic.AddUint64(&this.c, 1)
}
func (this *Pad) IncreaseA() {
	atomic.AddUint64(&this.a, 1)
}
func (this *Pad) IncreaseB() {
	atomic.AddUint64(&this.b, 1)
}

// func main() {
// 	myatomic := &Pad{}
// 	myatomic.IncreaseAllEles()
// 	fmt.Printf("%d\n", myatomic.a)
// }
