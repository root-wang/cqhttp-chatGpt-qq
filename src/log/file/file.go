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

const Prefix = "--------------------------------------------------------------------------------------\n" +
	"                       _oo0oo_\n" +
	"                      o8888888o\n" +
	"                      88\" . \"88\n" +
	"                      (| -_- |)\n" +
	"                      0\\  =  /0\n" +
	"                    ___/`---'\\___\n" +
	"                  .' \\\\|     |// '.\n" +
	"                 / \\\\|||  :  |||// \\\n" +
	"                / _||||| -:- |||||- \\\n" +
	"               |   | \\\\\\  - /// |   |\n" +
	"               | \\_|  ''\\---/''  |_/ |\n" +
	"               \\  .-\\__  '-'  ___/-. /\n" +
	"             ___'. .'  /--.--\\  `. .'___\n" +
	"          .\"\" '<  `.___\\_<|>_/___.' >' \"\".\n" +
	"         | | :  `- \\`.;`\\ _ /`;.`/ - ` : | |\n" +
	"         \\  \\ `_.   \\_ __\\ /__ _/   .-` /  /\n" +
	"     =====`-.____`.___ \\_____/___.-`___.-'=====\n" +
	"                       `=---='\n" +
	"\n" +
	"\n" +
	"     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n" +
	"\n" +
	"           佛祖保佑     永不宕机     永无BUG\n" +
	"--------------------------------------------------------------------------------------\n"

func InitLogFile() io.Writer {
	file := path.Join(time.Now().Format("2006y01m02d15h") + ".log")
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return logFile
}
