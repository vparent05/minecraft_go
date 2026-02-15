package level

import (
	"math"
	"sync"

	"github.com/vparent05/minecraft_go/internal/utils"
	"github.com/vparent05/minecraft_go/internal/utils/chanx"
)

const WORLD_GENERATOR_WORKER_COUNT = 2

type worldGenerator struct {
	mu         sync.Mutex
	wg         sync.WaitGroup
	stop       chan struct{}
	toGenerate chan *Chunk
}

func newWorldGenerator() *worldGenerator {
	w := &worldGenerator{
		mu:         sync.Mutex{},
		wg:         sync.WaitGroup{},
		stop:       make(chan struct{}),
		toGenerate: make(chan *Chunk),
	}

	return w
}

func (w *worldGenerator) enqueue(c *Chunk) {
	w.toGenerate <- c
}

func (w *worldGenerator) clear() {
	w.mu.Lock()
	close(w.toGenerate)
	w.toGenerate = make(chan *Chunk)
	w.mu.Unlock()
}

func (w *worldGenerator) start(n int) {
	for range n {
		w.newWorker()
	}
}

func (w *worldGenerator) newWorker() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			if _, ok := chanx.TryReceive(w.stop); ok {
				return
			}

			w.mu.Lock()
			toGenerate := w.toGenerate
			w.mu.Unlock()

			for chunk := range toGenerate {
				generateChunk(chunk)
			}

			// channel closed -> queue reset -> loop and pick up new channel
		}
	}()
}

func (m *worldGenerator) stopWorkers() {
	close(m.stop)
	m.wg.Wait()
}

func generateChunk(chunk *Chunk) {
	const WATER_LEVEL = 60
	var blocks [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_WIDTH]BlockId

	for i := range CHUNK_WIDTH {
		for j := range CHUNK_WIDTH {
			xBlock := i + chunk.coordinates.X*CHUNK_WIDTH
			zBlock := j + chunk.coordinates.Y*CHUNK_WIDTH
			topY := int(math.Floor(float64(utils.FractalNoise2(float32(xBlock), float32(zBlock), 6)+1)*30) + 35)

			for k := range topY {
				id := STONE // stone
				if k == topY-1 {
					if k <= WATER_LEVEL {
						id = SAND // sand
					} else {
						id = GRASS // grass
					}
				} else if k > topY-5 {
					if k <= WATER_LEVEL {
						id = SAND // sand
					} else {
						id = DIRT // dirt
					}
				}
				blocks[i][k][j] = id
			}

			for k := topY; k < WATER_LEVEL; k++ {
				blocks[i][k][j] = WATER
			}
			if topY <= WATER_LEVEL {
				blocks[i][WATER_LEVEL][j] = WATER
			}
		}
	}

	chunk.setBlocks(blocks)
}
