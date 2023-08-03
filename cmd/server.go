package main

import (
	"game-server/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := router.InitRouter()

	//r := gin.Default()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("webpage/*")
	r.GET("/", func(c *gin.Context) {
		//if pusher := c.Writer.Pusher(); pusher != nil {
		//	if err := pusher.Push("/assets/script.js", nil); err != nil {
		//		log.Printf("Failed to push: %v", err)
		//	}
		//	if err := pusher.Push("/assets/simple.css", nil); err != nil {
		//		log.Printf("Failed to push: %v", err)
		//	}
		//} else {
		//	log.Printf("Pusher is nil!!!!")
		//}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Posts",
		})
	})
	err := r.Run(":8080")
	//err := r.RunTLS(":8080", "./cert/localhost.crt", "./cert/localhost.key")
	if err != nil {
		return
	}
}
