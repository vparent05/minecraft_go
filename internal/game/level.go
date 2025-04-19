package game

import (
	"iter"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
)

const CHUNK_WIDTH = 15
const CHUNK_HEIGHT = 255

type Chunk struct {
	blocks          [CHUNK_WIDTH * CHUNK_HEIGHT * CHUNK_WIDTH]block // 15 * 255 * 15, index = x * 3825 + y * 15 + z
	solidMesh       []uint32
	transparentMesh []uint32

	Loaded         bool
	Dirty          bool
	SolidVBO       uint32
	TransparentVBO uint32
}

type level struct {
	chunks *utils.MutexMap[mgl32.Vec2, *Chunk]
}

func NewLevel() *level {
	return &level{
		chunks: utils.NewMutexMap[mgl32.Vec2, *Chunk](),
	}
}

func (l *level) Iterator() iter.Seq2[mgl32.Vec2, *Chunk] {
	return func(yield func(mgl32.Vec2, *Chunk) bool) {
		for _, pos := range l.chunks.Keys() {
			if chunk, ok := l.chunks.Get(pos); ok {
				if !yield(pos, chunk) {
					return
				}
			}
		}
	}
}

func t(from, offset, orientation float32) float32 {
	if orientation < 0 {
		return (float32(math.Ceil(float64(from-offset))) - from) / orientation
	} else {
		return (float32(math.Floor(float64(from+offset))) - from) / orientation
	}
}

func (l *level) castRay(position mgl32.Vec3, orientation mgl32.Vec3, length float32) (chunk *Chunk, index int, direction direction) {
	blockPos := mgl32.Vec3{
		float32(math.Floor(float64(position.X()))),
		float32(math.Floor(float64(position.Y()))),
		float32(math.Floor(float64(position.Z()))),
	}
	steps := mgl32.Vec3{
		orientation.X() / float32(math.Abs(float64(orientation.X()))),
		orientation.Y() / float32(math.Abs(float64(orientation.Y()))),
		orientation.Z() / float32(math.Abs(float64(orientation.Z()))),
	}

	curr := position
	dir := -1
	for curr.Sub(position).Len() < length {
		deltaIfX := t(curr.X(), 1, orientation.X())
		deltaIfY := t(curr.Y(), 1, orientation.Y())
		deltaIfZ := t(curr.Z(), 1, orientation.Z())

		if deltaIfX < deltaIfY {
			if deltaIfX < deltaIfZ {
				// move in x
				curr = curr.Add(orientation.Mul(deltaIfX))
				blockPos[0] += steps.X()
				if steps.X() < 0 {
					dir = _RIGHT
				} else {
					dir = _LEFT
				}
			} else {
				// move in z
				curr = curr.Add(orientation.Mul(deltaIfZ))
				blockPos[2] += steps.Z()
				if steps.Z() < 0 {
					dir = _FRONT
				} else {
					dir = _BACK
				}
			}
		} else {
			if deltaIfY < deltaIfZ {
				// move in y
				curr = curr.Add(orientation.Mul(deltaIfY))
				blockPos[1] += steps.Y()
				if steps.Y() < 0 {
					dir = _UP
				} else {
					dir = _DOWN
				}
			} else {
				// move in z
				curr = curr.Add(orientation.Mul(deltaIfZ))
				blockPos[2] += steps.Z()
				if steps.Z() < 0 {
					dir = _FRONT
				} else {
					dir = _BACK
				}
			}
		}

		blockX := ((int(blockPos.X()) % CHUNK_WIDTH) + CHUNK_WIDTH) % CHUNK_WIDTH
		blockZ := ((int(blockPos.Z()) % CHUNK_WIDTH) + CHUNK_WIDTH) % CHUNK_WIDTH
		chunkPos := mgl32.Vec2{
			(blockPos.X() - float32(blockX)) / float32(CHUNK_WIDTH),
			(blockPos.Z() - float32(blockZ)) / float32(CHUNK_WIDTH),
		}

		chunk, ok := l.chunks.Get(chunkPos)
		if !ok {
			continue
		}
		index := chunkIndex(blockX, int(blockPos.Y())%CHUNK_HEIGHT, blockZ)
		if index >= 0 && index < 57375 && chunk.blocks[index].id != 0 {
			return chunk, index, dir
		}
	}
	return nil, -1, -1
}

func (l *level) updateChunksAround(chunkCoords *mgl32.Vec2, renderDistance *int) {
	for chunkCoords != nil {
		for pos, chunk := range l.Iterator() {
			inRenderDistance := mgl32.Abs(pos.X()-chunkCoords.X()) <= float32(*renderDistance) &&
				mgl32.Abs(pos.Y()-chunkCoords.Y()) <= float32(*renderDistance)

			if chunk != nil && !inRenderDistance {
				if chunk.SolidVBO == 0 && chunk.TransparentVBO == 0 {
					l.chunks.Delete(pos)
					continue
				}
				chunk.Loaded = false
			}
		}
		for x := -*renderDistance; x <= *renderDistance; x++ {
			for z := -*renderDistance; z <= *renderDistance; z++ {
				pos := chunkCoords.Add(mgl32.Vec2{float32(x), float32(z)})
				if c, ok := l.chunks.Get(pos); !ok {
					l.chunks.Set(pos, generateChunk(pos))
				} else {
					c.Loaded = true
				}
			}
		}
	}
}

func (c *Chunk) SolidMesh() []uint32 {
	return c.solidMesh
}

func (c *Chunk) TransparentMesh() []uint32 {
	return c.transparentMesh
}

/*
Returns true if block id "id2" is visible through block id "id1"
*/
func visible(id1, id2 uint8) bool {
	return id1 == 0 ||
		BLOCK_TYPES[id1].isTransparent && id1 != id2
}

func chunkIndex(x, y, z int) int {
	return x*CHUNK_WIDTH*CHUNK_HEIGHT + y*CHUNK_WIDTH + z
}

func (c *Chunk) generateMesh() {
	c.solidMesh = make([]uint32, 0)
	c.transparentMesh = make([]uint32, 0)
	for i, b := range c.blocks {
		if b.id == 0 {
			continue
		}
		x := (i / (CHUNK_WIDTH * CHUNK_HEIGHT)) & 0xf
		y := ((i % (CHUNK_WIDTH * CHUNK_HEIGHT)) / CHUNK_WIDTH) & 0xff
		z := (i % CHUNK_WIDTH) & 0xf

		render := [6]bool{
			// (middle AND visible) OR edge
			y+1 < CHUNK_HEIGHT && visible(c.blocks[chunkIndex(x, y+1, z)].id, b.id) || y+1 >= CHUNK_HEIGHT,
			y-1 >= 0 && visible(c.blocks[chunkIndex(x, y-1, z)].id, b.id) || y-1 < 0,

			// (middle AND (visible OR different height level)) OR edge
			x-1 >= 0 && (visible(c.blocks[chunkIndex(x-1, y, z)].id, b.id) || c.blocks[chunkIndex(x-1, y, z)].level != b.level) || x-1 < 0,
			x+1 < CHUNK_WIDTH && (visible(c.blocks[chunkIndex(x+1, y, z)].id, b.id) || c.blocks[chunkIndex(x+1, y, z)].level != b.level) || x+1 >= CHUNK_WIDTH,
			z+1 < CHUNK_WIDTH && (visible(c.blocks[chunkIndex(x, y, z+1)].id, b.id) || c.blocks[chunkIndex(x, y, z+1)].level != b.level) || z+1 >= CHUNK_WIDTH,
			z-1 >= 0 && (visible(c.blocks[chunkIndex(x, y, z-1)].id, b.id) || c.blocks[chunkIndex(x, y, z-1)].level != b.level) || z-1 < 0,
		}

		if BLOCK_TYPES[b.id].isTransparent {
			c.transparentMesh = append(c.transparentMesh, b.mesh(x, y, z, render)...)
		} else {
			c.solidMesh = append(c.solidMesh, b.mesh(x, y, z, render)...)
		}
	}
	c.Dirty = true
}

func (c *Chunk) set(index int, value block) {
	c.blocks[index] = value
	c.generateMesh()
}
