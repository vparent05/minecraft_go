package game

type direction = int

const (
	_FRONT direction = iota // -Z
	_BACK  direction = iota // +Z
	_LEFT  direction = iota // -X
	_RIGHT direction = iota // +X
	_UP    direction = iota // +Y
	_DOWN  direction = iota // -Y
)
