package game

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/movement"
)

type Direction = int

const (
	_FRONT Direction = iota
	_BACK  Direction = iota
	_LEFT  Direction = iota
	_RIGHT Direction = iota
	_UP    Direction = iota
	_DOWN  Direction = iota
)

type player struct {
	game               *Game
	movementController movement.Controller
	chunkCoords        *mgl32.Vec2
	cameraOffsets      []mgl32.Vec3
	selectedCamera     int
	yaw                float32
	pitch              float32
	maxPitch           float32
	renderDistance     int
}

func NewPlayer(game *Game) *player {
	p := player{
		game,
		movement.NewController(mgl32.Vec3{0, 65, 0}, 100, 100, 500),
		&mgl32.Vec2{0, 0},
		[]mgl32.Vec3{{0, 0, 0}},
		0,
		0,
		0,
		4 * math.Pi / 9,
		16,
	}
	go game.Level.updateChunksAround(p.chunkCoords, &p.renderDistance)
	return &p
}

func (p *player) Position() mgl32.Vec3 {
	return p.movementController.Position()
}

func (p *player) CameraPosition() mgl32.Vec3 {
	return p.movementController.Position().Add(p.cameraOffsets[p.selectedCamera])
}

func (p *player) Orientation() mgl32.Vec3 {
	orientation := mgl32.Vec3{0, 0, -1}

	matrix := mgl32.Rotate3DY(p.yaw)
	matrix = matrix.Mul3(mgl32.Rotate3DX(p.pitch))

	return matrix.Mul3x1(orientation)
}

func (p *player) updatePosition(directions []Direction, deltaTime float32) {
	direction := mgl32.Vec3{0, 0, 0}

	for _, dir := range directions {
		switch dir {
		case _FRONT:
			direction = direction.Add(mgl32.Vec3{0, 0, -1})
		case _BACK:
			direction = direction.Add(mgl32.Vec3{0, 0, 1})
		case _LEFT:
			direction = direction.Add(mgl32.Vec3{-1, 0, 0})
		case _RIGHT:
			direction = direction.Add(mgl32.Vec3{1, 0, 0})
		case _UP:
			direction = direction.Add(mgl32.Vec3{0, 1, 0})
		case _DOWN:
			direction = direction.Add(mgl32.Vec3{0, -1, 0})
		}
	}

	direction = mgl32.Vec3{
		direction.Z()*float32(math.Sin(float64(p.yaw))) + direction.X()*float32(math.Cos(float64(p.yaw))),
		direction.Y(),
		direction.Z()*float32(math.Cos(float64(p.yaw))) - direction.X()*float32(math.Sin(float64(p.yaw))),
	}

	if (direction != mgl32.Vec3{0, 0, 0}) {
		direction = direction.Normalize()
	}

	p.movementController.Move(deltaTime, direction)

	*p.chunkCoords = mgl32.Vec2{
		float32(math.Floor(float64(p.movementController.Position().X() / 15))),
		float32(math.Floor(float64(p.movementController.Position().Z() / 15))),
	}
}

func (p *player) updateOrientation(mouseDisplacementRad mgl32.Vec2) {
	p.yaw = float32(math.Mod(float64(p.yaw+mouseDisplacementRad.X()), math.Pi*2))
	p.pitch = max(min(p.pitch+mouseDisplacementRad.Y(), p.maxPitch), -p.maxPitch)
}

func (p *player) ProcessInputs(deltaTime float32) {
	if glfw.GetCurrentContext().GetKey(glfw.KeyEscape) == glfw.Press {
		glfw.GetCurrentContext().SetShouldClose(true)
	}

	directions := make([]Direction, 0, 6)
	if glfw.GetCurrentContext().GetKey(glfw.KeyW) == glfw.Press {
		directions = append(directions, _FRONT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyA) == glfw.Press {
		directions = append(directions, _LEFT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyS) == glfw.Press {
		directions = append(directions, _BACK)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyD) == glfw.Press {
		directions = append(directions, _RIGHT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == glfw.Press {
		directions = append(directions, _UP)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftShift) == glfw.Press {
		directions = append(directions, _DOWN)
	}

	width, height := glfw.GetCurrentContext().GetSize()
	xpos, ypos := glfw.GetCurrentContext().GetCursorPos()
	glfw.GetCurrentContext().SetCursorPos(float64(width)/2, float64(height)/2)

	p.updateOrientation(mgl32.Vec2{(float32(width/2) - float32(xpos)) / 1000, (float32(height)/2 - float32(ypos)) / 1000})
	p.updatePosition(directions, deltaTime)
}
