// Package clt
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:34
//
package clt

import (
	cst "cqhttp-client/src/constant"
	msg "cqhttp-client/src/message"
	"cqhttp-client/src/module"
	"cqhttp-client/src/module/gpt"
	"github.com/gorilla/websocket"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestClient_ClearUserChat(t *testing.T) {
	type fields struct {
		c                  *websocket.Conn
		receiveMessageChan chan *msg.ReceiveMessage
		responseChan       chan *Response
		UserChatHistoryMap map[int64][]*gpt.ChatMessage
		UserChatTimerMap   map[int64]<-chan time.Time
		mu                 sync.Mutex
		Module             *module.Module
	}
	type args struct {
		receiveMsg *msg.ReceiveMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Client{
					c:                  tt.fields.c,
					receiveMessageChan: tt.fields.receiveMessageChan,
					responseChan:       tt.fields.responseChan,
					UserChatHistoryMap: tt.fields.UserChatHistoryMap,
					UserChatTimerMap:   tt.fields.UserChatTimerMap,
					mu:                 tt.fields.mu,
					Module:             tt.fields.Module,
				}
				c.ClearUserChat(tt.args.receiveMsg)
			},
		)
	}
}

func TestClient_ReceiveMessage(t *testing.T) {
	type fields struct {
		c                  *websocket.Conn
		receiveMessageChan chan *msg.ReceiveMessage
		responseChan       chan *Response
		UserChatHistoryMap map[int64][]*gpt.ChatMessage
		UserChatTimerMap   map[int64]<-chan time.Time
		mu                 sync.Mutex
		Module             *module.Module
	}
	type args struct {
		message *msg.ReceiveMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Client{
					c:                  tt.fields.c,
					receiveMessageChan: tt.fields.receiveMessageChan,
					responseChan:       tt.fields.responseChan,
					UserChatHistoryMap: tt.fields.UserChatHistoryMap,
					UserChatTimerMap:   tt.fields.UserChatTimerMap,
					mu:                 tt.fields.mu,
					Module:             tt.fields.Module,
				}
				if err := c.ReceiveMessage(tt.args.message); (err != nil) != tt.wantErr {
					t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestClient_ReplyGroupMessage(t *testing.T) {
	type fields struct {
		c                  *websocket.Conn
		receiveMessageChan chan *msg.ReceiveMessage
		responseChan       chan *Response
		UserChatHistoryMap map[int64][]*gpt.ChatMessage
		UserChatTimerMap   map[int64]<-chan time.Time
		mu                 sync.Mutex
		Module             *module.Module
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Client{
					c:                  tt.fields.c,
					receiveMessageChan: tt.fields.receiveMessageChan,
					responseChan:       tt.fields.responseChan,
					UserChatHistoryMap: tt.fields.UserChatHistoryMap,
					UserChatTimerMap:   tt.fields.UserChatTimerMap,
					mu:                 tt.fields.mu,
					Module:             tt.fields.Module,
				}
				c.ReplyGroupMessage()
			},
		)
	}
}
func TestClient_formatGroupMessage(t *testing.T) {
	type fields struct {
		c                  *websocket.Conn
		receiveMessageChan chan *msg.ReceiveMessage
		responseChan       chan *Response
		UserChatHistoryMap map[int64][]*gpt.ChatMessage
		UserChatTimerMap   map[int64]<-chan time.Time
		mu                 sync.Mutex
		Module             *module.Module
	}
	type args struct {
		groupId int64
		cqCode  *msg.CQCode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Client{
					c:                  tt.fields.c,
					receiveMessageChan: tt.fields.receiveMessageChan,
					responseChan:       tt.fields.responseChan,
					UserChatHistoryMap: tt.fields.UserChatHistoryMap,
					UserChatTimerMap:   tt.fields.UserChatTimerMap,
					mu:                 tt.fields.mu,
					Module:             tt.fields.Module,
				}
				c.formatGroupMessage(tt.args.groupId, tt.args.cqCode)
			},
		)
	}
}

func TestClient_saveUserChatHistory(t *testing.T) {
	type fields struct {
		c                  *websocket.Conn
		receiveMessageChan chan *msg.ReceiveMessage
		responseChan       chan *Response
		UserChatHistoryMap map[int64][]*gpt.ChatMessage
		UserChatTimerMap   map[int64]<-chan time.Time
		mu                 sync.Mutex
		Module             *module.Module
	}
	type args struct {
		user    int64
		message string
		role    cst.Role
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*gpt.ChatMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Client{
					c:                  tt.fields.c,
					receiveMessageChan: tt.fields.receiveMessageChan,
					responseChan:       tt.fields.responseChan,
					UserChatHistoryMap: tt.fields.UserChatHistoryMap,
					UserChatTimerMap:   tt.fields.UserChatTimerMap,
					mu:                 tt.fields.mu,
					Module:             tt.fields.Module,
				}
				if got := c.saveUserChatHistory(tt.args.user, tt.args.message, tt.args.role); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("saveUserChatHistory() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
