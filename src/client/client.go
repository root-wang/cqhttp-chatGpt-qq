// Package clt
// @Description
// @Author root_wang
// @Date 2022/12/10 20:20
package clt

import (
	"cqhttp-client/src/constant"
	"cqhttp-client/src/log"
	"cqhttp-client/src/message"
	"cqhttp-client/src/module"
	"cqhttp-client/src/module/gpt"
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
	// SendGroupMessage 响应cq-http字段 勿动
	SendGroupMessage ACTION = "send_group_msg"
)

type Response struct {
	Action ACTION      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

type Client struct {
	c                  *websocket.Conn
	receiveMessageChan chan *msg.ReceiveMessage
	responseChan       chan *Response
	// 同时为该QQ用户存储先前聊天记录 超时进行清除
	UserChatHistoryMap map[int64][]*gpt.ChatMessage
	// Chat模式计时器
	UserChatTimerMap map[int64]<-chan time.Time
	mu               sync.Mutex
	*module.Module
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		c:                  c,
		receiveMessageChan: make(chan *msg.ReceiveMessage),
		responseChan:       make(chan *Response),
		UserChatHistoryMap: make(map[int64][]*gpt.ChatMessage),
		UserChatTimerMap:   make(map[int64]<-chan time.Time),
		Module:             module.NewModule(),
	}
}

// saveUserChatHistory 为用户保存聊天时间内的聊天历史记录
func (c *Client) saveUserChatHistory(user int64, message string, role cst.Role) []*gpt.ChatMessage {
	chatMessage := gpt.NewChatMessage(message, role)
	defer c.mu.Unlock()
	c.mu.Lock()
	c.UserChatHistoryMap[user] = append(c.UserChatHistoryMap[user], chatMessage)

	return c.UserChatHistoryMap[user]
}

// formatGroupMessage 将要发送的群消息format并且送入管道中
func (c *Client) formatGroupMessage(groupId int64, cqCode *msg.CQCode) {
	resp := &Response{
		Action: SendGroupMessage,
		Params: msg.GroupResp{
			GroupId: groupId,
			Message: fmt.Sprintf("%s", cqCode),
			// client 自己转不需要server再转了
			AutoEscape: false,
		},
	}
	c.responseChan <- resp
}

// ClearUserChat 当有用户的聊天定时器到期时进行处理
func (c *Client) ClearUserChat(receiveMsg *msg.ReceiveMessage) {
	userId := receiveMsg.UserId
	for {
		select {
		case <-c.UserChatTimerMap[userId]:
			c.mu.Lock()
			// 删除用户定时器
			delete(c.UserChatTimerMap, userId)
			// 删除用户的聊天历史
			delete(c.UserChatHistoryMap, userId)
			c.mu.Unlock()
			log.Infof("user%v end\n", userId)

			// 向QQ群中提示用户聊天时间截至
			cqCode := msg.NewCQCode("你的对话聊天时间到了,如需要体验,则再次使用c&来开启", msg.AT)
			cqCode.SetKeyValue([]msg.CQKEY{msg.QQ}, userId)
			c.formatGroupMessage(receiveMsg.GroupId, cqCode)
		}
	}
}

// ReceiveMessage 从qq群中接收@机器人的消息
func (c *Client) ReceiveMessage(message *msg.ReceiveMessage) error {
	// 只处理群中有用户@机器人的消息
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
		case receiveMessage := <-c.receiveMessageChan:
			// 获取发送者qq号码
			userId := receiveMessage.Sender.UserId
			// 获取发送信息的qq群号
			groupId := receiveMessage.GroupId
			// 获取用户聊天计时器
			timer := c.UserChatTimerMap[userId]
			// 获取请求模式
			moduleType, prompt := c.CheckModuleType(receiveMessage.Message)
			switch {
			// 开启聊天模式
			case moduleType == cst.GptChat || (moduleType == cst.GptText && timer != nil):
				// 用户第一次的提问
				if timer == nil {
					// 为用户添加5分钟计时器
					const tenMin = time.Minute * 2
					c.mu.Lock()
					c.UserChatTimerMap[userId] = time.Tick(tenMin)
					c.mu.Unlock()

					go c.ClearUserChat(receiveMessage)
				}
				chatMessageSlice := c.saveUserChatHistory(userId, prompt, cst.USER)
				tick := time.Now()
				handlerMessage, err := c.Handler(cst.GptChat).HandlerMessage(chatMessageSlice)
				if err != nil {
					log.Errorf("reply group message failed with chat module: %v", err)
				}
				// 保存gpt返回第一次的回答
				if c.UserChatHistoryMap[userId] != nil {
					c.saveUserChatHistory(userId, handlerMessage, cst.ASSISTANT)
				}

				respMessage := fmt.Sprintf("%s\n\nfrom OPENAI", handlerMessage)
				log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())

				// 将响应消息转换为CQCode结构体 然后再变为字符串
				cqCode := msg.NewCQCode(respMessage, msg.REPLY)
				cqCode.SetKeyValue([]msg.CQKEY{msg.ID}, receiveMessage.MessageId)
				c.formatGroupMessage(groupId, cqCode)

			case moduleType == cst.GptImage || moduleType == cst.GptText:
				// 	生成图片或者QA模式
				go func(receiveMessage *msg.ReceiveMessage) {
					var handlerMessage string
					var err error
					// 使用第三方模块对消息进行处理
					tick := time.Now()

					handlerMessage, err = c.Handler(moduleType).HandlerMessage(prompt)
					if err != nil {
						log.Errorf("reply group message failed: %v", err)
					}
					respMessage := fmt.Sprintf("%s\n\nfrom OPENAI", handlerMessage)
					log.Infof("处理此消息共用时 %.2f s", time.Since(tick).Seconds())

					cqCode := msg.NewCQCode(respMessage, msg.REPLY)
					cqCode.SetKeyValue([]msg.CQKEY{msg.ID}, receiveMessage.MessageId)
					c.formatGroupMessage(groupId, cqCode)
				}(receiveMessage)
			}
		}
	}
}

func testGoroutine(message string) (handleMsg string) {
	handleMsg = "test" + message
	time.Sleep(time.Duration(len(message)) * time.Second)
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
