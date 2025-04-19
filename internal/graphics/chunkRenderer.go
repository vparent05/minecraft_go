package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	p_game "github.com/vparent05/minecraft_go/internal/game"
)

type chunkRenderer struct {
	game    *p_game.Game
	program *program
	_VAO    uint32
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

func (r *chunkRenderer) deleteVBOs(chunk *p_game.Chunk) {
	if chunk.SolidVBO != 0 {
		gl.DeleteBuffers(1, &chunk.SolidVBO)
		chunk.SolidVBO = 0
	}
	if chunk.TransparentVBO != 0 {
		gl.DeleteBuffers(1, &chunk.TransparentVBO)
		chunk.TransparentVBO = 0
	}
}

func (r *chunkRenderer) updateVBOs(chunk *p_game.Chunk) {
	gl.BindVertexArray(r._VAO)
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

func (r *chunkRenderer) draw(vbo uint32, pos *float32, count int) error {
	if vbo == 0 {
		return nil
	}

	chunkCoordinatesLocation, err := r.program.getUniformLocation("chunkCoordinates")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}

	gl.Uniform2fv(chunkCoordinatesLocation, 1, pos)

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
	gl.UniformMatrix4fv(viewLocation, 1, false, &r.game.View[0])
	gl.BindVertexArray(r._VAO)

	// draw solid geometry
	for pos, chunk := range r.game.Level.Iterator() {
		if !chunk.Loaded {
			r.deleteVBOs(chunk)
			continue
		}
		if chunk.SolidVBO == 0 {
			r.updateVBOs(chunk)
		}

		err = r.draw(chunk.SolidVBO, &pos[0], len(chunk.SolidMesh()))
		if err != nil {
			return fmt.Errorf("draw(): %w", err)
		}
	}

	// draw transparent geometry
	for pos, chunk := range r.game.Level.Iterator() {
		if !chunk.Loaded {
			r.deleteVBOs(chunk)
			continue
		}
		if chunk.SolidVBO == 0 {
			r.updateVBOs(chunk)
		}

		err = r.draw(chunk.TransparentVBO, &pos[0], len(chunk.TransparentMesh()))
		if err != nil {
			return fmt.Errorf("draw(): %w", err)
		}
	}
	return nil
}
