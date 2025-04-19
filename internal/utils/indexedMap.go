package utils

import (
	"iter"
)

type indexedMapElement[V any] struct {
	Value V
	Index int
}

type IndexedMap[K comparable, V any] struct {
	m    map[K]indexedMapElement[V]
	keys []K
}

func NewIndexedMap[K comparable, V any]() *IndexedMap[K, V] {
	return &IndexedMap[K, V]{
		m:    make(map[K]indexedMapElement[V]),
		keys: make([]K, 0),
	}
}

func (im *IndexedMap[K, V]) Iterator() iter.Seq2[K, indexedMapElement[V]] {
	return func(yield func(K, indexedMapElement[V]) bool) {
		for _, k := range im.keys {
			if !yield(k, im.m[k]) {
				return
			}
		}
	}
}

func (im *IndexedMap[K, V]) Insert(key K, value V) {
	im.m[key] = indexedMapElement[V]{value, len(im.keys)}
	im.keys = append(im.keys, key)
}

func (im *IndexedMap[K, V]) Remove(key K) {
	i := im.m[key].Index
	lastI := len(im.keys) - 1
	lastK := im.keys[lastI]

	im.keys[i] = lastK
	im.m[lastK] = indexedMapElement[V]{im.m[lastK].Value, i}

	im.keys = im.keys[:lastI]
	delete(im.m, key)
}

func (im *IndexedMap[K, V]) Get(key K) V {
	return im.m[key].Value
}

func (im *IndexedMap[K, V]) GetIndex(i int) (K, V) {
	return im.keys[i], im.m[im.keys[i]].Value
}

func (im *IndexedMap[K, V]) Len() int {
	return len(im.keys)
}
