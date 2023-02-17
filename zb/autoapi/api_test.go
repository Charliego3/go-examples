package autoapi

import (
	"github.com/shopspring/decimal"
	"github.com/whimthen/temp/logger"
	"net"
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

func TestName(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	logger.Info(listener.Addr().String())

}
