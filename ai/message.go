package myai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var MessageStore ChatMessages

func init() {
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear() // 清理和初始化
}

func K8sChat(prompt string) string {
	c := NewOpenAiClient()
	MessageStore.AddForUser(prompt)
	rsp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: MessageStore.ToMessage(),
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	MessageStore.AddForAssistant(rsp.Choices[0].Message.Content)

	return MessageStore.GetLast()
}

type (
	ChatMessages []*ChatMessage
	ChatMessage  struct {
		Msg openai.ChatCompletionMessage
	}
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

func (cm *ChatMessages) Clear() {
	*cm = make([]*ChatMessage, 0) // 重新初始化

	// cm.AddForSystem("You are a helpful K8S assistant. Use the provided text to form your answer. Keep your answer within 10 sentences. Accurate, helpful, concise and to the point")
	cm.AddForSystem("You are a helpful K8S assistant. I will provide you with text, including the question title or keywords, question description, and reference answer. If I provide a reference answer, please try to use it as much as possible.Try to answer in Chinese as much as possible, Keep your answer within 10 sentences. Accurate, useful, concise and to the point")
}

func (cm *ChatMessages) AddFor(msg string, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
}

const CommandPattern = "```\\s*(.*?)\\s*```"

func (cm *ChatMessages) Dump(file string) error {
	// 把内容json化后后 ，存到文件里
	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (cm *ChatMessages) AddForAssistant(msg string) {
	cm.AddFor(msg, RoleAssistant)
}

func (cm *ChatMessages) AddForSystem(msg string) {
	cm.AddFor(msg, RoleSystem)
}

func (cm *ChatMessages) AddForUser(msg string) {
	cm.AddFor(msg, RoleUser)
}

func (cm *ChatMessages) Apply(content string) error {
	content = strings.Trim(content, " ")
	if len(*cm) == 0 || content == "" {
		return fmt.Errorf("无需apply")
	}
	msg := (*cm)[len(*cm)-1]
	if msg.Msg.Role != RoleAssistant {
		fmt.Println("user/system内容无需apply")
		return fmt.Errorf("user/system内容无需apply")
	}
	msg.Msg.Content = content
	return nil
}

func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for index, c := range *cm {
		ret[index] = c.Msg
	}
	return ret
}

func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "什么都没找到"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}
