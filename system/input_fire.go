
package system

import (
	"math"
	"time"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/diegoxter/blockgame/component"
	"github.com/diegoxter/blockgame/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

func angle(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y1-y2, x1-x2)
}

type fireInputSystem struct {
	Position *component.Position
	Weapon   *component.Weapon
}

func NewFireInputSystem() *fireInputSystem {
	return &fireInputSystem{}
}

func (s *fireInputSystem) fire(fireAngle float64) {
	if time.Since(s.Weapon.LastFire) < s.Weapon.FireRate {
		return
	}

	s.Weapon.Ammo--
	s.Weapon.LastFire = time.Now()

	speedX := math.Cos(fireAngle) * -s.Weapon.BulletSpeed
	speedY := math.Sin(fireAngle) * -s.Weapon.BulletSpeed

	bullet := entity.NewBullet(s.Position.X, s.Position.Y, speedX, speedY)
	_ = bullet
}

func (s *fireInputSystem) Update(entity gohan.Entity) error {
	if s.Weapon.Ammo <= 0 {
		return nil
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		fireAngle := angle(s.Position.X, s.Position.Y, float64(cursorX), float64(cursorY))

		s.fire(fireAngle)
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		cursorX, cursorY := ebiten.CursorPosition()
		fireAngle := angle(s.Position.X, s.Position.Y, float64(cursorX), float64(cursorY))

		const div = 5
		s.Weapon.BulletSpeed /= div
		for i := 0.0; i < 24; i++ {
			s.fire(fireAngle + i*(math.Pi/12))
			s.Weapon.LastFire = time.Time{}
		}
		s.Weapon.BulletSpeed *= div
	}

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(math.Pi / 4)
	case ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(-math.Pi / 4)
	case ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(math.Pi * .75)
	case ebiten.IsKeyPressed(ebiten.KeyRight) && ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(-math.Pi * .75)
	case ebiten.IsKeyPressed(ebiten.KeyLeft):
		s.fire(0)
	case ebiten.IsKeyPressed(ebiten.KeyRight):
		s.fire(math.Pi)
	case ebiten.IsKeyPressed(ebiten.KeyUp):
		s.fire(math.Pi / 2)
	case ebiten.IsKeyPressed(ebiten.KeyDown):
		s.fire(-math.Pi / 2)
	}

	return nil
}

func (*fireInputSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}