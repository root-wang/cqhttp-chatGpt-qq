// Package log
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:31
//
package log

import "testing"


func TestErrorInside(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := ErrorInside(tt.args.e); (err != nil) != tt.wantErr {
					t.Errorf("ErrorInside() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestErrorInsidef(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := ErrorInsidef(tt.args.format, tt.args.v...); (err != nil) != tt.wantErr {
					t.Errorf("ErrorInsidef() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}