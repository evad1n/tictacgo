package main

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
)

func startLocal() {
	fmt.Println("Starting local game")
	go listenLocalGame()
	var keyEv keyboard.KeyEvent
	gameNumber = 0
	for {
		gameNumber++
		reset()
		clearScreen(localLog)
		drawBoard(localLog)
		fmt.Println("X's turn:")
		localGameLoop()
		clearScreen(localLog)
		fmt.Println("Press ENTER to play again")
		fmt.Println("Press ESC to go back to the main menu")

		for {
			keyEv = <-keyboardInput

			if keyEv.Key == keyboard.KeyEsc {
				return
			} else if keyEv.Key == keyboard.KeyEnter {
				break
			}
		}
	}
}

func listenLocalGame() {
	for m := range moves {
		if m.end {
			log.Println("Game session ended")
			return
		}
		if playing {
			clearScreen(localLog)
			m.do()
			// Draw new board
			drawBoard(localLog)
			// Check for winner
			if gameOver, endMsg := checkGameState(); gameOver {
				playing = false
				// Draw winning tiles
				clearScreen(localLog)
				drawBoard(localLog)
				fmt.Println("Game Over!")
				fmt.Print(endMsg)
				fmt.Println("Press any key to continue")
			} else {
				if xTurn {
					fmt.Println("X's turn:")
				} else {
					fmt.Println("O's turn")
				}
			}
		}
	}
}

func localGameLoop() {
	for playing {
		keyEv := <-keyboardInput
		// End current game session on escape
		if keyEv.Key == keyboard.KeyEsc {
			return
		}
		if validMoveFunc, exists := keyToMove[keyEv.Key]; exists {
			moves <- move{
				player: nil,
				do:     validMoveFunc,
			}
		} else if validMoveFunc, exists := runeToMove[keyEv.Rune]; exists {
			moves <- move{
				player: nil,
				do:     validMoveFunc,
			}
		}
	}
}
