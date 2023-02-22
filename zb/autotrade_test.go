package zb

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"github.com/whimthen/temp/logger"
	"github.com/whimthen/temp/websocket"
	"github.com/whimthen/temp/zb/autoapi"
	"gopkg.in/yaml.v3"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	EventAddChannel = "addChannel"

	ChannelIncrRecord = "push_user_incr_record"
)

type User struct {
	Name      string `yaml:"name"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`

	ch chan objx.Map
	ac *autoapi.Account
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

func (p *TradeProcessor) SetLogger(l *logger.Logger) {
	p.logger = l
}

func (p *TradeProcessor) OnReceive(frame *websocket.Frame) {
	buf, err := io.ReadAll(frame.Reader)
	if err != nil {
		p.logger.Errorf("读取响应失败: %v", err)
		return
	}

	content := string(buf)
	if !strings.HasPrefix(content, "{") || !strings.HasSuffix(content, "}") {
		p.logger.Warnf("收到非 JSON 类型的消息: %d, Msg: %s", frame.Type, buf)
		return
	}

	obj, err := objx.FromJSON(content)
	if err != nil {
		p.logger.Errorf("收到消息反序列化失败: %+v", err)
		return
	}

	channel := obj.Get("channel").String()
	switch channel {
	case ChannelIncrRecord:
		p.user.ch <- obj
	default:
		p.logger.Warnf("收到未处理的消息类型: %s", buf)
	}
}

func order(ctx context.Context, logger *logger.Logger, user User) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("服务已停止")
		case obj := <-user.ch:
			record := obj.Get("record").InterSlice()
			logger.Infof("Json: %+v", record)
			entrustId := record[0].(string)
			unitPrice := decimal.NewFromFloat(record[1].(float64))
			numbers := decimal.NewFromFloat(record[2].(float64))
			completeNumbers := decimal.NewFromFloat(record[3].(float64))
			types := cast.ToInt(record[5])

			if numbers.Equal(completeNumbers) {
				logger.Infof("本次成交: [%s:%d] = %s", entrustId, types, numbers)
				market := strings.TrimSuffix(obj.Get("market").String(), "default")
				var opUsers []User
				for _, u := range config.Users {
					if u.Name == user.Name {
						continue
					}

					opUsers = append(opUsers, u)
				}

				if len(opUsers) == 0 {
					logger.Errorf("没有对手用户, 使用自成交模式")
					opUsers = append(opUsers, user)
				}

				opu := opUsers[rand.Intn(len(opUsers))]
				autoapi.Order(market, numbers, unitPrice, autoapi.ReverseTradeType(types), autoapi.WithAccount(opu.ac))
				return
			}

			logger.Debugf("收到订单: ID: %s, Price: %s, Numbers: %s - %s, Types: %d",
				entrustId, unitPrice, numbers, completeNumbers, types)
		}
	}
}

type Websocket struct {
	*websocket.Client
	logger *logger.Logger
	user   User
	err    error
}

func NewWebsocket(ctx context.Context, user User) *Websocket {
	prefix := user.Name + "*" + config.WsapiURL
	log := logger.NewLogger(logger.WithPrefix(prefix))
	client := websocket.NewClient(
		ctx, config.WsapiURL, &TradeProcessor{user: user},
		websocket.WithPing(websocket.NewStringMessage("ping")),
		websocket.WithLogger(log),
	)
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	return &Websocket{
		Client: client,
		logger: log,
		user:   user,
	}
}

func (w *Websocket) SubscribeRecord(markets ...string) *Websocket {
	if w.err != nil {
		return w
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
	return w
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

	rand.Seed(time.Now().UnixMilli())
}

func TestAutoTrade(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	for _, user := range config.Users {
		user.ch = make(chan objx.Map, 10)
		user.ac = &autoapi.Account{
			Account:   user.Name,
			AccessKey: user.AccessKey,
			SecretKey: user.SecretKey,
			API:       config.ApiURL,
			Trade:     config.TradeRUL,
			KLine:     config.KlineURL,
			WSAPI:     config.WsapiURL,
		}

		client := NewWebsocket(ctx, user).
			SubscribeRecord(config.Markets...)

		go order(ctx, client.logger, user)
	}

	select {
	case <-ctx.Done():
		cancel()
	}
}
