// Package utils
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:33
//
package utils

import (
	"fmt"
	"testing"
)

type testStr string

func (t testStr) String() string {
	return string(t)
}

func TestStringEqual(t *testing.T) {
	type args struct {
		s1 string
		s2 fmt.Stringer
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{
				s1: "hello",
				s2: testStr("hello"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := StringEqual(tt.args.s1, tt.args.s2); got != tt.want {
					t.Errorf("StringEqual() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
