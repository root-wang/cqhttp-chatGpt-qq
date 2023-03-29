// Package utils
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:33
//
package utils

import "testing"

func TestAny2string(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		{
			name:    "string",
			args:    args{v: "hello"},
			wantStr: "hello",
		},
		{
			name:    "int64",
			args:    args{v: int64(298398235798124)},
			wantStr: "298398235798124",
		},
		{
			name:    "bool",
			args:    args{v: true},
			wantStr: "true",
		},
		{
			name:    "bool",
			args:    args{v: false},
			wantStr: "false",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotStr := Any2string(tt.args.v); gotStr != tt.wantStr {
					t.Errorf("Any2string() = %v, want %v", gotStr, tt.wantStr)
				}
			},
		)
	}
}
