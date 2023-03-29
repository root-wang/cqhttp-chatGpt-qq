// Package gpt
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:33
//
package gpt

import (
	"cqhttp-client/src/module"
	"reflect"
	"testing"
)

func TestAPI_APIByName(t *testing.T) {
	type fields struct {
		apis   map[SaasName]Parser
		urls   map[SaasName]string
		ApiKey string
	}
	type args struct {
		n SaasName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   module.Moduler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				A := &API{
					apis:   tt.fields.apis,
					urls:   tt.fields.urls,
					ApiKey: tt.fields.ApiKey,
				}
				if got := A.APIByName(tt.args.n); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("APIByName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestAPI_HandlerMessage(t *testing.T) {
	type fields struct {
		apis   map[SaasName]Parser
		urls   map[SaasName]string
		ApiKey string
	}
	type args struct {
		s   interface{}
		api SaasName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				A := &API{
					apis:   tt.fields.apis,
					urls:   tt.fields.urls,
					ApiKey: tt.fields.ApiKey,
				}
				got, err := A.HandlerMessage(tt.args.s, tt.args.api)
				if (err != nil) != tt.wantErr {
					t.Errorf("HandlerMessage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("HandlerMessage() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestAPI_MakeRequestStruct(t *testing.T) {
	type fields struct {
		apis   map[SaasName]Parser
		urls   map[SaasName]string
		ApiKey string
	}
	type args struct {
		api SaasName
		msg interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				A := &API{
					apis:   tt.fields.apis,
					urls:   tt.fields.urls,
					ApiKey: tt.fields.ApiKey,
				}
				if got := A.MakeRequestStruct(tt.args.api, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MakeRequestStruct() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestAPI_MakeRespStruct(t *testing.T) {
	type fields struct {
		apis   map[SaasName]Parser
		urls   map[SaasName]string
		ApiKey string
	}
	type args struct {
		api SaasName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				A := &API{
					apis:   tt.fields.apis,
					urls:   tt.fields.urls,
					ApiKey: tt.fields.ApiKey,
				}
				if got := A.MakeRespStruct(tt.args.api); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MakeRespStruct() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestAPI_ParseMessage(t *testing.T) {
	type fields struct {
		apis   map[SaasName]Parser
		urls   map[SaasName]string
		ApiKey string
	}
	type args struct {
		api  SaasName
		resp interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				A := &API{
					apis:   tt.fields.apis,
					urls:   tt.fields.urls,
					ApiKey: tt.fields.ApiKey,
				}
				got, err := A.ParseMessage(tt.args.api, tt.args.resp)
				if (err != nil) != tt.wantErr {
					t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("ParseMessage() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestChatApi_HandlerMessage(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := ChatApi{
					API: tt.fields.API,
				}
				got, err := c.HandlerMessage(tt.args.s)
				if (err != nil) != tt.wantErr {
					t.Errorf("HandlerMessage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("HandlerMessage() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestChatApi_Matcher(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantPrompt string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := ChatApi{
					API: tt.fields.API,
				}
				gotOk, gotPrompt := c.Matcher(tt.args.s)
				if gotOk != tt.wantOk {
					t.Errorf("Matcher() gotOk = %v, want %v", gotOk, tt.wantOk)
				}
				if gotPrompt != tt.wantPrompt {
					t.Errorf("Matcher() gotPrompt = %v, want %v", gotPrompt, tt.wantPrompt)
				}
			},
		)
	}
}

func TestChatApi_Parse(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		resp interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := ChatApi{
					API: tt.fields.API,
				}
				got, err := c.Parse(tt.args.resp)
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestImageApi_HandlerMessage(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				i := ImageApi{
					API: tt.fields.API,
				}
				got, err := i.HandlerMessage(tt.args.s)
				if (err != nil) != tt.wantErr {
					t.Errorf("HandlerMessage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("HandlerMessage() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestImageApi_Matcher(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOk     bool
		wantPrompt string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				i := ImageApi{
					API: tt.fields.API,
				}
				gotOk, gotPrompt := i.Matcher(tt.args.s)
				if gotOk != tt.wantOk {
					t.Errorf("Matcher() gotOk = %v, want %v", gotOk, tt.wantOk)
				}
				if gotPrompt != tt.wantPrompt {
					t.Errorf("Matcher() gotPrompt = %v, want %v", gotPrompt, tt.wantPrompt)
				}
			},
		)
	}
}

func TestImageApi_Parse(t *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		resp interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				i := ImageApi{
					API: tt.fields.API,
				}
				got, err := i.Parse(tt.args.resp)
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestInitApi(t *testing.T) {
	tests := []struct {
		name string
		want *API
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := InitApi(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("InitApi() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSaasName_String(t *testing.T) {
	tests := []struct {
		name string
		s    SaasName
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.s.String(); got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestTextApi_HandlerMessage(t1 *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(
			tt.name, func(t1 *testing.T) {
				t := TextApi{
					API: tt.fields.API,
				}
				got, err := t.HandlerMessage(tt.args.s)
				if (err != nil) != tt.wantErr {
					t1.Errorf("HandlerMessage() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t1.Errorf("HandlerMessage() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestTextApi_Matcher(t1 *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(
			tt.name, func(t1 *testing.T) {
				t := TextApi{
					API: tt.fields.API,
				}
				got, got1 := t.Matcher(tt.args.s)
				if got != tt.want {
					t1.Errorf("Matcher() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t1.Errorf("Matcher() got1 = %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}

func TestTextApi_Parse(t1 *testing.T) {
	type fields struct {
		API *API
	}
	type args struct {
		resp interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(
			tt.name, func(t1 *testing.T) {
				t := TextApi{
					API: tt.fields.API,
				}
				got, err := t.Parse(tt.args.resp)
				if (err != nil) != tt.wantErr {
					t1.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t1.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
