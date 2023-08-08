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

	chatRoom := chat.Init()
	root := r.Group("/api")
	{
		chatRouter := root.Group("/chat")
		{
			chatRouter.GET("/register", handler.EnterChat(chatRoom))
			chatRouter.POST("/message", handler.SendMessage(chatRoom))
			chatRouter.GET("/users", handler.GetUserList(chatRoom))
		}
	}
	return r
}
