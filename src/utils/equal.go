// Package utils
// @Description compare some build-in type with constant or other type
// @Author root_wang
// @Date 2022/12/11 20:49
package utils

import "fmt"

// StringEqual use for compare the string with constant string
func StringEqual(s1 string, s2 fmt.Stringer) bool {
	return s1 == fmt.Sprint(s2)
}
