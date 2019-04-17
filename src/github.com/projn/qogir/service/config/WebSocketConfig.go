package config

import "time"

type WebSocketConfig struct {
	_prefix              string        `prefix:"webSocket"`
	WebsocketContextPath string        `val:"websocket.contextPath"`
	HandshakeTimeout     time.Duration `val:"handshakeTimeout"`
	WriteTimeout         time.Duration `val:"writeTimeout"`
	ReadTimeout          time.Duration `val:"readTimeout"`
	MaxMessageSize       int64         `val:"maxMessageSize"`
	BinaryMessages       bool          `val:"binaryMessages"`
	ReadBufferSize       int           `val:"readBufferSize"`
	WriteBufferSize      int           `val:"writeBufferSize"`
}
