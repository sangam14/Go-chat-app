package main

import (
	"fmt"
	"net"
	"time"
)

type Chat struct {
	clients  []*Client
	connect  chan net.Conn
	outgoing chan string
}


func CreateChat() *Chat {
	chat := &Chat{
		clients:  make([]*Client, 0),
		connect:  make(chan net.Conn),
		outgoing: make(chan string),
	}

	chat.Listen()

	return chat
}


func (chat *Chat) Listen() {
	go func() {
		for {
			select {
			case conn := <-chat.connect:
				chat.Join(conn)
			case msg := <-chat.outgoing:
				chat.Broadcast(msg)
			}
		}
	}()
}

func (chat *Chat) Connect(conn net.Conn) {
	chat.connect <- conn
}

func (chat *Chat) Join(conn net.Conn) {
	client := CreateClient(conn)
	chat.clients = append(chat.clients, client)
	go func() {
		for {
			chat.outgoing <- <-client.incoming
		}
	}()
}

func (chat *Chat) Remove(i int) {
	chat.clients = append(chat.clients[:i], chat.clients[i+1:]...)
}


func (chat *Chat) UpdateClientsList() {
	connectedClients := "/clients>"
	for _, client := range chat.clients {
		connectedClients += client.name + " "
	}
	connectedClients += "\n"
	for _, client := range chat.clients {
		client.outgoing <- connectedClients
	}
}


func (chat *Chat) Broadcast(data string) {
	currentTime := time.Now().Format("15:04:05")
	msg := fmt.Sprintf("[%s] %s", currentTime, data)
	for i, client := range chat.clients {
		if client.status == 0 {
			chat.Remove(i)
		}
	}
	chat.UpdateClientsList()
	for _, client := range chat.clients {
		client.outgoing <- msg
	}
}
