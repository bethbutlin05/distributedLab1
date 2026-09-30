package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// Tell main if we stop accepting connections.
	defer close(conns)

	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			return
		}

		// Pass the new connection to main.
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	defer client.Close()

	reader := bufio.NewReader(client)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// EOF means the client has disconnected normally.
			if err != io.EOF {
				handleError(err)
			}

			fmt.Println("Client disconnected:", clientid)
			return
		}

		// Remove the line ending, keeping the message text.
		msg = strings.TrimRight(msg, "\r\n")

		msgs <- Message{
			sender:  clientid,
			message: msg,
		}
	}
}

func main() {
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError(err)
		return
	}
	defer ln.Close()

	conns := make(chan net.Conn)
	msgs := make(chan Message)
	clients := make(map[int]net.Conn)

	// Close remaining client connections when main finishes.
	defer func() {
		for _, conn := range clients {
			conn.Close()
		}
	}()

	go acceptConns(ln, conns)

	nextID := 0
	fmt.Println("Server listening on", *portPtr)

	for {
		select {
		case conn, ok := <-conns:
			if !ok {
				return // The connection-accepting function stopped.
			}

			id := nextID
			nextID++

			clients[id] = conn
			fmt.Println("Client connected:", id)

			go handleClient(conn, id, msgs)

		case msg := <-msgs:
			for id, conn := range clients {
				// Do not send the message back to its sender.
				if id == msg.sender {
					continue
				}

				// Add one newline so the receiving client
				// can identify the end of the message.
				_, err := fmt.Fprintln(conn, msg.message)
				if err != nil {
					handleError(err)
					conn.Close()
					delete(clients, id)
				}
			}
		}
	}
}
