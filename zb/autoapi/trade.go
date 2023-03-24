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

type OrderResp struct {
	Response `json:",inline"`
	ID       string `json:"ID,omitempty"`
}

// Order spot trade
//
//	Optional parameters:
//	WithAcctType: default AccountTypeMain
//	WithEnableExpress: default false
//	WithEnableRepay: default false
//	WithOrderType: default OrderTypeLimit
//	WithCustomerOrderId: default ""
func Order(market string, price, amount decimal.Decimal, tradeType TradeType, opts ...Option[*Values]) OrderResp {
	opts = append(opts, WithCurrencyMarket(market), WithAmount(amount), WithPrice(price), WithTradeType(tradeType), WithTrade())
	return request[OrderResp]("api/order", opts...)
}

func QueueOrder(market string, price, amount decimal.Decimal, tradeType TradeType, opts ...Option[*Values]) OrderResp {
	opts = append(opts, WithCurrencyMarket(market), WithAmount(amount), WithPrice(price), WithTradeType(tradeType), WithTrade())
	return request[OrderResp]("api/queueOrder", opts...)
}

func BatchOrder(market string, tradeType TradeType, tradeParams [][]decimal.Decimal, opts ...Option[*Values]) {
	opts = append(opts, WithMarket(market), WithTradeType(tradeType), WithObj("tradeParams", tradeParams), WithTrade())
	resp := request[any]("api/orderMoreV2", opts...)
	logger.Infof("OrderMoreV2 response: %+v", resp)
}

func CancelAllOrders(market string, opts ...Option[*Values]) any {
	opts = append(opts, WithCurrencyMarket(market), WithTrade())
	return request[any]("api/cancelAllOpenedOrders", opts...)
}

// Trade ended ==============================
