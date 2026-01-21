package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	p_game "github.com/vparent05/minecraft_go/internal/game"
	"github.com/vparent05/minecraft_go/internal/level"
	"github.com/vparent05/minecraft_go/internal/utils"
	"github.com/vparent05/minecraft_go/internal/utils/chanx"
)

type chunkData struct {
	solidCount       int
	transparentCount int
	solidVBO         uint32
	transparentVBO   uint32
}

type chunkRenderer struct {
	game       *p_game.Game
	program    *program
	_VAO       uint32
	chunksData []chunkData
}

func NewChunkRenderer(game *p_game.Game) (*chunkRenderer, error) {
	if game == nil {
		return nil, errors.New("game pointer is nil")
	}

	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("gl.Init(): %w", err)
	}

	level.BLOCK_TEXTURE_ATLAS, err = loadTextureAtlas("./textures/blocks", _BLOCKS_TEXTURE, 16)
	if err != nil {
		return nil, fmt.Errorf("loadTextureAtlas(): %w", err)
	}

	// create the block shader program
	blockProgram, err := NewProgram(
		NewShader("./shaders/block/Vertex.glsl", gl.VERTEX_SHADER),
		NewShader("./shaders/block/Fragment.glsl", gl.FRAGMENT_SHADER),
	)
	if err != nil {
		return nil, fmt.Errorf("NewProgram(): %w", err)
	}

	// create vertex array object
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	gl.EnableVertexAttribArray(0)

	blockProgram.use()
	projectionLocation, err := blockProgram.getUniformLocation("projection")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(projectionLocation, 1, false, &game.Projection[0]) // TODO separate game from graphic variables

	textureLocation, err := blockProgram.getUniformLocation("atlas")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.Uniform1i(textureLocation, _BLOCKS_TEXTURE)

	return &chunkRenderer{
		game,
		blockProgram,
		VAO,
		make([]chunkData, 33*33), // TODO actually link the render distance to the size of the slice
	}, nil
}

func (r *chunkRenderer) applyMeshUpdate(chunk *level.Chunk) {
	_, ok := chanx.TryReceive(chunk.MeshUpdates)
	if ok {
		newMesh := chunk.Mesh.Load()
		r.chunksData[chunk.Slot].solidCount = len(newMesh.Solid)
		r.chunksData[chunk.Slot].transparentCount = len(newMesh.Transparent)
		r.updateVBOs(chunk, newMesh)
	}
}

func (r *chunkRenderer) updateVBOs(chunk *level.Chunk, newMesh level.ChunkMesh) {
	gl.BindVertexArray(r._VAO)

	chunkData := r.chunksData[chunk.Slot]

	if chunkData.solidVBO == 0 {
		gl.GenBuffers(1, &chunkData.solidVBO)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, chunkData.solidVBO)
	vertices := newMesh.Solid
	if len(vertices) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}

	if chunkData.transparentVBO == 0 {
		gl.GenBuffers(1, &chunkData.transparentVBO)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, chunkData.transparentVBO)
	vertices = newMesh.Transparent
	if len(vertices) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}

	r.chunksData[chunk.Slot] = chunkData
}

func (r *chunkRenderer) deleteVBOs(chunk *level.Chunk) {
	chunkData := r.chunksData[chunk.Slot]
	if chunkData.solidVBO != 0 {
		gl.DeleteBuffers(1, &chunkData.solidVBO)
		chunkData.solidVBO = 0
	}
	if chunkData.transparentVBO != 0 {
		gl.DeleteBuffers(1, &chunkData.transparentVBO)
		chunkData.transparentVBO = 0
	}

	r.chunksData[chunk.Slot] = chunkData
}

func (r *chunkRenderer) draw(vbo uint32, pos utils.IntVector2, count int) error {
	if vbo == 0 {
		return nil
	}

	chunkCoordinatesLocation, err := r.program.getUniformLocation("chunkCoordinates")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}

	gl.Uniform2fv(chunkCoordinatesLocation, 1, &intVector2ToFloat32Slice(pos)[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(count))

	return nil
}

func (r *chunkRenderer) Draw() error {
	r.program.use()
	viewLocation, err := r.program.getUniformLocation("view")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(viewLocation, 1, false, &r.game.View[0]) // TODO separate game from graphic variables
	gl.BindVertexArray(r._VAO)

	// draw solid geometry
	for pos, chunk := range r.game.Level.Chunks() {
		r.applyMeshUpdate(chunk)
		if r.chunksData[chunk.Slot].solidVBO == 0 {
			r.updateVBOs(chunk, level.ChunkMesh{})
		}

		err = r.draw(r.chunksData[chunk.Slot].solidVBO, pos, r.chunksData[chunk.Slot].solidCount)
		if err != nil {
			return fmt.Errorf("draw(): %w", err)
		}
	}

	// draw transparent geometry
	for pos, chunk := range r.game.Level.Chunks() {
		err = r.draw(r.chunksData[chunk.Slot].transparentVBO, pos, r.chunksData[chunk.Slot].transparentCount)
		if err != nil {
			return fmt.Errorf("draw(): %w", err)
		}
	}
	return nil
}

func intVector2ToFloat32Slice(v utils.IntVector2) []float32 {
	return []float32{float32(v.X), float32(v.Y)}
}
