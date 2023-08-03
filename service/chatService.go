package service

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

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

func ChatService() (event *ChatRoom) {
	event = &ChatRoom{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *ChatRoom) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			log.Printf("BroadCast Message to client. %s", eventMsg)
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *ChatRoom) JoinUserNew() chan string {
	clientChan := make(ClientChan)
	stream.NewClients <- clientChan

	return clientChan
}

func (stream *ChatRoom) JoinUser(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Initialize client channel
	clientChan := make(ClientChan)

	//clientChan <- c.Param("user")
	//userName :=
	// Send new connection to event server
	stream.NewClients <- clientChan

	defer func() {
		// Send closed connection to event server
		stream.ClosedClients <- clientChan
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent("info", msg)
			c.Writer.Flush()
			return true
		}
		return false
	})
}

type Message struct {
	Msg string `json:"msg"`
}

func (stream *ChatRoom) BroadcastMessage(c *gin.Context) {
	var msg Message
	err := c.BindJSON(&msg)
	if err != nil {
		log.Printf("bind json Error: %s", err.Error())
		return
	}

	log.Printf("%s, %s\n", "testuser", msg.Msg)
	stream.Message <- msg.Msg

}
