// Package gpt
// @Description
// @Author root_wang
// @Date 2023/3/15 15:03
package gpt

import "cqhttp-client/src/constant"

type (
	ChatMessage struct {
		Role    cst.Role `json:"role"`
		Content string   `json:"content"`
	}

	ChatReq struct {
		Model    string         `json:"model"`
		Messages []*ChatMessage `json:"messages"`
		// default 1
		Temperature int `json:"temperature"`
		// default 0
		TopP int `json:"top_p"`
		// default 1 how many resp of msg
		N int `json:"n"`
		// max resp token
		MaxTokens        int `json:"max_tokens"`
		PresencePenalty  int `json:"presence_penalty"`
		FrequencyPenalty int `json:"frequency_penalty"`
	}

	chatChoices struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
	}

	ChatResp struct {
		Id      string        `json:"id"`
		Object  string        `json:"object"`
		Created int64         `json:"created"`
		Model   string        `json:"model"`
		Usage   usage         `json:"usage"`
		Choices []chatChoices `json:"choices"`
	}
)

func NewChatMessage(msg string, role cst.Role) *ChatMessage {
	return &ChatMessage{
		Role:    role,
		Content: msg,
	}
}

const (
	ChatCompletion SaasName = "chatcompletion"
	ChatBseURL              = "https://api.openai.com/v1/chat/completions"
)

const (
	ChatMaxTokens        = 1024
	ChatTemperature      = 1
	ChatTopP             = 0
	ChatN                = 1
	ChatStream           = false
	ChatPresencePenalty  = 0
	ChatFrequencyPenalty = 0
	ChatStop             = ""
)

// Model
const (
	// Gpt35Turbo Most capable GPT-3.5 model and optimized for chat at 1/10th the cost of text-davinci-003. Will be updated with our latest model iteration.
	Gpt35Turbo = "gpt-3.5-turbo"
	// Gpt350301 Snapshot of gpt-3.5-turbo from March 1st 2023. Unlike gpt-3.5-turbo, this model will not receive updates, and will only be supported for a three month period ending on June 1st 2023.
	Gpt350301 = "gpt-3.5-turbo-0301"
	// Gpt4 More capable than any GPT-3.5 model, able to do more complex tasks, and optimized for chat. Will be updated with our latest model iteration.
	Gpt4 = "gpt-4"
	// Gpt40314 Snapshot of gpt-4 from March 14th 2023. Unlike gpt-4, this model will not receive updates, and will only be supported for a three month period ending on June 14th 2023.
	Gpt40314 = "gpt-4-0314"
)
