package main

import (
	"game-server/api"
	"github.com/gin-contrib/static"
	"log"
)

func main() {
	r := api.InitRouter()

	r.Use(static.Serve("/", static.LocalFile("./web", false)))

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
