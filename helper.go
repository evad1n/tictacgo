package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

var (
	runeToMove map[rune]moveFunc
	keyToMove  map[keyboard.Key]moveFunc
)

// Initialize move maps
func loadMoves() {
	runeToMove = make(map[rune]moveFunc)
	runeToMove['w'] = moveUp
	runeToMove['s'] = moveDown
	runeToMove['a'] = moveLeft
	runeToMove['d'] = moveRight

	keyToMove = make(map[keyboard.Key]moveFunc)
	keyToMove[keyboard.KeyArrowUp] = moveUp
	keyToMove[keyboard.KeyArrowDown] = moveDown
	keyToMove[keyboard.KeyArrowLeft] = moveLeft
	keyToMove[keyboard.KeyArrowRight] = moveRight
	keyToMove[keyboard.KeyEnter] = moveSelect
}

// Pseudo state machine to process my limited characters from telnet
func checkSpecial(buffer []rune) (keyboard.Key, error) {
	// for i, r := range buffer {
	// 	fmt.Printf("%d: %d, ", i, r)
	// }

	if len(buffer) > 1 {
		switch buffer[0] {
		case 13:
			// Enter key start
			switch buffer[1] {
			case 0:
				return keyboard.KeyEnter, nil
			}
			break
		case 27:
			// Arrow key start
			switch buffer[1] {
			case 91:
				if len(buffer) > 2 {
					// the '[', part of arrow sequence
					switch buffer[2] {
					case 65:
						return keyboard.KeyArrowUp, nil
					case 66:
						return keyboard.KeyArrowDown, nil
					case 68:
						return keyboard.KeyArrowLeft, nil
					case 67:
						return keyboard.KeyArrowRight, nil
					}
				}
			}
			break
		}
	}
	return keyboard.KeyArrowDown, errors.New("Not a special key")
}

// Listen for keyboard input in background and send to keyboardInput channel
func listenKeyboard() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	defer keyboard.Close()

	loadMoves()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		switch key {
		case keyboard.KeyCtrlC:
			os.Exit(0)
		default:
			keyboardInput <- keyboard.KeyEvent{
				Key:  key,
				Rune: char,
				Err:  nil,
			}
		}
	}
}

// IO manip stuff

// Highlightable menu
func drawMenu(options []string, choice int) {
	for i, option := range options {
		if i == choice {
			fmt.Println(ansiWrap(option, "\x1b[30;46m"))
		} else {
			fmt.Println(option)
		}
	}
}

// Wrap some text in an ansi code
func ansiWrap(text string, code string) string {
	return fmt.Sprintf("%s%s\x1b[0m", code, text)
}

// Clears the screen
func clearScreen(log *log.Logger) {
	log.Print("\x1b[2J")
}

// Helper function for seeing if an int exists in a slice
func contains(set []int, item int) bool {
	for _, x := range set {
		if x == item {
			return true
		}
	}
	return false
}
