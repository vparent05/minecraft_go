package graphics

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/vparent05/minecraft_go/internal/game"
)

type chunkRenderer struct {
	game        *game.Game
	program     *program
	VAO         uint32
	VBOs        []uint32
	vertexCount []int32
}

func NewChunkRenderer(game *game.Game) (*chunkRenderer, error) {
	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("gl.Init(): %w", err)
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

	// create vertex buffer objects
	chunkCount := int32(len(game.Chunks))
	vertexCount := make([]int32, chunkCount)
	VBOs := make([]uint32, chunkCount)
	gl.GenBuffers(chunkCount, &VBOs[0])
	for i := range chunkCount {
		gl.BindBuffer(gl.ARRAY_BUFFER, VBOs[i])
		vertices := game.Chunks[i].SolidMesh()

		vertexCount[i] = int32(len(vertices))
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}

	gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
	gl.EnableVertexAttribArray(0)

	blockProgram.use()
	projectionLocation, err := blockProgram.getUniformLocation("projection")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(projectionLocation, 1, false, &game.Projection[0])

	return &chunkRenderer{
		game,
		blockProgram,
		VAO,
		VBOs,
		vertexCount,
	}, nil
}

func (r *chunkRenderer) Draw() error {
	r.program.use()
	viewLocation, err := r.program.getUniformLocation("view")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(viewLocation, 1, false, &r.game.View[0])
	gl.BindVertexArray(r.VAO)
	for i, VBO := range r.VBOs {
		gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
		gl.DrawArrays(gl.TRIANGLES, 0, r.vertexCount[i])
	}
	return nil
}

func (r *chunkRenderer) UpdateVBO(index int) {
	gl.BindBuffer(gl.ARRAY_BUFFER, r.VBOs[index])
	vertices := r.game.Chunks[index].SolidMesh()

	r.vertexCount[index] = int32(len(vertices))
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
}
