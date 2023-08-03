package main

import (
	"game-server/api"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := api.InitRouter()

	r.Static("/static", "./web/static")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")
	r.LoadHTMLFiles("web/static/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Posts",
		})
	})

	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panic("set trusted proxies fail")
		return
	}
	err = r.Run(":8080")
	if err != nil {
		log.Panic("error on running")
		return
	}
}
