package syncx

import (
	"sync/atomic"
	"time"
)

// is an implementation of atomic duration.
type AtomicDuration int64

func NewAtomicDuration() *AtomicDuration {
	return new(AtomicDuration)
}

// returns an AtomicDuration with given value
func ForAtomicDuration(val time.Duration) *AtomicDuration {
	d := NewAtomicDuration()
	d.Set(val)
	return d
}

// compares current value with old, if equals, set the value to val.
func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(old), int64(val))
}

func (d *AtomicDuration) Load() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

func (d *AtomicDuration) Set(val time.Duration) {
	atomic.StoreInt64((*int64)(d), (int64)(val))
}
