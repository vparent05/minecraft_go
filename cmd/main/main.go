package main

import (
	"fmt"
	"math"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	p_game "github.com/vparent05/minecraft_go/internal/game"
	"github.com/vparent05/minecraft_go/internal/graphics"
	"github.com/vparent05/minecraft_go/internal/level"
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
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	err := glfw.Init()
	if err != nil {
		panic(fmt.Errorf("glfw.Init(): %w", err))
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(1440, 810, "Testing", nil, nil)
	if err != nil {
		panic(fmt.Errorf("glfw.CreateWindow(): %w", err))
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("gl.Init(): %w", err))
	}
	checkGLError()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version: %s\n", version)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	level.LoadBlocks()
	game := p_game.NewGame(mgl32.Perspective(math.Pi/4, 16.0/9.0, 0.1, 2048))

	chunkRenderer, err := graphics.NewChunkRenderer(game)
	if err != nil {
		panic(fmt.Errorf("graphics.NewChunkRenderer(): %w", err))
	}

	skyboxRenderer, err := graphics.NewSkyboxRenderer(game)
	if err != nil {
		panic(fmt.Errorf("graphics.NewSkyboxRenderer(): %w", err))
	}

	lastFrame := glfw.GetTime()
	var deltaTime float64
	var currentTime float64

	for !window.ShouldClose() {
		currentTime = glfw.GetTime()
		deltaTime = currentTime - lastFrame
		lastFrame = currentTime

		fmt.Printf("FPS: %.2f\n", 1.0/deltaTime)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		game.View = mgl32.LookAtV(game.Player.CameraPosition(), game.Player.CameraPosition().Add(game.Player.Orientation()), mgl32.Vec3{0, 1, 0})

		err = skyboxRenderer.Draw()
		if err != nil {
			panic(fmt.Errorf("Draw(): %w", err))
		}

		err = chunkRenderer.Draw()
		if err != nil {
			panic(fmt.Errorf("Draw(): %w", err))
		}

		window.SwapBuffers()
		glfw.PollEvents()
		game.FrameTick(float32(deltaTime))
		checkGLError()
	}
}
