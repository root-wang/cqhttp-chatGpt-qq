// Package module
// @Description
// @Author root_wang
// @Date 2022/12/11 14:37
package module

import "cqhttp-client/src/log"

type Moduler interface {
	// HandlerMessage 每个模块需要实现对输入进行处理后的消息 处理后也就是回复用户的是一条消息
	HandlerMessage(interface{}) (string, error)
	// Matcher 每个模块实现对自己模块的识别处理
	Matcher(string) (ok bool, prompt string)
}

// Module 响应信息的模块
type Module struct {
	m map[string]Moduler
}

func NewModule() *Module {
	return &Module{m: make(map[string]Moduler)}
}

func (m *Module) Handler(name string) Moduler {
	return m.m[name]
}

func (m *Module) AddModule(name string, mer Moduler) {
	m.m[name] = mer
}

// CheckModuleType 调用每个模块实现 Moduler.Matcher 接口的方法
func (m *Module) CheckModuleType(msg string) (module string, prompt string) {
	for t, moduler := range m.m {
		var ok bool
		if ok, prompt = moduler.Matcher(msg); ok {
			module = t
			return
		}
	}
	log.Errorf("can't find appropriate module")
	return "", ""
}
