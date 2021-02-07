package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
)

func startLocal() {
	fmt.Println("Starting local game")
	go handleLocalPlayer(nil)
	xWins := 0
	oWins := 0
	gameNumber = 0
	for {
		gameNumber++
		reset()
		clearScreen(localLog)
		drawBoard(localLog)
		fmt.Println("X's turn:")
		if completed := listenLocalGame(&xWins, &oWins); !completed {
			gameNumber--
		} else {
			clearScreen(localLog)
		}
		fmt.Printf("X wins: %d\n", xWins)
		fmt.Printf("O wins: %d\n", oWins)
		fmt.Printf("Ties: %d\n", gameNumber-(xWins+oWins))
		fmt.Println("Press ENTER to play again")
		fmt.Println("Press ESC to go back to the main menu")

		for {
			keyEv := <-keyboardInput

			if keyEv.Key == keyboard.KeyEsc {
				return
			} else if keyEv.Key == keyboard.KeyEnter {
				break
			}
		}
	}
}

// Listen for one game. Returns if the game successfully completed
func listenLocalGame(xWins *int, oWins *int) bool {
	for m := range moves {
		if m.end {
			playing = false
			clearScreen(localLog)
			fmt.Print("Game ended early\n\n")
			return false
		}
		if playing {
			clearScreen(localLog)
			m.do()
			// Draw new board
			drawBoard(localLog)
			// Check for winner
			if gameOver, winner := checkGameState(); gameOver {
				playing = false
				// Draw winning tiles
				clearScreen(localLog)
				drawBoard(localLog)
				fmt.Println("Game Over!")
				switch winner {
				case 'X':
					fmt.Println("XWins!")
					*xWins++
					break
				case 'O':
					fmt.Println("O Wins!")
					*oWins++
					break
				default:
					fmt.Println("It's a tie")
					break
				}
				return true
			}

			if xTurn {
				fmt.Println("X's turn:")
			} else {
				fmt.Println("O's turn")
			}
		}
	}
	return false
}

func handleLocalPlayer(p *player) {
	for {
		keyEv := <-keyboardInput
		// Put back menu events if we aren't playing
		if !playing {
			keyboardInput <- keyEv
			continue
		}
		// End current game session on escape
		if keyEv.Key == keyboard.KeyEsc {
			moves <- move{
				player: p,
				do:     nil,
				end:    true,
			}
		}
		if validMoveFunc, exists := keyToMove[keyEv.Key]; exists {
			moves <- move{
				player: p,
				do:     validMoveFunc,
			}
		} else if validMoveFunc, exists := runeToMove[keyEv.Rune]; exists {
			moves <- move{
				player: p,
				do:     validMoveFunc,
			}
		}
	}
}
