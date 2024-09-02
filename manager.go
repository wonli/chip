package chip

import (
	"sync"
)

type ActionManager struct {
	handlerMap map[string]HandlersChain
}

var msy sync.Once
var manager *ActionManager

func InitManager() *ActionManager {
	msy.Do(func() {
		manager = &ActionManager{
			handlerMap: map[string]HandlersChain{},
		}
	})

	return manager
}

func (m *ActionManager) Add(name string, router HandlersChain) {
	m.handlerMap[name] = router
}

func (m *ActionManager) Has(name string) bool {
	_, ok := m.handlerMap[name]
	return ok
}

func (m *ActionManager) Handlers(name string) HandlersChain {
	return m.handlerMap[name]
}
