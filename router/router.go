package router

import (
	"game-server/service"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//r.NoRoute(func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "/")
	//})

	chatService := service.ChatService()

	root := r.Group("/api")
	{
		chat := root.Group("/chat")
		{
			chat.GET("/register", chatService.JoinUser)
			chat.POST("/message", chatService.BroadcastMessage)
		}
	}
	return r
}
