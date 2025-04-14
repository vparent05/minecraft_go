package game

type blockType struct {
	name          string
	isTransparent bool
	isLiquid      bool
	viscosity     float32
	textureRight  string
	textureLeft   string
	textureTop    string
	textureBottom string
	textureFront  string
	textureBack   string
}

type block struct {
	id      uint8
	visible bool
}

var BLOCK_TEXTURE_ATLAS map[string]uint8

var BLOCK_TYPES map[int]blockType

const atlasWidth = 16

func LoadBlocks() {
	BLOCK_TYPES = map[int]blockType{
		0: {
			"grass",
			false,
			false,
			1.0,
			"grass_side.png",
			"grass_side.png",
			"grass_top.png",
			"grass_bottom.png",
			"grass_side.png",
			"grass_side.png",
		},
	}
}

func (b *block) mesh(x, y, z int) []uint32 {
	// one vertex is encoded as: x (4bits) | y (8bits) | z (4bits) | orientation (4bits) | texture coordinate (x + atlasWidth*y) (12bits)

	id := int(b.id)

	// top
	v11 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]+atlasWidth+1))
	v12 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]))
	v13 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]+atlasWidth))

	v21 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]+atlasWidth+1))
	v22 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]+1))
	v23 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x0<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureTop]))

	// bottom
	v31 := uint32(x<<28 | y<<20 | z<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]+1))
	v32 := uint32((x+1)<<28 | y<<20 | z<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]))
	v33 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]+atlasWidth))

	v41 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]+atlasWidth))
	v42 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]+1+atlasWidth))
	v43 := uint32(x<<28 | y<<20 | z<<16 | 0x1<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBottom]+1))

	// left
	v51 := uint32(x<<28 | y<<20 | z<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]+atlasWidth))
	v52 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]+1+atlasWidth))
	v53 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]))

	v61 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]+1+atlasWidth))
	v62 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]+1))
	v63 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x2<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureLeft]))

	// right
	v71 := uint32((x+1)<<28 | y<<20 | z<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]+1+atlasWidth))
	v72 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]+1))
	v73 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]+atlasWidth))

	v81 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]+atlasWidth))
	v82 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]+1))
	v83 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x3<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureRight]))

	// front
	v91 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]+atlasWidth))
	v92 := uint32((x+1)<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]+atlasWidth+1))
	v93 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]+1))

	vA1 := uint32(x<<28 | y<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]+atlasWidth))
	vA2 := uint32((x+1)<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]+1))
	vA3 := uint32(x<<28 | (y+1)<<20 | (z+1)<<16 | 0x4<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureFront]))

	// back
	vB1 := uint32(x<<28 | y<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]+1+atlasWidth))
	vB2 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]))
	vB3 := uint32((x+1)<<28 | y<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]+atlasWidth))

	vC1 := uint32(x<<28 | y<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]+1+atlasWidth))
	vC2 := uint32(x<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]+1))
	vC3 := uint32((x+1)<<28 | (y+1)<<20 | z<<16 | 0x5<<12 | int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[id].textureBack]))

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
