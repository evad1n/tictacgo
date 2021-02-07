package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/eiannone/keyboard"
)

func awaitPlayer() {
	clearScreen()
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
		stdout: conn,
		wins:   0,
	}
	// Local player
	p1 := &player{
		number: 0,
		stdout: os.Stdout,
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
				moves <- move{
					player: p,
					do:     keyToMove[key],
				}
				// Buffer used, so discard
				buffer = nil
				sequence = false
			}
		} else {
			// Treat as single character
			moves <- move{
				player: p,
				do:     runeToMove[r],
			}
		}
	}

	if err := scanner.Err(); err != nil {

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
	drawBoard(p1.stdout)
	drawBoard(p2.stdout)
	for m := range moves {
		if m.end {
			log.Println("Game session ended")
			return
		}
		if playing {
			// If it is their turn
			if m.player.number == turn%2 {
				m.do()
				// Check for winner
				drawBoard(p1.stdout)
				drawBoard(p2.stdout)
			}
		}
	}
}
