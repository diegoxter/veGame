package component

import (
	"time"
)

type Weapon struct {
	Ammo int

	Damage int

	FireRate time.Duration
	LastFire time.Time

	BulletSpeed float64
}