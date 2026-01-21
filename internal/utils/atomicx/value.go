package atomicx

import (
	"sync/atomic"
)

type Value[T any] struct {
	atomic.Value
}

// Load atomically loads and returns the value stored in v.
func (v *Value[T]) Load() T {
	return v.Value.Load().(T)
}

// Store atomically stores val into v.
func (v *Value[T]) Store(val T) {
	v.Value.Store(val)
}

// Swap atomically stores new into v and returns the previous value.
func (v *Value[T]) Swap(new T) (old T) {
	return v.Value.Swap(new).(T)
}

// CompareAndSwap executes the compare-and-swap operation for v.
func (v *Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.Value.CompareAndSwap(old, new)
}
