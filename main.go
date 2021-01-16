package main

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

var xTurn bool
var highlighted int
var board []byte
var turn int
var playing bool

// Winning sets
var winSets [][]int = [][]int{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},

	{0, 4, 8},
	{6, 4, 2},

	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
}

var winCells []int

// The cell characters
// Width 9
var xCell []string = []string{
	"X       X",
	"  X   X  ",
	"    X    ",
	"  X   X  ",
	"X       X"}
var oCell []string = []string{
	"  OOOOO  ",
	" O     O ",
	" O     O ",
	" O     O ",
	"  OOOOO  "}

func clear() {
	fmt.Print("\x1b[2J")
}

// Drawn board is 29x17 characters (9x5 for each cell)
func drawBoard() {
	clear()
	if xTurn {
		fmt.Println("X's turn:")
	} else {
		fmt.Println("O's turn")
	}
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
							line += "\x1b[1;41m\x1b[1;37m"
						} else {
							// Add red bg and white fg to unavailable cell
							line += "\x1b[1;42m\x1b[1;37m"
						}
					}
				} else if contains(winCells, (row*3)+col) {
					// Add blue bg to winning 3 cells
					line += "\x1b[1;44m\x1b[1;37m"
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
				line += "\x1b[1;00m"
				// Vertical divider
				if col != 2 {
					line += "|"
				}
			}
			fmt.Println(line)
		}
		// Horizontal divider
		if row != 2 {
			fmt.Println("---------+---------+---------")
		}
	}
}

func nextTurn() {
	// Default highlight to middle cell
	highlighted = 4
	xTurn = !xTurn
}

func reset() {
	playing = true
	turn = 0
	board = []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	winCells = []int{}
	xTurn = true
	highlighted = 4
	gameLoop()
}

// Why do I have to write this; why can't go have this built-in? :C
func contains(set []int, item int) bool {
	for _, x := range set {
		if x == item {
			return true
		}
	}
	return false
}

func getWinner() byte {
	set := []int{}
	// Whose turn decides the valid winning(vw) byte
	var vw byte
	if xTurn {
		vw = 'X'
	} else {
		vw = 'O'
	}
	for i, v := range board {
		if v == vw {
			set = append(set, i)
		}
	}

	for _, winSet := range winSets {
		if contains(set, winSet[0]) && contains(set, winSet[1]) && contains(set, winSet[2]) {
			winCells = []int{winSet[0], winSet[1], winSet[2]}
			return vw
		}
	}
	return ' '
}

func gameOver() {
	playing = false
	drawBoard()
	fmt.Println("Game over!")
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
	turn++

	// Check for game over conditions
	if winner := getWinner(); winner != ' ' {
		gameOver()
		fmt.Printf("%c wins!\n", winner)
		return
	}

	if turn == 9 {
		gameOver()
		fmt.Println("It's a tie!")
		return
	}
	nextTurn()
	return
}

func gameLoop() {
	for playing {
		drawBoard()
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		// Otherwise Ctrl+C gets eaten
		case keyboard.KeyCtrlC:
			os.Exit(0)
		case keyboard.KeyEsc:
			return
		case keyboard.KeyEnter, keyboard.KeySpace:
			selectCell()
		case keyboard.KeyArrowUp:
			if highlighted > 2 {
				highlighted -= 3
			}
		case keyboard.KeyArrowDown:
			if highlighted < 6 {
				highlighted += 3
			}
		case keyboard.KeyArrowLeft:
			if highlighted != 0 && highlighted != 3 && highlighted != 6 {
				highlighted--
			}
		case keyboard.KeyArrowRight:
			if highlighted != 2 && highlighted != 5 && highlighted != 8 {
				highlighted++
			}
		}
	}
}

func main() {
	fmt.Println("Welcome to Tic-Tac-Go!")
	fmt.Println("Press ENTER to begin")

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return
		case keyboard.KeyEnter:
			reset()
			fmt.Println("Press ENTER to play again!")
		}
	}
}
