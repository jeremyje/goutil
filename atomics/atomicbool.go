package atomics

import "sync/atomic"

type AtomicBool struct {
	val int32
}

func (b *AtomicBool) Set(value bool) {
	var i int32
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.val), int32(i))
}

func (b *AtomicBool) Get() bool {
	if atomic.LoadInt32(&(b.val)) != 0 {
		return true
	}
	return false
}

func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}
