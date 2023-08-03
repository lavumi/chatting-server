package chat

import (
	"log"
)

type Message struct {
	RoomId string `json:"room_id"`
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

type IChatRoom interface {
	JoinRoom() chan string
	ExitRoom(chan string)
	SendMessage(string)
}

type Room struct {
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

//var room *Room

func Init() (room *Room) {
	room = &Room{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go room.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (room *Room) listen() {
	for {
		select {
		// Add new available client
		case client := <-room.NewClients:
			room.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(room.TotalClients))

		// Remove closed client
		case client := <-room.ClosedClients:
			delete(room.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(room.TotalClients))

		// Broadcast message to client
		case eventMsg := <-room.Message:
			log.Printf("BroadCast Message to client. %s", eventMsg)
			for clientMessageChan := range room.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (room *Room) JoinRoom() chan string {
	clientChan := make(ClientChan)
	room.NewClients <- clientChan

	return clientChan
}

func (room *Room) ExitRoom(clientChan chan string) {
	room.ClosedClients <- clientChan
}

func (room *Room) SendMessage(msg string) {
	room.Message <- msg
}
