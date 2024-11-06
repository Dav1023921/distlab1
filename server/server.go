package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	if err != nil {
		// print error
		fmt.Println("error")
	}
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		handleError(err)
		conns <- conn
	}
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		r := new(Message)
		r.message = "Message from: " + strconv.Itoa(clientid) + msg
		r.sender = clientid
		msgs <- *r
	}
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	ln, _ := net.Listen("tcp", *portPtr)
	//TODO Create a Listener for TCP connections on the port given above.

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	clid := 0
	for {
		select {
		case conn := <-conns:
			clients[clid] = conn
			go handleClient(conn, clid, msgs)
			clid += 1
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client

		case msg := <-msgs:
			for id, connection := range clients {
				if id != msg.sender {
					fmt.Fprint(connection, msg.message)
				}
			}
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		}
	}
}
