package handler

import (
	"game-server/internal/chat"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

type ReqRoomUri struct {
	RoomId string `uri:"roomId" binding:"required"`
}

func EnterChat(room *chat.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		userId := c.GetHeader("UserName")

		var roomUri ReqRoomUri
		if err := c.ShouldBindUri(&roomUri); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		client := room.JoinRoom(roomUri.RoomId, userId)
		defer room.ExitRoom(roomUri.RoomId, userId, client)
		clientDone := c.Request.Context().Done()
		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientDone:
				return false
			case message := <-client:
				c.SSEvent("msg", message)
				return true
			}
		})
	}
}

func SendMessage(room *chat.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		var roomUri ReqRoomUri
		if err := c.ShouldBindUri(&roomUri); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		userId := c.GetHeader("UserName")
		data, err := c.GetRawData()
		if err != nil {
			log.Printf("bind json Error: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		room.SendMessage(roomUri.RoomId, userId, string(data))
		c.String(http.StatusOK, "")
	}
}

func GetRoomList(room *chat.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"roomList": room.GetRoomList(),
		})
	}
}

func GetRoomInfo(room *chat.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		var roomUri ReqRoomUri
		if err := c.ShouldBindUri(&roomUri); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"roomInfo": room.GetRoomInfo(roomUri.RoomId),
		})
	}
}
