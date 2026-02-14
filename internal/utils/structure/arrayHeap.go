package structure

type Item[T, S any] struct {
	value T
	score S
	index int
}

type Heap[T, S any] struct {
	values  []*Item[T, S]
	compare func(a, b *Item[T, S]) int
	score   func(a T) S
}

// cmp(a, b) > 0 if a has higher priority than b
func NewHeap[T, S any](score func(a T) S, compare func(a, b S) int, capacity int) *Heap[T, S] {
	return &Heap[T, S]{
		values: make([]*Item[T, S], 0, capacity),
		compare: func(a, b *Item[T, S]) int {
			return compare(a.score, b.score)
		},
		score: score,
	}
}

func (h *Heap[T, S]) Size() int {
	return len(h.values)
}

func (h *Heap[T, S]) Add(new T) *Item[T, S] {
	newItem := &Item[T, S]{value: new, index: len(h.values), score: h.score(new)}
	h.values = append(h.values, newItem)
	h.siftUp(newItem)

	return newItem
}

func (h *Heap[T, S]) Fix(item *Item[T, S]) {
	item.score = h.score(item.value)
	h.siftUp(item)
	h.siftDown(item)
}

func (h *Heap[T, S]) Pop() T {
	result := h.values[0]

	last := len(h.values) - 1
	toPlace := h.values[last]
	h.values = h.values[:last]

	if last == 0 {
		return result.value
	}

	toPlace.index = 0
	h.siftDown(toPlace)

	return result.value
}

func (h *Heap[T, S]) siftUp(item *Item[T, S]) {
	current := item.index
	for current > 0 {
		p := parent(current)
		if h.compare(item, h.values[p]) <= 0 {
			break
		}

		h.values[current] = h.values[p]
		h.values[current].index = current
		current = p
	}

	h.values[current] = item
	item.index = current
}

func (h *Heap[T, S]) siftDown(item *Item[T, S]) {
	current := item.index
	n := len(h.values)
	for {
		child := left(current)
		if child >= n {
			break
		}

		rChild := child + 1
		if rChild < n && h.compare(h.values[rChild], h.values[child]) > 0 {
			child = rChild
		}

		if h.compare(item, h.values[child]) >= 0 {
			break
		}

		h.values[current] = h.values[child]
		h.values[current].index = current
		current = child
	}

	h.values[current] = item
	h.values[current].index = current
}

func left(i int) int {
	return i*2 + 1
}

func parent(i int) int {
	return (i - 1) / 2
}
