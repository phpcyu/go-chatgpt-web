package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	apiChat "qiming-server/src/api/chat"
	"qiming-server/src/config"
	"qiming-server/src/util"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	conf := config.GetConfig()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Cookie", "X-CSRF-TOKEN", "Accept", "Authorization", "X-XSRF-TOKEN", "Access-Control-Allow-Origin", "X-Token", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "authenticated"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/user/login/test", func(c *gin.Context) {
		util.ApiSuc(c, map[string]interface{}{
			"id": time.Now().Unix(),
		})
	})
	// 接收用户发送的消息
	r.POST("/ai/chat/completion", JWTAuthMiddleware(), apiChat.Completion)

	// 保持用户端的长链接
	r.GET("/ai/chat/stream", HeadersMiddleware(), JWTAuthMiddleware(), apiChat.Stream)

	err := r.Run(fmt.Sprintf("%s", conf.App.Addr))
	log.Fatalf("http_server start fail: %s", err)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.GetHeader("Token"))
		// 检测合法性等等...
		c.Set("user_id", id)
		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
