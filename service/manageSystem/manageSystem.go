package manageSystem

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"sherlock/client"
	"sherlock/log"
	"strings"
)

type (
	// 管理系统定义
	ManageSystem interface {
		// 服务信息
		Info() string
		// 初始化
		Init(client.Client) error
		// 运行 （必须阻塞该执行线程）
		Run() error
		// 销毁
		Destroy() error

		// 注册路由
		RegisterRoute(subject string, handler nats.MsgHandler, middleware ...client.HandleFunc) error
		// 主动通知关闭
		Close()
		// 使用全局中间件
		UseMiddleware(middlewareList ...client.HandleFunc) error
	}

	// 管理系统路由项
	manageSystemRoute struct {
		subject    string              // 主题
		queue      string              // 组
		handler    nats.MsgHandler     // 处理函数
		middleware []client.HandleFunc // 中间件
	}

	// 管理系统实现
	manageSystem struct {
		name          string                       // 名称
		version       string                       // 版本
		routes        map[string]manageSystemRoute // 路由组 map[subject]manageSystemRoute
		subscriptions []*nats.Subscription         // 订阅记录
		closeChan     chan struct{}                // 关闭通知通道
		middleware    []client.HandleFunc          // 中间件组
	}
)

const (
	QueuePrefix = "ManageSystem"
)

var ()

func init() {}

func NewManageSystem(name, version string) ManageSystem {
	return &manageSystem{
		name:          name,
		version:       version,
		routes:        map[string]manageSystemRoute{},
		subscriptions: []*nats.Subscription{},
		closeChan:     make(chan struct{}),
		middleware:    []client.HandleFunc{},
	}
}

// 注册路由
func (ms *manageSystem) RegisterRoute(subject string, handler nats.MsgHandler, middleware ...client.HandleFunc) error {
	if subject == "" {
		return errors.New("subject can't be nil")
	}
	if handler == nil {
		return errors.New("handler can't be nil")
	}
	if _, exit := ms.routes[ms.SubscribeSubject(subject)]; exit {
		return errors.New("register route exist")
	}

	ms.routes[ms.SubscribeSubject(subject)] = manageSystemRoute{
		subject:    ms.SubscribeSubject(subject),
		queue:      ms.SubscribeQueue(),
		handler:    handler,
		middleware: append(ms.middleware, middleware...),
	}

	return nil
}

// 主动通知关闭
func (ms *manageSystem) Close() {
	if ms.closeChan != nil {
		ms.closeChan <- struct{}{}
	}
}

// 使用全局中间件
func (ms *manageSystem) UseMiddleware(middlewareList ...client.HandleFunc) error {
	if middlewareList == nil || len(middlewareList) == 0 {
		return errors.New("middleware function list can't be nil or empty")
	}

	ms.middleware = append(ms.middleware, middlewareList...)

	return nil
}

// 服务信息
func (ms *manageSystem) Info() string {
	return fmt.Sprintf("%s %s", ms.name, ms.version)
}

// 初始化
func (ms *manageSystem) Init(c client.Client) error {
	// 注册所有订阅
	for _, route := range ms.routes {
		if sp, err := c.Response(route.subject, route.queue, route.handler, route.middleware...); err != nil {
			return err
		} else {
			ms.subscriptions = append(ms.subscriptions, sp)
			log.DebugF("Subscribe [%s]-[%s] success", sp.Subject, sp.Queue)
		}
	}

	return nil
}

// 运行 （必须阻塞该执行线程）
func (ms *manageSystem) Run() error {
	select {
	case <-ms.closeChan:
		log.DebugF("%s close", ms.Info())
	}

	return nil
}

// 销毁
func (ms *manageSystem) Destroy() error {
	// 取消所有订阅
	for _, sp := range ms.subscriptions {
		if err := sp.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

// 构建订阅主题
func (ms *manageSystem) SubscribeSubject(subject string) string {
	return strings.Join([]string{ms.SubscribeQueue(), subject}, ".")
}

// 构建订阅组 ManageSystem.Boss
func (ms *manageSystem) SubscribeQueue() string {
	return strings.Join([]string{QueuePrefix, ms.name}, ".")
}
