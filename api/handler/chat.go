package handler

import (
	"fmt"
	"game-server/internal/chat"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type Message struct {
	RoomId string `json:"room_id"`
	Sender string `json:"sender"`
	Msg    string `json:"msg"`
}

func EnterChat(room chat.IChatRoom) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		client := room.JoinRoom()
		defer room.ExitRoom(client)

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
}

func SendMessage(room chat.IChatRoom) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg Message
		err := c.BindJSON(&msg)
		if err != nil {
			log.Printf("bind json Error: %s", err.Error())
			return
		}

		room.SendMessage(fmt.Sprintf("%s : %s", msg.Sender, msg.Msg))
		c.String(http.StatusOK, "")
	}
}

func GetUserList(room chat.IChatRoom) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userList": []string{"a", "b", "c"},
		})
	}
}
