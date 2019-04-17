package _struct

import "github.com/kataras/iris/websocket"

type WebSocketSessionInfo struct {
	session             websocket.Connection
	CreateTime          int64
	LastReceiveDataTime int64
	LastSendDataTime    int64
}
