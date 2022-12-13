// Package file
// @Description
// @Author root_wang
// @Date 2022/12/13 19:29
package file

import (
	"io"
	"os"
	"path"
	"time"
)

func InitLogFile() io.Writer {
	file := path.Join(time.Now().Format("2006y01m02d15h04m") + ".log")
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return logFile
}
