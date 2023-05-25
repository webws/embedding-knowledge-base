package myai

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	SocksProxy = "socks5://127.0.0.1:1080"
	AIKey      = "your open api key"
)

// 设置自己的科学代理地址
func myProxyTransport() *http.Transport {
	// SocksProxy := "socks5://127.0.0.1:1080"

	uri, err := url.Parse(SocksProxy)
	if err != nil {
		log.Fatalln(err)
	}
	return &http.Transport{
		Proxy: http.ProxyURL(uri),
	}
}

func NewOpenAiClient() *openai.Client {
	token := os.Getenv("AIKey")
	if token == "" {
		token = AIKey
	}
	config := openai.DefaultConfig(token)
	config.HTTPClient.Transport = myProxyTransport()
	return openai.NewClientWithConfig(config)
}

// 把 搜索词变成向量
func SimpleGetVec(prompt string) ([]float32, error) {
	c := NewOpenAiClient()
	req := openai.EmbeddingRequest{
		Input: []string{prompt},
		Model: openai.AdaEmbeddingV2,
	}
	rsp, err := c.CreateEmbeddings(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	// 打印日志,向量总数量,不打印向量具体内容,其他都打印
	// log.Printf("len(rsp.Data): %d", len(rsp.Data))
	// 取向量的平均值
	var v []float32
	for _, d := range rsp.Data {
		v = append(v, d.Embedding...)
	}
	// 取平均值
	for i := 0; i < len(v); i++ {
		v[i] = v[i] / float32(len(rsp.Data))
	}
	// 打印向量
	// log.Printf("v: %v", v)
	return v, nil
	// return rsp.Data[0].Embedding, nil
}
