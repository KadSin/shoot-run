package entities

import (
	"time"

	term "github.com/nsf/termbox-go"
)

type Shooter struct {
	Person  Object
	Bullets []*Object
}

func (shooter *Shooter) Shoot() {
	bullet := &Object{
		Shape:     '*',
		Direction: shooter.Person.Direction,
		Color:     term.ColorLightGray,
		Location:  shooter.Person.Location,
		Screen:    shooter.Person.Screen,
	}

	shooter.Bullets = append(shooter.Bullets, bullet)

	bullet.UpdateLocation(2)
}

func (shooter *Shooter) ListenToBullets(speed int) {
	ticker := time.NewTicker(time.Second / time.Duration(speed))

	for range ticker.C {
		for _, b := range shooter.Bullets {
			error := b.UpdateLocation(1)

			if error != nil {
				shooter.RemoveBullet(b)
			}
		}
	}
}

func (shooter *Shooter) RemoveBullet(bullet *Object) {
	for id, b := range shooter.Bullets {
		if b == bullet {
			shooter.Bullets[id] = nil

			if id == 0 {
				shooter.Bullets = shooter.Bullets[id+1:]
			} else if id == len(shooter.Bullets)-1 {
				shooter.Bullets = shooter.Bullets[:id-1]
			} else {
				shooter.Bullets = append(shooter.Bullets[id-1:], shooter.Bullets[:id]...)
			}

			break
		}
	}

}

func (shooter *Shooter) Run(speed int) {
	ticker := time.NewTicker(time.Second / time.Duration(speed))

	for range ticker.C {
		shooter.Person.UpdateLocation(1)
	}
}
