package client

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	Incoming       chan string
	Outgoing       chan string
	reader         *bufio.Reader
	writer         *bufio.Writer
	Conn           net.Conn
	onDisconnected func(client *Client)
}

func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\n')

		if err != nil {
			if err := client.Conn.Close(); err != nil {
				log.Println("Error while closing connection for client", client.Conn.RemoteAddr().String(), ":", err.Error())
			}

			client.onDisconnected(client)
			break
		}

		client.Incoming <- line
	}
}

func (client *Client) Write() {
	for data := range client.Outgoing {
		_, err := client.writer.WriteString(data)

		if err != nil {
			log.Println("Error while writing:", err.Error())
		}

		if err := client.writer.Flush(); err != nil {
			log.Println("Error while flushing:", err.Error())
		}
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn, onClientDisconnected func(client *Client)) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		Incoming:       make(chan string),
		Outgoing:       make(chan string),
		reader:         reader,
		writer:         writer,
		Conn:           connection,
		onDisconnected: onClientDisconnected,
	}

	// Starts a couple threads to keep updating Client state
	client.Listen()

	return client
}
