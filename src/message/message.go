// Package message
// @Description
// @Author root_wang
// @Date 2022/12/10 17:58
package message

import (
	"regexp"
	"strconv"
	"strings"
)

// CQCode 包含一个CQType和一系列键值对 不能确定有哪些键值对采取懒加载
type CQCode struct {
	rawMessage string
	keyValue   map[string]interface{}
}

func (c *CQCode) ParseKey(keys ...string) {
	for _, key := range keys {
		reg := key + `=(.*)[\],]{1}`
		if key == "CQ" {
			reg = key + `:(.*),`
		}
		keyReg := regexp.MustCompile(reg)
		matches := keyReg.FindStringSubmatch(c.rawMessage)
		if matches != nil {
			c.keyValue[key] = matches[1]
		}
	}
}

func (c *CQCode) ValueByKey(key string) interface{} {
	if value, ok := c.keyValue[key]; ok {
		return value
	}
	return nil
}

func (c *CQCode) IsAt() bool {
	c.ParseKey("qq")
	if c.ValueByKey("CQ") == "at" && c.ValueByKey("qq").(string) == strconv.FormatInt(QQ, 10) {
		return true
	}
	return false
}

// CQMessage 包含了CQCode和消息
type CQMessage struct {
	cqCode *CQCode
	msg    string
}

func (c *CQMessage) CQCode() *CQCode {
	return c.cqCode
}

func (c *CQMessage) Message() string {
	return c.msg
}

// RawMessage 原始信息主要包含了CQCode
type RawMessage string

func (m RawMessage) IsPlainMessage() bool {
	if strings.HasPrefix(string(m), "[CQ:") {
		return false
	}
	return true
}

func (m RawMessage) IsEmpty() bool {
	return string(m) == ""
}

func (m RawMessage) ToCQCode() (cqMsg *CQMessage, err error) {
	cqMsg = new(CQMessage)
	reg := `(\[.*\])\s(.+)`
	cqReg := regexp.MustCompile(reg)
	matches := cqReg.FindStringSubmatch(string(m))
	cqMsg.cqCode = &CQCode{
		rawMessage: matches[1],
		keyValue:   make(map[string]interface{}),
	}
	cqMsg.msg = matches[2]
	// 初始化CQMsg必须指明CQ类型
	cqMsg.cqCode.ParseKey("CQ")
	return
}

// Sender 发送消息的发送者
type Sender struct {
	Age      int64  `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserId   int64  `json:"user_id"`
}

type ReceiveMessage struct {
	PostType    string     `json:"post_type"`
	MessageType string     `json:"message_type"`
	Time        int64      `json:"time"`
	SelfId      int64      `json:"self_id"`
	SubType     string     `json:"sub_type"`
	UserId      int64      `json:"user_id"`
	MessageId   int64      `json:"message_id"`
	Font        int64      `json:"font"`
	GroupId     int64      `json:"group_id"`
	MessageSeq  int64      `json:"message_seq"`
	RawMessage  RawMessage `json:"raw_message"`
	Message     string     `json:"message"`
	Sender      Sender     `json:"sender"`
}
