package movement

type Direction = int

const (
	FRONT Direction = iota // -Z
	BACK                   // +Z
	LEFT                   // -X
	RIGHT                  // +X
	UP                     // +Y
	DOWN                   // -Y
)
