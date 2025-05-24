package game

import (
	"iter"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
)

const CHUNK_WIDTH = 15
const CHUNK_HEIGHT = 255

type Chunk struct {
	coordinates     mgl32.Vec2
	blocks          [CHUNK_WIDTH * CHUNK_HEIGHT * CHUNK_WIDTH]block // 15 * 255 * 15, index = x * 3825 + y * 15 + z
	solidMesh       []uint32
	transparentMesh []uint32

	Loaded         bool
	Dirty          bool
	SolidVBO       uint32
	TransparentVBO uint32
}

type level struct {
	observer *mgl32.Vec2
	chunks   [][]*Chunk
}

func (l *level) originIndex() (int, int) {
	return utils.Mod(int(l.observer.X()), len(l.chunks[0])), utils.Mod(int(l.observer.Y()), len(l.chunks))
}

func (l *level) index(coordinates mgl32.Vec2) (int, int) {
	delta := coordinates.Sub(*l.observer)
	originX, originZ := l.originIndex()
	return utils.Mod(int(delta.X())+originX, len(l.chunks)), utils.Mod(int(delta.Y())+originZ, len(l.chunks[0]))
}

func (l *level) get(chunkCoordinates mgl32.Vec2) *Chunk {
	i, j := l.index(chunkCoordinates)
	return l.chunks[i][j]
}

func (l *level) set(chunkCoordinates mgl32.Vec2, value *Chunk) {
	i, j := l.index(chunkCoordinates)
	l.chunks[i][j] = value
}

func NewLevel(renderDistance int, observer *mgl32.Vec2) *level {
	l := &level{
		chunks:   make([][]*Chunk, renderDistance*2+1),
		observer: observer,
	}
	for i := range l.chunks {
		l.chunks[i] = make([]*Chunk, renderDistance*2+1)
	}
	return l
}

func (l *level) Iterator() iter.Seq2[mgl32.Vec2, *Chunk] {
	return func(yield func(mgl32.Vec2, *Chunk) bool) {
		originX, originY := l.originIndex()
		for _, chunk := range utils.FromOriginIterator2(l.chunks, originX, originY) {
			if chunk == nil {
				continue
			}
			if !yield(chunk.coordinates, chunk) {
				return
			}
		}
	}
}

func (l *level) positionInChunk(position mgl32.Vec3) (*Chunk, int) {
	blockX := utils.Mod(int(position.X()), CHUNK_WIDTH)
	blockZ := utils.Mod(int(position.Z()), CHUNK_WIDTH)
	chunkPos := mgl32.Vec2{
		(position.X() - float32(blockX)) / float32(CHUNK_WIDTH),
		(position.Z() - float32(blockZ)) / float32(CHUNK_WIDTH),
	}

	chunk := l.get(chunkPos)
	if chunk == nil {
		return nil, -1
	}
	index := indexInChunk(blockX, utils.Mod(int(position.Y()), CHUNK_HEIGHT), blockZ)
	return chunk, index
}

func t(from, offset, orientation float32) float32 {
	if orientation < 0 {
		return (float32(math.Ceil(float64(from-offset))) - from) / orientation
	} else {
		return (float32(math.Floor(float64(from+offset))) - from) / orientation
	}
}

func (l *level) castRay(position mgl32.Vec3, orientation mgl32.Vec3, length float32) (targetedChunk *Chunk, targeted int, frontChunk *Chunk, front int) {
	previousBlockPos := mgl32.Vec3{
		float32(math.Floor(float64(position.X()))),
		float32(math.Floor(float64(position.Y()))),
		float32(math.Floor(float64(position.Z()))),
	}
	blockPos := previousBlockPos
	steps := mgl32.Vec3{
		orientation.X() / float32(math.Abs(float64(orientation.X()))),
		orientation.Y() / float32(math.Abs(float64(orientation.Y()))),
		orientation.Z() / float32(math.Abs(float64(orientation.Z()))),
	}

	curr := position
	for curr.Sub(position).Len() < length {
		previousBlockPos = blockPos
		deltaIfX := t(curr.X(), 1, orientation.X())
		deltaIfY := t(curr.Y(), 1, orientation.Y())
		deltaIfZ := t(curr.Z(), 1, orientation.Z())

		if deltaIfX < deltaIfY {
			if deltaIfX < deltaIfZ {
				// move in x
				curr = curr.Add(orientation.Mul(deltaIfX))
				blockPos[0] += steps.X()
			} else {
				// move in z
				curr = curr.Add(orientation.Mul(deltaIfZ))
				blockPos[2] += steps.Z()
			}
		} else {
			if deltaIfY < deltaIfZ {
				// move in y
				curr = curr.Add(orientation.Mul(deltaIfY))
				blockPos[1] += steps.Y()
			} else {
				// move in z
				curr = curr.Add(orientation.Mul(deltaIfZ))
				blockPos[2] += steps.Z()
			}
		}

		// if targeted block exists and isn't air, return targeted block
		targetedChunk, targeted := l.positionInChunk(blockPos)
		if targetedChunk != nil && targeted >= 0 && targeted < CHUNK_WIDTH*CHUNK_HEIGHT*CHUNK_WIDTH &&
			targetedChunk.blocks[targeted].id != 0 {

			// if front block exists, return targeted block and front block
			frontChunk, front := l.positionInChunk(previousBlockPos)
			if frontChunk != nil && front >= 0 && front < CHUNK_WIDTH*CHUNK_HEIGHT*CHUNK_WIDTH {
				return targetedChunk, targeted, frontChunk, front
			}

			return targetedChunk, targeted, nil, -1
		}
	}
	return nil, -1, nil, -1
}

func (l *level) updateChunksAround(renderDistance *int) {
	for l.observer != nil {
		if len(l.chunks) != (*renderDistance)*2+1 {
			l.chunks = make([][]*Chunk, (*renderDistance)*2+1)
			for i := range l.chunks {
				l.chunks[i] = make([]*Chunk, (*renderDistance)*2+1)
			}
		}

		for x := -*renderDistance; x <= *renderDistance; x++ {
			for z := -*renderDistance; z <= *renderDistance; z++ {
				pos := l.observer.Add(mgl32.Vec2{float32(x), float32(z)})

				if c := l.get(pos); c == nil || pos != c.coordinates {
					l.set(pos, generateChunk(pos))
				} else {
					c.Loaded = true
				}
			}
		}
	}
}

func (c *Chunk) SolidMesh() []uint32 {
	return c.solidMesh
}

func (c *Chunk) TransparentMesh() []uint32 {
	return c.transparentMesh
}

/*
Returns true if block id "id2" is visible through block id "id1"
*/
func visible(id1, id2 uint8) bool {
	return id1 == 0 ||
		BLOCK_TYPES[id1].isTransparent && id1 != id2
}

func indexInChunk(x, y, z int) int {
	return x*CHUNK_WIDTH*CHUNK_HEIGHT + y*CHUNK_WIDTH + z
}

func (c *Chunk) generateMesh() {
	c.solidMesh = make([]uint32, 0)
	c.transparentMesh = make([]uint32, 0)
	for i, b := range c.blocks {
		if b.id == 0 {
			continue
		}
		x := (i / (CHUNK_WIDTH * CHUNK_HEIGHT)) & 0xf
		y := ((i % (CHUNK_WIDTH * CHUNK_HEIGHT)) / CHUNK_WIDTH) & 0xff
		z := (i % CHUNK_WIDTH) & 0xf

		render := [6]bool{
			// (middle AND visible) OR edge
			y+1 < CHUNK_HEIGHT && visible(c.blocks[indexInChunk(x, y+1, z)].id, b.id) || y+1 >= CHUNK_HEIGHT,
			y-1 >= 0 && visible(c.blocks[indexInChunk(x, y-1, z)].id, b.id) || y-1 < 0,

			// (middle AND (visible OR different height level)) OR edge
			x-1 >= 0 && (visible(c.blocks[indexInChunk(x-1, y, z)].id, b.id) || c.blocks[indexInChunk(x-1, y, z)].level != b.level) || x-1 < 0,
			x+1 < CHUNK_WIDTH && (visible(c.blocks[indexInChunk(x+1, y, z)].id, b.id) || c.blocks[indexInChunk(x+1, y, z)].level != b.level) || x+1 >= CHUNK_WIDTH,
			z+1 < CHUNK_WIDTH && (visible(c.blocks[indexInChunk(x, y, z+1)].id, b.id) || c.blocks[indexInChunk(x, y, z+1)].level != b.level) || z+1 >= CHUNK_WIDTH,
			z-1 >= 0 && (visible(c.blocks[indexInChunk(x, y, z-1)].id, b.id) || c.blocks[indexInChunk(x, y, z-1)].level != b.level) || z-1 < 0,
		}

		if BLOCK_TYPES[b.id].isTransparent {
			c.transparentMesh = append(c.transparentMesh, b.mesh(x, y, z, render)...)
		} else {
			c.solidMesh = append(c.solidMesh, b.mesh(x, y, z, render)...)
		}
	}
	c.Dirty = true
}

func (c *Chunk) set(index int, value block) {
	c.blocks[index] = value
	c.generateMesh()
}
