package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

const (
	port = "9090"
)

// Game variables
var (
	xTurn       bool
	highlighted int // The currently highlighted tile
	board       []byte
	turn        int // Keep track of turns, mostly for ties
	playing     bool
)

// Channels
var (
	keyboardInput chan keyboard.KeyEvent // Keyboard events from the local user
	moves         chan move              // The move channel for game events
)

// I have to use log because telnet is shitty and doesn't understand newlines
var (
	localLog *log.Logger // Just an alias for os.Stdout so that it fits my function signatures
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
