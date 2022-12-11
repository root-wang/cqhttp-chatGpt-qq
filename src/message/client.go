// Package message
// @Description
// @Author root_wang
// @Date 2022/12/10 20:20
package message

import (
	"cqhttp-client/src/log"
	module "cqhttp-client/src/module"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const BotQQ int64 = 1966962723

type ACTION string

func (A ACTION) String() string {
	return string(A)
}

const (
	SendGroupMsg ACTION = "send_group_msg"
)

type Response struct {
	Action ACTION      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type Client struct {
	c                  *websocket.Conn
	receiveMessageChan chan *ReceiveMessage
	responseChan       chan *Response
	*module.Module
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		c:                  c,
		receiveMessageChan: make(chan *ReceiveMessage),
		responseChan:       make(chan *Response),
		Module:             module.NewModule(),
	}
}

func (c *Client) ReceiveMessage(message *ReceiveMessage) error {
	if !message.RawMessage.IsPlainMessage() {
		cqMsg, err := message.RawMessage.ToCQCode()
		if err != nil {
			return log.ErrorInsidef("transform raw message 2 cqcode failed: %v", err)
		}
		if cqMsg.CQCode().IsAt() {
			log.Infof("get quote message:%s from group:%d", cqMsg.Message(), message.GroupId)
			message.Message = cqMsg.Message()
			c.receiveMessageChan <- message
		}
	}
	return nil
}

// ReplyGroupMessage if some group member at bot reply the message
func (c *Client) ReplyGroupMessage() {
	for {
		select {
		case receiveMsg := <-c.receiveMessageChan:
			go func(receiveMsg *ReceiveMessage) {
				msg := receiveMsg.Message
				// 使用第三方模块对消息进行处理
				tick := time.Now()
				// handleMsg := testGoroutine(msg)
				handleMsg, err := c.Handler("gpt").HandlerMessage(msg)
				if err != nil {
					log.Errorf("reply group message failed: %v", err)
				}
				respMsg := fmt.Sprintf(
					"%s\n\nfrom OPENAI \n\n处理此消息共用时 %.2f s", handleMsg,
					time.Since(tick).Seconds(),
				)
				log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())
				// 将响应消息转换为CQCode结构体 然后再变为字符串
				cqCode := &CQCode{
					rawMessage: respMsg,
					keyValue:   make(map[CQKEY]string),
					cqtype:     "",
				}
				cqCode.SetType(REPLY)
				cqCode.SetKeyValue([]CQKEY{ID}, receiveMsg.MessageId)

				resp := &Response{
					Action: SendGroupMsg,
					Params: groupResp{
						GroupId: receiveMsg.GroupId,
						Message: fmt.Sprintf("%s", cqCode),
						// client 自己转不需要server再转了
						AutoEscape: false,
					},
				}
				c.responseChan <- resp
			}(receiveMsg)
		}
	}
}

func testGoroutine(msg string) (handleMsg string) {
	handleMsg = "test" + msg
	time.Sleep(time.Duration(len(msg)) * time.Second)
	return
}

func (c *Client) Run() {
	for {
		select {
		case resp := <-c.responseChan:
			err := c.c.WriteJSON(resp)
			log.Infof("send message to %s", resp.Action)
			if err != nil {
				log.Errorf("write response message 2 JSON failed: %v", err)
				return
			}
		}
	}
}
