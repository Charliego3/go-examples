package autoapi

import (
	"github.com/whimthen/temp/logger"
)

// Market start ==============================

func Markets() map[string]MarketConfig {
	configs := request[map[string]MarketConfig]("data/v1/markets")
	logger.Infof("MarketConfig: %+v", configs)
	return configs
}

func AllTicker() {
	tickers := request[Tickers]("data/v1/allTicker")
	logger.Infof("Tickers: %+v", tickers)
}

func TickerData(market string) {
	ticker := request[Ticker]("data/v1/ticker", WithMarket(market))
	logger.Infof("Ticker: %s = %+v", market, ticker)
}

// Market ended ==============================
