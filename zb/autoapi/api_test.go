package autoapi

import (
	"net"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/shopspring/decimal"
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
	Order("trx_usdt", decimal.NewFromFloat32(0.062), decimal.NewFromInt(20), TradeTypeBuy)
}

func TestName(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	log.Info(listener.Addr().String())

}
