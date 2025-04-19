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
	id    uint8
	level int
}

var BLOCK_TEXTURE_ATLAS map[string]uint8

var BLOCK_TYPES map[uint8]blockType

const atlasWidth = 16

func LoadBlocks() {
	BLOCK_TYPES = map[uint8]blockType{
		1: {
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
		2: {
			"glass",
			true,
			false,
			1.0,
			"glass.png",
			"glass.png",
			"glass.png",
			"glass.png",
			"glass.png",
			"glass.png",
		},
		3: {
			"water",
			true,
			true,
			0.5,
			"water.png",
			"water.png",
			"water.png",
			"water.png",
			"water.png",
			"water.png",
		},
		4: {
			"sand",
			false,
			false,
			1.0,
			"sand.png",
			"sand.png",
			"sand.png",
			"sand.png",
			"sand.png",
			"sand.png",
		},
		5: {
			"dirt",
			false,
			false,
			1.0,
			"dirt.png",
			"dirt.png",
			"dirt.png",
			"dirt.png",
			"dirt.png",
			"dirt.png",
		},
		6: {
			"stone",
			false,
			false,
			1.0,
			"stone.png",
			"stone.png",
			"stone.png",
			"stone.png",
			"stone.png",
			"stone.png",
		},
	}
}

func (b *block) mesh(x, y, z int, render [6]bool) []uint32 {
	// one vertex is encoded as: x (4bits) | y (8bits) | z (4bits) | orientation (4bits) | texture coordinate (x + atlasWidth*y) (8bits) | level (height: (level+1) / 16) (4bits)
	if b.id == 0 {
		return []uint32{}
	}

	mesh := make([]uint32, 0, 36)

	// top
	if render[0] {
		mesh = append(mesh,
			uint32(x<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop]+atlasWidth+1)<<4|b.level),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop])<<4|b.level),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop]+atlasWidth)<<4|b.level),

			uint32(x<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop]+atlasWidth+1)<<4|b.level),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop]+1)<<4|b.level),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureTop])<<4|b.level),
		)
	}

	// bottom
	if render[1] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom]+1)<<4|15),
			uint32((x+1)<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom])<<4|15),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom]+atlasWidth)<<4|15),

			uint32((x+1)<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom]+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom]+1+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBottom]+1)<<4|15),
		)
	}

	// left
	if render[2] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft]+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft])<<4|b.level),

			uint32(x<<28|y<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft]+1)<<4|b.level),
			uint32(x<<28|(y+1)<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureLeft])<<4|b.level),
		)
	}

	// right
	if render[3] {
		mesh = append(mesh,
			uint32((x+1)<<28|y<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight]+1+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight]+1)<<4|b.level),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight]+atlasWidth)<<4|15),

			uint32((x+1)<<28|y<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight]+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight]+1)<<4|b.level),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureRight])<<4|b.level),
		)
	}

	// front
	if render[4] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront]+atlasWidth)<<4|15),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront]+atlasWidth+1)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront]+1)<<4|b.level),

			uint32(x<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront]+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront]+1)<<4|b.level),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureFront])<<4|b.level),
		)
	}

	// back
	if render[5] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack]+1+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack])<<4|b.level),
			uint32((x+1)<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack]+atlasWidth)<<4|15),

			uint32(x<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack]+1)<<4|b.level),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b.id].textureBack])<<4|b.level),
		)
	}
	return mesh
}
