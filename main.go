package main

import (
	"kadsin/shoot-run/entities"
	"time"

	term "github.com/nsf/termbox-go"
)

var shooter = entities.Object{X: 0, Y: 0, Shape: '●', Direction: entities.DIRECTION_RIGHT, Speed: 2}

var exit = false

func main() {
	term.Init()
	term.HideCursor()
	defer term.Close()

	startGame()
}

func startGame() {
	var width, height = term.Size()
	shooter.X = width / 2
	shooter.Y = height / 2

	ticker := time.NewTicker(time.Second / time.Duration(shooter.Speed))

	go listenToKeyboard()

	for range ticker.C {
		if exit {
			break
		}

		term.Clear(term.ColorDefault, term.ColorDefault)
		term.SetChar(shooter.X, shooter.Y, shooter.Shape)
		term.Sync()

		shooter.UpdateLocation()
	}
}

func listenToKeyboard() {
	for {
		if exit {
			break
		}

		var event = term.PollEvent()

		if event.Type == term.EventKey {
			switch event.Key {
			case term.KeyArrowLeft:
				shooter.MoveLeft()
			case term.KeyArrowRight:
				shooter.MoveRight()
			case term.KeyArrowUp:
				shooter.MoveUp()
			case term.KeyArrowDown:
				shooter.MoveDown()
			case term.KeyCtrlC:
				exit = true
			}
		}
	}
}
