package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Player     *player
	Level      *level
	Projection mgl32.Mat4
	View       mgl32.Mat4
}

func NewGame(projection mgl32.Mat4) *Game {
	game := Game{}
	game.Projection = projection

	game.Player = NewPlayer(&game)
	game.Level = NewLevel(game.Player.renderDistance, game.Player.chunkCoords)
	go game.Level.updateChunksAround(&game.Player.renderDistance)

	return &game
}
