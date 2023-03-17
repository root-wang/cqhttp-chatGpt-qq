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
	KeyAccessToken                   = "accessToken"
	UserAgent                        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	InitURL                          = "https://chat.openai.com/backend-api/conversation"
)

type Role string

const (
	SYSTEM    Role = "system"
	USER           = "user"
	ASSISTANT      = "assistant"
)
