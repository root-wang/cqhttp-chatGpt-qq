// Package cst
// @Description
// @Author root_wang
// @Date 2022/12/11 20:46
package cst

type GPTError string

func (G GPTError) String() string {
	return string(G)
}

const (
	RefreshAccessTokenError GPTError = "RefreshAccessTokenError"
)

type Role string

const (
	SYSTEM    Role = "system"
	USER           = "user"
	ASSISTANT      = "assistant"
)

const (
	GptText  = "gptText"
	GptImage = "gptImage"
	GptChat  = "gptChat"
)

const (
	ErrorOpenAIResponse      = "请求OpenAI出错啦"
	ErrorOpenAIParseResponse = "处理OpenAI的响应出错啦"
)
