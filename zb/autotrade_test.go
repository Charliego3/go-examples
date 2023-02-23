package zb

import (
	"context"
	"flag"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/stretchr/objx"
	"github.com/whimthen/temp/logger"
	"github.com/whimthen/temp/websocket"
	"github.com/whimthen/temp/zb/autoapi"
	"gopkg.in/yaml.v3"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync/atomic"
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
	ApiURL   string        `yaml:"apiURL"`
	TradeRUL string        `yaml:"tradeRUL"`
	KlineURL string        `yaml:"klineURL"`
	WsapiURL string        `yaml:"wsapiURL"`
	Interval time.Duration `yaml:"interval"`
	Markets  []string      `yaml:"markets"`
	Users    []User        `yaml:"accounts"`
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

	dataType := obj.Get("dataType").String()
	switch dataType {
	case "quickDepth":
		p.logger.Debugf("收到快速行情: %s", content)
		cprice.Store(decimal.NewFromFloat(obj.Get("currentPrice").Float64()))

	case "userIncrRecord":
		p.user.ch <- obj
	default:
		p.logger.Warnf("收到未处理的消息类型: %s", content)
	}
}

func receiveOrder(ctx context.Context, logger *logger.Logger, user User) {
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
				logger.Infof("本次成交: [%s:%d] = %s / %s", entrustId, types, unitPrice, numbers)
				if !*rOrder {
					return
				}

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
				market := strings.TrimSuffix(obj.Get("market").String(), "default")
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

var (
	config Config
	rOrder = flag.Bool("rOrder", false, "成交后是否下 taker 吃单")
	buyer  atomic.Value
	seller atomic.Value
	cprice atomic.Value
)

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

func listenDepth(ctx context.Context) {
	group := getGroupMarkets()
	logger.Debugf("GroupMarkets: %+v", group)
	zoneRegex := regexp.MustCompile("(" + strings.Join(group["zone"], "|") + ")$")

	clients := make(map[string]*websocket.Client)

	for _, market := range config.Markets {
		websocketURL := config.WsapiURL
		if !strings.HasSuffix(websocketURL, "/") {
			websocketURL += "/"
		}
		websocketURL += zoneRegex.ReplaceAllString(market, "")

		var client *websocket.Client
		if c, ok := clients[websocketURL]; ok {
			client = c
		} else {
			c := websocket.NewClient(
				ctx, websocketURL, &TradeProcessor{},
				websocket.WithPing(websocket.NewStringMessage("ping")),
			)
			err := c.Connect()
			if err != nil {
				panic(err)
			}
			client = c
			clients[websocketURL] = client
		}

		err := client.SendMessage(&websocket.RequestMsg{
			Event:   EventAddChannel,
			Channel: market + "_quick_depth",
		})
		if err != nil {
			logger.Fatalf("订阅深度失败: %+v", err)
		}
	}
}

func TestAutoTrade(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	listenDepth(ctx)

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

		go receiveOrder(ctx, client.logger, user)

		var isBuy bool
		for _, market := range config.Markets {
			//go makeOrder(ctx, user, market, isBuy)
			_ = market
			isBuy = !isBuy
		}
	}

	select {
	case <-ctx.Done():
		cancel()
	}
}

func makeOrder(ctx context.Context, user User, market string, isBuy bool) {
	log := logger.NewLogger(logger.WithPrefix("ORDER:%t:%s:%s", isBuy, strings.ToUpper(market), user.Name))
	ticker := time.NewTicker(config.Interval)
	for {
		select {
		case <-ticker.C:
			tradeType := autoapi.TradeTypeSell
			if isBuy {
				tradeType = autoapi.TradeTypeBuy
			}

			var numbers, price decimal.Decimal // numbers = minNumber * 2
			if isBuy {

			}

			autoapi.Order(market, numbers, price, tradeType, autoapi.WithAccount(user.ac))
		case <-ctx.Done():
			log.Infof("退出布单流程....")
			return
		}
	}
}

func getGroupMarkets() map[string][]string {
	resp, err := http.Get(config.ApiURL + "data/v1/getGroupMarkets")
	if err != nil {
		logger.Fatalf("获取 GroupMarkets 失败: %+v", err)
	}

	group := make(map[string][]string)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&group)
	if err != nil {
		logger.Fatalf("GroupMarkets 反序列化失败: %s", err)
	}
	return group
}
