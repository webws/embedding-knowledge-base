package ai

import (
	openai "github.com/sashabaranov/go-openai"
)

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
	*cm = make([]*ChatMessage, 0)
}

func InitChatMessages() ChatMessages {
	cm := make(ChatMessages, 0)
	// cm.AddForSystem("You are a helpful K8S assistant. Use the provided text to form your answer. Keep your answer within 10 sentences. Accurate, helpful, concise and to the point")
	// cm.AddForSystem("You are a helpful K8S assistant. I will provide you with text, including the question title or keywords, question description, and reference answer. If I provide a reference answer, please try to use it as much as possible.Try to answer in Chinese as much as possible, Keep your answer within 10 sentences. Accurate, useful, concise and to the point")
	return cm
}

func (cm *ChatMessages) AddFor(msg string, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
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

func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for index, c := range *cm {
		ret[index] = c.Msg
	}
	return ret
}

func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "nothing"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}
