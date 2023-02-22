package zb

import (
	"context"
	"github.com/whimthen/temp/logger"
	"github.com/whimthen/temp/websocket"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"testing"
)

const (
	EventAddChannel = "addChannel"

	ChannelIncrRecord = "push_user_incr_record"
)

type User struct {
	Name      string `yaml:"name"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`

	ch chan struct{}
}

type Config struct {
	ApiURL   string   `yaml:"apiURL"`
	TradeRUL string   `yaml:"tradeRUL"`
	KlineURL string   `yaml:"klineURL"`
	WsapiURL string   `yaml:"wsapiURL"`
	Markets  []string `yaml:"markets"`
	Users    []User   `yaml:"accounts"`
}

type TradeProcessor struct {
	logger *logger.Logger
	user   User
}

func (p *TradeProcessor) OnReceive(frame *websocket.Frame) {
	buf, err := io.ReadAll(frame.Reader)
	if err != nil {
		p.logger.Errorf("读取响应失败: %v", err)
		return
	}

	p.logger.Infof("Recived Type: %d, Msg: %s", frame.Type, buf)
}

func (p *TradeProcessor) SetLogger(l *logger.Logger) {
	p.logger = l
}

type Websocket struct {
	*websocket.Client
	user User
	err  error
}

func NewWebsocket(ctx context.Context, user User) *Websocket {
	client := websocket.NewClient(
		ctx, config.WsapiURL, &TradeProcessor{user: user},
		websocket.WithPing(websocket.NewStringMessage("ping")),
		websocket.WithPrefix(user.Name),
	)
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	return &Websocket{
		Client: client,
		user:   user,
	}
}

func (w *Websocket) SubscribeRecord(markets ...string) {
	if w.err != nil {
		return
	}

	for _, market := range markets {
		market := market
		w.err = w.Client.SendMessage((&websocket.RequestMsg{
			Event:     EventAddChannel,
			Channel:   ChannelIncrRecord,
			AccessKey: w.user.AccessKey,
			Market:    &market,
		}).Signed(w.user.SecretKey))
		if w.err != nil {
			panic(w.err)
		}
	}
}

var config Config

func init() {
	logger.SetFormatter(&logger.Formatter{})
	buf, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}
}

func TestAutoTrade(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	for _, user := range config.Users {
		user.ch = make(chan struct{})

		NewWebsocket(ctx, user).
			SubscribeRecord(config.Markets...)
	}

	select {
	case <-ctx.Done():
		cancel()
	}
}
