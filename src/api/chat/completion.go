package apiChat

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net/http"
	"qiming-server/src/config"
	openaiServ "qiming-server/src/package/openai"
	strBase "qiming-server/src/struct"
)

var Clients = make(map[int]chan string)

type CompletionRequest struct {
	Text    string `json:"text"`
	Context []struct {
		Role int    `form:"role"`
		Text string `form:"text"`
	}
}

// Completion 接收聊天消息
func Completion(c *gin.Context) {
	var form CompletionRequest
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("BindParams Failed: %s\n", err)
		return
	}
	if form.Text == "" {
		return
	}
	userId, _ := c.Get("user_id")
	uid := userId.(int)

	conf := config.GetConfig()

	aiServ := openaiServ.NewService(conf.ChatGPT.Keys, conf.ChatGPT.HttpProxy)
	ai, err := aiServ.NewClient()
	if err != nil {
		log.Printf("openai server Failed: %s\n", err)
		return
	}

	// 检测是否包含违规词
	_, isOk, err := aiServ.Moderation(form.Text)
	if err != nil {
		log.Printf("Moderation Failed: %s\n", err)
		return
	}
	if !isOk {
		c.JSON(http.StatusOK, strBase.ApiResult{
			Code: 1001,
			Msg:  "输入内容包含违规词",
		})
		return
	}

	var chatList []openai.ChatCompletionMessage
	for _, v := range form.Context {
		chatList = append(chatList, openai.ChatCompletionMessage{
			Role:    parseRole(v.Role),
			Content: v.Text,
		})
	}
	go func() {
		resp, err := ai.CreateChatCompletionStream(c, openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo16K, //gpt-3.5-turbo-16k 这个模型的请求频率和上下文要高一些
			Messages:  chatList,
			MaxTokens: 2500,
			Stream:    true,
			User:      fmt.Sprintf("%d", userId),
		})
		if err != nil {
			log.Printf("CreateChatCompletion failed: %s", err)
			Clients[uid] <- "<!limit!>"
			return
		}
		defer resp.Close()

		for {
			r, err := resp.Recv()
			if errors.Is(err, io.EOF) {
				Clients[uid] <- "<!end!>"
				return
			}
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return
			}
			Clients[uid] <- r.Choices[0].Delta.Content
		}
	}()

	// TODO
	// 检查用户状态、余额
	// 调用openai检测关键词是否违规
	// 从key池中取一个key
	// 调用...
	// 扣费

	c.JSON(http.StatusOK, strBase.ApiResult{
		Code: http.StatusOK,
	})
}

func parseRole(r int) string {
	switch r {
	case 1:
		return openai.ChatMessageRoleSystem
	case 2:
		return openai.ChatMessageRoleUser
	case 3:
		return openai.ChatMessageRoleAssistant
	}
	return openai.ChatMessageRoleUser
}
