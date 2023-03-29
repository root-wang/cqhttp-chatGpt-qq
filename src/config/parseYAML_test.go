// Package cfg
// @Description 
// @Author root_wang
// @Date 2023/3/22 12:34
//
package cfg

import (
	"reflect"
	"testing"
)

func TestGetYAMLFiled(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name            string
		args            args
		wantFieldValMap map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if gotFieldValMap := GetYAMLFiled(tt.args.fields...); !reflect.DeepEqual(
					gotFieldValMap, tt.wantFieldValMap,
				) {
					t.Errorf("GetYAMLFiled() = %v, want %v", gotFieldValMap, tt.wantFieldValMap)
				}
			},
		)
	}
}
