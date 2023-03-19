package event_default

import (
	"errors"
	"sync"

	"github.com/infrago/event"
	"github.com/infrago/util"
)

// mo

var (
	errRunning    = errors.New("Event is running")
	errNotRunning = errors.New("Event is not running")
)

type (
	defaultDriver  struct{}
	defaultConnect struct {
		mutex   sync.RWMutex
		running bool
		actives int64

		instance *event.Instance

		runner *util.Runner
		events map[string]chan *defaultMsg
	}

	defaultMsg struct {
		name string
		data []byte
	}
)

// 连接
func (driver *defaultDriver) Connect(inst *event.Instance) (event.Connect, error) {
	return &defaultConnect{
		instance: inst, runner: util.NewRunner(),
		events: make(map[string]chan *defaultMsg, 0),
	}, nil
}

// 打开连接
func (connect *defaultConnect) Open() error {
	return nil
}
func (connect *defaultConnect) Health() (event.Health, error) {
	connect.mutex.RLock()
	defer connect.mutex.RUnlock()
	return event.Health{Workload: connect.actives}, nil
}

// 关闭连接
func (connect *defaultConnect) Close() error {
	return nil
}

func (connect *defaultConnect) Register(name, group string) error {
	connect.mutex.Lock()
	defer connect.mutex.Unlock()

	connect.events[name] = make(chan *defaultMsg, 10)

	return nil
}

// 开始订阅者
func (connect *defaultConnect) Start() error {
	if connect.running {
		return errRunning
	}

	for _, cccc := range connect.events {
		connect.runner.Run(func() {
			for {
				select {
				case msg := <-cccc:
					connect.instance.Serve(msg.name, msg.data)
				case <-connect.runner.Stop():
					return
				}
			}
		})
	}

	connect.running = true
	return nil
}

// 停止订阅
func (connect *defaultConnect) Stop() error {
	if false == connect.running {
		return errNotRunning
	}

	connect.runner.End()

	connect.running = false
	return nil
}

func (connect *defaultConnect) Notify(name string, data []byte) error {
	if qqq, ok := connect.events[name]; ok {
		qqq <- &defaultMsg{name, data}
	}
	return nil
}

//------------------------- 默认事件驱动 end --------------------------
