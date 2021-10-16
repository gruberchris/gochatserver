package chatroom

import (
	"fmt"
	"github.com/gruberchris/gochatserver/chatclient"
	"net"
	"strconv"
)

type ChatRoom struct {
	clients  []*chatclient.ChatClient
	Joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (chatRoom *ChatRoom) Broadcast(data string) {
	connectionsCount := len(chatRoom.clients)

	if connectionsCount > 0 {
		fmt.Println("Broadcasting to", strconv.Itoa(connectionsCount), "clients:", data)

		for _, c := range chatRoom.clients {
			c.Outgoing <- data
		}
	}
}

func (chatRoom *ChatRoom) Join(connection net.Conn) {
	c := chatclient.NewClient(connection, chatRoom.Remove)
	chatRoom.clients = append(chatRoom.clients, c)

	go func() {
		for {
			chatRoom.incoming <- <-c.Incoming
		}
	}()
}

func (chatRoom *ChatRoom) Listen() {
	go func() {
		for {
			select {
			case data := <-chatRoom.incoming:
				chatRoom.Broadcast(data)
			case conn := <-chatRoom.Joins:
				chatRoom.Join(conn)
			}
		}
	}()
}

func (chatRoom *ChatRoom) Remove(disconnectedClient *chatclient.ChatClient) {
	newClients := make([]*chatclient.ChatClient, len(chatRoom.clients))
	index := 0

	for _, c := range chatRoom.clients {
		if c != disconnectedClient {
			newClients[index] = c
			index++
		}
	}

	// remove the chatclient from the chat room
	chatRoom.clients = newClients[:index]

	// broadcast the chatclient left to any remaining active clients
	disconnectedMessage := fmt.Sprintf("%s left", disconnectedClient.Conn.RemoteAddr().String())
	chatRoom.Broadcast(disconnectedMessage)

	fmt.Println(len(chatRoom.clients), "clients connected.")
}

func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{
		clients:  make([]*chatclient.ChatClient, 0),
		Joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	chatRoom.Listen()

	return chatRoom
}
