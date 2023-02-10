package websocket

import (
	"github.com/whimthen/temp/logger"
	"io"
)

type Frame struct {
	Type   int
	Reader io.Reader
}

type IMessage interface {
	IsSubscribe() bool
	ToUnsubscribe() IMessage
	IsPing() bool
}

func NewStringMessage(msg string) *StringMessage {
	sm := new(StringMessage)
	*sm = StringMessage(msg)
	return sm
}

type JsonMessage struct{}
type StringMessage string

func (j *JsonMessage) IsPing() bool              { return false }
func (j *JsonMessage) IsSubscribe() bool         { return true }
func (j *JsonMessage) ToUnsubscribe() IMessage   { return j }
func (s *StringMessage) IsPing() bool            { return true }
func (s *StringMessage) IsSubscribe() bool       { return false }
func (s *StringMessage) ToUnsubscribe() IMessage { return s }

type IWebsocket interface {
	Connect() error
	Shutdown() error
	SendMessage(message IMessage) error
}

type IWebsocketProcessor interface {
	OnReceive(frame *Frame)
	SetLogger(l *logger.Logger)
}
