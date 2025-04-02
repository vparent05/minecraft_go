package movement

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Controller struct {
	position     mgl32.Vec3
	speed        float32
	maxSpeed     float32
	acceleration float32
}

func NewController(position mgl32.Vec3, maxSpeed, acceleration float32) Controller {
	return Controller{
		position,
		0,
		maxSpeed,
		acceleration,
	}
}

func (c *Controller) Position() mgl32.Vec3 {
	return c.position
}

func (c *Controller) Move(deltaTime float32, direction mgl32.Vec3) {
	speedVec := direction.Mul(c.speed)
	accelerationVec := direction.Mul(c.acceleration)
	displacement := speedVec.Mul(deltaTime).Add(accelerationVec.Mul(deltaTime * deltaTime))
	if maxDisplacement := direction.Mul(c.maxSpeed * deltaTime); displacement.Len() > maxDisplacement.Len() {
		displacement = maxDisplacement
	}

	c.position = c.position.Add(displacement)
	c.speed = min(c.acceleration*deltaTime+c.speed, c.maxSpeed)
}
