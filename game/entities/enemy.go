package entities

import "math/rand"

type Enemy struct {
	Person Object
	Target *Object
	Speed  int
	OnKill func()
}

func (enemy *Enemy) Chase() {
	canChaseTwoDirections := rand.Float32() > 0.5
	canChaseVertical := rand.Float32() > 0.5

	if canChaseVertical || canChaseTwoDirections {
		enemy.chaseVertical()
	}

	if !canChaseVertical || canChaseTwoDirections {
		enemy.chaseHorizontal()
	}
}

func (enemy *Enemy) chaseVertical() {
	if enemy.Target.Location.Y > enemy.Person.Location.Y {
		enemy.Person.MoveDown()
	} else if enemy.Target.Location.Y < enemy.Person.Location.Y {
		enemy.Person.MoveUp()
	}

	enemy.Person.UpdateLocation(1)
}

func (enemy *Enemy) chaseHorizontal() {
	if enemy.Target.Location.X > enemy.Person.Location.X {
		enemy.Person.MoveRight()
	} else if enemy.Target.Location.X < enemy.Person.Location.X {
		enemy.Person.MoveLeft()
	}

	enemy.Person.UpdateLocation(1)
}
