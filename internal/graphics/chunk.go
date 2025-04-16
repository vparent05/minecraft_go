package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	p_game "github.com/vparent05/minecraft_go/internal/game"
)

type chunkRenderer struct {
	game        *p_game.Game
	program     *program
	VAO         uint32
	VBOs        []uint32
	vertexCount []int32
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
		[]uint32{},
		[]int32{},
	}, nil
}

func (r *chunkRenderer) UpdateVBOs() error {
	gl.BindVertexArray(r.VAO)

	// create vertex buffer objects
	chunkCount := int32(len(r.game.Chunks))

	if len(r.VBOs) > 0 {
		gl.DeleteBuffers(int32(len(r.VBOs)), &r.VBOs[0])
	}

	if chunkCount == 0 {
		r.VBOs = []uint32{}
		r.vertexCount = []int32{}
		return nil
	}

	vertexCount := make([]int32, chunkCount*2)
	VBOs := make([]uint32, chunkCount*2)
	gl.GenBuffers(chunkCount*2, &VBOs[0])
	for i := range chunkCount {
		// solid geometry
		gl.BindBuffer(gl.ARRAY_BUFFER, VBOs[i])
		vertices := r.game.Chunks[i].SolidMesh()

		vertexCount[i] = int32(len(vertices))
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		// transparent geometry
		gl.BindBuffer(gl.ARRAY_BUFFER, VBOs[i+chunkCount])
		vertices = r.game.Chunks[i].TransparentMesh()

		vertexCount[i+chunkCount] = int32(len(vertices))
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	}

	r.VBOs = VBOs
	r.vertexCount = vertexCount

	return nil
}

func (r *chunkRenderer) UpdateVBO(index int) {
	gl.BindVertexArray(r.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.VBOs[index])
	vertices := r.game.Chunks[index].SolidMesh()

	r.vertexCount[index] = int32(len(vertices))
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
}

func (r *chunkRenderer) Draw() error {
	r.program.use()
	viewLocation, err := r.program.getUniformLocation("view")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(viewLocation, 1, false, &r.game.View[0])
	gl.BindVertexArray(r.VAO)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	for i, VBO := range r.VBOs {
		chunkCoordinatesLocation, err := r.program.getUniformLocation("chunkCoordinates")
		if err != nil {
			return fmt.Errorf("getUniformLocation(): %w", err)
		}
		chunkPosition := r.game.Chunks[i%len(r.game.Chunks)].Position()
		gl.Uniform2fv(chunkCoordinatesLocation, 1, &chunkPosition[0])

		gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
		gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
		gl.DrawArrays(gl.TRIANGLES, 0, r.vertexCount[i])
	}
	return nil
}
