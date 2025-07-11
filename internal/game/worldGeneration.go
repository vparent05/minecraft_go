package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
)

func generateChunk(pos mgl32.Vec2) *Chunk {
	const WATER_LEVEL = 60
	chunk := Chunk{
		pos,
		[CHUNK_WIDTH * CHUNK_HEIGHT * CHUNK_WIDTH]block{},
		nil,
		nil,
		true,
		false,
		0,
		0,
	}
	x := int(pos.X())
	z := int(pos.Y())
	for i := range CHUNK_WIDTH {
		for j := range CHUNK_WIDTH {
			xBlock := i + x*CHUNK_WIDTH
			zBlock := j + z*CHUNK_WIDTH
			topY := int(math.Floor(float64(utils.FractalNoise2(float32(xBlock), float32(zBlock), 6)+1)*30) + 35)

			for k := range topY {
				id := 6 // stone
				if k == topY-1 {
					if k <= WATER_LEVEL {
						id = 4 // sand
					} else {
						id = 1 // grass
					}
				} else if k > topY-5 {
					if k <= WATER_LEVEL {
						id = 4 // sand
					} else {
						id = 5 // dirt
					}
				}
				chunk.blocks[indexInChunk(i, k, j)] = block{uint8(id), 15}
			}

			for k := topY; k < WATER_LEVEL; k++ {
				chunk.blocks[indexInChunk(i, k, j)] = block{3, 15}
			}
			if topY <= WATER_LEVEL {
				chunk.blocks[indexInChunk(i, WATER_LEVEL, j)] = block{3, 13}
			}
		}
	}
	chunk.generateMesh()
	return &chunk
}
