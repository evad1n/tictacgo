package main

import (
	"fmt"
	"log"
	"os"

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
		log    *log.Logger
		wins   int
	}
)

var (
	xTurn       bool
	gameNumber  int
	highlighted int
	board       []byte
	turn        int
	playing     bool
	listening   bool
)

const (
	port = "9090"
)

var (
	keyboardInput chan keyboard.KeyEvent
	moves         chan move
)

// I have to use log because telnet is shitty and doesn't understand newlines
var (
	localLog *log.Logger
)

func main() {
	// Set up keyboard listening
	keyboardInput = make(chan keyboard.KeyEvent)
	go listenKeyboard()

	// Move event channel
	moves = make(chan move)

	// Just for consistency
	localLog = log.New(os.Stdout, "", 0)

	// Menu selection
	mode := 0
	modes := []string{"1. Single Keyboard", "2. Remote Connection"}

	for {
		clearScreen(localLog)
		fmt.Println("Welcome to Tic-Tac-Go!")
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
