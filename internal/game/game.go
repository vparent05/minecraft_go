package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/level"
)

type Game struct {
	Player     *player
	Chunks     []*level.Chunk
	Projection mgl32.Mat4
}
