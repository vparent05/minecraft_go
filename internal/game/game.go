package game

import "github.com/vparent05/minecraft_go/internal/level"

type Game struct {
	Player *player
	Chunks []*level.Chunk
}
