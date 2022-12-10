// Package message
// @Description
// @Author root_wang
// @Date 2022/12/10 20:20
package message

import (
	"cqhttp-client/src/log"
	"github.com/gorilla/websocket"
)

const QQ int64 = 1966962723

type Response struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type Client struct {
	c                  *websocket.Conn
	receiveMessageChan chan *ReceiveMessage
	responseChan       chan *Response
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		c:                  c,
		receiveMessageChan: make(chan *ReceiveMessage),
		responseChan:       make(chan *Response),
	}
}

func (c *Client) ReceiveMessage(message *ReceiveMessage) error {
	if !message.RawMessage.IsPlainMessage() {
		cqMsg, err := message.RawMessage.ToCQCode()
		if err != nil {
			return err
		}
		if cqMsg.CQCode().IsAt() {
			log.Infof("get quote message:%s from group:%d", cqMsg.Message(), message.GroupId)
			message.Message = cqMsg.Message()
			c.receiveMessageChan <- message
		}

	}
	return nil
}

func (c *Client) ReplyGroupMessage() {
	for {
		select {
		case receiveMsg := <-c.receiveMessageChan:
			msg := receiveMsg.Message
			respMsg := "this is test for " + msg
			// 将响应消息转换为CQCode结构体 然后再变为字符串
			cqCode := &CQCode{
				rawMessage: respMsg,
				keyValue:   make(map[string]interface{}),
			}
			cqCode.SetType("reply")
			cqCode.SetKeyValue([]string{"id"}, receiveMsg.MessageId)
			resp := &Response{
				Action: "send_group_msg",
				Params: groupResp{
					GroupId: receiveMsg.GroupId,
					Message: cqCode.String(),
					// client 自己转不需要server再转了
					AutoEscape: false,
				},
			}
			c.responseChan <- resp
		}
	}
}

func (c *Client) Run() {
	for {
		select {
		case resp := <-c.responseChan:
			err := c.c.WriteJSON(resp)
			log.Infof("send message to %s", resp.Action)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}
