package views

import (
	"context"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 创建WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	channel := c.Param("channel")
	requestID := requestid.Get(c)

	global.App.Log.Infof("WebSocket connected: channel=%s, request_id=%s", channel, requestID)

	// 升级为 WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.App.Log.Errorf("WebSocket upgrade failed: channel=%s, request_id=%s, error=%v", channel, requestID, err)
		return
	}
	defer func() {
		ws.Close()
		global.App.Log.Infof("WebSocket closed: channel=%s, request_id=%s", channel, requestID)
	}()

	// 创建可取消的上下文，用于管理 goroutine 生命周期
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// 监听客户端是否断开（读协程）
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				global.App.Log.Infof("WebSocket read message failed: channel=%s, request_id=%s, error=%v", channel, requestID, err)
				return
			}

			// 心跳处理
			if msgType == websocket.TextMessage && string(msg) == "ping" {
				ws.WriteMessage(websocket.TextMessage, []byte("pong"))
				global.App.Log.Infof("WebSocket heartbeat pong: channel=%s, request_id=%s", channel, requestID)
				continue
			}
		}
	}()

	// Redis 订阅
	sub := global.App.Redis.Subscribe(ctx, channel)
	if sub == nil {

		global.App.Log.Errorf("WebSocket Redis subscribe failed: channel=%s, request_id=%s", channel, requestID)
		return
	}
	defer func() {
		if err := sub.Close(); err != nil {
			global.App.Log.Errorf("WebSocket Redis subscribe close failed: channel=%s, request_id=%s, error=%v", channel, requestID, err)
		}
	}()

	global.App.Log.Infof("WebSocket Redis subscribe successful: channel=%s, request_id=%s", channel, requestID)

	// 缓冲消息，避免阻塞
	msgChan := make(chan string, 100)

	// Redis 消息读取 goroutine
	go func() {
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				global.App.Log.Infof("WebSocket Redis context done: channel=%s, request_id=%s", channel, requestID)
				return
			case m, ok := <-ch:
				if !ok {
					global.App.Log.Infof("WebSocket Redis channel closed: channel=%s, request_id=%s", channel, requestID)
					return
				}
				msgChan <- m.Payload
			}
		}
	}()

	// 主循环：把 redis 消息推送给 websocket
	for {
		select {
		case <-done:
			// websocket 被关闭
			global.App.Log.Infof("WebSocket disconnected: channel=%s, request_id=%s", channel, requestID)
			return

		case m := <-msgChan:
			if err := ws.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
				// 写入失败 = 连接断开
				global.App.Log.Infof("WebSocket write message failed: channel=%s, request_id=%s, error=%v", channel, requestID, err)
				return
			}
		}
	}
}
