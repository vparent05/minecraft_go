package movement

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Controller struct {
	position     mgl32.Vec3
	direction    mgl32.Vec3
	velocity     float32
	acceleration float32
	drag         float32
	maxVelocity  float32
}

func NewController(position mgl32.Vec3, acceleration, drag, maxVelocity float32) *Controller {
	return &Controller{
		position,
		mgl32.Vec3{0, 0, 0},
		0,
		acceleration,
		drag,
		maxVelocity,
	}
}

func (c *Controller) Position() mgl32.Vec3 {
	return c.position
}

func (c *Controller) Move(deltaTime float32, accelerationDirection mgl32.Vec3) bool {
	dragVec := c.direction.Mul(c.drag)
	dragVec = dragVec.Sub(accelerationDirection.Mul(dragVec.Dot(accelerationDirection)))
	accelerationVec := accelerationDirection.Mul(c.acceleration).Sub(dragVec)
	velocityVec := c.direction.Mul(c.velocity)

	displacement := velocityVec.Mul(deltaTime).Add(accelerationVec.Mul(0.5 * deltaTime * deltaTime))
	if maxDisplacement := c.direction.Mul(c.maxVelocity * deltaTime); displacement.Len() > maxDisplacement.Len() {
		displacement = maxDisplacement
	}
	if velocityVec.X() < 5 && velocityVec.Y() < 5 && velocityVec.Z() < 5 && velocityVec.X() > -5 && velocityVec.Y() > -5 && velocityVec.Z() > -5 {
		displacement = mgl32.Vec3{0, 0, 0}
	}
	c.position = c.position.Add(displacement)

	newVelocityVec := accelerationVec.Mul(deltaTime).Add(velocityVec)
	c.velocity = max(min(newVelocityVec.Len(), c.maxVelocity), 0)
	if newDirection := newVelocityVec; (newDirection != mgl32.Vec3{0, 0, 0}) {
		c.direction = newDirection.Normalize()
	} else {
		c.direction = mgl32.Vec3{0, 0, 0}
	}

	return displacement.X() != 0 || displacement.Y() != 0 || displacement.Z() != 0
}

type EntityController struct {
	*Controller
	yaw      float32
	pitch    float32
	maxPitch float32
}

func NewEntityController(position mgl32.Vec3, acceleration, drag, maxVelocity, yaw, pitch, maxPitch float32) *EntityController {
	return &EntityController{
		NewController(position, acceleration, drag, maxVelocity),
		yaw,
		pitch,
		maxPitch,
	}
}

func (e *EntityController) Orientation() mgl32.Vec3 {
	orientation := mgl32.Vec3{0, 0, -1}

	matrix := mgl32.Rotate3DY(e.yaw)
	matrix = matrix.Mul3(mgl32.Rotate3DX(e.pitch))

	return matrix.Mul3x1(orientation)
}

func (e *EntityController) UpdatePosition(directions []Direction, deltaTime float32) bool {
	directionVec := mgl32.Vec3{0, 0, 0}

	for _, dir := range directions {
		switch dir {
		case FRONT:
			directionVec = directionVec.Add(mgl32.Vec3{0, 0, -1})
		case BACK:
			directionVec = directionVec.Add(mgl32.Vec3{0, 0, 1})
		case LEFT:
			directionVec = directionVec.Add(mgl32.Vec3{-1, 0, 0})
		case RIGHT:
			directionVec = directionVec.Add(mgl32.Vec3{1, 0, 0})
		case UP:
			directionVec = directionVec.Add(mgl32.Vec3{0, 1, 0})
		case DOWN:
			directionVec = directionVec.Add(mgl32.Vec3{0, -1, 0})
		}
	}

	directionVec = mgl32.Vec3{
		directionVec.Z()*float32(math.Sin(float64(e.yaw))) + directionVec.X()*float32(math.Cos(float64(e.yaw))),
		directionVec.Y(),
		directionVec.Z()*float32(math.Cos(float64(e.yaw))) - directionVec.X()*float32(math.Sin(float64(e.yaw))),
	}

	if (directionVec != mgl32.Vec3{0, 0, 0}) {
		directionVec = directionVec.Normalize()
	}

	return e.Move(deltaTime, directionVec)
}

func (e *EntityController) UpdateOrientation(mouseDisplacementRad mgl32.Vec2) bool {
	e.yaw = float32(math.Mod(float64(e.yaw+mouseDisplacementRad.X()), math.Pi*2))
	e.pitch = max(min(e.pitch+mouseDisplacementRad.Y(), e.maxPitch), -e.maxPitch)
	return mouseDisplacementRad.X() != 0 || mouseDisplacementRad.Y() != 0
}
