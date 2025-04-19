package utils

import (
	"iter"
	"sync"
)

type IndexedMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func NewIndexedMap[K comparable, V any]() *IndexedMap[K, V] {
	return &IndexedMap[K, V]{
		m: make(map[K]V),
	}
}

func (im *IndexedMap[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		im.mutex.RLock()
		defer im.mutex.RUnlock()

		for k, v := range im.m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (im *IndexedMap[K, V]) Keys() []K {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	keys := make([]K, len(im.m))

	i := 0
	for k := range im.m {
		keys[i] = k
		i++
	}
	return keys
}

func (im *IndexedMap[K, V]) Insert(key K, value V) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	im.m[key] = value
}

func (im *IndexedMap[K, V]) Remove(key K) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	delete(im.m, key)
}

func (im *IndexedMap[K, V]) Get(key K) (V, bool) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	v, ok := im.m[key]
	return v, ok
}
