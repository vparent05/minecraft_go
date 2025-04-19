package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	p_game "github.com/vparent05/minecraft_go/internal/game"
)

type chunkRenderer struct {
	game    *p_game.Game
	program *program
	VAO     uint32
}

func NewChunkRenderer(game *p_game.Game) (*chunkRenderer, error) {
	if game == nil {
		return nil, errors.New("game pointer is nil")
	}

	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("gl.Init(): %w", err)
	}

	p_game.BLOCK_TEXTURE_ATLAS, err = loadTextureAtlas("./textures/blocks", BLOCKS_TEXTURE, 16)
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
	gl.UniformMatrix4fv(projectionLocation, 1, false, &game.Projection[0])

	textureLocation, err := blockProgram.getUniformLocation("atlas")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.Uniform1i(textureLocation, BLOCKS_TEXTURE)

	return &chunkRenderer{
		game,
		blockProgram,
		VAO,
	}, nil
}

func (r *chunkRenderer) DeleteVBO(buf uint32) {
	gl.DeleteBuffers(1, &buf)
}

func (r *chunkRenderer) UpdateVBO(pos mgl32.Vec2) {
	chunk, ok := r.game.Level.Chunks.Get(pos)
	if !ok {
		return
	}

	gl.BindVertexArray(r.VAO)
	if chunk.SolidVBO == 0 {
		gl.GenBuffers(1, &chunk.SolidVBO)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, chunk.SolidVBO)
	vertices := chunk.SolidMesh()
	if len(vertices) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}

	if chunk.TransparentVBO == 0 {
		gl.GenBuffers(1, &chunk.TransparentVBO)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, chunk.TransparentVBO)
	vertices = chunk.TransparentMesh()
	if len(vertices) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}
}

func (r *chunkRenderer) Draw() error {
	r.program.use()
	viewLocation, err := r.program.getUniformLocation("view")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(viewLocation, 1, false, &r.game.View[0])
	gl.BindVertexArray(r.VAO)

	for pos, chunk := range r.game.Level.Chunks.Iterator() {
		if chunk.SolidVBO == 0 {
			continue
		}

		chunkCoordinatesLocation, err := r.program.getUniformLocation("chunkCoordinates")
		if err != nil {
			return fmt.Errorf("getUniformLocation(): %w", err)
		}

		gl.Uniform2fv(chunkCoordinatesLocation, 1, &pos[0])

		gl.BindBuffer(gl.ARRAY_BUFFER, chunk.SolidVBO)
		gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(chunk.SolidMesh())))
	}

	for pos, chunk := range r.game.Level.Chunks.Iterator() {
		if chunk.SolidVBO == 0 {
			continue
		}

		chunkCoordinatesLocation, err := r.program.getUniformLocation("chunkCoordinates")
		if err != nil {
			return fmt.Errorf("getUniformLocation(): %w", err)
		}

		gl.Uniform2fv(chunkCoordinatesLocation, 1, &pos[0])

		gl.BindBuffer(gl.ARRAY_BUFFER, chunk.TransparentVBO)
		gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(chunk.TransparentMesh())))
	}
	return nil
}
