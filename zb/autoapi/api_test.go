package autoapi

import (
	"github.com/shopspring/decimal"
	"testing"
)

func init() {
	loadConfig()
}

func TestApis(t *testing.T) {
	Markets()
	AllTicker()
	TickerData("trx_usdt")
}

func TestFunds(t *testing.T) {
	UserInfo()
	DepositAddress("usdt")
	Order("trx_usdt", decimal.NewFromInt(20), decimal.NewFromFloat32(0.062), TradeTypeBuy)
}
