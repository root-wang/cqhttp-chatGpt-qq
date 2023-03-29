// Package msg
// @Description 
// @Author root_wang
// @Date 2023/3/21 23:51
//
package msg

import (
	"reflect"
	"testing"
)

func TestCQCode_String(t *testing.T) {
	type fields struct {
		rawMessage string
		keyValue   map[CQKEY]string
		cqtype     CQTYPE
	}
	tests := []struct {
		name   string
		fields *fields
		want   string
	}{
		{
			name: "at someone",
			fields: &fields{
				rawMessage: "早上好啊",
				keyValue:   map[CQKEY]string{QQ: "114514"},
				cqtype:     AT,
			},
			want: "[CQ:at,qq=114514]早上好啊",
		}, {
			name: "reply someone's message by message Id",
			fields: &fields{
				rawMessage: "早上好啊",
				keyValue:   map[CQKEY]string{ID: "123456"},
				cqtype:     REPLY,
			},
			want: "[CQ:reply,id=123456]早上好啊",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &CQCode{
					rawMessage: tt.fields.rawMessage,
					keyValue:   tt.fields.keyValue,
					cqtype:     tt.fields.cqtype,
				}
				if got := c.String(); got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCQCode_IsAt(t *testing.T) {
	type fields struct {
		rawMessage string
		keyValue   map[CQKEY]string
		cqtype     CQTYPE
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:at,qq=3223893903]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     AT,
			},
			want: true,
		}, {
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:at,qq=233403]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     AT,
			},
			want: false,
		}, {
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:reply,id=123456]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     REPLY,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &CQCode{
					rawMessage: tt.fields.rawMessage,
					keyValue:   tt.fields.keyValue,
					cqtype:     tt.fields.cqtype,
				}
				if got := c.IsAt(); got != tt.want {
					t.Errorf("IsAt() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCQCode_ParseKey(t *testing.T) {
	type fields struct {
		rawMessage string
		keyValue   map[CQKEY]string
		cqtype     CQTYPE
	}
	type args struct {
		keys []CQKEY
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:at,qq=3223893903]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     "",
			},
			args: args{keys: []CQKEY{QQ}},
		}, {
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:at,qq=233403]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     "",
			},
			args: args{keys: []CQKEY{QQ}},
		}, {
			name: "at bot",
			fields: fields{
				rawMessage: "[CQ:reply,id=123456]hello",
				keyValue:   map[CQKEY]string{},
				cqtype:     "",
			},
			args: args{keys: []CQKEY{ID}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &CQCode{
					rawMessage: tt.fields.rawMessage,
					keyValue:   tt.fields.keyValue,
					cqtype:     tt.fields.cqtype,
				}
				// todo:完成测试判断
				c.ParseKey(tt.args.keys...)
			},
		)
	}
}

func TestCQCode_SetKeyValue(t *testing.T) {
	type fields struct {
		rawMessage string
		keyValue   map[CQKEY]string
		cqtype     CQTYPE
	}
	type args struct {
		keys   []CQKEY
		values []interface{}
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
				c := &CQCode{
					rawMessage: tt.fields.rawMessage,
					keyValue:   tt.fields.keyValue,
					cqtype:     tt.fields.cqtype,
				}
				c.SetKeyValue(tt.args.keys, tt.args.values...)
			},
		)
	}
}

func TestCQMessage_IsEmpty(t *testing.T) {
	type fields struct {
		cqCode *CQCode
		msg    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				cqCode: &CQCode{
					rawMessage: "[CQ:at,qq=3223893903]hello",
					keyValue:   map[CQKEY]string{},
					cqtype:     AT,
				},
				msg: "hello",
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				cqCode: &CQCode{
					rawMessage: "[CQ:at,qq=3223893903]",
					keyValue:   map[CQKEY]string{},
					cqtype:     AT,
				},
				msg: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &CQMessage{
					cqCode: tt.fields.cqCode,
					msg:    tt.fields.msg,
				}
				if got := c.IsEmpty(); got != tt.want {
					t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRawMessage_IsPlainMessage(t *testing.T) {
	tests := []struct {
		name string
		m    RawMessage
		want bool
	}{
		{
			name: "video",
			m:    RawMessage("[CQ:video,file=http://baidu.com/1.mp4]"),
			want: false,
		}, {
			name: "plain",
			m:    RawMessage("hello"),
			want: false,
		}, {
			name: "contact",
			m:    RawMessage("[CQ:contact,type=group,id=100100]hello"),
			want: false,
		}, {
			name: "at",
			m:    RawMessage("[CQ:at,qq=10001000]hello"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.m.IsAtMessage(); got != tt.want {
					t.Errorf("IsAtMessage() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRawMessage_ToCQCode(t *testing.T) {
	tests := []struct {
		name      string
		m         RawMessage
		wantCqMsg *CQMessage
		wantErr   bool
	}{
		{
			name: "at",
			m:    RawMessage("[CQ:at,qq=10001000]hello"),
			wantCqMsg: &CQMessage{
				cqCode: &CQCode{
					rawMessage: "[CQ:at,qq=10001000]hello",
					keyValue:   map[CQKEY]string{QQ: "10001000"},
					cqtype:     AT,
				},
				msg: "hello",
			},
			wantErr: false,
		},
		{
			name: "reply",
			m:    RawMessage("[CQ:reply,text=Hello World,qq=10086,time=3376656000,seq=5123]hll"),
			wantCqMsg: &CQMessage{
				cqCode: &CQCode{
					rawMessage: "[CQ:reply,text=Hello World,qq=10086,time=3376656000,seq=5123]hll",
					keyValue:   map[CQKEY]string{TEXT: "Hello World", QQ: "10086", TIME: "3376656000", SEQ: "5123"},
					cqtype:     REPLY,
				},
				msg: "hll",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotCqMsg, err := tt.m.ToCQCode()
				if (err == nil) == tt.wantErr {
					t.Errorf("ToCQCode() error = %v, wantErr %v", (err == nil) == tt.wantErr, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(gotCqMsg.Message(), tt.wantCqMsg.Message()) {
					t.Errorf("ToCQCode() gotCqMsg = %v, want %v", gotCqMsg, tt.wantCqMsg)
				}
			},
		)
	}
}
