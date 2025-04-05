package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/level"
	"github.com/vparent05/minecraft_go/internal/player"
	"github.com/vparent05/minecraft_go/internal/shader"
)

func checkGLError() {
	for {
		err := gl.GetError()
		if err == gl.NO_ERROR {
			break
		}
		fmt.Println("OpenGL Error:", err)
	}
}

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(960, 540, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version:", version)

	program, err := shader.Program(
		shader.NewShader("./shaders/Vertex.glsl", gl.VERTEX_SHADER),
		shader.NewShader("./shaders/Fragment.glsl", gl.FRAGMENT_SHADER),
	)

	if err != nil {
		panic(err)
	}

	gl.LinkProgram(program)

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
	gl.EnableVertexAttribArray(0)

	gl.Enable(gl.DEPTH_TEST)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Enable(gl.CULL_FACE)

	gl.UseProgram(program)

	proj := mgl32.Perspective(math.Pi/4, 16.0/9.0, 0.1, 100)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &proj[0])

	player := player.NewPlayer()
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))

	vertices := level.GetTestMesh()

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var view mgl32.Mat4
	lastFrame := glfw.GetTime()
	var deltaTime float64
	var currentTime float64
	for !window.ShouldClose() {
		currentTime = glfw.GetTime()
		deltaTime = currentTime - lastFrame
		lastFrame = currentTime

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		view = mgl32.LookAtV(player.CameraPosition(), player.CameraPosition().Add(player.Orientation()), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

		window.SwapBuffers()
		glfw.PollEvents()
		player.ProcessInputs(float32(deltaTime))
		checkGLError()
	}
}
