package game

import (
	"time"

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
	g.Level = level.NewLevel(g.Player.updates)

	t := time.NewTicker(50 * time.Millisecond)
	go func() {
		for range t.C {
			g.Tick()
		}
	}()

	return &g
}

func (g *Game) Tick() {
	g.Level.GameTick()
}
