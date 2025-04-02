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
)

const vertex = `#version 330 core
layout (location = 0) in int ver;

uniform mat4 view;
uniform mat4 projection;

void main()
{	
		int x = (ver>>28) & 0xF;
		int y = (ver>>20) & 0xFF;
		int z = (ver>>16) & 0xF;
    vec4 homogeneous = projection * view * vec4(float(x), float(y), float(-z), 1.0);
    gl_Position = homogeneous / homogeneous.w;
}` + "\x00"

const fragment = `#version 330 core
out vec4 FragColor;

void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 0.5f);
}` + "\x00"

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

	window, err := glfw.CreateWindow(500, 500, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertex)
	gl.ShaderSource(vShader, 1, csources, nil)
	free()
	gl.CompileShader(vShader)
	gl.AttachShader(program, vShader)
	gl.DeleteShader(vShader)

	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragment)
	gl.ShaderSource(fShader, 1, csources, nil)
	free()
	gl.CompileShader(fShader)
	gl.AttachShader(program, fShader)
	gl.DeleteShader(fShader)

	gl.LinkProgram(program)

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.VertexAttribIPointer(0, 1, gl.INT, 4, nil)
	gl.EnableVertexAttribArray(0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.UseProgram(program)

	proj := mgl32.Perspective(math.Pi/4, 1, 0.1, 100)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &proj[0])

	player := player.NewPlayer()
	viewLocation := gl.GetUniformLocation(program, gl.Str("view\x00"))

	vertices := level.GetTestMesh()

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var view mgl32.Mat4
	lastFrame := float32(glfw.GetTime())
	var deltaTime float32
	var currentTime float32
	for !window.ShouldClose() {
		currentTime = float32(glfw.GetTime())
		deltaTime = currentTime - lastFrame
		lastFrame = currentTime

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)
		view = mgl32.LookAtV(player.CameraPosition(), player.CameraPosition().Add(player.Orientation()), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])

		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))

		window.SwapBuffers()
		glfw.PollEvents()
		player.ProcessInputs(deltaTime)
		checkGLError()
	}
}
