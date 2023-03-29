// Package utils
// @Description
// @Author root_wang
// @Date 2022/12/11 13:42
package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func Any2string(v interface{}) (str string) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		str = v.(string)
	case reflect.Int64:
		str = strconv.FormatInt(v.(int64), 10)
	case reflect.Bool:
		if v.(bool) == true {
			str = "true"
		} else {
			str = "false"
		}
	default:
		panic(fmt.Sprintf("invalid type: %s 2 string failed", rv.Kind().String()))
	}
	return
}
