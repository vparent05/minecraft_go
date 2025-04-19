package graphics

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/vparent05/minecraft_go/internal/game"
)

var skyboxVertices = []float32{
	-1, 1, -1,
	-1, -1, -1,
	1, -1, -1,
	1, -1, -1,
	1, 1, -1,
	-1, 1, -1,

	-1, -1, 1,
	-1, -1, -1,
	-1, 1, -1,
	-1, 1, -1,
	-1, 1, 1,
	-1, -1, 1,

	1, -1, -1,
	1, -1, 1,
	1, 1, 1,
	1, 1, 1,
	1, 1, -1,
	1, -1, -1,

	-1, -1, 1,
	-1, 1, 1,
	1, 1, 1,
	1, 1, 1,
	1, -1, 1,
	-1, -1, 1,

	-1, 1, -1,
	1, 1, -1,
	1, 1, 1,
	1, 1, 1,
	-1, 1, 1,
	-1, 1, -1,

	-1, -1, -1,
	-1, -1, 1,
	1, -1, -1,
	1, -1, -1,
	-1, -1, 1,
	1, -1, 1,
}

type skyboxRenderer struct {
	game    *game.Game
	program *program
	VAO     uint32
}

func NewSkyboxRenderer(game *game.Game) (*skyboxRenderer, error) {
	if game == nil {
		return nil, errors.New("game pointer is nil")
	}

	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("gl.Init(): %w", err)
	}

	err = loadCubemap("./textures/skybox", SKYBOX_TEXTURE)
	if err != nil {
		return nil, fmt.Errorf("loadCubemap(): %w", err)
	}

	// create the skybox shader program
	skyboxProgram, err := NewProgram(
		NewShader("./shaders/skybox/Vertex.glsl", gl.VERTEX_SHADER),
		NewShader("./shaders/skybox/Fragment.glsl", gl.FRAGMENT_SHADER),
	)
	if err != nil {
		return nil, fmt.Errorf("NewProgram(): %w", err)
	}

	// create vertex array object
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	// create vertex buffer objects
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.BufferData(gl.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(skyboxVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	skyboxProgram.use()
	projectionLocation, err := skyboxProgram.getUniformLocation("projection")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.UniformMatrix4fv(projectionLocation, 1, false, &game.Projection[0])

	skyboxLocation, err := skyboxProgram.getUniformLocation("skybox")
	if err != nil {
		return nil, fmt.Errorf("getUniformLocation(): %w", err)
	}
	gl.Uniform1i(skyboxLocation, SKYBOX_TEXTURE)

	return &skyboxRenderer{
		game,
		skyboxProgram,
		VAO,
	}, nil
}

func (r *skyboxRenderer) Draw() error {
	r.program.use()
	viewLocation, err := r.program.getUniformLocation("view")
	if err != nil {
		return fmt.Errorf("getUniformLocation(): %w", err)
	}
	rotationOnlyView := r.game.View.Mat3().Mat4()
	gl.UniformMatrix4fv(viewLocation, 1, false, &rotationOnlyView[0])

	gl.DepthMask(false)
	gl.DepthFunc(gl.LEQUAL)
	gl.BindVertexArray(r.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LESS)

	return nil
}
