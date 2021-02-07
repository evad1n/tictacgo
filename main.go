package main

import (
	"fmt"
	"io"

	"github.com/eiannone/keyboard"
)

type (
	move struct {
		player *player  // The player sending the move
		do     moveFunc // The move to perform
		end    bool     // Signifies game end
	}

	// Player is
	player struct {
		number int
		stdout io.Writer
		wins   int
	}
)

var (
	xTurn         bool
	gameNumber    int
	highlighted   int
	board         []byte
	turn          int
	playing       bool
	listening     bool
	remoteIsAlive bool
)

const (
	port = "9090"
)

var (
	keyboardInput chan keyboard.KeyEvent
	moves         chan move
)

func main() {
	// Set up keyboard listening
	keyboardInput = make(chan keyboard.KeyEvent)
	go listenKeyboard()

	// Move event channel
	moves = make(chan move)

	// Menu selection
	mode := 0
	modes := []string{"1. Single Keyboard", "2. Remote Connection"}

	for {
		clearScreen()
		fmt.Print("Welcome to Tic-Tac-Go!")
		fmt.Println("Select game mode")
		drawMenu(modes, mode)

		keyEv := <-keyboardInput

		switch keyEv.Key {
		case keyboard.KeyArrowUp, keyboard.KeyArrowDown:
			mode++
			mode %= 2
		case keyboard.KeyEnter:
			switch mode {
			case 1:
				awaitPlayer()
			default:
				startLocal()
			}
			break
		case keyboard.KeyEsc:
			return
		}
	}
}
