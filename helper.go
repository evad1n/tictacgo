package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

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
func clearScreen() {
	fmt.Print("\x1b[2J")
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

func keyboardPrintLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	// One char at a time
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		r := scanner.Text()
		// fmt.Println(rune(r[0]))
		fmt.Println(rune(r[0]), r)
		fmt.Println("HELLO")

		// moves <- move{p, key}
	}
}

type moveFunc func()

func moveUp() {
	if highlighted > 2 {
		highlighted -= 3
	}
}
func moveDown() {
	if highlighted < 6 {
		highlighted += 3
	}
}
func moveLeft() {
	if highlighted != 0 && highlighted != 3 && highlighted != 6 {
		highlighted--
	}
}
func moveRight() {
	if highlighted != 2 && highlighted != 5 && highlighted != 8 {
		highlighted++
	}
}
func moveSelect() {
	if playing {
		selectCell()
	}
}

var (
	runeToMove map[rune]moveFunc
	keyToMove  map[keyboard.Key]moveFunc
)

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
