package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/configs"
	"github.com/transerver/commons/dbs"
	"github.com/transerver/commons/logger"
	"github.com/xo/dburl"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	tradingPath = "/api/fake/"
)

var (
	minProfitRate = decimal.NewFromFloat(0.002)
	two           = decimal.NewFromInt(2)
	db            *dbs.Database
)

type Trading struct {
	user   User
	market Market

	logger *logger.Logger
	robot  Robot
}

func initDatabase() {
	u, err := dburl.Parse(Settings.Database[0].URL)
	if err != nil {
		logger.Fatal(err)
	}

	db = dbs.NewDatabase(dbs.WithConfig(&configs.DBConfig{
		Driver:         "mysql",
		DSN:            u.DSN,
		URL:            u.RawPath,
		DBName:         u.Path[1:],
		DesensitiseDSN: u.Redacted(),
		Options:        Settings.Database[0].Options,
	}))
	err = db.Connect()
	if err != nil {
		logger.Fatal("获取不到数据库连接", err)
	}
	// db.SetDatabaseLoggerHook()
}

func StartTradingRobot(user User, market Market) *Trading {
	t := &Trading{
		user:   user,
		market: market,
		logger: logger.NewLogger(logger.WithPrefix("TRADING:%d:%s", user.ID, user.Username)),
	}

	t.Start()
	return t
}

func (t *Trading) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	row := db.QueryRowxContext(ctx, "SELECT id, strategyId, userId, status, userName, initialAsset, coinAmount, faitAmount, income, extractedIncome, totalIncome, isBuy, createTime, startTime, params FROM robot WHERE userId = ? AND status <= 1 AND marketName = 'btc_qc' LIMIT 1", t.user.ID)
	if row.Err() != nil {
		t.logger.Errorf("查询机器人失败: %v", row.Err())
		return
	}

	var robot Robot
	err := row.StructScan(&robot)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := t.CreateRobot()
			if err != nil {
				t.logger.Errorf("重来!!!, %v", err)
				time.Sleep(time.Second * 10)
				t.Start()
				return
			}
			t.logger.Debugf("没有记录")
		} else {
			t.logger.Errorf("反序列化机器人失败: %v", err)
		}
	}
	t.robot = robot
}

func (t *Trading) CreateRobot() (err error) {
	count := Random(int64(101))
	_, upper := t.getPrice(count)
	if upper.IsZero() {
		t.logger.Warnf("生成网格价格异常 -> Upper: %s", upper)
		return nil
	}

	diff := minProfitRate.Add(minProfitRate.Mul(two)).Mul(upper).Sub(minProfitRate)
	lower := upper.Sub(diff.Mul(decimal.NewFromInt(count)))

	strategyId := RandomRange(15, 18)
	var currency string
	var isBuy bool
	if strategyId == 16 {
		currency = t.market.Currency
		isBuy = true
	} else if strategyId == 17 {
		currency = t.market.Symbol
		isBuy = false
	} else {
		currency = t.randomCurrency()
		isBuy = currency == t.market.Currency
	}

	_ = isBuy

	amount := t.getTotalAmount(upper, lower, count, isBuy)
	err = Recharge(t.user.ID, currency, amount)
	if err != nil {
		return err
	}

	params := url.Values{
		"userId":                 []string{strconv.Itoa(t.user.ID)},
		"username":               []string{t.user.Username},
		"market":                 []string{t.market.Name},
		"currency":               []string{currency},
		"strategyId":             []string{strconv.Itoa(strategyId)},
		"gridAmount":             []string{strconv.FormatInt(count, 10)},
		"type":                   []string{"1"},
		"totalAmount":            []string{amount.String()},
		"lowerPrice":             []string{lower.String()},
		"upperPrice":             []string{upper.String()},
		"triggerPrice":           []string{},
		"exchangeWithStopLoss":   []string{t.randomExchange()},
		"exchangeWithStopProfit": []string{t.randomExchange()},
	}

	t.logger.Debugf("创建机器人参数: %s", params.Encode())

	var resp Response[Robot]
	err = timeout.GET(t.getRequestURL("saveStrategy")).SetQuery(params.Encode()).BindJSON(&resp).Do()
	if err != nil || !resp.Success() {
		return fmt.Errorf("创建机器人失败 -> %s%v", resp.ResMsg.Message, err)
	}

	t.robot = resp.Payload
	return nil
}

func (t *Trading) Shutdown() error {
	params := url.Values{
		"userId":           []string{strconv.Itoa(t.user.ID)},
		"id":               []string{t.user.Username},
		"exchangeWithStop": []string{t.market.Name},
	}

	var resp Response[any]
	err := timeout.GET(t.getRequestURL("shutdown")).SetQuery(params.Encode()).BindJSON(&resp).Do()
	if err != nil {
		t.logger.Errorf("创建机器人失败 -> %v", err)
		return err
	}
	return nil
}

func (t *Trading) getTotalAmount(upper, lower decimal.Decimal, count int64, isBuy bool) decimal.Decimal {
	// diff := minProfitRate.Add(minProfitRate.Mul(decimal.NewFromInt(2))).Mul(upper)
	// diff := upper.Sub(lower).Div(decimal.NewFromInt(count))
	// amount := upper.Sub(lower).Div(diff)
	c := decimal.NewFromInt(int64(t.market.CurrencyBix))
	minexchange := t.market.MinExchange.Add(c)
	if isBuy {
		return minexchange.Mul(decimal.NewFromInt(count))
	}

	minAmount := minexchange.Div(lower)
	return decimal.Max(minAmount, t.market.MinAmount.Add(decimal.NewFromInt(int64(t.market.SymbolBix)))).Mul(c)
}

func (t *Trading) getPrice(count int64) (min, max decimal.Decimal) {
	depth := Depth()
	if depth.NotValid() {
		return
	}

	downLen := len(depth.ListDown)
	upLen := len(depth.ListUp)

	if count <= 30 {
		return t.prices(depth, downLen, upLen, 5)
	} else if count <= 60 {
		return t.prices(depth, downLen, upLen, 10)
	} else if count <= 80 {
		return t.prices(depth, downLen, upLen, 15)
	} else {
		return t.prices(depth, downLen, upLen, 20)
	}

	// if count <= 30 {
	// 	return t.prices(depth, downLen, upLen, 5)
	// } else if count <= 70 {
	// 	return t.prices(depth, downLen, upLen, 10)
	// } else if count <= 110 {
	// 	return t.prices(depth, downLen, upLen, 15)
	// } else {
	// 	return t.prices(depth, downLen, upLen, 20)
	// }
}

func (t *Trading) prices(depth QuickDepth, downLen, upLen, index int) (min, max decimal.Decimal) {
	if downLen > index {
		min = depth.ListDown[downLen-index][0]
	} else {
		min = depth.ListDown[downLen-1][0]
	}
	if upLen > index {
		max = depth.ListUp[index][0]
	} else {
		max = depth.ListUp[upLen-1][0]
	}
	return
}

func (t *Trading) randomCurrency() string {
	if Random(2) == 1 {
		return t.market.Currency
	}
	return t.market.Symbol
}

func (t *Trading) randomExchange() string {
	if Random(2) == 1 {
		return "true"
	}
	return "false"
}

func (t *Trading) getRequestURL(method string) string {
	requestURL := Settings.Basic.TradingUrl
	if strings.HasSuffix(requestURL, "/") {
		requestURL = requestURL[:len(requestURL)-1]
	}
	return requestURL + tradingPath + method
}
