package main

import (
	"github.com/guonaihong/gout"
	"github.com/transerver/commons/logger"
	"sync"
)

const (
	getMarketEndpoint = "/api/fake/V1_0_0/getMarket"
)

var (
	markets     = make(map[string]Market)
	marketMutex sync.RWMutex
)

func FetchMarket(market string) (Market, bool) {
	marketMutex.RLock()
	m, ok := markets[market]
	marketMutex.RUnlock()

	if !ok {
		var err error
		m, err = getMarket(market)
		if err != nil {
			return Market{}, false
		}
	}

	return m, true
}

func getMarket(market string) (Market, error) {
	var resp Response[Market]
	err := timeout.GET(getRequestURL(getMarketEndpoint)).SetQuery(gout.H{"market": market}).BindJSON(&resp).Do()
	if err != nil {
		logger.Error("获取市场失败: ", err)
		return Market{}, err
	}

	marketMutex.Lock()
	defer marketMutex.Unlock()

	markets[resp.Payload.Name] = resp.Payload
	return resp.Payload, nil
}
