// Package websocket 通过 github.com/olahol/melody 实现 websocket
package websocket

import (
	"net/http"

	"github.com/olahol/melody"
)

var (
	// Server websocket 服务器
	Server *melody.Melody

	// ConnectHandler 连接 websocket 时的处理器
	ConnectHandler http.HandlerFunc
)

func init() {
	Server = melody.New()
	ConnectHandler = Server.HandleRequest
}
