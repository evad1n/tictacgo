package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/eiannone/keyboard"
)

type (
	move struct {
		player *player  // The player sending the move
		do     moveFunc // The move to perform
		end    bool     // Signifies game end
	}

	// Player is
	player struct {
		number int
		stdout io.Writer
		wins   int
	}
)

var (
	xTurn         bool
	gameNumber    int
	highlighted   int
	board         []byte
	turn          int
	playing       bool
	listening     bool
	remoteIsAlive bool
)

const (
	port = "9090"
)

var (
	keyboardInput chan keyboard.KeyEvent
	moves         chan move
)

// Listen for keyboard input in background
func listenKeyboard() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	defer keyboard.Close()

	loadMoves()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)

		switch key {
		case keyboard.KeyCtrlC:
			os.Exit(0)
		default:
			keyboardInput <- keyboard.KeyEvent{
				Key:  key,
				Rune: char,
				Err:  nil,
			}
		}
	}
}

func main() {
	// Set up keyboard listening
	keyboardInput = make(chan keyboard.KeyEvent)
	go listenKeyboard()

	// Move event channel
	moves = make(chan move)

	// Menu selection
	mode := 0
	modes := []string{"1. Single Keyboard", "2. Remote Connection"}

	for {
		clearScreen()
		fmt.Print("Welcome to Tic-Tac-Go!")
		fmt.Println("Select game mode")
		drawMenu(modes, mode)

		keyEv := <-keyboardInput

		switch keyEv.Key {
		case keyboard.KeyArrowUp, keyboard.KeyArrowDown:
			mode++
			mode %= 2
		case keyboard.KeyEnter:
			switch mode {
			case 1:
				awaitPlayer()
			default:
				startLocal()
			}
		}
	}
}

func startLocal() {
	fmt.Println("Starting local game")
	go listenLocalGame()
	var keyEv keyboard.KeyEvent
	gameNumber = 0
	for replay := true; replay; replay = (keyEv.Key == keyboard.KeyEnter) {
		gameNumber++
		reset()
		localGameLoop()
		fmt.Println("Press ENTER to play again!")

		keyEv = <-keyboardInput
	}
}

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

func listenLocalGame() {
	for m := range moves {
		fmt.Println("move")
		if m.end {
			log.Println("Game session ended")
			return
		}
		if playing {
			fmt.Println("doing move")
			m.do()
			// Check for winner
			drawBoard(os.Stdout)
		}
	}
}

func localGameLoop() {
	for playing {
		drawBoard(os.Stdout)
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

// Drawn board is 29x17 characters (9x5 for each cell)
func drawBoard(w io.Writer) {
	// clearScreen()
	if xTurn {
		fmt.Fprintln(w, "X's turn:")
	} else {
		fmt.Fprintln(w, "O's turn")
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
			fmt.Fprintln(w, line)
		}
		// Horizontal divider
		if row != 2 {
			fmt.Fprintln(w, "---------+---------+---------")
		}
	}
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

func getWinner() byte {
	set := []int{}
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

func reset() {
	playing = true
	turn = 0
	board = []byte{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	winCells = []int{}
	// Switch start player every other game
	xTurn = true
	highlighted = 4
}

func nextTurn() {
	// Default highlight to middle cell
	highlighted = 4
	xTurn = !xTurn
}

func gameOver() {
	playing = false
	drawBoard(os.Stdout)
	fmt.Println("Game over!")
}
