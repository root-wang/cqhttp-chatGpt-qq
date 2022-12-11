// Package gpt
// @Description 此模块由https://github.com/m1guelpf/chatgpt-telegram开发
// @Author
// @Date 2022/12/11 14:31
package gpt

import (
	"cqhttp-client/src/log"
	"cqhttp-client/src/module/gpt/expirymap"
	"cqhttp-client/src/module/gpt/sse"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"time"
)

const KeyAccessToken = "accessToken"
const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

type ChatGPT struct {
	SessionToken   string
	AccessTokenMap expirymap.ExpiryMap
}

type SessionResult struct {
	Error       string `json:"error"`
	Expires     string `json:"expires"`
	AccessToken string `json:"accessToken"`
}

type MessageResponse struct {
	ConversationId string `json:"conversation_id"`
	Error          string `json:"error"`
	Message        struct {
		ID      string `json:"id"`
		Content struct {
			Parts []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
}

type ChatResponse struct {
	Message        string
	MessageId      string
	ConversationId string
}

// func NewGpt(config config.Config) ChatGPT {
//     return ChatGPT{
//         AccessTokenMap: expirymap.New(),
//         SessionToken:   config.OpenAISession,
//     }
// }

func NewGpt(token string) *ChatGPT {
	return &ChatGPT{
		AccessTokenMap: expirymap.New(),
		SessionToken:   token,
	}
}

func (c *ChatGPT) IsAuthenticated() bool {
	_, err := c.refreshAccessToken()
	return err == nil
}

func (c *ChatGPT) EnsureAuth() error {
	_, err := c.refreshAccessToken()
	return err
}

func (c *ChatGPT) SendMessage(message string, conversationId string, messageId string) (chan ChatResponse, error) {
	r := make(chan ChatResponse)
	accessToken, err := c.refreshAccessToken()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't get access token: %v", err))
	}

	client := sse.Init("https://chat.openai.com/backend-api/conversation")

	client.Headers = map[string]string{
		"User-Agent":    UserAgent,
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	err = client.Connect(message, conversationId, messageId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't connect to ChatGPT: %v", err))
	}

	go func() {
		defer close(r)
	mainLoop:
		for {
			select {
			case chunk, ok := <-client.EventChannel:
				if !ok {
					break mainLoop
				}

				var res MessageResponse
				err := json.Unmarshal([]byte(chunk), &res)
				if err != nil {
					log.Infof("Couldn't unmarshal message response: %v", err)
					continue
				}

				if len(res.Message.Content.Parts) > 0 {
					r <- ChatResponse{
						MessageId:      res.Message.ID,
						ConversationId: res.ConversationId,
						Message:        res.Message.Content.Parts[0],
					}
				}
			}
		}
	}()

	return r, nil
}

func (c *ChatGPT) refreshAccessToken() (string, error) {
	cachedAccessToken, ok := c.AccessTokenMap.Get(KeyAccessToken)
	if ok {
		return cachedAccessToken, nil
	}

	req, err := http.NewRequest("GET", "https://chat.openai.com/api/auth/session", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cookie", fmt.Sprintf("__Secure-next-auth.session-token=%s", c.SessionToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer res.Body.Close()

	var result SessionResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	accessToken := result.AccessToken
	if accessToken == "" {
		return "", errors.New("unauthorized")
	}

	if result.Error != "" {
		if result.Error == "RefreshAccessTokenError" {
			return "", errors.New("Session token has expired")
		}

		return "", errors.New(result.Error)
	}

	expiryTime, err := time.Parse(time.RFC3339, result.Expires)
	if err != nil {
		return "", fmt.Errorf("failed to parse expiry time: %v", err)
	}
	c.AccessTokenMap.Set(KeyAccessToken, accessToken, expiryTime.Sub(time.Now()))

	return accessToken, nil
}

func (c *ChatGPT) HandlerMessage(s string) string {
	feed, err := c.SendMessage(s, "", "")
	if err != nil {
		log.Error(err.Error())
	}
	var msg string
pollResponse:
	for {
		select {
		case response, ok := <-feed:
			if !ok {
				break pollResponse
			}
			msg = response.Message
		}
	}
	return msg
}
