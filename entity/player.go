package entity

import (
	"github.com/diegoxter/blockgame/asset"
	"github.com/diegoxter/blockgame/component"

	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.Position{})
	player.AddComponent(&component.Velocity{})

	weapon := &component.Weapon{
		Damage:      1,
		FireRate:    144 / 16,
		BulletSpeed: 8,
	}
	player.AddComponent(weapon)

	player.AddComponent(&component.Sprite{
		Image: asset.ImgPlayer,
	})

	player.AddComponent(&component.Rail{})

	return player
}