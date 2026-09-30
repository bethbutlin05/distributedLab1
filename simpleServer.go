package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		fmt.Println(msg)
		fmt.Fprintln(conn, "hi bethany")
	}
}

func main() {
	//start listening for connections on port 8030
	ln, _ := net.Listen("tcp", ":8030")
	for {
		//wait for client to connect, return its connection
		conn, _ := ln.Accept()
		//goroutine makes it handle many connections at a time. as soon as we call it we go back round the for loop again, the function doesn't need to complete first.
		go handleConnection(conn)
	}
}
