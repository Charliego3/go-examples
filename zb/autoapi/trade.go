package autoapi

import (
	"github.com/shopspring/decimal"
	"github.com/whimthen/temp/logger"
)

// Fund start ==============================

func UserInfo() {
	info := request[any]("api/getAccountInfo", WithTrade())
	logger.Infof("UserInfo: %+v", info)
}

func DepositAddress(currency string) {
	address := request[any]("api/getUserAddress", WithCurrency(currency), WithTrade())
	logger.Infof("Deposit address: %+v", address)
}

// Fund ended ==============================

// Trade start ==============================

// Order spot trade
//
//	Optional parameters:
//	WithAcctType: default AccountTypeMain
//	WithEnableExpress: default false
//	WithEnableRepay: default false
//	WithOrderType: default OrderTypeLimit
//	WithCustomerOrderId: default ""
func Order(market string, amount, price decimal.Decimal, tradeType TradeType, opts ...Option[*Values]) {
	opts = append(opts, WithCurrencyMarket(market), WithAmount(amount), WithPrice(price), WithTradeType(tradeType), WithTrade())
	resp := request[any]("api/order", opts...)
	logger.Infof("Order: %+v", resp)
}

// Trade ended ==============================
