package service

import (
	"log"
)

type Message struct {
	RoomId string `json:"room_id"`
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

type ChatRoom struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

type ClientChan chan string

var chatRoom *ChatRoom

func init() {
	chatRoom = &ChatRoom{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go listen()
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func listen() {
	for {
		select {
		// Add new available client
		case client := <-chatRoom.NewClients:
			chatRoom.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(chatRoom.TotalClients))

		// Remove closed client
		case client := <-chatRoom.ClosedClients:
			delete(chatRoom.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(chatRoom.TotalClients))

		// Broadcast message to client
		case eventMsg := <-chatRoom.Message:
			log.Printf("BroadCast Message to client. %s", eventMsg)
			for clientMessageChan := range chatRoom.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func JoinRoom() chan string {
	clientChan := make(ClientChan)
	chatRoom.NewClients <- clientChan

	return clientChan
}

func ExitRoom(clientChan chan string) {
	chatRoom.ClosedClients <- clientChan
}

func SendMessage(msg string) {
	chatRoom.Message <- msg
}
