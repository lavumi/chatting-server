package handler

import (
	"game-server/internal/chat"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type RoomInfo struct {
	RoomId string `uri:"roomId" binding:"required"`
}

func EnterChat(room *chat.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		userId := c.GetHeader("UserName")

		var roomInfo RoomInfo
		if err := c.ShouldBindUri(&roomInfo); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		client := room.JoinRoom(roomInfo.RoomId, userId)
		defer room.ExitRoom(roomInfo.RoomId, userId, client)
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

func SendMessage(room *chat.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		var roomInfo RoomInfo
		if err := c.ShouldBindUri(&roomInfo); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		data, err := c.GetRawData()
		if err != nil {
			log.Printf("bind json Error: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room.SendMessage(roomInfo.RoomId, string(data))
		c.String(http.StatusOK, "")
	}
}

func GetRoomList(room *chat.Service) gin.HandlerFunc {

	//roomList := room.GetRoomList()
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"roomList": []string{"123", "chat", "test"},
		})
	}
}

func GetUserList(room *chat.Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		var roomInfo RoomInfo
		if err := c.ShouldBindUri(&roomInfo); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"userList": []string{"aaa", "bbb", "ccc"},
		})
	}
}
