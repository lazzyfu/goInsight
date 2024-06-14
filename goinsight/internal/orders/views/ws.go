package views

import (
	"goInsight/global"
	"goInsight/pkg/utils"
	"net/http"

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
	// 升级HTTP请求为WebSocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// 监听客户端是否断开连接
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				global.App.Log.Error(err.Error())
				return
			}
		}
	}()
	msg := make(chan string)
	// 订阅Redis频道
	sub := global.App.Redis.Subscribe(c.Request.Context(), channel)
	defer sub.Close()
	// 等待读取消息
	go utils.Subscribe(sub, c.Request.Context(), channel, msg)
	// 循环
	for {
		select {
		case <-done:
			global.App.Log.Info("websocket客户端断开连接")
			return
		case t := <-msg:
			// 将读取到的消息写入到ws，推送给客户端
			err = ws.WriteMessage(1, []byte(t))
			if err != nil {
				global.App.Log.Error("Error when writing message: ", err)
				continue
			}
			go utils.Subscribe(sub, c.Request.Context(), channel, msg)
		}
	}
}
