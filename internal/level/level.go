package level

type block struct {
	id      uint8
	visible bool
}

type mesh struct {
	vertices []uint8
	indices  []uint16
	normals  []bool
	uv       []uint8
}

type chunk struct {
	blocks [65536]block // 16 * 256 * 16, index = x * 4096 + y * 16 + z

	transparentMesh mesh
	solidMesh       mesh
}

func (c *chunk) generateMesh() {
	vertices := map[(float32, float32, float32)]struct{}{}
	for i, b := range c.blocks {
		if b == nil || !b.visible {
			continue
		}
		x := i / 4096
		y := (i % 4096) / 16
		z := i % 16

		vertices[(x, y, z)] = struct{}{}
		vertices[(x, y, z+1)] = struct{}{}
		vertices[(x, y+1, z)] = struct{}{}
		vertices[(x, y+1, z+1)] = struct{}{}
		vertices[(x+1, y, z)] = struct{}{}
		vertices[(x+1, y, z+1)] = struct{}{}
		vertices[(x+1, y+1, z)] = struct{}{}
		vertices[(x+1, y+1, z+1)] = struct{}{}


	}
}
