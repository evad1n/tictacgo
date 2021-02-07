package main

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
