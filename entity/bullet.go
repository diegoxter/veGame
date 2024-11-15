
package entity

import (
	"code.rocketnine.space/tslocum/gohan"
	"github.com/diegoxter/blockgame/component"
)

func NewBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
	bullet := gohan.NewEntity()

	bullet.AddComponent(&component.Position{
		X: x,
		Y: y,
	})

	bullet.AddComponent(&component.Velocity{
		X: xSpeed,
		Y: ySpeed,
	})

	bullet.AddComponent(&component.Bullet{})

	return bullet
}