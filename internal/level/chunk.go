package level

import (
	"iter"
	"math"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
	"github.com/vparent05/minecraft_go/internal/utils/atomicx"
	"github.com/vparent05/minecraft_go/internal/utils/chanx"
)

const CHUNK_WIDTH = 15
const CHUNK_HEIGHT = 255

type ChunkMesh struct {
	Solid       []uint32
	Transparent []uint32
}

type chunkSnapshot struct {
	coordinates utils.IntVector2
	blocks      [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_WIDTH]BlockId
	observer    utils.IntVector3
}

// Using the exported members of Chunk is fully thread safe
type Chunk struct {
	mu            sync.Mutex
	generator     *meshGenerator
	coordinates   utils.IntVector2
	observer      *atomicx.Value[LevelObserver] // Coordinates of the block closest to the level observer in the chunk
	observerCache utils.IntVector3
	blocks        [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_WIDTH]BlockId
	Slot          int

	Mesh        *atomicx.Value[ChunkMesh]
	MeshUpdates chan struct{}
}

func newChunk(generator *meshGenerator, observer *atomicx.Value[LevelObserver]) *Chunk {
	c := &Chunk{
		generator:   generator,
		observer:    observer,
		Mesh:        &atomicx.Value[ChunkMesh]{},
		MeshUpdates: make(chan struct{}, 1),
	}

	go func() {
		for {
			c.updateObserverCache()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return c
}

func (c *Chunk) snapshot() chunkSnapshot {
	c.mu.Lock()
	defer c.mu.Unlock()

	return chunkSnapshot{
		coordinates: c.coordinates,
		blocks:      c.blocks,
		observer:    c.observerCache,
	}
}

func (c *chunkSnapshot) iter() iter.Seq2[utils.IntVector3, BlockId] {
	return func(yield func(utils.IntVector3, BlockId) bool) {
		for pos, b := range utils.UnsafeFromOriginIterator3(
			&c.blocks[0][0][0],
			utils.IntVector3{X: c.coordinates.X * CHUNK_WIDTH, Y: 0, Z: c.coordinates.Y * CHUNK_WIDTH},                        // Start
			utils.IntVector3{X: (c.coordinates.X + 1) * CHUNK_WIDTH, Y: CHUNK_HEIGHT, Z: (c.coordinates.Y + 1) * CHUNK_WIDTH}, // End
			utils.IntVector3{X: CHUNK_WIDTH, Y: CHUNK_HEIGHT, Z: CHUNK_WIDTH},                                                 // Size
			c.observer, // Origin
		) {

			if !yield(pos, b) {
				return
			}
		}
	}
}

func (c *Chunk) setContent(coordinates utils.IntVector2, blocks [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_WIDTH]BlockId) {
	c.mu.Lock()
	c.coordinates = coordinates
	c.blocks = blocks
	c.mu.Unlock()

	c.generator.enqueue(c)
}

func visibleOrDifferentHeightLevel(a, b BlockId) bool {
	return visible(a, b) || BLOCK_TYPES[a].height != BLOCK_TYPES[b].height
}

func (c *Chunk) generateMesh(level *Level) {
	snap := c.snapshot()
	mesh := ChunkMesh{make([]uint32, 0), make([]uint32, 0)}

	for pos, b := range snap.iter() {
		if b == AIR {
			continue
		}
		x := pos.X
		y := pos.Y
		z := pos.Z

		render := [6]bool{}

		// (middle AND visible) OR edge
		render[0] = y+1 < CHUNK_HEIGHT && visible(snap.blocks[x][y+1][z], b) || y+1 >= CHUNK_HEIGHT
		render[1] = y-1 >= 0 && visible(snap.blocks[x][y-1][z], b) || y-1 < 0

		// visible OR different height level
		levelX := snap.coordinates.X*CHUNK_WIDTH + x
		levelZ := snap.coordinates.Y*CHUNK_WIDTH + z
		if x-1 >= 0 {
			render[2] = visibleOrDifferentHeightLevel(snap.blocks[x-1][y][z], b)
		} else if a, ok := level.getBlockPosition(mgl32.Vec3{float32(levelX - 1), float32(y), float32(levelZ)}).Get(); ok {
			render[2] = visibleOrDifferentHeightLevel(a, b)
		}

		if x+1 < CHUNK_WIDTH {
			render[3] = visibleOrDifferentHeightLevel(snap.blocks[x+1][y][z], b)
		} else if a, ok := level.getBlockPosition(mgl32.Vec3{float32(levelX + 1), float32(y), float32(levelZ)}).Get(); ok {
			render[3] = visibleOrDifferentHeightLevel(a, b)
		}

		if z+1 < CHUNK_WIDTH {
			render[4] = visibleOrDifferentHeightLevel(snap.blocks[x][y][z+1], b)
		} else if a, ok := level.getBlockPosition(mgl32.Vec3{float32(levelX), float32(y), float32(levelZ + 1)}).Get(); ok {
			render[4] = visibleOrDifferentHeightLevel(a, b)
		}

		if z-1 >= 0 {
			render[5] = visibleOrDifferentHeightLevel(snap.blocks[x][y][z-1], b)
		} else if a, ok := level.getBlockPosition(mgl32.Vec3{float32(levelX), float32(y), float32(levelZ - 1)}).Get(); ok {
			render[5] = visibleOrDifferentHeightLevel(a, b)
		}

		if BLOCK_TYPES[b].isTransparent {
			mesh.Transparent = append(mesh.Transparent, b.mesh(x, y, z, render)...)
		} else {
			mesh.Solid = append(mesh.Solid, b.mesh(x, y, z, render)...)
		}
	}

	c.Mesh.Store(mesh)
	chanx.TrySend(c.MeshUpdates, struct{}{})
}

func (c *Chunk) updateObserverCache() {
	observer := c.observer.Load().Vec3

	newObserverCache := utils.IntVector3{
		X: int(math.Floor(math.Max(math.Min(float64(observer.X()), float64((c.coordinates.X+1)*CHUNK_WIDTH)), float64(c.coordinates.X*CHUNK_WIDTH)))),
		Y: int(math.Floor(math.Max(math.Min(float64(observer.Y()), CHUNK_HEIGHT), 0))),
		Z: int(math.Floor(math.Max(math.Min(float64(observer.Z()), float64((c.coordinates.Y+1)*CHUNK_WIDTH)), float64(c.coordinates.Y*CHUNK_WIDTH)))),
	}

	if c.observerCache != newObserverCache {
		c.observerCache = newObserverCache
		c.generator.enqueue(c)
	}
}

func (c *Chunk) getBlock(coordinates utils.IntVector3) BlockId {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.blocks[coordinates.X][coordinates.Y][coordinates.Z]
}

func (c *Chunk) setBlock(coordinates utils.IntVector3, value BlockId) {
	c.mu.Lock()
	c.blocks[coordinates.X][coordinates.Y][coordinates.Z] = value
	c.mu.Unlock()

	c.generator.enqueue(c)
}
