package main

import (
	"encoding/json"
	"io"
	"sync/atomic"

	"github.com/charliego93/websocket"
	"github.com/charmbracelet/log"
	"github.com/gookit/goutil/strutil"
	"github.com/shopspring/decimal"
)

type OkexResponse struct {
	Event string     `json:"event"`
	Arg   OkexArg    `json:"arg"`
	Data  []OkexData `json:"data"`
}

type OkexData struct {
	InstId  string          `json:"instId"`  // 市场
	TradeId string          `json:"tradeId"` // 委托ID
	Px      decimal.Decimal `json:"px"`      // 成交价格
	Sz      decimal.Decimal `json:"sz"`      // 成交数量
	Side    string          `json:"side"`    // 成交方向，buy sell
	Ts      string          `json:"ts"`      // 成交时间，Unix时间戳的毫秒数格式，如 1597026383085
}

type OwnTicker struct {
	Buy  decimal.Decimal `json:"buy"`
	Sell decimal.Decimal `json:"sell"`
	Last decimal.Decimal `json:"last"`
}

type OwnProcessor struct {
	logger  *log.Logger
	Ticker  atomic.Value
	Markets map[string]*MarketConfig
}

type MarketConfig struct {
	AmountScale int32           `json:"amountScale"`
	PriceScale  int32           `json:"priceScale"`
	MinAmount   decimal.Decimal `json:"minAmount"`
	MinSize     decimal.Decimal `json:"minSize"`
}

func (b *OwnProcessor) OnReceive(frame *websocket.Frame) {
	defer func() {
		if err := recover(); err != nil {
			b.logger.Error("OnReceive", "err", err)
		}
	}()

	bs, _ := io.ReadAll(frame.Reader)
	content := string(bs)
	if content == "pong" {
		b.logger.Info("收到心跳pong")
		return
	}

	var resp = struct {
		Type    string    `json:"dataType"`
		Ticker  OwnTicker `json:"ticker"`
		Channel string    `json:"channel"`
		Data    any       `json:"data"`
	}{}
	err := json.Unmarshal(bs, &resp)
	if err != nil {
		b.logger.Error("解析响应数据失败", "err", err, "data", content)
		return
	}

	if resp.Type == "ticker" {
		b.Ticker.Store(&resp.Ticker)
	} else if resp.Channel == "markets" {
		bs, _ := json.Marshal(resp.Data)
		m := make(map[string]*MarketConfig, len(resp.Data.(map[string]any)))
		json.Unmarshal(bs, &m)
		b.Markets = make(map[string]*MarketConfig, len(m))
		for k, v := range m {
			b.Markets[ClearMarket(k)] = v
		}
		// b.logger.Debug("Receive", "markets", b.Markets)
	} else {
		b.logger.Warn("未处理的消息", "msg", content)
	}
}

func (b *OwnProcessor) SetLogger(log *log.Logger) {
	b.logger = log
}

type OkexProcessor struct {
	logger *log.Logger
	ch     chan OkexData
}

func NewOkexProcessor(ch chan OkexData) *OkexProcessor {
	return &OkexProcessor{ch: ch}
}

func (b *OkexProcessor) OnReceive(frame *websocket.Frame) {
	var resp OkexResponse
	err := json.NewDecoder(frame.Reader).Decode(&resp)
	if err != nil {
		b.logger.Error("解析响应失败", "err", err)
		return
	}

	if strutil.IsNotBlank(resp.Event) {
		b.logger.Info("订阅成功", "Channel", resp.Arg.Channel, "InstId", resp.Arg.InstId)
		return
	}

	if resp.Arg.Channel == "trades" {
		b.ch <- resp.Data[0]
	} else {
		b.logger.Warn("未处理的通道类型", "Channel", resp.Arg.Channel, "InstId", resp.Arg.InstId)
	}
}

func (b *OkexProcessor) SetLogger(log *log.Logger) {
	b.logger = log
}
