package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Chunk struct {
	position mgl32.Vec2
	blocks   [65536]block // 16 * 256 * 16, index = x * 4096 + y * 16 + z

	solidMesh []uint32
}

func (c *Chunk) generateMesh() {
	solidMesh := make([]uint32, 0)
	for i, b := range c.blocks {
		if !b.visible {
			continue
		}
		x := (i / 4096) & 0xf
		y := ((i % 4096) / 16) & 0xff
		z := (i % 16) & 0xf

		solidMesh = append(solidMesh, b.mesh(x, y, z)...)
	}

	c.solidMesh = solidMesh
}

func (c *Chunk) SolidMesh() []uint32 {
	return c.solidMesh
}

func GetTestChunk() *Chunk {
	chunk := Chunk{
		mgl32.Vec2{0, 0},
		[65536]block{},
		nil,
	}
	for i := range 15 {
		for j := range 15 {
			b := block{0, true}
			chunk.blocks[i*4096+j] = b
		}
	}
	chunk.generateMesh()
	return &chunk
}
