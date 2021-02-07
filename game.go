package main

import (
	"fmt"
	"log"
)

// Drawn board is 29x17 characters (9x5 for each cell)
func drawBoard(log *log.Logger) {
	for row := 0; row < 3; row++ {
		// Each cell is 5 lines tall
		for cellLine := 0; cellLine < 5; cellLine++ {
			var line string
			for col := 0; col < 3; col++ {
				if playing {
					// Add bg color to selected cells
					if (row*3)+col == highlighted {
						// Add green bg and white fg to available cell
						if board[highlighted] != ' ' {
							line += "\x1b[1;37;41m"
						} else {
							// Add red bg and white fg to unavailable cell
							line += "\x1b[1;37;42m"
						}
					}
				} else if contains(winCells, (row*3)+col) {
					// Add blue bg to winning 3 cells
					line += "\x1b[1;37;44m"
				}
				switch board[(row*3)+col] {
				case 'X':
					line += xCell[cellLine]
				case 'O':
					line += oCell[cellLine]
				default:
					line += "         "
				}
				// Reset bg/fg colors
				line += "\x1b[0m"
				// Vertical divider
				if col != 2 {
					line += "|"
				}
			}
			log.Println(line)
		}
		// Horizontal divider
		if row != 2 {
			log.Println("---------+---------+---------")
		}
	}
	log.Println()
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
	selectCell()
}

func selectCell() {
	if board[highlighted] != ' ' {
		return
	}
	if xTurn {
		board[highlighted] = 'X'
	} else {
		board[highlighted] = 'O'
	}
	// Default highlight to middle cell
	turn++
	highlighted = 4
	xTurn = !xTurn
}

func checkGameState() (bool, string) {
	// Check for game over conditions
	if winner := getWinner(); winner != ' ' {
		return true, fmt.Sprintf("%c wins!\n", winner)
	}

	if turn == 9 {
		return true, fmt.Sprintf("It's a tie!\n")
	}

	return false, ""
}

func getWinner() byte {
	set := []int{}
	var validWinner byte
	// The winner to check is the previously played turn
	if xTurn {
		validWinner = 'O'
	} else {
		validWinner = 'X'
	}
	for i, v := range board {
		if v == validWinner {
			set = append(set, i)
		}
	}

	for _, winSet := range winSets {
		if contains(set, winSet[0]) && contains(set, winSet[1]) && contains(set, winSet[2]) {
			winCells = []int{winSet[0], winSet[1], winSet[2]}
			return validWinner
		}
	}
	return ' '
}

func reset() {
	playing = true
	turn = 0
	board = []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	winCells = []int{}
	// Switch start player every other game
	xTurn = true
	highlighted = 4
}
