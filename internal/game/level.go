package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type block struct {
	id      uint8
	visible bool
}

func (b *block) mesh(x, y, z int) []uint32 {
	id := int(b.id << 4)

	// top
	v11 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | id)
	v12 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | id)
	v13 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | id)

	v21 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | id)
	v22 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | id)
	v23 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | id)

	// bottom
	v31 := uint32(x<<28 | (y)<<20 | z<<16 | 0x1<<12 | id)
	v32 := uint32((x+1)<<28 | (y)<<20 | z<<16 | 0x1<<12 | id)
	v33 := uint32((x+1)<<28 | (y)<<20 | (z+1)<<16 | 0x1<<12 | id)

	v41 := uint32((x+1)<<28 | (y)<<20 | (z+1)<<16 | 0x1<<12 | id)
	v42 := uint32(x<<28 | (y)<<20 | (z+1)<<16 | 0x1<<12 | id)
	v43 := uint32(x<<28 | (y)<<20 | z<<16 | 0x1<<12 | id)

	// left
	v51 := uint32(x<<28 | y<<20 | z<<16 | 0x2<<12 | id)
	v52 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x2<<12 | id)
	v53 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x2<<12 | id)

	v61 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x2<<12 | id)
	v62 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x2<<12 | id)
	v63 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x2<<12 | id)

	// right
	v71 := uint32((x+1)<<28 | y<<20 | z<<16 | 0x3<<12 | id)
	v72 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x3<<12 | id)
	v73 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x3<<12 | id)

	v81 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x3<<12 | id)
	v82 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x3<<12 | id)
	v83 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x3<<12 | id)

	// front
	v91 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | id)
	v92 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | id)
	v93 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | id)

	vA1 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | id)
	vA2 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | id)
	vA3 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | id)

	// back
	vB1 := uint32(x<<28 | y<<20 | z<<16 | 0x5<<12 | id)
	vB2 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | id)
	vB3 := uint32((x+1)<<28 | y<<20 | z<<16 | 0x5<<12 | id)

	vC1 := uint32(x<<28 | y<<20 | z<<16 | 0x5<<12 | id)
	vC2 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | id)
	vC3 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | id)

	return []uint32{
		v11, v12, v13,
		v21, v22, v23,
		v31, v32, v33,
		v41, v42, v43,
		v51, v52, v53,
		v61, v62, v63,
		v71, v72, v73,
		v81, v82, v83,
		v91, v92, v93,
		vA1, vA2, vA3,
		vB1, vB2, vB3,
		vC1, vC2, vC3,
	}
}

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
