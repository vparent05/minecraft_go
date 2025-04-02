package player

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/movement"
)

type Direction = int

const (
	FRONT Direction = iota
	BACK  Direction = iota
	LEFT  Direction = iota
	RIGHT Direction = iota
	UP    Direction = iota
	DOWN  Direction = iota
)

type player struct {
	movementController movement.Controller
	cameraOffsets      []mgl32.Vec3
	selectedCamera     int
	yaw                float32
	pitch              float32
}

func (p *player) ProcessInputs(deltaTime float32) {
	if glfw.GetCurrentContext().GetKey(glfw.KeyEscape) == glfw.Press {
		glfw.GetCurrentContext().SetShouldClose(true)
	}

	directions := make([]Direction, 0, 6)
	if glfw.GetCurrentContext().GetKey(glfw.KeyW) == glfw.Press {
		directions = append(directions, FRONT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyA) == glfw.Press {
		directions = append(directions, LEFT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyS) == glfw.Press {
		directions = append(directions, BACK)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyD) == glfw.Press {
		directions = append(directions, RIGHT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == glfw.Press {
		directions = append(directions, UP)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftShift) == glfw.Press {
		directions = append(directions, DOWN)
	}

	p.UpdatePosition(directions, deltaTime)
}

func NewPlayer() player {
	return player{
		movement.NewController(mgl32.Vec3{0, 0, 0}, 15, 15, 5),
		[]mgl32.Vec3{{0, 0, 0}},
		0,
		0,
		0,
	}
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
	matrix = mgl32.HomogRotate3D(p.pitch, matrix.Mul3x1(mgl32.Vec3{1, 0, 0})).Mat3()

	return matrix.Mul3x1(orientation)
}

func (p *player) UpdatePosition(directions []Direction, deltaTime float32) {
	direction := mgl32.Vec3{0, 0, 0}

	for _, dir := range directions {
		switch dir {
		case FRONT:
			direction = direction.Add(mgl32.Vec3{0, 0, -1})
		case BACK:
			direction = direction.Add(mgl32.Vec3{0, 0, 1})
		case LEFT:
			direction = direction.Add(mgl32.Vec3{-1, 0, 0})
		case RIGHT:
			direction = direction.Add(mgl32.Vec3{1, 0, 0})
		case UP:
			direction = direction.Add(mgl32.Vec3{0, 1, 0})
		case DOWN:
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
}
