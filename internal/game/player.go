package game

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/level"
	"github.com/vparent05/minecraft_go/internal/movement"
	"github.com/vparent05/minecraft_go/internal/utils"
)

type player struct {
	*movement.EntityController
	game           *Game
	cameraOffsets  []mgl32.Vec3
	selectedCamera int
	renderDistance int
	reach          float32
	updates        *utils.UpdateChannel[level.LevelObserver]
}

func NewPlayer(game *Game) *player {
	p := &player{
		movement.NewEntityController(mgl32.Vec3{0, 65, 0}, 1000, 1000, 5000, 0, 0, 4*math.Pi/9),
		game,
		[]mgl32.Vec3{{0, 0, 0}},
		0,
		16,
		5,
		utils.NewUpdateChannel[level.LevelObserver](),
	}

	p.updates.Send(p.asLevelObserver())

	return p
}

func (p *player) CameraPosition() mgl32.Vec3 {
	return p.Position().Add(p.cameraOffsets[p.selectedCamera])
}

func (p *player) asLevelObserver() level.LevelObserver {
	return level.LevelObserver{
		Vec3:           p.CameraPosition(),
		RenderDistance: p.renderDistance,
	}
}

func (p *player) FrameTick(deltaTime float32) {
	if glfw.GetCurrentContext().GetKey(glfw.KeyEscape) == glfw.Press {
		glfw.GetCurrentContext().SetShouldClose(true)
	}

	directions := make([]movement.Direction, 0, 6)
	if glfw.GetCurrentContext().GetKey(glfw.KeyW) == glfw.Press {
		directions = append(directions, movement.FRONT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyA) == glfw.Press {
		directions = append(directions, movement.LEFT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyS) == glfw.Press {
		directions = append(directions, movement.BACK)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyD) == glfw.Press {
		directions = append(directions, movement.RIGHT)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == glfw.Press {
		directions = append(directions, movement.UP)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftShift) == glfw.Press {
		directions = append(directions, movement.DOWN)
	}

	width, height := glfw.GetCurrentContext().GetSize()
	xpos, ypos := glfw.GetCurrentContext().GetCursorPos()
	glfw.GetCurrentContext().SetCursorPos(float64(width)/2, float64(height)/2)

	p.UpdateOrientation(mgl32.Vec2{(float32(width/2) - float32(xpos)) / 1000, (float32(height)/2 - float32(ypos)) / 1000})
	p.UpdatePosition(directions, deltaTime)

	p.updates.Send(p.asLevelObserver())

	if glfw.GetCurrentContext().GetMouseButton(glfw.MouseButton1) == glfw.Press {
		targeted, _ := p.game.Level.CastRay(p.CameraPosition(), p.Orientation(), p.reach)
		targeted.Set(level.AIR)
	}
	if glfw.GetCurrentContext().GetMouseButton(glfw.MouseButton2) == glfw.Press {
		_, front := p.game.Level.CastRay(p.CameraPosition(), p.Orientation(), p.reach)
		front.Set(level.STONE)
	}
}
