package api

import (
	"game-server/api/handler"
	"game-server/internal/chat"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	chatService := chat.Init()

	root := r.Group("/api")
	{
		chatRouter := root.Group("/chat")
		{
			chatRouter.GET("/register", handler.EnterChat(chatService))
			chatRouter.POST("/message", handler.SendMessage(chatService))
			chatRouter.GET("/users", handler.GetUserList(chatService))
		}
	}
	return r
}
