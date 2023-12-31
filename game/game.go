package game

import (
	"fmt"
	"kadsin/shoot-run/game/assets"
	"kadsin/shoot-run/game/entities"
	"time"

	term "github.com/nsf/termbox-go"
)

type Game struct {
	Screen             assets.Screen
	Exited             bool
	Shooter            entities.Shooter
	Enemies            []*entities.Enemy
	Blocks             []entities.Object
	KilledEnemiesCount uint
	LastTimeActions    LastActionAt
	StartedAt          int64
}

type LastActionAt struct {
	BlocksGenerator int64
	Enemies         map[*entities.Enemy]int64
	EnemyGenerator  int64
	Shooter         int64
	Bullets         int64
	Kill            int64
}

func (game *Game) Start() {
	game.showStoryReady()

	game.StartedAt = time.Now().Unix()

	game.LastTimeActions.Enemies = make(map[*entities.Enemy]int64)

	game.Shooter.Speed = assets.SPEED_SHOOTER
	game.Shooter.Person.Location = assets.Coordinate{
		X: game.Screen.End.X / 2,
		Y: game.Screen.End.Y / 2,
	}

	go game.listenToKeyboard()

	game.update()
}

func (game *Game) listenToKeyboard() {
	for {
		var event = term.PollEvent()

		if event.Type == term.EventKey {
			switch event.Key {
			case term.KeyArrowLeft:
				game.Shooter.Person.MoveLeft()
			case term.KeyArrowRight:
				game.Shooter.Person.MoveRight()
			case term.KeyArrowUp:
				game.Shooter.Person.MoveUp()
			case term.KeyArrowDown:
				game.Shooter.Person.MoveDown()
			case term.KeySpace:
				go game.Shooter.Shoot()
			case term.KeyCtrlC:
				game.Exited = true
			}
		}
	}
}

func (game Game) ScreenTime() string {
	screenTime := time.Now().Unix() - game.StartedAt

	minutes := screenTime / 60
	seconds := screenTime % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (game *Game) ScreenCircumference() int {
	return 2*game.Screen.End.X + 2*game.Screen.End.Y
}
