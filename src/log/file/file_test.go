// Package file
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:32
//
package file

import (
	"io"
	"reflect"
	"testing"
)

func TestInitLogFile(t *testing.T) {
	tests := []struct {
		name string
		want io.Writer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := InitLogFile(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("InitLogFile() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
