package utils

import (
	"iter"
	"sync"
)

type indexedMapElement[V any] struct {
	Value V
	Index int
}

type IndexedMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]indexedMapElement[V]
	keys  []K
}

func NewIndexedMap[K comparable, V any]() *IndexedMap[K, V] {
	return &IndexedMap[K, V]{
		m:    make(map[K]indexedMapElement[V]),
		keys: make([]K, 0),
	}
}

func (im *IndexedMap[K, V]) Iterator() iter.Seq2[K, indexedMapElement[V]] {
	return func(yield func(K, indexedMapElement[V]) bool) {
		im.mutex.RLock()
		defer im.mutex.RUnlock()

		for _, k := range im.keys {
			if !yield(k, im.m[k]) {
				return
			}
		}
	}
}

func (im *IndexedMap[K, V]) Keys() []K {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	return im.keys
}

func (im *IndexedMap[K, V]) Has(key K) bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	_, ok := im.m[key]
	return ok
}

func (im *IndexedMap[K, V]) Insert(key K, value V) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	im.m[key] = indexedMapElement[V]{value, len(im.keys)}
	im.keys = append(im.keys, key)
}

func (im *IndexedMap[K, V]) Remove(key K) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	i := im.m[key].Index
	lastI := len(im.keys) - 1
	lastK := im.keys[lastI]

	im.keys[i] = lastK
	im.m[lastK] = indexedMapElement[V]{im.m[lastK].Value, i}

	im.keys = im.keys[:lastI]
	delete(im.m, key)
}

func (im *IndexedMap[K, V]) Get(key K) (int, V) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	return im.m[key].Index, im.m[key].Value
}

func (im *IndexedMap[K, V]) GetIndex(i int) (K, V) {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	k := im.keys[i]
	return k, im.m[k].Value
}

func (im *IndexedMap[K, V]) Len() int {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	return len(im.keys)
}
