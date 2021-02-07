package main

import (
	"fmt"
	"net"
)

var (
	// The sequence to disable telnet line mode
	noLineMode = [][]byte{
		{255, 251, 1},
		{255, 251, 3},
		{255, 252, 34},
	}
)

// Disables telnet line mode
func disableLineMode(conn net.Conn) {
	for _, line := range noLineMode {
		conn.Write(line)
	}
}

// Handle a client connection with their own command loop
func handleConnection(conn net.Conn, in chan int) {
	defer conn.Close()
	fmt.Fprintln(conn)
	// Log connection to server
	fmt.Printf("Client connected from %s", conn.RemoteAddr().String())
	fmt.Fprintln(conn, "Welcome to Tic-Tac-Go!")
}

func getLocalAddress() string {
	var localaddress string

	ifaces, err := net.Interfaces()
	if err != nil {
		panic("getLocalAddress: failed to find network interfaces")
	}

	// find the first non-loopback interface with an IPv4 address
	for _, elt := range ifaces {
		if elt.Flags&net.FlagLoopback == 0 && elt.Flags&net.FlagUp != 0 {
			addrs, err := elt.Addrs()
			if err != nil {
				panic("getLocalAddress: failed to get addresses for network interface")
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok {
					if ip4 := ipnet.IP.To4(); len(ip4) == net.IPv4len {
						localaddress = ip4.String()
						break
					}
				}
			}
		}
	}
	if localaddress == "" {
		panic("localaddress: failed to find non-loopback interface with valid IPv4 address")
	}

	return localaddress
}
