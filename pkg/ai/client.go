package ai

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/webws/go-moda/logger"
)

type AIClient struct {
	openai *openai.Client
	proxy  string
	apiKey string
}

func NewAiClient(proxy, apiKey string) (*AIClient, error) {
	if apiKey == "" {
		apiKey = os.Getenv("AIKey")
	}
	config := openai.DefaultConfig(apiKey)
	if proxy != "" {
		uri, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
			logger.Errorw("NewAiClient", "err", err)
			return nil, err

		}
		config.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	return &AIClient{
		openai: openai.NewClientWithConfig(config),
		proxy:  proxy,
		apiKey: apiKey,
	}, nil
}

// SimpleGetVec  Get vector from prompt
func (c *AIClient) SimpleGetVec(prompt string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{prompt},
		Model: openai.AdaEmbeddingV2,
	}
	rsp, err := c.openai.CreateEmbeddings(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	var v []float32
	for _, d := range rsp.Data {
		v = append(v, d.Embedding...)
	}
	for i := 0; i < len(v); i++ {
		v[i] = v[i] / float32(len(rsp.Data))
	}
	return v, nil
}

func (c *AIClient) Chat(prompt string) (string, error) {
	messageStore := InitChatMessages()
	messageStore.AddForUser(prompt)
	defer messageStore.Clear()
	rsp, err := c.openai.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messageStore.ToMessage(),
	})
	if err != nil {
		return "", err
	}
	messageStore.AddForAssistant(rsp.Choices[0].Message.Content)
	return messageStore.GetLast(), nil
}
