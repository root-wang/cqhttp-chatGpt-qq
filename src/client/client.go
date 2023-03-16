// Package client
// @Description
// @Author root_wang
// @Date 2022/12/10 20:20
package client

import (
	"cqhttp-client/src/constant"
	"cqhttp-client/src/log"
	"cqhttp-client/src/message"
	"cqhttp-client/src/module"
	"cqhttp-client/src/module/gpt/API"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

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
	receiveMessageChan chan *message.ReceiveMessage
	responseChan       chan *Response
	// 同时为该QQ用户存储先前聊天记录 超时进行清除
	UserChatHistoryMap map[int64][]*API.ChatMessage
	// Chat模式计时器
	UserChatTimerMap map[int64]<-chan time.Time
	mu               sync.Mutex
	*module.Module
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		c:                  c,
		receiveMessageChan: make(chan *message.ReceiveMessage),
		responseChan:       make(chan *Response),
		UserChatHistoryMap: make(map[int64][]*API.ChatMessage),
		UserChatTimerMap:   make(map[int64]<-chan time.Time),
		Module:             module.NewModule(),
	}
}

func (c *Client) saveUserChatHistory(user int64, msg string, role constant.Role) []*API.ChatMessage {
	chatMessage := API.NewChatMessage(msg, role)
	defer c.mu.Unlock()
	c.mu.Lock()
	c.UserChatHistoryMap[user] = append(c.UserChatHistoryMap[user], chatMessage)

	return c.UserChatHistoryMap[user]
}

func (c *Client) formatGroupMessageReply(receiveMsg *message.ReceiveMessage, respMsg string) {

	// 将响应消息转换为CQCode结构体 然后再变为字符串
	cqCode := message.NewCQCode(respMsg, message.REPLY)
	cqCode.SetKeyValue([]message.CQKEY{message.ID}, receiveMsg.MessageId)

	resp := &Response{
		Action: SendGroupMsg,
		Params: message.GroupResp{
			GroupId: receiveMsg.GroupId,
			Message: fmt.Sprintf("%s", cqCode),
			// client 自己转不需要server再转了
			AutoEscape: false,
		},
	}
	c.responseChan <- resp
}

func (c *Client) ClearUserChat(receiveMsg *message.ReceiveMessage) {
	userId := receiveMsg.UserId
	groupId := receiveMsg.GroupId
	for {
		select {
		case <-c.UserChatTimerMap[userId]:
			c.mu.Lock()
			delete(c.UserChatTimerMap, userId)
			delete(c.UserChatHistoryMap, userId)
			c.mu.Unlock()
			log.Infof("user%v end\n", userId)

			cqCode := message.NewCQCode("你的对话聊天时间到了,如需要体验,则再次使用c&来开启", message.AT)
			cqCode.SetKeyValue([]message.CQKEY{message.QQ}, userId)

			resp := &Response{
				Action: SendGroupMsg,
				Params: message.GroupResp{
					GroupId: groupId,
					Message: fmt.Sprintf("%s", cqCode),
					// client 自己转不需要server再转了
					AutoEscape: false,
				},
			}
			c.responseChan <- resp
		}
	}
}

func (c *Client) ReceiveMessage(message *message.ReceiveMessage) error {
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
			// 获取发送者qq号码
			UserId := receiveMsg.Sender.UserId
			// 获取发送的经过处理的文本消息
			msg := receiveMsg.Message
			moduleType, prompt := message.CheckModuleType(msg)
			// 开启聊天模式
			if moduleType == constant.GPTCHAT && c.UserChatTimerMap[UserId] == nil {
				// 为用户添加5分钟计时器
				const tenMin = time.Minute * 2
				c.mu.Lock()
				c.UserChatTimerMap[UserId] = time.Tick(tenMin)
				c.mu.Unlock()

				go c.ClearUserChat(receiveMsg)
				// 保存用户第一次的提问
				chatMessageSlice := c.saveUserChatHistory(UserId, prompt, constant.USER)
				tick := time.Now()
				handlerMessage, err := c.Handler(constant.GPTCHAT).HandlerMessage(chatMessageSlice)
				if err != nil {
					log.Errorf("reply group message failed with chat module: %v", err)
				}
				// 保存gpt返回第一次的回答
				c.saveUserChatHistory(UserId, handlerMessage, constant.ASSISTANT)

				respMsg := fmt.Sprintf("%s\n\nfrom OPENAI", handlerMessage)
				log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())

				c.formatGroupMessageReply(receiveMsg, respMsg)
			} else if moduleType != constant.GPTIMAGE && c.UserChatTimerMap[UserId] != nil {
				chatMessageSlice := c.saveUserChatHistory(UserId, prompt, constant.USER)
				tick := time.Now()
				handlerMessage, err := c.Handler(constant.GPTCHAT).HandlerMessage(chatMessageSlice)
				if err != nil {
					log.Errorf("reply group message failed with chat module: %v", err)
				}
				c.saveUserChatHistory(UserId, handlerMessage, constant.ASSISTANT)

				respMsg := fmt.Sprintf("%s\n\nfrom OPENAI", handlerMessage)
				log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())

				c.formatGroupMessageReply(receiveMsg, respMsg)
			} else {
				// 	不是聊天模式
				go func(receiveMsg *message.ReceiveMessage) {
					var handlerMessage string
					var err error
					// 使用第三方模块对消息进行处理
					tick := time.Now()

					handlerMessage, err = c.Handler(moduleType).HandlerMessage(prompt)
					if err != nil {
						log.Errorf("reply group message failed: %v", err)
					}
					respMsg := fmt.Sprintf("%s\n\nfrom OPENAI", handlerMessage)
					log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())

					c.formatGroupMessageReply(receiveMsg, respMsg)
				}(receiveMsg)
			}
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
