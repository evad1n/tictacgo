package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

var xTurn bool
var highlighted int
var board []byte
var turn int
var playing bool

// Win sets
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

// Width 10
var xCell []string = []string{"X       X", "  X   X  ", "    X    ", "  X   X  ", "X       X"}
var oCell []string = []string{"  OOOOO  ", " O     O ", " O     O ", " O     O ", "  OOOOO  "}

// Drawn board is 17x17 characters (5x5 for each cell)
func drawBoard() {
	clear()
	if xTurn {
		fmt.Println("X's turn:")
	} else {
		fmt.Println("O's turn")
	}
	for row := 0; row < 3; row++ {
		for cellLine := 0; cellLine < 5; cellLine++ {
			var line string
			for col := 0; col < 3; col++ {
				// Add green bg and white fg to selected cell
				if (row*3)+col == highlighted {
					line += "\x1b[1;42m\x1b[1;37m"
				}
				switch board[(row*3)+col] {
				case 'X':
					line += xCell[cellLine]
				case 'O':
					line += oCell[cellLine]
				default:
					line += "         "
				}
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

func reset() int {
	playing = true
	turn = 0
	board = []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	xTurn = true
	highlighted = 4
	return gameLoop()
}

func contains(set map[int]struct{}, item int) bool {
	_, ok := set[item]
	return ok
}

func gameOver() byte {
	set := make(map[int]struct{})
	var vw byte
	if xTurn {
		vw = 'X'
	} else {
		vw = 'O'
	}
	for i, v := range board {
		if v == vw {
			set[i] = struct{}{}
		}
	}

	for _, winSet := range winSets {
		if contains(set, winSet[0]) && contains(set, winSet[1]) && contains(set, winSet[2]) {
			return vw
		}
	}
	return ' '
}

func selectCell() bool {
	if board[highlighted] != ' ' {
		return false
	}
	if xTurn {
		board[highlighted] = 'X'
	} else {
		board[highlighted] = 'O'
	}
	// Check for game over conditions
	if winner := gameOver(); winner != ' ' {
		drawBoard()
		fmt.Printf("%c wins!\n", winner)
		return true
	}

	if turn == 9 {
		drawBoard()
		fmt.Println("It's a tie!")
		return true
	}
	nextTurn()
	return false
}

func gameLoop() int {
	for playing {
		drawBoard()
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return 0
		case keyboard.KeyEnter, keyboard.KeySpace:
			if gameOver := selectCell(); gameOver {
				return 1
			}
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
	return 0
}

func clear() {
	fmt.Print("\x1b[2J")
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
			if r := reset(); r == 0 {
				return
			}
			fmt.Println("Press ENTER to play again!")
		}
	}
}
