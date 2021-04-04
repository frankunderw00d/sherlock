package lobby

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"sherlock/client"
	"sherlock/log"
	"strings"
)

type (
	// 子游戏房间定义
	Room interface {
		// 房间ID
		RoomID() string
	}

	// 子游戏大厅定义
	Lobby interface {
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
		// 通知关闭
		Close()
		// 使用全局中间件
		UseMiddleware(middlewareList ...client.HandleFunc) error
	}

	// 子游戏大厅路由项
	lobbyRoute struct {
		subject    string              // 主题
		queue      string              // 组
		handler    nats.MsgHandler     // 处理函数
		middleware []client.HandleFunc // 中间件
	}

	// 子游戏大厅实现
	lobby struct {
		platformID    string                // 业主 ID
		gameID        string                // 游戏 ID
		name          string                // 游戏名称
		routes        map[string]lobbyRoute // 路由组	map[subject]lobbyRoute
		subscriptions []*nats.Subscription  // 订阅记录
		closeChan     chan struct{}         // 关闭通知通道
		middleware    []client.HandleFunc   // 中间件组
	}
)

const (
	QueuePrefix = "Lobby"
)

var ()

func init() {}

func NewLobby(pid, gid, name string) Lobby {
	return &lobby{
		platformID:    pid,
		gameID:        gid,
		name:          name,
		routes:        map[string]lobbyRoute{},
		subscriptions: []*nats.Subscription{},
		closeChan:     make(chan struct{}),
		middleware:    []client.HandleFunc{},
	}
}

// 注册路由
func (l *lobby) RegisterRoute(subject string, handler nats.MsgHandler, middleware ...client.HandleFunc) error {
	if subject == "" {
		return errors.New("subject can't be nil")
	}
	if handler == nil {
		return errors.New("handler can't be nil")
	}
	if _, exit := l.routes[l.SubscribeSubject(subject)]; exit {
		return errors.New("register route exist")
	}

	l.routes[l.SubscribeSubject(subject)] = lobbyRoute{
		subject:    l.SubscribeSubject(subject),
		queue:      l.SubscribeQueue(),
		handler:    handler,
		middleware: append(l.middleware, middleware...),
	}

	return nil
}

// 主动通知关闭
func (l *lobby) Close() {
	if l.closeChan != nil {
		l.closeChan <- struct{}{}
	}
}

// 使用全局中间件
func (l *lobby) UseMiddleware(middlewareList ...client.HandleFunc) error {
	if middlewareList == nil || len(middlewareList) == 0 {
		return errors.New("middleware function list can't be nil or empty")
	}

	l.middleware = append(l.middleware, middlewareList...)

	return nil
}

// 服务信息
func (l *lobby) Info() string {
	return fmt.Sprintf("%s-%s-%s", l.platformID, l.gameID, l.name)
}

// 初始化
func (l *lobby) Init(c client.Client) error {
	// 注册所有订阅
	for _, route := range l.routes {
		if sp, err := c.Subscribe(route.subject, route.queue, route.handler, route.middleware...); err != nil {
			return err
		} else {
			l.subscriptions = append(l.subscriptions, sp)
			log.DebugF("Subscribe [%s]-[%s] success", sp.Subject, sp.Queue)
		}
	}

	return nil
}

// 运行 （必须阻塞该执行线程）
func (l *lobby) Run() error {
	select {
	case <-l.closeChan:
		log.DebugF("%s close", l.Info())
	}

	return nil
}

// 销毁
func (l *lobby) Destroy() error {
	// 取消所有订阅
	for _, sp := range l.subscriptions {
		if err := sp.Unsubscribe(); err != nil {
			return err
		}
	}

	return nil
}

// 构建订阅主题
func (l *lobby) SubscribeSubject(subject string) string {
	return strings.Join([]string{l.SubscribeQueue(), subject}, ".")
}

// 构建订阅组 Lobby.PID.GID
func (l *lobby) SubscribeQueue() string {
	return strings.Join([]string{QueuePrefix, l.platformID, l.gameID}, ".")
}
