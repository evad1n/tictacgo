package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/eiannone/keyboard"
)

func awaitPlayer() {
	clearScreen(localLog)
	fmt.Println("Waiting for player to connect...")
	fmt.Printf("Host IP: %s\n", getLocalAddress())
	fmt.Printf("Listening for connections on port %s\n", port)
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("cant't start server on port %s: %v", port, err)
	}
	defer server.Close()
	// Blocking
	listening = true
	conn, err := server.Accept()
	if err != nil {
		log.Fatalf("error accepting connection: %v", err)
	}
	defer conn.Close()

	log.Printf("Client connected from %s\n", conn.RemoteAddr().String())
	fmt.Fprintf(conn, "\nWelcome to Tic-Tac-Go!\n\n")
	disableLineMode(conn)
	fmt.Fprintln(conn, "\x1b[?25l")

	// Remote player
	p2 := &player{
		number: 1,
		log:    log.New(conn, "\x1b[0G", 0),
		wins:   0,
	}
	// Local player
	p1 := &player{
		number: 0,
		log:    localLog,
		wins:   0,
	}
	// Player input loops
	go handleRemotePlayer(p2, conn)
	go handleLocalPlayer(p1)
	startRemote(p1, p2)

	fmt.Println("Game session ended")
}

func handleRemotePlayer(p *player, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	// One char at a time
	scanner.Split(bufio.ScanRunes)
	buffer := []rune{}

	sequence := false
	for scanner.Scan() {
		r := rune(scanner.Text()[0])

		// Check if it is a control character
		if sequence || r < rune(40) {
			sequence = true
			// Update list
			buffer = append(buffer, r)
			if len(buffer) > 3 {
				// Buffer overflow, so reset
				buffer = nil
				sequence = false
			}
			// Check if the buffer has a valid char
			if key, err := checkSpecial(buffer); err == nil {
				if m, valid := keyToMove[key]; valid {
					moves <- move{
						player: p,
						do:     m,
					}
				}
				// Buffer used, so discard
				buffer = nil
				sequence = false
			}
		} else {
			// Treat as single character
			if m, valid := runeToMove[r]; valid {
				moves <- move{
					player: p,
					do:     m,
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("remote scanner: %v", err)
	}

	// Send game over
	moves <- move{
		player: p,
		do:     nil,
		end:    true,
	}
}

func startRemote(p1 *player, p2 *player) {
	fmt.Println("Starting remote game")

	gameNumber = 0
	p1Ready := true
	p2Ready := true
	for replay := true; replay; replay = (p1Ready && p2Ready) {
		gameNumber++
		reset()
		clearScreen(p1.log)
		clearScreen(p2.log)
		drawBoard(p1.log)
		drawBoard(p2.log)
		// Switch start player every other game
		if gameNumber > 1 {
			p1.number = (p1.number + 1) % 2
			p2.number = (p2.number + 1) % 2
		}

		printTurn(p1, p2)
		if completed := listenRemoteGame(p1, p2); !completed {
			gameNumber--
		} else {
			clearScreen(p1.log)
			clearScreen(p2.log)
		}

		printScores(p1, p2, gameNumber)
		p2.log.Println("Waiting for host to restart...")

		p1.log.Println("Press ENTER to play again")
		p1.log.Println("Press ESC to go back to the main menu (this will end the connection)")

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

func listenRemoteGame(p1 *player, p2 *player) bool {
	for m := range moves {
		if m.end {
			playing = false
			clearScreen(p1.log)
			clearScreen(p2.log)
			log.Println("Game ended early")
			p2.log.Println("Game ended early by host")
			return false
		}
		if playing {
			// If it is their turn
			if m.player.number == turn%2 {
				clearScreen(p1.log)
				clearScreen(p2.log)
				m.do()
				// Draw new board
				drawBoard(p1.log)
				drawBoard(p2.log)
				// Check for winner
				if gameOver, winner := checkGameState(); gameOver {
					playing = false
					// Draw winning tiles
					clearScreen(p1.log)
					clearScreen(p2.log)
					drawBoard(p1.log)
					drawBoard(p2.log)
					p1.log.Println("Game Over!")
					p2.log.Println("Game Over!")
					printWinner(winner, p1, p2)
					return true
				}
				printTurn(p1, p2)
			}
		}
	}
	return false
}

func printTurn(p1 *player, p2 *player) {
	if turn%2 == p1.number {
		p1.log.Println("Your turn")
	} else {
		p1.log.Println("Opponent's turn")
	}

	if turn%2 == p2.number {
		p2.log.Println("Your turn")
	} else {
		p2.log.Println("Opponent's turn")
	}
}

func printWinner(winner byte, p1 *player, p2 *player) {
	switch winner {
	case 'X':
		if p1.number == 0 {
			p1.log.Println("You win!")
			p2.log.Println("You lose!")
			p1.wins++
		} else {
			p1.log.Println("You lose!")
			p2.log.Println("You win!")
			p2.wins++
		}
		break
	case 'O':
		if p1.number == 1 {
			p1.log.Println("You win!")
			p2.log.Println("You lose!")
			p1.wins++
		} else {
			p1.log.Println("You lose!")
			p2.log.Println("You win!")
			p2.wins++
		}
		break
	default:
		p1.log.Println("It's a tie!")
		p2.log.Println("It's a tie!")
		break
	}
}

func printScores(p1 *player, p2 *player, totalGames int) {
	p1.log.Printf("Your wins: %d\n", p1.wins)
	p1.log.Printf("Opponent wins: %d\n", p2.wins)
	p1.log.Printf("Ties: %d\n", totalGames-(p1.wins+p2.wins))

	p2.log.Printf("Your wins: %d\n", p2.wins)
	p2.log.Printf("Opponent wins: %d\n", p1.wins)
	p2.log.Printf("Ties: %d\n", totalGames-(p1.wins+p2.wins))
}
