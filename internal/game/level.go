package game

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	position mgl32.Vec2
	blocks   [57375]block // 15 * 255 * 15, index = x * 3825 + y * 15 + z

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
		x := (i / 3825) & 0xf
		y := ((i % 3825) / 15) & 0xff
		z := (i % 15) & 0xf

		render := [6]bool{
			// (middle AND visible) OR edge
			i+15 < 57375 && visible(c.blocks[i+15].id, b.id) || i+15 >= 57375,
			i-15 >= 0 && visible(c.blocks[i-15].id, b.id) || i-15 < 0,

			// (middle AND (visible OR different height level)) OR edge
			i-3825 >= 0 && (visible(c.blocks[i-3825].id, b.id) || c.blocks[i-3825].level != b.level) || i-3825 < 0,
			i+3825 < 57375 && (visible(c.blocks[i+3825].id, b.id) || c.blocks[i+3825].level != b.level) || i+3825 >= 57375,
			i+1 < 57375 && (visible(c.blocks[i+1].id, b.id) || c.blocks[i+1].level != b.level) || i+1 >= 57375,
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

func (c *Chunk) Position() mgl32.Vec2 {
	return c.position
}

func GetTestChunk(i int) *Chunk {
	chunk := Chunk{
		mgl32.Vec2{float32(i % 16), float32(i / 16)},
		[57375]block{},
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
				b = block{id, 15}
			}
			chunk.blocks[i*3825+j] = b
		}
	}
	chunk.generateMesh()
	return &chunk
}
