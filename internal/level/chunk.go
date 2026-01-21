package level

import (
	"iter"
	"math"
	"sync"
	"sync/atomic"

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

// Using the exported members of Chunk is fully thread safe
type Chunk struct {
	mu          sync.Mutex
	coordinates utils.IntVector2
	observer    utils.IntVector3 // Coordinates of the block closest to the level observer in the chunk TODO use it
	blocks      [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_WIDTH]BlockId
	Slot        int

	Dirty       atomic.Bool
	Mesh        *atomicx.Value[ChunkMesh]
	MeshUpdates chan struct{}
}

func newChunk(coordinates utils.IntVector2) *Chunk {
	c := &Chunk{
		coordinates: coordinates,
		observer:    utils.IntVector3{}, // TODO properly set
		Mesh:        &atomicx.Value[ChunkMesh]{},
		MeshUpdates: make(chan struct{}, 1),
	}
	c.Dirty.Store(true)
	return c
}

func (c *Chunk) iter() iter.Seq2[utils.IntVector3, BlockId] {
	return func(yield func(utils.IntVector3, BlockId) bool) {
		for pos, b := range utils.UnsafeFromOriginIterator3(&c.blocks[0][0][0], utils.IntVector3{X: CHUNK_WIDTH, Y: CHUNK_HEIGHT, Z: CHUNK_WIDTH}, c.observer) {
			if !yield(pos, b) {
				return
			}
		}
	}
}

func (c *Chunk) generateMesh() {
	c.mu.Lock()
	defer c.mu.Unlock()

	mesh := ChunkMesh{make([]uint32, 0), make([]uint32, 0)}

	for pos, b := range c.iter() {
		if b == AIR {
			continue
		}
		x := pos.X
		y := pos.Y
		z := pos.Z

		render := [6]bool{
			// (middle AND visible) OR edge
			y+1 < CHUNK_HEIGHT && visible(c.blocks[x][y+1][z], b) || y+1 >= CHUNK_HEIGHT,
			y-1 >= 0 && visible(c.blocks[x][y-1][z], b) || y-1 < 0,

			// (middle AND (visible OR different height level)) OR edge
			x-1 >= 0 && (visible(c.blocks[x-1][y][z], b) || BLOCK_TYPES[c.blocks[x-1][y][z]].height != BLOCK_TYPES[b].height) || x-1 < 0,
			x+1 < CHUNK_WIDTH && (visible(c.blocks[x+1][y][z], b) || BLOCK_TYPES[c.blocks[x+1][y][z]].height != BLOCK_TYPES[b].height) || x+1 >= CHUNK_WIDTH,
			z+1 < CHUNK_WIDTH && (visible(c.blocks[x][y][z+1], b) || BLOCK_TYPES[c.blocks[x][y][z+1]].height != BLOCK_TYPES[b].height) || z+1 >= CHUNK_WIDTH,
			z-1 >= 0 && (visible(c.blocks[x][y][z-1], b) || BLOCK_TYPES[c.blocks[x][y][z-1]].height != BLOCK_TYPES[b].height) || z-1 < 0,
		}

		if BLOCK_TYPES[b].isTransparent {
			mesh.Transparent = append(mesh.Transparent, b.mesh(x, y, z, render)...)
		} else {
			mesh.Solid = append(mesh.Solid, b.mesh(x, y, z, render)...)
		}
	}

	c.Mesh.Store(mesh)
	chanx.TrySend(c.MeshUpdates, struct{}{})

	c.Dirty.Store(false)
}

func (c *Chunk) setObserver(observer mgl32.Vec3) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.observer = utils.IntVector3{
		X: int(math.Floor(math.Max(math.Min(float64(observer.X()), float64((c.coordinates.X+1)*CHUNK_WIDTH)), float64(c.coordinates.X*CHUNK_WIDTH)))),
		Y: int(math.Floor(float64(observer.Y()))),
		Z: int(math.Floor(math.Max(math.Min(float64(observer.Z()), float64((c.coordinates.Y+1)*CHUNK_WIDTH)), float64(c.coordinates.Y*CHUNK_WIDTH)))),
	}
	c.Dirty.Store(true)
}

func (c *Chunk) getBlock(coordinates utils.IntVector3) BlockId {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.blocks[coordinates.X][coordinates.Y][coordinates.Z]
}

func (c *Chunk) setBlock(coordinates utils.IntVector3, value BlockId) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.blocks[coordinates.X][coordinates.Y][coordinates.Z] = value
	c.Dirty.Store(true)
}
