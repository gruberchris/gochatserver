package main

import (
	"fmt"
	"github.com/gruberchris/gochatserver/chatroom"
	"log"
	"net"
)

func main() {
	const port = "5000"

	chatRoom := chatroom.NewChatRoom()

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal("Error while listening for connections:", err.Error())
	}

	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("Error while closing listening for connections: ", err.Error())
		}
	}()

	fmt.Println("Listening on", ":"+port, "...")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting a TCP connection:", err.Error())
			continue
		}

		fmt.Println("Accepted connection from remote client:", conn.RemoteAddr().String())

		chatRoom.Joins <- conn
	}
}
