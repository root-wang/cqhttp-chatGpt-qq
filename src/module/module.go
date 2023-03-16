// Package module
// @Description
// @Author root_wang
// @Date 2022/12/11 14:37
package module

type Moduler interface {
	HandlerMessage(interface{}) (string, error)
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
