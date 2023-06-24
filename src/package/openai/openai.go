package openaiServ

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type OpenAi struct {
	Logger    *zap.Logger
	Token     string
	HttpAgent string
}

func NewService(token, agent string) *OpenAi {
	logger, _ := zap.NewDevelopment()
	return &OpenAi{
		Logger:    logger,
		Token:     token,
		HttpAgent: agent,
	}
}

func (o *OpenAi) NewClient() (*openai.Client, error) {
	config := openai.DefaultConfig(o.Token)
	if o.HttpAgent != "" {
		proxyUrl, err := url.Parse(o.HttpAgent)
		if err != nil {
			return nil, err
		}
		config.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
	}
	return openai.NewClientWithConfig(config), nil
}

func (o *OpenAi) Moderation(text string) (r *openai.ModerationResponse, isOk bool, e error) {
	ai, err := o.NewClient()
	if err != nil {
		return nil, false, err
	}
	resp, err := ai.Moderations(context.Background(), openai.ModerationRequest{
		Input: text,
		Model: "text-moderation-latest",
	})
	for _, v := range resp.Results {
		if v.Flagged {
			return &resp, false, nil
		}
	}
	return &resp, true, nil
}
