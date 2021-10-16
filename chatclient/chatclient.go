package chatclient

import (
	"bufio"
	"log"
	"net"
)

type ChatClient struct {
	Incoming       chan string
	Outgoing       chan string
	reader         *bufio.Reader
	writer         *bufio.Writer
	Conn           net.Conn
	onDisconnected func(client *ChatClient)
}

func (chatClient *ChatClient) Read() {
	for {
		line, err := chatClient.reader.ReadString('\n')

		if err != nil {
			if err := chatClient.Conn.Close(); err != nil {
				log.Println("Error while closing connection for chatclient", chatClient.Conn.RemoteAddr().String(), ":", err.Error())
			}

			chatClient.onDisconnected(chatClient)
			break
		}

		chatClient.Incoming <- line
	}
}

func (chatClient *ChatClient) Write() {
	for data := range chatClient.Outgoing {
		_, err := chatClient.writer.WriteString(data)

		if err != nil {
			log.Println("Error while writing:", err.Error())
		}

		if err := chatClient.writer.Flush(); err != nil {
			log.Println("Error while flushing:", err.Error())
		}
	}
}

func (chatClient *ChatClient) Listen() {
	go chatClient.Read()
	go chatClient.Write()
}

func NewClient(connection net.Conn, onClientDisconnected func(client *ChatClient)) *ChatClient {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	chatClient := &ChatClient{
		Incoming:       make(chan string),
		Outgoing:       make(chan string),
		reader:         reader,
		writer:         writer,
		Conn:           connection,
		onDisconnected: onClientDisconnected,
	}

	// Starts a couple threads to keep updating ChatClient state
	chatClient.Listen()

	return chatClient
}
