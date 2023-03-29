// Package log
// @Description
// @Author root_wang
// @Date 2022/11/20 16:20
package log

import (
	"cqhttp-client/src/log/file"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

var (
	errorLog = log.New(
		io.MultiWriter(os.Stdout, file.InitLogFile()), "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile,
	)
	infoLog = log.New(
		io.MultiWriter(os.Stdout, file.InitLogFile()), "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile,
	)
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func ErrorInside(e string) error {
	_, f, line, ok := runtime.Caller(1)
	if ok {
		short := f
		for i := len(f) - 1; i > 0; i-- {
			if f[i] == '/' {
				short = f[i+1:]
				break
			}
		}
		f = short
	} else {
		f = "??"
		line = 0
	}
	return fmt.Errorf("%s:%d: %s", f, line, e)
}

func ErrorInsidef(format string, v ...any) error {
	_, f, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("runtime.Caller() failed")
	}
	short := f
	for i := len(f) - 1; i > 0; i-- {
		if f[i] == '/' {
			short = f[i+1:]
			break
		}
	}
	f = short
	prefix := []any{f, line}
	prefix = append(prefix, v...)
	return fmt.Errorf("%s:%d: "+format, prefix...)
}
