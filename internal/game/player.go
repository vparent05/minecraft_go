package game

import (
	"math"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/vparent05/minecraft_go/internal/level"
	"github.com/vparent05/minecraft_go/internal/movement"
	"github.com/vparent05/minecraft_go/internal/utils/atomicx"
	"github.com/vparent05/minecraft_go/internal/utils/chanx"
	"github.com/vparent05/minecraft_go/internal/utils/debounce"
)

type player struct {
	*movement.EntityController
	game                 *Game
	cameraOffsets        []mgl32.Vec3
	selectedCamera       int
	renderDistance       int
	reach                float32
	levelObserver        *atomicx.Value[level.LevelObserver]
	levelObserverUpdates chan struct{}

	blockAction *debounce.Debounce
}

func NewPlayer(game *Game) *player {
	p := &player{
		EntityController:     movement.NewEntityController(mgl32.Vec3{0, 65, 0}, 1000, 1000, 5000, 0, 0, 4*math.Pi/9),
		game:                 game,
		cameraOffsets:        []mgl32.Vec3{{0, 0, 0}},
		selectedCamera:       0,
		renderDistance:       16,
		reach:                16,
		levelObserver:        &atomicx.Value[level.LevelObserver]{},
		levelObserverUpdates: make(chan struct{}, 1),
		blockAction:          debounce.NewDebounce(100 * time.Millisecond),
	}

	p.levelObserver.Store(p.asLevelObserver())
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

	p.levelObserver.Store(p.asLevelObserver())
	chanx.TrySend(p.levelObserverUpdates, struct{}{})

	if glfw.GetCurrentContext().GetMouseButton(glfw.MouseButton1) == glfw.Press {
		targeted, _ := p.game.Level.CastRay(p.CameraPosition(), p.Orientation(), p.reach)
		if targeted != nil {
			p.blockAction.Do(func() { targeted.Set(level.AIR) })
		}
	}

	if glfw.GetCurrentContext().GetMouseButton(glfw.MouseButton2) == glfw.Press {
		_, front := p.game.Level.CastRay(p.CameraPosition(), p.Orientation(), p.reach)
		if front != nil {
			p.blockAction.Do(func() { front.Set(level.STONE) })
		}
	}
}
