package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"net"
	nHttp "net/http"
	"sherlock/client"
	"sherlock/log"
	"sherlock/util/encrypt"
	"strings"
	"time"
)

type (
	Gateway interface {
		// 服务信息
		Info() string
		// 初始化
		Init(client.Client) error
		// 运行 （必须阻塞该执行线程）
		Run() error
		// 销毁
		Destroy() error

		// 主动关闭
		Close() error
	}

	baseGateway struct {
		address       string
		server        *nHttp.Server
		subscriptions []*nats.Subscription
		bl            BlackList
	}

	http struct {
		baseGateway
		engine *gin.Engine
	}

	webSocket struct {
		baseGateway
		engine   *gin.Engine
		upgrader *websocket.Upgrader
		conn     *websocket.Conn
	}

	Message struct {
		Subject string `json:"subject"`
		Data    []byte `json:"data"`
	}
)

const (
	HTTPGatewayName      = "HTTP-GATEWAY"
	WebSocketGatewayName = "WEBSOCKET-GATEWAY"
)

var ()

func init() { gin.SetMode(gin.ReleaseMode) }

func NewHTTPGateway(address string) Gateway {
	return &http{
		baseGateway: baseGateway{
			address:       address,
			server:        nil,
			subscriptions: []*nats.Subscription{},
			bl:            NewBlackList(),
		},
	}
}

func NewWebSocketGateway(address string) Gateway {
	return &webSocket{
		baseGateway: baseGateway{
			address:       address,
			server:        nil,
			subscriptions: []*nats.Subscription{},
			bl:            NewBlackList(),
		},
	}
}

func (h *http) Close() error { return h.server.Shutdown(context.Background()) }
func (h *http) Info() string { return HTTPGatewayName }
func (h *http) Init(c client.Client) error {
	// 网关订阅黑名单开关
	if sp, err := c.Subscribe(BlackListSwitchSubject, "", h.bl.OnlineSwitch); err != nil {
		log.ErrorF("HTTP gateway subscribe [%s] error : %s", BlackListSwitchSubject, err.Error())
	} else {
		log.DebugF("HTTP gateway subscribe [%s] success", sp.Subject)
		h.subscriptions = append(h.subscriptions, sp)
	}

	// 网关订阅黑名单更新
	if sp, err := c.Subscribe(BlackListUpdateSubject, "", h.bl.OnlineUpdate); err != nil {
		log.ErrorF("HTTP gateway subscribe [%s] error : %s", BlackListUpdateSubject, err.Error())
	} else {
		log.DebugF("HTTP gateway subscribe [%s] success", sp.Subject)
		h.subscriptions = append(h.subscriptions, sp)
	}

	// 初始化 HTTP 引擎
	h.engine = gin.New()
	// 添加IP黑名单中间件
	h.engine.Use(h.baseGateway.FilterIPMiddleware)
	// 只支持 POST 请求
	h.engine.POST("/:module/:path", func(context *gin.Context) {
		// 模块.路径 指定了发布主题，请求 Body 指定了数据
		module := context.Param("module")
		path := context.Param("path")
		data, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			context.String(nHttp.StatusBadRequest, "get body error :"+err.Error())
			return
		}

		message := &Message{
			Subject: strings.Join([]string{module, path}, "."),
			Data:    data,
		}

		log.DebugF("HTTP get new message from [%s] : %s", context.Request.RemoteAddr, message.Subject)
		log.DebugF("HTTP get new message : %s", string(message.Data))

		// 通过请求的方式发布及接收响应
		response, err := c.Request(message.Subject, "", message.Data, time.Duration(12)*time.Second)
		if err != nil {
			log.ErrorF("HTTP request to [%s] subject error : %s", message.Subject, err.Error())
			context.String(nHttp.StatusInternalServerError, err.Error())
			return
		}

		context.String(nHttp.StatusOK, string(response.Data))
	})
	// 初始化服务
	h.server = &nHttp.Server{
		Addr:    h.address,
		Handler: h.engine,
	}
	return nil
}
func (h *http) Run() error {
	if err := h.server.ListenAndServe(); err != nil {
		if err == nHttp.ErrServerClosed { // 主动关闭
			return nil
		} else {
			return err
		}
	}
	return nil
}
func (h *http) Destroy() error {
	// 取消订阅
	for _, sp := range h.subscriptions {
		if err := sp.Unsubscribe(); err != nil {
			return err
		} else {
			log.DebugF("HTTP gateway unsubscribe [%s] success", sp.Subject)
		}
	}

	h.engine = nil
	return nil
}

func (ws *webSocket) Close() error { return ws.server.Shutdown(context.Background()) }
func (ws *webSocket) Info() string { return WebSocketGatewayName }
func (ws *webSocket) Init(c client.Client) error {
	// 网关订阅黑名单开关
	if sp, err := c.Subscribe(BlackListSwitchSubject, "", ws.bl.OnlineSwitch); err != nil {
		log.ErrorF("WebSocket gateway subscribe [%s] error : %s", BlackListSwitchSubject, err.Error())
	} else {
		ws.subscriptions = append(ws.subscriptions, sp)
		log.DebugF("WebSocket gateway subscribe [%s] success", sp.Subject)
	}

	// 网关订阅黑名单更新
	if sp, err := c.Subscribe(BlackListUpdateSubject, "", ws.bl.OnlineUpdate); err != nil {
		log.ErrorF("WebSocket gateway subscribe [%s] error : %s", BlackListUpdateSubject, err.Error())
	} else {
		ws.subscriptions = append(ws.subscriptions, sp)
		log.DebugF("WebSocket gateway subscribe [%s] success", sp.Subject)
	}

	// 初始化 HTTP 引擎
	ws.engine = gin.New()
	// 添加IP黑名单中间件
	ws.engine.Use(ws.baseGateway.FilterIPMiddleware)
	// 初始化 WebSocket 升级件
	ws.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *nHttp.Request) bool {
			return true
		},
	}
	ws.engine.GET("/ws", func(context *gin.Context) {
		// 对 [ws://address/ws] 路径上的请求统一升级为长连接
		conn, err := ws.upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.FatalF("upgrade http connection to WebSocket error : %s", err.Error())
			return
		}
		// 同时开启一个 [WS_CONN.远程地址摘要] 主题的订阅，用于接收回复消息
		sp, err := c.Subscribe(ws.connSubject(conn.RemoteAddr().String()), "", func(msg *nats.Msg) {
			if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				log.ErrorF("Write message to [%s] error : %s", conn.RemoteAddr().String(), err.Error())
			}
		})
		if err != nil {
			log.ErrorF("Subscribe [%s] according to the remote address [%s] error : %s", "WS_"+encrypt.MD5(conn.RemoteAddr().String()), conn.RemoteAddr().String())
			// 如果订阅失败，直接关闭连接
			if err := conn.Close(); err != nil {
				log.ErrorF("WebSocket connection [%s] close error : %s", conn.RemoteAddr().String(), err.Error())
			}
			return
		}

		defer func() {
			if err := conn.Close(); err != nil {
				log.ErrorF("WebSocket connection [%s] close error : %s", conn.RemoteAddr().String(), err.Error())
			}

			if err := sp.Unsubscribe(); err != nil {
				log.ErrorF("WebSocket connection [%s] close error : %s", conn.RemoteAddr().String(), err.Error())
			}
		}()

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.ErrorF("WebSocket connection [%s] read message error : %s", conn.RemoteAddr().String(), err.Error())
				break
			}
			// 接收的消息必须以 Message 的形式指定
			message := &Message{}
			if err := json.Unmarshal(data, message); err != nil {
				log.ErrorF("WebSocket connection [%s] read message error : %s", conn.RemoteAddr().String(), err.Error())
				continue
			}

			log.DebugF("WebSocket get new message from [%s] : %s", context.Request.RemoteAddr, message.Subject)
			log.DebugF("WebSocket get new message : %s", string(message.Data))

			// 通过指定 Reply 为 [WS_CONN.远程地址摘要] ，由上面的订阅接收并且回复给用户
			if err := c.Publish(message.Subject, ws.connSubject(conn.RemoteAddr().String()), message.Data); err != nil {
				log.ErrorF("WebSocket publish a message to [%s] subject error : %s", message.Subject, err.Error())
				continue
			}
		}
	})
	// 初始化服务
	ws.server = &nHttp.Server{
		Addr:    ws.address,
		Handler: ws.engine,
	}
	return nil
}
func (ws *webSocket) Run() error {
	if err := ws.server.ListenAndServe(); err != nil {
		if err == nHttp.ErrServerClosed { // 主动关闭
			return nil
		} else {
			return err
		}
	}
	return nil
}
func (ws *webSocket) Destroy() error {
	// 取消订阅
	for _, sp := range ws.subscriptions {
		if err := sp.Unsubscribe(); err != nil {
			return err
		} else {
			log.DebugF("WebSocket gateway unsubscribe [%s] success", sp.Subject)
		}
	}

	ws.engine = nil
	return nil
}
func (ws *webSocket) connSubject(address string) string {
	return fmt.Sprintf("WS_CONN.%s", encrypt.MD5(address))
}

// 黑名单过滤
func (bg *baseGateway) FilterIPMiddleware(context *gin.Context) {
	h, _, err := net.SplitHostPort(context.Request.Host)
	if err != nil {
		log.ErrorF("Parse connection host [%s] error : %s", context.Request.Host, err.Error())
		context.Abort()
	}
	log.DebugF("Split host : %s", h)

	if !bg.bl.Filter(h) {
		context.String(nHttp.StatusBadRequest, "You are block by blacklist")
		context.Abort()
	}
}
