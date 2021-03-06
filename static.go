package main

import "log"

type (
	move struct {
		player *player  // The player sending the move
		do     moveFunc // The move to perform
		end    bool     // Signifies game end
	}

	player struct {
		number int         // The player number for keeping track of turns
		log    *log.Logger // The output log to print to
		wins   int         // The total wins for this sessions
	}

	moveFunc func() // A move that affects the board
)

var winCells []int

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
