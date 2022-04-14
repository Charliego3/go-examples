package main

import (
	"context"
	"flag"
	"github.com/guonaihong/gout"
	"github.com/jmoiron/sqlx"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/logger"
	"github.com/transerver/commons/utils"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

var (
	Settings ConfigSettings
	cgPath   = "./tradingbot/config.json"

	timeout = gout.NewWithOpt(gout.WithTimeout(time.Second * 10))
)

func (r Response[T]) Success() bool {
	return r.ResMsg.Code == 1000
}

func loadConfig() {
	buf, err := ioutil.ReadFile(cgPath)
	if err != nil {
		logger.Fatal("加载配置文件出错", err)
	}

	err = json.Unmarshal(buf, &Settings)
	if err != nil {
		logger.Fatal("反序列化配置失败", err)
	}
}

func main() {
	var err error
	cf := flag.String("cf", cgPath, "config.json path")
	tf := flag.String("tf", "./tradingbot/index.html", "template path")
	al := flag.Float64("maxPriceLimit", 1.4, "max order price rate with current price")
	il := flag.Float64("minPriceLimit", 0.6, "min order price rate with current price")
	tu := flag.Int("tu", 2, "trade users count")
	// ru := flag.Int("ru", 5, "trading users count")
	td := flag.Duration("td", time.Second, "trade interval duration")
	flag.Parse()

	cgPath = *cf
	loadConfig()

	maxPriceLimitRate = decimal.NewFromFloat(*al)
	minPriceLimitRate = decimal.NewFromFloat(*il)

	if len(Settings.TradeUsers) < *tu {
		for i := 0; i < *tu-len(Settings.TradeUsers); i++ {
			user, err := Register(false)
			if err != nil {
				return
			}
			logger.Debugf("成功注册用户 -> U: %d, N: %s", user.UserId, user.Username)
		}
	}

	// if len(Settings.TradingUsers) < *ru {
	// 	logger.Warn(*ru - len(Settings.TradingUsers))
	// 	for i := 0; i < *ru-len(Settings.TradingUsers); i++ {
	// 		user, err := Register(true)
	// 		if err != nil {
	// 			return
	// 		}
	// 		logger.Debugf("成功注册用户 -> U: %d, N: %s", user.UserId, user.Username)
	// 	}
	// }

	market, ok := FetchMarket("btcqc")
	if !ok {
		return
	}

	done, cancel := context.WithCancel(context.Background())
	err = SubscribeQuickDepth(done, market)
	if err != nil {
		return
	}

	initDatabase()

	go serve(*tf)
	time.Sleep(time.Second)
	// <-make(chan struct{})
	// return

	// StartTradingRobot(*Settings.TradingUsers[1], market)
	// return

	// for _, user := range Settings.TradingUsers[:*ru] {
	// 	StartTradingRobot(*user, market)
	// }

	for _, user := range Settings.TradeUsers[:*tu] {
		trader := NewTrader(*user, market, *td)
		trader.Run(done)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGQUIT)

	for {
		select {
		case <-exit:
			cancel()
			break
		case <-done.Done():
			logger.Debugf("正常结束.....")
			return
		}
	}
}

func serve(tf string) {
	type Order struct {
		GridIndex  int             `db:"gridIndex"`
		OrderPrice decimal.Decimal `db:"orderPrice"`
		IsBuy      Bool            `db:"isBuy"`
		Count      int             `db:"count"`
		Rate       decimal.Decimal `db:"-"`
	}

	type Params struct {
		AccountId      int             `json:"accountId,omitempty"`
		UpperPrice     decimal.Decimal `json:"upperPrice,omitempty"`
		LowerPrice     decimal.Decimal `json:"lowerPrice,omitempty"`
		StopLowerPrice decimal.Decimal `json:"stopLowerPrice,omitempty"`
		StopUpperPrice decimal.Decimal `json:"stopUpperPrice,omitempty"`
		StopUpper      bool            `json:"stopUpper,omitempty"`
		StopLower      bool            `json:"stopLower,omitempty"`
	}

	queryOrderList := `SELECT g.gridIndex gridIndex, g.orderPrice orderPrice, g.isBuy isBuy, IF(SUM(c.hedgeCount) IS NULL, 0, SUM(c.hedgeCount)) count FROM gridrecordv2 g
LEFT JOIN (SELECT gridIndex, COUNT(1) hedgeCount FROM gridrecordv2 a WHERE userId = ? AND robotId = ?
AND isOrignOrder = FALSE AND status = 4 AND orderStatus = 2 GROUP BY gridIndex) c ON g.gridIndex = c.gridIndex
WHERE userId = ? AND robotId = ? AND status = 1 AND (orderStatus = 1 OR orderStatus = 3) GROUP BY isBuy, orderPrice, gridIndex;`

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, err := template.ParseFiles(tf)
		if err != nil {
			logger.Errorf("解析页面模版失败: %v", err)
			return
		}

		var userIds []interface{}
		// for _, user := range Settings.TradingUsers {
		// 	userIds = append(userIds, user.ID)
		// }
		for _, userId := range Settings.ViewRobotUsers {
			userIds = append(userIds, userId)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		querySQL := "SELECT id, strategyId, userId, status, userName, initialAsset, coinAmount, faitAmount, income, extractedIncome, totalIncome, isBuy, createTime, startTime, params FROM robot WHERE userId IN (?) AND marketName = 'btc_qc' ORDER BY status"
		// querySQL := "SELECT id, strategyId, userId, status, userName, initialAsset, coinAmount, faitAmount, income, extractedIncome, totalIncome, isBuy, createTime, startTime FROM robot ORDER BY id DESC LIMIT 10"
		querySQL, args, err := sqlx.In(querySQL, userIds)
		if err != nil {
			_ = page.Execute(w, err.Error())
			return
		}
		querySQL = db.Rebind(querySQL)
		rows, err := db.QueryxContext(ctx, querySQL, args...)
		// rows, err := db.Queryx(querySQL)
		if err != nil {
			_ = page.Execute(w, err.Error())
			return
		}

		one := decimal.NewFromInt(1)
		var currentPrice decimal.Decimal
		depth := Depth()
		if depth.NotValid() {
			logger.Warn("获取不到当前最新价格")
			currentPrice = one
		} else {
			currentPrice = depth.CurrentPrice
		}

		var rtn []map[string]interface{}
		for rows.Next() {
			var robot Robot
			err := rows.StructScan(&robot)
			if err != nil {
				_ = page.Execute(w, err.Error())
				return
			}

			orderRows, err := db.QueryxContext(ctx, queryOrderList, robot.UserID, robot.ID, robot.UserID, robot.ID)
			if err != nil {
				_ = page.Execute(w, err.Error())
				return
			}

			var buys []Order
			var sells []Order
			for orderRows.Next() {
				var order Order
				err := orderRows.StructScan(&order)
				if err != nil {
					_ = page.Execute(w, err.Error())
					return
				}

				if currentPrice == one {
					order.Rate = decimal.Decimal{}
				} else {
					order.Rate = order.OrderPrice.Div(currentPrice).RoundDown(4).Sub(one).Shift(2)
				}
				if order.IsBuy {
					buys = append(buys, order)
				} else {
					sells = append(sells, order)
				}
			}

			sort.Slice(buys, func(i, j int) bool {
				return buys[i].GridIndex < buys[j].GridIndex
			})
			sort.Slice(sells, func(i, j int) bool {
				return sells[i].GridIndex > sells[j].GridIndex
			})

			var params Params
			index := strings.Index(robot.Params, "\"gridInfo\":{")
			lastIndex := strings.LastIndex(robot.Params, "}},")
			robotParams := robot.Params[:index] + robot.Params[lastIndex+3:]
			err = json.Unmarshal(utils.Bytes(robotParams), &params)
			if err != nil {
				logger.Errorf("反序列化失败: %v", err)
				_ = page.Execute(w, err.Error())
				return
			}

			var subUserId int
			err = db.GetContext(ctx, &subUserId, "SELECT subUserId FROM account WHERE id = ?", params.AccountId)
			if err != nil {
				_ = page.Execute(w, err.Error())
				return
			}

			funds, err := GetFunds(subUserId, "btc", "qc")
			if err != nil {
				_ = page.Execute(w, err.Error())
				return
			}

			// logger.Warnf("UserFund: %+v", funds)

			m := make(map[string]interface{})
			m["robot"] = robot
			m["buys"] = buys
			m["sells"] = sells
			m["price"] = depth.CurrentPrice
			m["btc"] = funds.BTC
			m["qc"] = funds.QC
			m["params"] = params
			rtn = append(rtn, m)
		}

		_ = page.Execute(w, rtn)
	})

	if err := http.ListenAndServe(":9090", nil); err != nil {
		logger.Fatal(err)
	}
}

func overrideConfig() {
	buf, err := json.MarshalIndent(Settings, "", "    ")
	if err != nil {
		logger.Errorf("覆盖配置文件失败: %v", err)
		return
	}

	err = ioutil.WriteFile(cgPath, buf, os.ModePerm)
	if err != nil {
		logger.Error("写入配置文件失败", err)
	}
}
