package utils

import (
	"iter"
	"sync"
)

type MutexMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
	return &MutexMap[K, V]{
		m: make(map[K]V),
	}
}

func (mm *MutexMap[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		mm.mutex.RLock()
		defer mm.mutex.RUnlock()

		for k, v := range mm.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (mm *MutexMap[K, V]) Keys() []K {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	keys := make([]K, len(mm.m))

	i := 0
	for k := range mm.m {
		keys[i] = k
		i++
	}
	return keys
}

func (mm *MutexMap[K, V]) Set(key K, value V) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	mm.m[key] = value
}

func (mm *MutexMap[K, V]) Get(key K) (V, bool) {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()

	v, ok := mm.m[key]
	return v, ok
}

func (mm *MutexMap[K, V]) Delete(key K) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	delete(mm.m, key)
}
