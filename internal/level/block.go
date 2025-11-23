package level

type BlockId uint8

const (
	AIR BlockId = iota
	GRASS
	GLASS
	WATER
	SAND
	DIRT
	STONE
)

type blockType struct {
	name          string
	height        int
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

var BLOCK_TEXTURE_ATLAS map[string]BlockId

var BLOCK_TYPES map[BlockId]blockType

const atlasWidth = 16

func LoadBlocks() {
	BLOCK_TYPES = map[BlockId]blockType{
		GRASS: {
			"grass",
			15,
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
		GLASS: {
			"glass",
			15,
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
		WATER: {
			"water",
			13,
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
		SAND: {
			"sand",
			15,
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
		DIRT: {
			"dirt",
			15,
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
		STONE: {
			"stone",
			15,
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

func (b BlockId) mesh(x, y, z int, render [6]bool) []uint32 {
	// one vertex is encoded as: x (4bits) | y (8bits) | z (4bits) | orientation (4bits) | texture coordinate (x + atlasWidth*y) (8bits) | height (real height: (height+1) / 16) (4bits)
	if b == AIR {
		return []uint32{}
	}

	mesh := make([]uint32, 0, 36)

	// top
	if render[0] {
		mesh = append(mesh,
			uint32(x<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop]+atlasWidth+1)<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop])<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop]+atlasWidth)<<4|BLOCK_TYPES[b].height),

			uint32(x<<28|(y+1)<<20|z<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop]+atlasWidth+1)<<4|BLOCK_TYPES[b].height),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop]+1)<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x0<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureTop])<<4|BLOCK_TYPES[b].height),
		)
	}

	// bottom
	if render[1] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom]+1)<<4|15),
			uint32((x+1)<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom])<<4|15),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom]+atlasWidth)<<4|15),

			uint32((x+1)<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom]+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|(z+1)<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom]+1+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|z<<16|0x1<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBottom]+1)<<4|15),
		)
	}

	// left
	if render[2] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft]+atlasWidth)<<4|15),
			uint32(x<<28|y<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft])<<4|BLOCK_TYPES[b].height),

			uint32(x<<28|y<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft]+1)<<4|BLOCK_TYPES[b].height),
			uint32(x<<28|(y+1)<<20|z<<16|0x2<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureLeft])<<4|BLOCK_TYPES[b].height),
		)
	}

	// right
	if render[3] {
		mesh = append(mesh,
			uint32((x+1)<<28|y<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight]+1+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight]+1)<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight]+atlasWidth)<<4|15),

			uint32((x+1)<<28|y<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight]+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight]+1)<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x3<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureRight])<<4|BLOCK_TYPES[b].height),
		)
	}

	// front
	if render[4] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront]+atlasWidth)<<4|15),
			uint32((x+1)<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront]+atlasWidth+1)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront]+1)<<4|BLOCK_TYPES[b].height),

			uint32(x<<28|y<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront]+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront]+1)<<4|BLOCK_TYPES[b].height),
			uint32(x<<28|(y+1)<<20|(z+1)<<16|0x4<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureFront])<<4|BLOCK_TYPES[b].height),
		)
	}

	// back
	if render[5] {
		mesh = append(mesh,
			uint32(x<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack]+1+atlasWidth)<<4|15),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack])<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack]+atlasWidth)<<4|15),

			uint32(x<<28|y<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack]+1+atlasWidth)<<4|15),
			uint32(x<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack]+1)<<4|BLOCK_TYPES[b].height),
			uint32((x+1)<<28|(y+1)<<20|z<<16|0x5<<12|int(BLOCK_TEXTURE_ATLAS[BLOCK_TYPES[b].textureBack])<<4|BLOCK_TYPES[b].height),
		)
	}
	return mesh
}

// Visible returns true if block id back is visible through block id front
func visible(front, back BlockId) bool {
	return front == 0 ||
		BLOCK_TYPES[front].isTransparent && front != back
}
