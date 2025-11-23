package level

import (
	"math"

	"github.com/vparent05/minecraft_go/internal/utils"
)

func generateChunk(pos utils.IntVector2) *Chunk {
	const WATER_LEVEL = 60
	chunk := newChunk(pos)

	for i := range CHUNK_WIDTH {
		for j := range CHUNK_WIDTH {
			xBlock := i + pos.X*CHUNK_WIDTH
			zBlock := j + pos.Y*CHUNK_WIDTH
			topY := int(math.Floor(float64(utils.FractalNoise2(float32(xBlock), float32(zBlock), 6)+1)*30) + 35)

			for k := range topY {
				id := STONE // stone
				if k == topY-1 {
					if k <= WATER_LEVEL {
						id = SAND // sand
					} else {
						id = GRASS // grass
					}
				} else if k > topY-5 {
					if k <= WATER_LEVEL {
						id = SAND // sand
					} else {
						id = DIRT // dirt
					}
				}
				chunk.blocks[i][k][j] = id
			}

			for k := topY; k < WATER_LEVEL; k++ {
				chunk.blocks[i][k][j] = WATER
			}
			if topY <= WATER_LEVEL {
				chunk.blocks[i][WATER_LEVEL][j] = WATER
			}
		}
	}
	chunk.generateMesh()
	return chunk
}
