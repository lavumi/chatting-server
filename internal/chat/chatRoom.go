package chat

import (
	"log"
)

type IChatRoom interface {
	joinRoom(chan string)
	exitRoom(chan string)
	sendMessage(string)
	closeRoom() error
	//GetMembers() []string
}

type chatRoom struct {
	message chan string
	open    chan chan string
	close   chan chan string
	streams map[chan string]bool
}

//var chatRoom *chatRoom

func makeRoom() IChatRoom {
	room := &chatRoom{
		message: make(chan string),
		open:    make(chan chan string),
		close:   make(chan chan string),
		streams: make(map[chan string]bool),
	}

	go room.listen()

	return room
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (room *chatRoom) listen() {
	for {
		select {
		// Add new available client
		case client := <-room.open:
			room.streams[client] = true
			log.Printf("Client added. %d registered clients", len(room.streams))

		// Remove closed client
		case client := <-room.close:
			delete(room.streams, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(room.streams))

		// Broadcast message to client
		case eventMsg := <-room.message:
			log.Printf("BroadCast Message to client. %s", eventMsg)
			for clientMessageChan := range room.streams {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (room *chatRoom) joinRoom(channel chan string) {
	room.open <- channel
}

func (room *chatRoom) exitRoom(clientChan chan string) {
	log.Printf("close room registered")
	room.close <- clientChan
}

func (room *chatRoom) sendMessage(msg string) {
	room.message <- msg
}

func (room *chatRoom) closeRoom() error {
	close(room.open)
	close(room.close)
	return nil
}
