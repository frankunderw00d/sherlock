package client

import (
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type (
	Client interface {
		// 关闭
		Close()

		// 订阅
		// subject , queue , handler
		Subscribe(string, string, nats.MsgHandler, ...HandleFunc) (*nats.Subscription, error)

		// 发布
		// subject , reply , data
		Publish(string, string, []byte) error

		// 同步请求
		// subject , reply , data , timeout
		// 当 reply 为空时，由内部自动生成一个回复地址
		Request(string, string, []byte, time.Duration) (*nats.Msg, error)

		// 响应(其实就是订阅)
		// subject , queue , handler
		Response(string, string, nats.MsgHandler, ...HandleFunc) (*nats.Subscription, error)

		// 回复(其实就是发布)
		// subject , reply , data
		Reply(string, string, []byte) error

		// 加入中间件
		UseMiddleware(HandleFunc)
	}

	client struct {
		conn         *nats.Conn
		rmw          nats.MsgHandler
		handlerMaker func(nats.MsgHandler, nats.MsgHandler) nats.MsgHandler
		mw           Middleware
	}
)

const (
	DefaultMaxReconnects = 3  // 默认重联3次
	DefaultTimeout       = 10 // 默认客户端超时时间10分钟
	DefaultReconnectWait = 1  // 默认重联等待间隔时间1秒
)

var ()

func init() {}

func NewClient(name, address, token string) Client {
	opts := []nats.Option{
		nats.Name(name),
		nats.MaxReconnects(DefaultMaxReconnects),
		nats.Timeout(time.Minute * time.Duration(DefaultTimeout)),
		nats.ReconnectWait(DefaultReconnectWait * time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Println("nats client reconnected")
		}),
		nats.Token(token),
	}

	conn, err := nats.Connect(address, opts...)
	if err != nil {
		log.Printf("nats connect error : %s", err.Error())
		return nil
	}

	return &client{
		conn: conn,
		rmw:  func(msg *nats.Msg) {},
		handlerMaker: func(front nats.MsgHandler, back nats.MsgHandler) nats.MsgHandler {
			return func(msg *nats.Msg) {
				front(msg)
				back(msg)
			}
		},
		mw: NewMiddleware(),
	}
}

func (c *client) Close() {
	c.conn.Close()
}

func (c *client) Subscribe(subject, queue string, handler nats.MsgHandler, middleware ...HandleFunc) (*nats.Subscription, error) {
	// 派生新的中间件
	nmw := c.mw.Derive()

	// 加入特设中间件
	for _, mw := range middleware {
		nmw.Use(mw)
	}

	if queue == "" {
		//								  链路出最终执行函数
		return c.conn.Subscribe(subject, nmw.End(handler))
	}

	return c.conn.QueueSubscribe(subject, queue, nmw.End(handler))
}

func (c *client) Publish(subject, reply string, data []byte) error {
	return c.conn.PublishMsg(&nats.Msg{
		Subject: subject,
		Reply:   reply,
		Data:    data,
	})
}

func (c *client) Request(subject, reply string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	if reply == "" {
		return c.conn.Request(subject, data, timeout)
	}

	return c.conn.RequestMsg(&nats.Msg{
		Subject: subject,
		Reply:   reply,
		Data:    data,
	}, timeout)
}

func (c *client) Response(subject, queue string, handler nats.MsgHandler, middleware ...HandleFunc) (*nats.Subscription, error) {
	return c.Subscribe(subject, queue, handler, middleware...)
}

func (c *client) Reply(subject, reply string, data []byte) error {
	return c.Publish(subject, reply, data)
}

func (c *client) UseMiddleware(mw HandleFunc) {
	c.mw.Use(mw)
}
