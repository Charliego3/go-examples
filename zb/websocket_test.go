package zb

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/charliego93/websocket"
	logger "github.com/charmbracelet/log"
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

func TestLogger(t *testing.T) {
	log := logger.NewWithOptions(os.Stdout, logger.Options{
		ReportTimestamp: true,
		// ReportCaller:    true,
		TimeFormat: "3:04:05PM",
		Prefix:     "Baking 👀",
	})

	var group sync.WaitGroup
	group.Add(100)
	go func(log *logger.Logger) {
		for i := 0; i < 100; i++ {
			log.Info(fmt.Sprintf("%d/100...", i+1))
			group.Done()
		}
	}(log)

	group.Wait()
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
