package main

import (
	json "github.com/json-iterator/go"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/logger"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	oneHundred = decimal.NewFromInt(100)
)

type OrderResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

func Order(u User, m Market, logger *logger.Logger, price, number decimal.Decimal, isBuy int) {
	nanoid, err := gonanoid.Generate(alphabet, 6)
	if err != nil {
		logger.Errorf("生成自定义ID失败, 忽略本次下单 -> 价格: %s, 数量: %s, 买入: %t", price, number, isBuy == 1)
		return
	}
	params := url.Values{
		"tradeType":       []string{strconv.Itoa(isBuy)},
		"acctType":        []string{"0"},
		"currency":        []string{m.Name},
		"amount":          []string{number.String()},
		"price":           []string{price.String()},
		"customerOrderId": []string{nanoid},
	}

	// logger.Warnf("开始下单 -> 价格: %s, 数量: %s, 买入: %t", price, number, isBuy == 1)
	buf, err := u.Request("order", params)
	if err != nil {
		logger.Errorf("下单失败 -> 价格: %s, 数量: %s, 买入: %t >>> %v", price.StringFixed(int32(m.CurrencyBix)), number.StringFixed(int32(m.SymbolBix)), isBuy == 1, err)
		return
	}

	var resp OrderResp
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		logger.Errorf("反序列化下单响应失败: %s, err: %v", buf, err)
		return
	}

	if resp.Code == 2009 {
		coin := m.Symbol
		if isBuy == 1 {
			coin = m.Currency
		}
		amount := number.Mul(oneHundred)
		err := Recharge(u.ID, coin, amount)
		if err == nil {
			logger.Infof("充值成功 -> %s: %s", coin, amount.String())
		}
		return
	} else if resp.Code != 1000 {
		logger.Errorf("下单失败 -> 价格: %s, 数量: %s, 买入: %t >>> %s", price.StringFixed(int32(m.CurrencyBix)), number.StringFixed(int32(m.SymbolBix)), isBuy == 1, resp.Message)
	} else {
		logger.Debugf("下单成功 -> 价格: %s, 数量: %s, ID: %s, 买入: %t", price.StringFixed(int32(m.CurrencyBix)), number.StringFixed(int32(m.SymbolBix)), resp.ID, isBuy == 1)
	}
}

func (u User) Request(method string, params url.Values) (buf []byte, err error) {
	params["accesskey"] = []string{u.APIKey}
	params["method"] = []string{method}
	sign := HmacMD5(params.Encode(), u.APISecret)
	params["sign"] = []string{sign}
	params["reqTime"] = []string{strconv.Itoa(int(time.Now().Unix() * 1000))}

	var resp []byte
	requestURL := Settings.Basic.TradeUrl
	if !strings.HasSuffix(requestURL, "/") {
		requestURL += "/"
	}
	err = timeout.POST(requestURL + method).SetQuery(params.Encode()).BindBody(&resp).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
