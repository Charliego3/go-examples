package autoapi

import (
	"github.com/charmbracelet/log"
)

// Market start ==============================

func Markets() map[string]MarketConfig {
	configs := request[map[string]MarketConfig]("data/v1/markets")
	log.Infof("MarketConfig: %+v", configs)
	return configs
}

func AllTicker() {
	tickers := request[Tickers]("data/v1/allTicker")
	log.Infof("Tickers: %+v", tickers)
}

func TickerData(market string) {
	ticker := request[Ticker]("data/v1/ticker", WithMarket(market))
	log.Infof("Ticker: %s = %+v", market, ticker)
}

// Market ended ==============================
