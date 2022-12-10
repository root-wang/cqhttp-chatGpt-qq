// Package message
// @Description
// @Author root_wang
// @Date 2022/12/10 20:20
package message

import (
	"github.com/gorilla/websocket"
	"strconv"
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
	if !message.RawMessage.IsEmpty() {
		if message.RawMessage.IsPlainMessage() {
			message.Message = string(message.RawMessage)
		} else {
			cqMsg, err := message.RawMessage.ToCQCode()
			if err != nil {
				return err
			}
			if cqMsg.CQCode().IsAt() {
				println(cqMsg.Message())
				message.Message = cqMsg.Message()
			}
		}
	}
	c.receiveMessageChan <- message
	return nil
}

func (c *Client) SendGroupMessage() {
	for {
		select {
		case receiveMsg := <-c.receiveMessageChan:
			msg := receiveMsg.Message
			respMsg := "[CQ:reply,id=" + strconv.FormatInt(
				receiveMsg.MessageId, 10,
			) + "]" + "this is test for" + msg
			resp := &Response{
				Action: "send_group_msg",
				Params: groupResp{
					GroupId:    receiveMsg.GroupId,
					Message:    respMsg,
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
			if err != nil {
				return
			}
		}
	}
}
