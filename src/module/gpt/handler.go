// Package gpt
// @Description
// @Author root_wang
// @Date 2022/12/12 20:59
package gpt

import "cqhttp-client/src/log"

func (c *ChatGPT) HandlerMessage(s string) (string, error) {
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
	return msg, nil
}
