package api

import (
	"fmt"
	chatService "game-server/service"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func EnterChat(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	client := chatService.JoinRoom()
	defer chatService.ExitRoom(client)

	clientDone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientDone:
			return false
		case message := <-client:
			c.SSEvent("info", message)
			return true
		}
	})
}

type Message struct {
	RoomId string `json:"room_id"`
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

func SendMessage(c *gin.Context) {
	var msg Message
	err := c.BindJSON(&msg)
	if err != nil {
		log.Printf("bind json Error: %s", err.Error())
		return
	}

	chatService.SendMessage(fmt.Sprintf("%s : %s", msg.Sender, msg.Msg))
}
