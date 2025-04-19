package game

import (
	"iter"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
)

type Chunk struct {
	blocks          [57375]block // 15 * 255 * 15, index = x * 3825 + y * 15 + z
	solidMesh       []uint32
	transparentMesh []uint32

	Loaded         bool
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
	return x*3825 + y*15 + z
}

func (c *Chunk) generateMesh() {
	c.solidMesh = make([]uint32, 0)
	c.transparentMesh = make([]uint32, 0)
	for i, b := range c.blocks {
		if b.id == 0 {
			continue
		}
		x := (i / 3825) & 0xf
		y := ((i % 3825) / 15) & 0xff
		z := (i % 15) & 0xf

		render := [6]bool{
			// (middle AND visible) OR edge
			y+1 < 255 && visible(c.blocks[chunkIndex(x, y+1, z)].id, b.id) || y+1 >= 255,
			y-1 >= 0 && visible(c.blocks[chunkIndex(x, y-1, z)].id, b.id) || y-1 < 0,

			// (middle AND (visible OR different height level)) OR edge
			x-1 >= 0 && (visible(c.blocks[chunkIndex(x-1, y, z)].id, b.id) || c.blocks[chunkIndex(x-1, y, z)].level != b.level) || x-1 < 0,
			x+1 < 15 && (visible(c.blocks[chunkIndex(x+1, y, z)].id, b.id) || c.blocks[chunkIndex(x+1, y, z)].level != b.level) || x+1 >= 15,
			z+1 < 15 && (visible(c.blocks[chunkIndex(x, y, z+1)].id, b.id) || c.blocks[chunkIndex(x, y, z+1)].level != b.level) || z+1 >= 15,
			z-1 >= 0 && (visible(c.blocks[chunkIndex(x, y, z-1)].id, b.id) || c.blocks[chunkIndex(x, y, z-1)].level != b.level) || z-1 < 0,
		}

		if BLOCK_TYPES[b.id].isTransparent {
			c.transparentMesh = append(c.transparentMesh, b.mesh(x, y, z, render)...)
		} else {
			c.solidMesh = append(c.solidMesh, b.mesh(x, y, z, render)...)
		}
	}
}
