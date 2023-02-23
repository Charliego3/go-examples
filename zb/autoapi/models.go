package autoapi

import "github.com/shopspring/decimal"

type MarketConfig struct {
	AmountScale decimal.Decimal `json:"amountScale"`
	MinAmount   decimal.Decimal `json:"minAmount"`
	MinSize     decimal.Decimal `json:"minSize"`
	PriceScale  decimal.Decimal `json:"priceScale"`
}

type Tickers map[string]TickerModel

type Ticker struct {
	Date   string      `json:"date"`
	Ticker TickerModel `json:"ticker"`
}

type TickerModel struct {
	High     decimal.Decimal `json:"high"`
	Vol      decimal.Decimal `json:"vol"`
	Last     decimal.Decimal `json:"last"`
	Low      decimal.Decimal `json:"low"`
	Buy      decimal.Decimal `json:"buy"`
	Sell     decimal.Decimal `json:"sell"`
	Turnover decimal.Decimal `json:"turnover"`
	Open     decimal.Decimal `json:"open"`
	RiseRate decimal.Decimal `json:"riseRate"`
}

type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
