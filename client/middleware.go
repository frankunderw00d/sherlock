package client

import "github.com/nats-io/nats.go"

type (
	HandleFunc func(*nats.Msg) bool

	Middleware interface {
		// 加入中间件函数到调用链路中
		Use(HandleFunc)
		// 链路出最终执行函数
		End(nats.MsgHandler) nats.MsgHandler
		// 派生新的中间件
		Derive() Middleware
	}

	middleware struct {
		root  HandleFunc
		maker func(HandleFunc, HandleFunc) HandleFunc
	}
)

const ()

var ()

func NewMiddleware() Middleware {
	return &middleware{
		root: func(msg *nats.Msg) bool { return true },
		maker: func(front HandleFunc, behind HandleFunc) HandleFunc {
			return func(msg *nats.Msg) bool {
				if !front(msg) {
					return false
				}
				return behind(msg)
			}
		},
	}
}

func (mw *middleware) Use(handler HandleFunc) {
	mw.root = mw.maker(mw.root, handler)
}

func (mw *middleware) End(handler nats.MsgHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		if !mw.root(msg) {
			return
		}
		handler(msg)
	}
}

func (mw *middleware) Derive() Middleware {
	return &middleware{
		root: mw.root,
		maker: func(front HandleFunc, behind HandleFunc) HandleFunc {
			return func(msg *nats.Msg) bool {
				if !front(msg) {
					return false
				}
				return behind(msg)
			}
		},
	}
}
