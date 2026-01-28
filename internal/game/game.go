package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/level"
)

type Game struct {
	Player     *player
	Level      *level.Level
	Projection mgl32.Mat4
	View       mgl32.Mat4
}

func NewGame(projection mgl32.Mat4) *Game {
	g := Game{}
	g.Projection = projection

	g.Player = NewPlayer(&g)
	g.Level = level.NewLevel(g.Player.levelObserver)

	return &g
}

func (g *Game) Start() {
	go g.Level.GenerateAround() // TODO manage this goroutine (don't leave it hanging)
}

func (g *Game) FrameTick(deltaTime float32) {
	g.Player.FrameTick(deltaTime)
}
