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

	log.Printf("Client connected from %s\n", conn.RemoteAddr().String())
	fmt.Fprintf(conn, "\nWelcome to Tic-Tac-Go!\n\n")
	disableLineMode(conn)

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
	remoteIsAlive = true
	go handleRemotePlayer(p2, conn)
	go handleLocalPlayer(p1)
	startRemote(p1, p2)

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

func handleLocalPlayer(p *player) {
	for remoteIsAlive {
		keyEv := <-keyboardInput
		// End current game session on escape
		if keyEv.Key == keyboard.KeyEsc {

		}
		if validMoveFunc, exists := keyToMove[keyEv.Key]; exists {
			moves <- move{
				player: p,
				do:     validMoveFunc,
			}
		}
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

		listenRemoteGame(p1, p2)

		fmt.Println("Press ENTER to play again!")
	}
}

func listenRemoteGame(p1 *player, p2 *player) {
	for m := range moves {
		if m.end {
			log.Println("Game session ended")
			return
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
				if gameOver, endMsg := checkGameState(); gameOver {
					playing = false
					// Draw winning tiles
					clearScreen(p1.log)
					clearScreen(p2.log)
					drawBoard(p1.log)
					drawBoard(p2.log)
					// fmt.Println("Game Over!")
					fmt.Print(endMsg)
					// fmt.Println("Press any key to continue")
				} else {
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

			}
		}
	}
}
