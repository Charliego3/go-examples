package main

import (
	"github.com/guonaihong/gout"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/logger"
	"io/ioutil"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestCreateRobot(t *testing.T) {
	buf, err := ioutil.ReadFile("./config.json")
	if err != nil {
		logger.Fatal("加载配置文件出错", err)
	}

	err = json.Unmarshal(buf, &Settings)
	if err != nil {
		logger.Fatal("反序列化配置失败", err)
	}

	userId := 362682
	currency := "qc"
	amount := decimal.NewFromInt(20045)
	err = Recharge(userId, currency, amount)
	if err != nil {
		logger.Fatalf("充值失败: %v", err)
	}

	params := url.Values{
		"userId":                 []string{strconv.Itoa(userId)},
		"username":               []string{"15200000047"},
		"market":                 []string{"btc_qc"},
		"currency":               []string{currency},
		"strategyId":             []string{"16"},
		"gridAmount":             []string{"150"},
		"type":                   []string{"1"},
		"totalAmount":            []string{amount.String()},
		"lowerPrice":             []string{"30201"},
		"upperPrice":             []string{"36030"},
		"triggerPrice":           []string{},
		"exchangeWithStopLoss":   []string{"true"},
		"exchangeWithStopProfit": []string{"true"},
	}

	logger.Debugf("创建机器人参数: %s", params.Encode())

	timeout := gout.NewWithOpt(gout.WithTimeout(time.Minute * 3))
	var resp []byte
	err = timeout.GET("http://127.0.0.1:48620/api/fake/saveStrategy").SetQuery(params.Encode()).BindBody(&resp).Do()
	logger.Debugf("创建机器人响应 -> %s, %v", json.Get(resp, "resMsg").ToString(), err)
}
