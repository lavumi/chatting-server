package api

import (
	"github.com/gin-gonic/gin"
)

func EnterChat(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	//
	//username, err := c.Cookie("user")
	//log.Printf("User Request: %v", username)
	//if err != nil {
	//	return
	//}
	//
	//v, ok := c.Get("clientChan")
	//if !ok {
	//	return
	//}
	//
	//c.Stream(func(w io.Writer) bool {
	//	if msg, ok := <-clientChan; ok {
	//		c.SSEvent("message", msg)
	//		return true
	//	}
	//	return false
	//})

}
