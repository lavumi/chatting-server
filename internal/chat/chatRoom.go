package chat

import (
	"log"
)

type IChatRoom interface {
	JoinRoom() chan string
	ExitRoom(chan string)
	SendMessage(string)
	//GetMembers() []string
}

type Room struct {
	message      chan string
	newClient    chan chan string
	closedClient chan chan string
	totalClient  map[chan string]bool
}

type ClientChan chan string

//var room *Room

func Init() (room *Room) {
	room = &Room{
		message:      make(chan string),
		newClient:    make(chan chan string),
		closedClient: make(chan chan string),
		totalClient:  make(map[chan string]bool),
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
		case client := <-room.newClient:
			room.totalClient[client] = true
			log.Printf("Client added. %d registered clients", len(room.totalClient))

		// Remove closed client
		case client := <-room.closedClient:
			delete(room.totalClient, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(room.totalClient))

		// Broadcast message to client
		case eventMsg := <-room.message:
			log.Printf("BroadCast Message to client. %s", eventMsg)
			for clientMessageChan := range room.totalClient {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (room *Room) JoinRoom() chan string {
	clientChan := make(ClientChan)
	room.newClient <- clientChan
	return clientChan
}

func (room *Room) ExitRoom(clientChan chan string) {

	room.closedClient <- clientChan
}

func (room *Room) SendMessage(msg string) {
	room.message <- msg
}

func (room *Room) Close() error {
	close(room.newClient)
	close(room.closedClient)
	return nil
}
