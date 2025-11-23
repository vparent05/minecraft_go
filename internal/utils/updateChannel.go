package utils

import (
	"iter"
	"sync"
)

type UpdateChannel[T any] struct {
	m  sync.Mutex
	ch chan T
}

func NewUpdateChannel[T any]() *UpdateChannel[T] {
	return &UpdateChannel[T]{ch: make(chan T, 1)}
}

func (u *UpdateChannel[T]) TryReceive() (T, bool) {
	select {
	case t := <-u.ch:
		return t, true
	default:
		return *new(T), false
	}
}

func (u *UpdateChannel[T]) Receive() T {
	return <-u.ch
}

func (u *UpdateChannel[T]) Send(v T) {
	u.m.Lock()
	defer u.m.Unlock()

	select {
	case u.ch <- v:
		return
	case <-u.ch:
		u.ch <- v
	}
}

func (u *UpdateChannel[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for t := range u.ch {
			if !yield(t) {
				return
			}
		}
	}
}
