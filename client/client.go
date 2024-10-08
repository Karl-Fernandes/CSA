package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
			break
		}
		// Print the incoming message
		fmt.Print(msg) // Print the message directly without an extra newline
	}
}

func write(conn net.Conn) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter Text: ") // Prompt for input without newline
		msg, err := stdin.ReadString('\n')
		handleError(err) // Handle potential errors

		msg = strings.TrimSpace(msg)
		if msg == "exit" {
			fmt.Println("Exiting...")
			return
		}
		fmt.Fprintln(conn, msg) // Send the message to the server
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	conn, err := net.Dial("tcp", *addrPtr)
	handleError(err) // Handle potential connection errors

	go read(conn) // Start reading messages in a goroutine
	write(conn)   // Start writing messages

	// The program will exit when write() returns
}
