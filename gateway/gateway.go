package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	nHttp "net/http"
	"sherlock/client"
	"sherlock/log"
	"sherlock/util/encrypt"
	"strings"
	"time"
)

type (
	Gateway interface {
		Name() string
		Address() string
		Init(client.Client) error
		Run() error // Run function have to block the thread
		Destroy() error
	}

	http struct {
		address string
		engine  *gin.Engine
	}

	webSocket struct {
		address  string
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
	return &http{address: address}
}

func NewWebSocketGateway(address string) Gateway {
	return &webSocket{address: address}
}

func (h *http) Name() string    { return HTTPGatewayName }
func (h *http) Address() string { return h.address }
func (h *http) Init(c client.Client) error {
	h.engine = gin.New()
	h.engine.POST("/:module/:path", func(context *gin.Context) {
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

		response, err := c.Request(message.Subject, "", message.Data, time.Duration(12)*time.Second)
		if err != nil {
			log.ErrorF("HTTP request to [%s] subject error : %s", message.Subject, err.Error())
			context.String(nHttp.StatusInternalServerError, err.Error())
			return
		}

		context.String(nHttp.StatusOK, string(response.Data))
	})
	return nil
}
func (h *http) Run() error { return h.engine.Run(h.address) }
func (h *http) Destroy() error {
	h.engine = nil
	return nil
}

func (ws *webSocket) Name() string    { return WebSocketGatewayName }
func (ws *webSocket) Address() string { return ws.address }
func (ws *webSocket) Init(c client.Client) error {
	ws.engine = gin.New()
	ws.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *nHttp.Request) bool {
			return true
		},
	}
	ws.engine.GET("/ws", func(context *gin.Context) {
		conn, err := ws.upgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.FatalF("upgrade http connection to WebSocket error : %s", err.Error())
			return
		}
		sp, err := c.Subscribe(ws.connSubject(conn.RemoteAddr().String()), "", func(msg *nats.Msg) {
			if err := conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				log.ErrorF("Write message to [%s] error : %s", conn.RemoteAddr().String(), err.Error())
			}
		})
		if err != nil {
			log.ErrorF("Subscribe [%s] according to the remote address [%s] error : %s", "WS_"+encrypt.MD5(conn.RemoteAddr().String()), conn.RemoteAddr().String())
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

			message := &Message{}
			if err := json.Unmarshal(data, message); err != nil {
				log.ErrorF("WebSocket connection [%s] read message error : %s", conn.RemoteAddr().String(), err.Error())
				continue
			}

			log.DebugF("WebSocket get new message from [%s] : %s", context.Request.RemoteAddr, message.Subject)
			log.DebugF("WebSocket get new message : %s", string(message.Data))

			if err := c.Publish(message.Subject, ws.connSubject(conn.RemoteAddr().String()), message.Data); err != nil {
				log.ErrorF("WebSocket publish a message to [%s] subject error : %s", message.Subject, err.Error())
				continue
			}
		}
	})
	return nil
}
func (ws *webSocket) Run() error     { return ws.engine.Run(ws.address) }
func (ws *webSocket) Destroy() error { return nil }
func (ws *webSocket) connSubject(address string) string {
	return fmt.Sprintf("WS_CONN.%s", encrypt.MD5(address))
}
