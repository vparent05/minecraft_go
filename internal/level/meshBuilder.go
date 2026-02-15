package level

import (
	"sync"
	"time"

	"github.com/vparent05/minecraft_go/internal/utils/chanx"
	"github.com/vparent05/minecraft_go/internal/utils/debounce"
	"github.com/vparent05/minecraft_go/internal/utils/structure"
)

const MESH_BUILDING_WORKER_COUNT = 4

type meshBuilder struct {
	mu             sync.Mutex
	wg             sync.WaitGroup
	stop           chan struct{}
	new            sync.Cond
	queue          *structure.Heap[*Chunk, int]
	queueItems     map[*Chunk]*structure.Item[*Chunk, int]
	level          *Level
	fixAllDebounce *debounce.Debounce
}

func newMeshBuilder(level *Level, renderDistance int) *meshBuilder {
	m := &meshBuilder{
		mu:   sync.Mutex{},
		wg:   sync.WaitGroup{},
		stop: make(chan struct{}),
		queue: structure.NewHeap(
			func(a *Chunk) int {
				observer := level.observer.Load()
				observerChunkCoords := LevelToChunkCoords(observer.Vec3)
				diffA := a.coordinates.Sub(observerChunkCoords)
				dx := diffA.X
				dz := diffA.Y
				return dx*dx + dz*dz

			},
			func(a, b int) int {
				return b - a
			},
			(2*renderDistance+1)*(2*renderDistance+1)),
		queueItems:     make(map[*Chunk]*structure.Item[*Chunk, int], (2*renderDistance+1)*(2*renderDistance+1)),
		level:          level,
		fixAllDebounce: debounce.NewDebounce(50 * time.Millisecond),
	}

	m.new = *sync.NewCond(&m.mu)

	return m
}

func (m *meshBuilder) movedChunk() {
	m.fixAllDebounce.Do(func() {
		m.mu.Lock()
		m.queue.FixAll()
		m.mu.Unlock()
	})
}

func (m *meshBuilder) enqueue(c *Chunk) {
	m.mu.Lock()
	if item, ok := m.queueItems[c]; ok {
		m.queue.Fix(item)
	} else {
		m.queueItems[c] = m.queue.Add(c)
	}
	m.mu.Unlock()
	m.new.Signal()
}

func (m *meshBuilder) start(n int) {
	for range n {
		m.newWorker()
	}
}

func (m *meshBuilder) newWorker() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		for {
			if _, ok := chanx.TryReceive(m.stop); ok {
				return
			}

			m.mu.Lock()
			for m.queue.Size() == 0 {
				m.new.Wait()
			}

			chunk := m.queue.Pop()
			delete(m.queueItems, chunk)
			m.mu.Unlock()

			chunk.generateMesh(m.level)
		}
	}()

}

func (m *meshBuilder) stopWorkers() {
	close(m.stop)
	m.wg.Wait()
}
