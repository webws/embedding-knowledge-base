package ai

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
	"github.com/webws/embedding-knowledge-base/util"
	"github.com/webws/go-moda/logger"
)

type AIClient struct {
	openai *openai.Client
	proxy  string
	apiKey string
}

func NewAiClient(proxy, apiKey string, needCheckProxy bool) (*AIClient, error) {
	c := &AIClient{proxy: proxy, apiKey: apiKey}
	config := openai.DefaultConfig(c.apiKey)
	config.HTTPClient, _ = newHttpClientWithProxy(c.proxy)
	pass := util.ProxyCheck(config.HTTPClient, c.proxy, "https://www.google.com", 5)
	if needCheckProxy && !pass {
		config.HTTPClient, _ = newHttpClientWithProxy("")
	}
	c.openai = openai.NewClientWithConfig(config)
	return c, nil
}

func newHttpClientWithProxy(proxy string) (*http.Client, error) {
	client := &http.Client{}
	if proxy != "" {
		uri, err := url.Parse(proxy)
		if err != nil {
			log.Fatalln(err)
			logger.Errorw("NewAiClient", "err", err)
			return nil, err

		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	return client, nil
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
	if len(rsp.Data) > 0 {
		return rsp.Data[0].Embedding, nil
	}
	return nil, errors.New("embedding length is 0")
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
