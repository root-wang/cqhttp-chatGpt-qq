// Package cfg
// @Description
// @Author root_wang
// @Date 2023/3/16 17:01
package cfg

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetYAMLFiled(fields ...string) (fieldValMap map[string]interface{}) {
	fieldValMap = make(map[string]interface{})
	for _, field := range fields {
		fieldValMap[field] = viper.Get(field)
	}
	return
}

func init() {
	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")          // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./src/config")  // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()          // 查找并读取配置文件
	if err != nil {                      // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
