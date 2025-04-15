package game

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	position mgl32.Vec2
	blocks   [65536]block // 16 * 256 * 16, index = x * 4096 + y * 16 + z

	solidMesh       []uint32
	transparentMesh []uint32
}

/*
Returns true if block id "id2" is visible through block id "id1"
*/
func visible(id1, id2 uint8) bool {
	return id1 == 0 ||
		BLOCK_TYPES[id1].isTransparent && id1 != id2
}

func (c *Chunk) generateMesh() {
	c.solidMesh = make([]uint32, 0)
	c.transparentMesh = make([]uint32, 0)
	for i, b := range c.blocks {
		if b.id == 0 {
			continue
		}
		x := (i / 4096) & 0xf
		y := ((i % 4096) / 16) & 0xff
		z := (i % 16) & 0xf

		render := [6]bool{
			// (middle AND visible) OR edge
			i+16 < 65536 && visible(c.blocks[i+16].id, b.id) || i+16 >= 65536,
			i-16 >= 0 && visible(c.blocks[i-16].id, b.id) || i-16 < 0,

			// (middle AND (visible OR different height level)) OR edge
			i-4096 >= 0 && (visible(c.blocks[i-4096].id, b.id) || c.blocks[i-4096].level != b.level) || i-4096 < 0,
			i+4096 < 65536 && (visible(c.blocks[i+4096].id, b.id) || c.blocks[i+4096].level != b.level) || i+4096 >= 65536,
			i+1 < 65536 && (visible(c.blocks[i+1].id, b.id) || c.blocks[i+1].level != b.level) || i+1 >= 65536,
			i-1 >= 0 && (visible(c.blocks[i-1].id, b.id) || c.blocks[i-1].level != b.level) || i-1 < 0,
		}

		if BLOCK_TYPES[b.id].isTransparent {
			c.transparentMesh = append(c.transparentMesh, b.mesh(x, y, z, render)...)
		} else {
			c.solidMesh = append(c.solidMesh, b.mesh(x, y, z, render)...)
		}
	}
}

func (c *Chunk) SolidMesh() []uint32 {
	return c.solidMesh
}

func (c *Chunk) TransparentMesh() []uint32 {
	return c.transparentMesh
}

func GetTestChunk() *Chunk {
	chunk := Chunk{
		mgl32.Vec2{0, 0},
		[65536]block{},
		nil,
		nil,
	}
	for i := range 15 {
		for j := range 15 {
			id := uint8(rand.Int()%3 + 1)
			var b block
			if id == 3 {
				b = block{id, 13}
			} else {
				b = block{id, rand.Int() % 16}
			}
			chunk.blocks[i*4096+j] = b
		}
	}
	chunk.generateMesh()
	return &chunk
}
