package zb

import (
	"context"
	"github.com/whimthen/temp/logger"
	"github.com/whimthen/temp/websocket"
	"io"
	"testing"
)

type LoggedProcessor struct {
	logger *logger.Logger
}

func (p *LoggedProcessor) OnReceive(frame *websocket.Frame) {
	buf, err := io.ReadAll(frame.Reader)
	if err != nil {
		p.logger.Errorf("读取响应失败: %v", err)
		return
	}

	p.logger.Infof("Type: %d, Msg: %s", frame.Type, buf)
}

func (p *LoggedProcessor) SetLogger(l *logger.Logger) {
	p.logger = l
}

func TestKLineWebsocket(t *testing.T) {
	const url = "wss://kline.bw6.com/websocket"
	ctx := context.Background()
	client := websocket.NewClient(ctx, url, &LoggedProcessor{},
		websocket.WithPing(websocket.NewStringMessage("ping")))
	err := client.Connect()
	if err != nil {
		return
	}

	<-make(chan struct{})
}

func TestWSAPI(t *testing.T) {
	const url = "wss://api.bw6.com/websocket"
	ctx := context.Background()
	client := websocket.NewClient(ctx, url, &LoggedProcessor{},
		websocket.WithPing(websocket.NewStringMessage("ping")))
	err := client.Connect()
	if err != nil {
		return
	}

	<-make(chan struct{})
}
