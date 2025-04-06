package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Player     *player
	Chunks     []*Chunk
	Projection mgl32.Mat4
	View       mgl32.Mat4
}
