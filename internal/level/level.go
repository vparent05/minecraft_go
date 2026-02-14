package level

import (
	"iter"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/utils"
	"github.com/vparent05/minecraft_go/internal/utils/atomicx"
	"github.com/vparent05/minecraft_go/internal/utils/chanx"
)

type blockPosition struct {
	c *Chunk
	i utils.IntVector3
}

func (b *blockPosition) Set(value BlockId) {
	if b.c == nil {
		return
	}
	b.c.setBlock(b.i, value)
}

func (b *blockPosition) Get() (BlockId, bool) {
	if b.c == nil {
		return 0, false
	}
	return b.c.getBlock(b.i), true
}

type LevelObserver struct {
	mgl32.Vec3
	RenderDistance int
}

type Level struct {
	observer      *atomicx.Value[LevelObserver]
	observerCache LevelObserver
	chunks        [][]*Chunk
	generateOrder [][2]int
}

func NewLevel(observer *atomicx.Value[LevelObserver]) *Level {
	return &Level{
		observer:      observer,
		observerCache: observer.Load(),
	}
}

// originIndex returns the current index of the 0, 0 chunk
func (l *Level) originIndex() (int, int) {
	observerChunkCoords := LevelToChunkCoords(l.observerCache.Vec3)
	return utils.Mod(observerChunkCoords.X, len(l.chunks[0])),
		utils.Mod(observerChunkCoords.Y, len(l.chunks))
}

// chunkIndex returns the current index of the chunk described by chunkCoordinates
func (l *Level) chunkIndex(chunkCoordinates utils.IntVector2) (int, int) {
	delta := chunkCoordinates.Sub(LevelToChunkCoords(l.observerCache.Vec3))
	originX, originZ := l.originIndex()
	return utils.Mod(delta.X+originX, len(l.chunks[0])),
		utils.Mod(delta.Y+originZ, len(l.chunks))
}

func (l *Level) getChunk(chunkCoordinates utils.IntVector2) *Chunk {
	i, j := l.chunkIndex(chunkCoordinates)
	return l.chunks[i][j]
}

func (l *Level) setChunk(chunkCoordinates utils.IntVector2, value *Chunk) {
	i, j := l.chunkIndex(chunkCoordinates)
	value.Slot = i*len(l.chunks) + j
	l.chunks[i][j] = value
}

func (l *Level) Chunks() iter.Seq2[utils.IntVector2, *Chunk] {
	return func(yield func(utils.IntVector2, *Chunk) bool) {
		if len(l.chunks) == 0 || len(l.chunks[0]) == 0 {
			return
		}

		originX, originZ := l.originIndex()
		for _, chunk := range utils.FromOriginIterator2(l.chunks, utils.IntVector2{X: originX, Y: originZ}) {
			if chunk == nil {
				continue
			}
			if !yield(chunk.coordinates, chunk) {
				return
			}
		}
	}
}

func (l *Level) getBlockPosition(position mgl32.Vec3) *blockPosition {
	blockX := utils.Mod(int(position.X()), CHUNK_WIDTH)
	blockZ := utils.Mod(int(position.Z()), CHUNK_WIDTH)

	chunk := l.getChunk(LevelToChunkCoords(position))
	if chunk == nil {
		return &blockPosition{nil, utils.IntVector3{}}
	}

	return &blockPosition{chunk, utils.IntVector3{X: blockX, Y: utils.Mod(int(position.Y()), CHUNK_HEIGHT), Z: blockZ}}
}

func t(from, offset, orientation float32) float32 {
	if orientation < 0 {
		return (float32(math.Ceil(float64(from-offset))) - from) / orientation
	} else {
		return (float32(math.Floor(float64(from+offset))) - from) / orientation
	}
}

func (l *Level) CastRay(position mgl32.Vec3, orientation mgl32.Vec3, length float32) (targeted *blockPosition, front *blockPosition) {
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
		targeted := l.getBlockPosition(blockPos)
		if tBlock, ok := targeted.Get(); ok && tBlock != AIR {

			// if front block exists, return targeted block and front block
			front := l.getBlockPosition(previousBlockPos)
			if _, ok = front.Get(); ok {
				return targeted, front
			}

			return targeted, nil
		}
	}
	return nil, nil
}

func (l *Level) updateGenerateOrder() {
	size := l.observerCache.RenderDistance*2 + 1
	l.generateOrder = make([][2]int, 0, size*size)
	for L := range l.observerCache.RenderDistance + 1 {
		for i := range size {
			for j := range size {
				if int(math.Max(math.Abs(float64(i-l.observerCache.RenderDistance)), math.Abs(float64(j-l.observerCache.RenderDistance)))) == L {
					l.generateOrder = append(l.generateOrder, [2]int{i, j})
				}
			}
		}
	}
}

func (l *Level) GenerateAround() {
	l.observerCache = l.observer.Load()
	generator := newMeshGenerator(l, l.observerCache.RenderDistance) // TODO update render distance dynamically
	generator.start(WORKER_COUNT)

	for {
		if len(l.chunks) != l.observerCache.RenderDistance*2+1 {
			l.updateGenerateOrder()
			l.chunks = make([][]*Chunk, l.observerCache.RenderDistance*2+1)
			for i := range l.chunks {
				l.chunks[i] = make([]*Chunk, l.observerCache.RenderDistance*2+1)
			}
		}

		observerChunkCoords := LevelToChunkCoords(l.observerCache.Vec3)
		for _, i := range l.generateOrder {
			l.observerCache = l.observer.Load()
			newObserverChunkCoords := LevelToChunkCoords(l.observerCache.Vec3)
			if newObserverChunkCoords != observerChunkCoords {
				// Observer changed chunks, cut our losses to regenerate around the new center
				break
			}

			xOffset := i[0] - l.observerCache.RenderDistance
			zOffset := i[1] - l.observerCache.RenderDistance
			pos := utils.IntVector2{X: xOffset + observerChunkCoords.X, Y: zOffset + observerChunkCoords.Y}
			if c := l.getChunk(pos); c == nil || pos != c.coordinates {
				if c == nil {
					c = newChunk(generator, l.observer)
				}

				c.Mesh.Store(ChunkMesh{make([]uint32, 0), make([]uint32, 0)})
				chanx.TrySend(c.MeshUpdates, struct{}{})
				generateChunk(c, pos)
				l.setChunk(pos, c)
			}
		}
	}
}

func LevelToChunkCoords(level mgl32.Vec3) utils.IntVector2 {
	return utils.IntVector2{
		X: int(math.Floor(float64(level.X() / CHUNK_WIDTH))),
		Y: int(math.Floor(float64(level.Z() / CHUNK_WIDTH))),
	}
}
