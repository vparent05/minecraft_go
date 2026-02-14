package level

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/vparent05/minecraft_go/internal/utils/chanx"
	"github.com/vparent05/minecraft_go/internal/utils/structure"
)

const WORKER_COUNT = 4

type meshGenerator struct {
	mu         sync.Mutex
	wg         sync.WaitGroup
	stop       chan struct{}
	new        sync.Cond
	queue      *structure.Heap[*Chunk, float64]
	queueItems map[*Chunk]*structure.Item[*Chunk, float64]
	level      *Level
}

func newMeshGenerator(level *Level, renderDistance int) *meshGenerator {
	m := &meshGenerator{
		mu:   sync.Mutex{},
		wg:   sync.WaitGroup{},
		stop: make(chan struct{}),
		queue: structure.NewHeap(
			func(a *Chunk) float64 {
				observer := level.observer.Load()
				observerChunkCoords := LevelToChunkCoords(observer.Vec3)
				diffA := a.coordinates.Sub(observerChunkCoords)
				return math.Abs(float64(diffA.X)) + math.Abs(float64(diffA.Y))
			},
			func(a, b float64) int {
				return int(b - a)
			},
			(2*renderDistance+1)*(2*renderDistance+1)),
		queueItems: make(map[*Chunk]*structure.Item[*Chunk, float64], (2*renderDistance+1)*(2*renderDistance+1)),
		level:      level,
	}

	m.new = *sync.NewCond(&m.mu)

	return m
}

var callerMap = make(map[struct {
	file string
	line int
}]int)

var o = sync.Once{}

func (m *meshGenerator) enqueue(c *Chunk) {
	_, file, line, _ := runtime.Caller(1)

	o.Do(func() {
		go func() {
			time.Sleep(10 * time.Second)
			for key, value := range callerMap {
				fmt.Printf("%s:%d : %d times\n", key.file, key.line, value)
			}
		}()
	})

	m.mu.Lock()
	if _, ok := callerMap[struct {
		file string
		line int
	}{file, line}]; !ok {
		callerMap[struct {
			file string
			line int
		}{file, line}] = 0
	}
	callerMap[struct {
		file string
		line int
	}{file, line}]++
	if item, ok := m.queueItems[c]; ok {
		m.queue.Fix(item)
	} else {
		m.queueItems[c] = m.queue.Add(c)
	}
	m.mu.Unlock()
	m.new.Signal()
}

func (m *meshGenerator) start(n int) {
	for range n {
		m.newWorker()
	}
}

func (m *meshGenerator) newWorker() {
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

func (m *meshGenerator) stopWorkers() {
	close(m.stop)
	m.wg.Wait()
}
