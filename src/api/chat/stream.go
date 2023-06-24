package apiChat

import (
	"github.com/gin-gonic/gin"
	"io"
)

func Stream(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		return
	}
	uid := userId.(int)
	Clients[uid] = make(chan string)
	defer func() {
		delete(Clients, uid)
		close(Clients[uid])
	}()
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-Clients[uid]; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
