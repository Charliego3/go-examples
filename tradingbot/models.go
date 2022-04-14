package main

import (
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/configs"
)

type ResMsg struct {
	Code    int    `json:"code"`
	Method  string `json:"method"`
	Message string `json:"message"`
}

type Response[T any] struct {
	ResMsg  ResMsg `json:"resMsg"`
	Payload T      `json:"datas"`
}

// ========= Settings start =========

type ConfigSettings struct {
	Basic          *Basic             `json:"basic,omitempty"`
	ViewRobotUsers []int              `json:"viewRobotUsers"`
	Database       []configs.DBConfig `json:"databases,omitempty"`
	TradeUsers     []*User            `json:"tradeUsers,omitempty"`
	TradingUsers   []*User            `json:"tradingUsers,omitempty"`
	Websocket      *Websocket         `json:"websocket,omitempty"`
}

type Database struct {
	URL     string `json:"url,omitempty"`
	Options struct {
		MaxOpenConns    int    `json:"maxOpenConns"`
		MaxIdleConns    int    `json:"maxIdleConns"`
		ConnMaxIdleTime string `json:"connMaxIdleTime"`
		ConnMaxLifeTime string `json:"connMaxLifeTime"`
	} `json:"options"`
}

type Basic struct {
	Domain       string `json:"domain,omitempty"`
	LastRegister int    `json:"lastRegister,omitempty"`
	TradeUrl     string `json:"tradeUrl,omitempty"`
	TradingUrl   string `json:"tradingUrl,omitempty"`
}

type User struct {
	ID        int    `json:"userId"`
	Username  string `json:"username"`
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type Websocket struct {
	Address string `json:"address,omitempty"`
}

// ========= Settings end =========

type Market struct {
	Symbol      string          `json:"symbol"`
	MinAmount   decimal.Decimal `json:"minAmount"`
	SymbolBix   int             `json:"symbolBix"`
	CurrencyBix int             `json:"currencyBix"`
	Name        string          `json:"name"`
	Currency    string          `json:"currency"`
	MaxPrice    decimal.Decimal `json:"maxPrice"`
	MinExchange decimal.Decimal `json:"minExchange"`
	MarketID    int             `json:"marketId"`
}

type UserIncrAsset struct {
}

type UserAsset struct {
	Coins    []Coin `json:"coins"`
	DataType string `json:"dataType"`
	Channel  string `json:"channel"`
	Version  int64  `json:"version"`
	Usdtcny  string `json:"usdtcny"`
}

type Coin struct {
	IsCanWithdraw bool            `json:"isCanWithdraw"`
	CanLoan       bool            `json:"canLoan"`
	Fundstype     int             `json:"fundstype"`
	ShowName      string          `json:"showName"`
	IsCanRecharge bool            `json:"isCanRecharge"`
	CnName        string          `json:"cnName"`
	EnName        string          `json:"enName"`
	Available     decimal.Decimal `json:"available"`
	Freez         decimal.Decimal `json:"freez"`
	UnitTag       string          `json:"unitTag"`
	Key           string          `json:"key"`
	UnitDecimal   int             `json:"unitDecimal"`
}

func (c Coin) ToFund() *Fund {
	return &Fund{
		ID:          c.Fundstype,
		Name:        c.ShowName,
		CanRecharge: c.IsCanRecharge,
		Available:   c.Available,
		Freeze:      c.Freez,
		Unit:        c.UnitDecimal,
	}
}

type Fund struct {
	ID          int             `json:"fundstype,omitempty"`
	Name        string          `json:"showName,omitempty"`
	CanRecharge bool            `json:"isCanRecharge,omitempty"`
	Available   decimal.Decimal `json:"available"`
	Freeze      decimal.Decimal `json:"freez"`
	Unit        int             `json:"unitDecimal,omitempty"`
}

type QuickDepth struct {
	LastTime     int64               `json:"lastTime"`
	DataType     string              `json:"dataType"`
	Channel      string              `json:"channel"`
	CurrentPrice decimal.Decimal     `json:"currentPrice"`
	ListDown     [][]decimal.Decimal `json:"listDown"`
	Market       string              `json:"market"`
	ListUp       [][]decimal.Decimal `json:"listUp"`
	High         decimal.Decimal     `json:"high"`
	Rate         string              `json:"rate"`
	Low          decimal.Decimal     `json:"low"`
	CurrentIsBuy bool                `json:"currentIsBuy"`
	DayNumber    decimal.Decimal     `json:"dayNumber"`
	TotalBtc     decimal.Decimal     `json:"totalBtc"`
	ShowMarket   string              `json:"showMarket"`
}

func (q QuickDepth) NotValid() bool {
	if q.CurrentPrice.IsZero() {
		return true
	}

	if len(q.ListDown) <= 0 && len(q.ListUp) <= 0 {
		return true
	}

	return false
}

func (q QuickDepth) BuyOne() decimal.Decimal {
	if q.NotValid() {
		return decimal.Decimal{}
	}

	return q.ListDown[0][0]
}

func (q QuickDepth) SellOne() decimal.Decimal {
	if q.NotValid() {
		return decimal.Decimal{}
	}

	return q.ListUp[0][0]
}

func (q QuickDepth) MinBuy() decimal.Decimal {
	if q.NotValid() {
		return decimal.Decimal{}
	}

	return q.ListDown[len(q.ListDown)-1][0]
}

func (q QuickDepth) MaxSell() decimal.Decimal {
	if q.NotValid() {
		return decimal.Decimal{}
	}

	return q.ListUp[len(q.ListUp)-1][0]
}

type Robot struct {
	ID         int64           `json:"id,omitempty" db:"id"`
	StrategyId int             `json:"strategyId,omitempty" db:"strategyId"`
	UserID     int             `json:"userId,omitempty" db:"userId"`
	Status     int             `json:"status,omitempty" db:"status"`
	Username   string          `json:"userName,omitempty" db:"userName"`
	Asset      decimal.Decimal `json:"initialAsset,omitempty" db:"initialAsset"`
	Coin       decimal.Decimal `json:"coinAmount,omitempty" db:"coinAmount"`
	Fait       decimal.Decimal `json:"faitAmount,omitempty" db:"faitAmount"`
	Income     decimal.Decimal `json:"income,omitempty" db:"income"`
	EIncome    decimal.Decimal `json:"extractedIncome,omitempty" db:"extractedIncome"`
	TIncome    decimal.Decimal `json:"totalIncome,omitempty" db:"totalIncome"`
	Buy        Bool            `json:"buy,omitempty" db:"isBuy"`
	CreateAt   DateTime        `json:"createTime,omitempty" db:"createTime"`
	StartAt    DateTime        `json:"startTime,omitempty" db:"startTime"`

	Params string `json:"params,omitempty" db:"params"`
}
