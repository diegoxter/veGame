package system

import (
	"github.com/diegoxter/blockgame/component"
	"github.com/diegoxter/blockgame/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const CameraMoveSpeed = 0.132

type CameraSystem struct {
	Weapon   *component.Weapon
	Position *component.Position
}

func NewCameraSystem() *CameraSystem {
	s := &CameraSystem{}

	return s
}

func (s *CameraSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	world.World.CamMoving = world.World.CamY > 0
	if world.World.CamMoving {
		world.World.CamY -= CameraMoveSpeed
	}
	// else {
	// 	world.World.GameOver = true
	// }
	return nil
}

func (*CameraSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}