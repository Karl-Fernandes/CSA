package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		conns <- conn
	}
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
			client.Close()
			break
		}
		msg = strings.TrimSpace(msg)
		msgs <- Message{clientid, msg}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError(err)
		return
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
			clientID := len(clients) + 1
			clients[clientID] = conn
			go handleClient(clients[clientID], clientID, msgs)
		case msg := <-msgs:
			fmt.Printf("Received message from client %d: %s\n", msg.sender, msg.message) // Print the sender's ID and the message
			formattedMessage := fmt.Sprintf("Client %d: %s", msg.sender, msg.message)
			for id, client := range clients {
				if id != msg.sender { // Check if the client ID is not the sender
					_, err := fmt.Fprintln(client, formattedMessage) // Send the message to the client
					if err != nil {
						fmt.Println("Error sending message to client:", err)
						client.Close()
					}
				}
			}

		}
	}
}
