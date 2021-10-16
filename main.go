package main

import (
	"fmt"
	"github.com/gruberchris/gochatserver/chatroom"
	"net"
	"os"
)

func main() {
	const (
		host = "localhost"
		port = "5000"
	)

	chatRoom := chatroom.NewChatRoom()

	listener, err := net.Listen("tcp", host + ":" + port)

	if err != nil {
		fmt.Println("Error while listening for connections: ", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("Error while closing listening for connections: ", err.Error())
		}
	}()

	fmt.Println("Listening on " + host + ":" + port + "...")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting a TCP connection: ", err.Error())
			continue
		}

		fmt.Println("Accepted connection from remote client: " + conn.RemoteAddr().String())

		chatRoom.Joins <- conn
	}
}
