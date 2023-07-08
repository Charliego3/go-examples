package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/charliego93/websocket"
	"github.com/charmbracelet/log"
	"github.com/shopspring/decimal"
	"github.com/whimthen/temp/zb/autoapi"
)

type Orderer struct {
	ctx       context.Context
	logger    *log.Logger
	opts      log.Options
	ForClient *websocket.Client
	OwnClient *websocket.Client
}

func (o *Orderer) Start(c OrderConfig) {
	o.logger = log.NewWithOptions(os.Stdout, o.opts)
	ch := make(chan OkexData, 10)
	o.ForClient = websocket.NewClient(
		o.ctx, c.ForURL, NewOkexProcessor(ch),
		websocket.WithPrefix(c.Name),
		websocket.WithLoggerOptions(&o.opts),
		// websocket.WithLogger(o.logger),
	)
	err := o.ForClient.Connect()
	if err != nil {
		o.logger.Fatal("外盘websocket连接失败", "err", err)
	}
	o.ForClient.Subscribe(&OkexReqMsg{
		Op: "subscribe",
		Args: []OkexArg{
			{"trades", c.Market},
		},
	})

	ownURL := c.OwnURL
	// quick depth
	// regex := regexp.MustCompile("(usdt|qc|btc|usdc)$")
	// ownURL = ownURL + "/" + regex.ReplaceAllString(ClearMarket(c.Market), "")
	processor := &OwnProcessor{}
	o.OwnClient = websocket.NewClient(
		o.ctx, ownURL, processor,
		websocket.WithPrefix(c.Name),
		websocket.WithPing(websocket.NewStringMessage("ping")),
		websocket.WithLoggerOptions(&o.opts),
		// websocket.WithLogger(o.logger),
	)
	err = o.OwnClient.Connect()
	if err != nil {
		o.logger.Fatal("websocket连接失败", "err", err)
	}
	o.OwnClient.Subscribe(NewOwnReq("markets"))
	o.OwnClient.Subscribe(NewOwnTickerReq(c.Market))

	rand.Seed(time.Now().UnixNano())
	for data := range ch {
		pt := processor.Ticker.Load()
		// o.logger.Info("收到外盘成交", "data", data)
		if pt == nil {
			o.logger.Error("receive", "data", data)
			continue
		}

		config := processor.Markets[ClearMarket(c.Market)]
		if config == nil {
			o.logger.Error("can not fetch market config", "market", c.Market)
			continue
		}

		ticker := pt.(*OwnTicker)
		price := o.getPrice(ticker, config, data.Px)
		if price.IsZero() {
			o.logger.Warn("price warning", "data", data.Px, "ticker", ticker)
			continue
		}

		numbers := getNumbers(config, price, data.Sz)
		if numbers.IsZero() {
			o.logger.Warn("random number is zero")
			continue
		}

		tradeType := autoapi.TradeTypeBuy
		if data.Side == "sell" {
			tradeType = autoapi.TradeTypeSell
		}
		account := &autoapi.Account{
			Account:   c.Name,
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
			API:       c.ApiURL,
			Trade:     c.TrdURL,
			WSAPI:     c.OwnURL,
		}
		resp := autoapi.QueueOrder(ClearMarket(c.Market), price, numbers, tradeType, autoapi.WithAccount(account))
		o.logger.Info("下单结果", "price", price, "numbers", numbers, "resp", resp, "ticker", ticker, "for", data.Px)
	}
}

func getNumbers(config *MarketConfig, price, num decimal.Decimal) decimal.Decimal {
	rate := 0.5 + rand.Float64()
	exponent := config.AmountScale
	f, _ := num.Shift(exponent).Float64()
	numbers := decimal.NewFromFloat(f*rate + f).Shift(-exponent).Round(exponent)
	if config.MinAmount.GreaterThan(numbers) || config.MinSize.GreaterThan(numbers.Mul(price)) {
		numbers = config.MinSize.Div(price).RoundUp(config.AmountScale)
	}
	return numbers
}

func (o *Orderer) getPrice(ticker *OwnTicker, config *MarketConfig, forPrice decimal.Decimal) decimal.Decimal {
	var lower, upper decimal.Decimal
	var rate float64
	if ticker.Buy.GreaterThan(forPrice) { // 买一大于第三方价格
		lower = ticker.Buy
		upper = ticker.Last
		rate = rand.Float64() * 0.5
	} else if ticker.Sell.LessThan(forPrice) { // 卖一小于第三方价格
		lower = ticker.Last
		upper = ticker.Sell
		rate = rand.Float64()*(1-0.5) + 0.5
	}

	if rate > 0 {
		exponent := config.PriceScale
		low, _ := lower.Shift(exponent).Float64()
		upp, _ := upper.Shift(exponent).Float64()
		number := rate*(upp-low) + low
		price := decimal.NewFromFloat(number).Shift(-exponent).Round(exponent)
		// o.logger.Debug("生成买一卖一间价格", "price", price, "rate", rate, "lower", lower, "upper", upper, "exponent", exponent)
		return price
	}

	o.logger.Info("本盘价格", "buy", ticker.Buy, "sell", ticker.Sell)
	return decimal.Zero
}
