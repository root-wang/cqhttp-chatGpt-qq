// Package log
// @Description
// @Author root_wang
// @Date 2022/11/20 16:20
package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func ErrorInside(e string) error {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	return fmt.Errorf("%s:%d: %s", file, line, e)
}

func ErrorInsidef(format string, v ...any) error {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	prefix := []any{file, line}
	prefix = append(prefix, v...)
	return fmt.Errorf("%s:%d: "+format, prefix...)
}
