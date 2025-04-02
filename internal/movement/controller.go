package movement

import (
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

func NewController(position mgl32.Vec3, acceleration, drag, maxVelocity float32) Controller {
	return Controller{
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

func (c *Controller) Move(deltaTime float32, accelerationDirection mgl32.Vec3) {
	dragVec := c.direction.Mul(c.drag)
	dragVec = dragVec.Sub(accelerationDirection.Mul(dragVec.Dot(accelerationDirection)))
	accelerationVec := accelerationDirection.Mul(c.acceleration).Sub(dragVec)
	velocityVec := c.direction.Mul(c.velocity)

	displacement := velocityVec.Mul(deltaTime).Add(accelerationVec.Mul(0.5 * deltaTime * deltaTime))
	if maxDisplacement := c.direction.Mul(c.maxVelocity * deltaTime); displacement.Len() > maxDisplacement.Len() {
		displacement = maxDisplacement
	}
	c.position = c.position.Add(displacement)

	newVelocityVec := accelerationVec.Mul(deltaTime).Add(velocityVec)
	c.velocity = min(newVelocityVec.Len(), c.maxVelocity)
	if newDirection := newVelocityVec; (newDirection != mgl32.Vec3{0, 0, 0}) {
		c.direction = newDirection.Normalize()
	} else {
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}
