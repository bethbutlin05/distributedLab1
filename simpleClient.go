package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	msg, _ := reader.ReadString('\n')
	fmt.Println(msg)
}

func main() {
	stdin := bufio.NewReader(os.Stdin)
	//connect to the server on your own computer
	conn, _ := net.Dial("tcp", "127.0.0.1:8030")

	for {
		fmt.Printf("Enter text: ")
		msg, _ := stdin.ReadString('\n')

		//send message followed by a newline through that connection
		fmt.Fprintln(conn, msg)

		read(conn)
	}
}
