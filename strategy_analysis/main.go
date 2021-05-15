package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/golog"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	args = Args{}

	strategyDBURL string
	entrustDBURL  string

	robot *Robot
	grid  = Grid{}

	market   string
	currency string
	symbol   string

	isProd bool
)

const (
	strategyDBName = "strategy"
	entrustDBName  = "zb_%sentrust"
)

func main() {
	golog.SetLevel("debug")
	golog.SetTimeFormat("")

	cmd := os.Args[0]

	rootCmd := cobra.Command{
		Use:     cmd,
		Example: fmt.Sprintf(`  %s -e 130 --robot 14196%s  %s --robot 14196 --sd "user:password@tcp(ip:port)/dbname" --ed "user:password@tcp(ip:port)/dbname"`, cmd, "\n", cmd),
		Run:     analysisFunc,
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&args.Env, "env", "e", "", "环境, 测试环境有效值: 130, 123, 129, 218")
	flags.Int64Var(&args.RobotId, "robot", 0, "机器人ID")
	flags.StringVar(&args.ProdStrategyDBURL, "sd", "", "网格服务的数据库连接地址, eg: user:password@tcp(ip:port)/dbname")
	flags.StringVar(&args.ProdEntrustDBURL, "ed", "", "盘口服务的数据库连接地址, eg: user:password@tcp(ip:port)/dbname")

	if err := rootCmd.Execute(); err != nil {
		golog.Error("Execute Error:", err)
	}
}

func analysisFunc(*cobra.Command, []string) {
	if !parseArgs() {
		return
	}

	connected := connect(strategyDBURL, false)
	if !connected {
		return
	}

	r, ok := loadRobot(args.RobotId)
	if !ok {
		return
	}

	market = r.Name
	currency = strings.ToUpper(market[strings.Index(market, "_")+1:])
	symbol = strings.ToUpper(market[:strings.Index(market, "_")])

	if entrustDBURL == "" {
		env := Envs[args.Env]
		entrustDBURL = env + fmt.Sprintf(entrustDBName, strings.Replace(market, "_", "", -1))
	}

	if isProd {
		golog.Debugf("网格DB URL: %s", strategyDBURL[strings.Index(strategyDBURL, "@"):])
		golog.Debugf("盘口DB URL: %s", entrustDBURL[strings.Index(entrustDBURL, "@"):])
	}

	connected = connect(entrustDBURL, true)
	if !connected {
		return
	}

	robot = r
	parseGrid()
	rand.Seed(time.Now().UnixNano())

	if args.RobotId > 0 {
		withRobot(args.RobotId)
	} else {
		withGridRecord(args.ID)
	}
}

func withRobot(robotId int64) {
	records, ok := getGridRecordByRobotId(robotId)
	if !ok {
		golog.Warn("未查询到网格记录, RobotId: ", robotId)
		return
	}

	for _, record := range records {
		colorStatus = reset
		analysisWarp(record)
		println()
	}
}

func withGridRecord(id int64) {
	record, ok := getGridRecordById(id)
	if !ok {
		if id == args.ID {
			golog.Warn("未查询到网格记录, ID: ", id)
		}
		return
	}

	analysisWarp(record)
}

func analysisWarp(record *GridRecord) {
	color.Green("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")

	analysis(record, true)

	color.Green("↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑")
}

func analysis(record *GridRecord, first bool) {
	var recordType = "对冲单"
	if record.IsOriginOrder {
		recordType = "原始单"
		if !first {
			colorStatus = change
		}
	} else {
		colorStatus = noChange
	}
	first = false
	var prefix string
	if record.OriginRecordId.Valid {
		prefix = fmt.Sprintf("[%s] GridIndex:[%d] -> 原始ID:[%s], 订单ID:[%d]", recordType, record.GridIndex, randomColor(record.OriginRecordId.Int64), record.Id)
	} else {
		prefix = fmt.Sprintf("[%s] GridIndex:[%d] -> 订单ID:[%s]", recordType, record.GridIndex, randomColor(record.Id))
	}
	if record.OrderId.Valid {
		prefix += fmt.Sprintf(" -> 委托ID:[%s] ", record.OrderId.String)
	}

	var iocType = "限价"
	var entrustType = "买入"
	if !record.IsBuy {
		entrustType = "卖出"
	}
	if record.IsIocOrder {
		iocType = "IOC"
	}
	entrustType = iocType + entrustType
	info := grid.GridInfo[record.GridIndex]
	prefix += fmt.Sprintf("- 价格:[%s], 数量:[%s], 委托类型:[%s]", info.Price.String(), info.Amount.String(), entrustType)

	var entrust *Entrust
	var entrustOK bool
	if !record.OrderId.Valid { // 已委托
		golog.Error("网格记录以生成, 委托还未成功....")
	} else {
		entrust, entrustOK = getEntrustById(record.OrderId.String)
		if !entrustOK {
			golog.Error("在盘口中找不到委托记录")
		} else {
			prefix += fmt.Sprintf(", 挂单价格:[%s]", decimal.Decimal(entrust.UnitPrice).String())
		}
	}
	prefix += " >> "
	golog.SetPrefix(prefix)

	switch record.Status {
	case 1: // 挂单中
		{
			if record.OrderStatus != 1 {
				golog.Error("记录异常")
				break
			}
			if entrustOK {
				isPart := false
				completeNumber := decimal.Decimal(entrust.CompleteNumber)
				numbers := decimal.Decimal(entrust.Numbers)
				if completeNumber.Cmp(decimal.Zero) > 0 && completeNumber.Cmp(numbers) < 0 {
					isPart = true
				}

				if entrust.Status == 2 {
					golog.Warn("委托在盘口已成交")
				} else if entrust.Status == 3 {
					if isPart {
						golog.Info("网格记录正在挂单中, 盘口中委托记录正在委托[部分成交]....")
					} else {
						golog.Info("网格记录正在挂单中, 盘口中委托记录正在委托[待成交]....")
					}
				} else if entrust.Status == 1 {
					if isPart {
						golog.Info("网格记录正在挂单中, 盘口中委托记录已撤销[部分成交]....")
					} else {
						golog.Info("网格记录正在挂单中, 盘口中委托记录已撤销[待成交]....")
					}
				}
			}
		}
	case 2: // 待对冲
		{
			golog.Warn("待对冲")
		}
	case 3: // 对冲中
		{
			if entrustOK {
				if entrust.Status == 2 {
					golog.Infof("该订单对冲中, 委托已成交, 成交数量:[%s]", decimal.Decimal(entrust.CompleteNumber).String())
				} else if entrust.Status == 1 {
					golog.Infof("该订单对冲中, 委托已撤销")
				} else if entrust.Status == 3 {
					golog.Infof("该订单对冲中, 委托已部分成交[%s]", decimal.Decimal(entrust.CompleteNumber).String())
				}
			} else {
				golog.Info("对冲中")
			}

			hedgeRecord, ok := getHedgeRecord(record.Id)
			if !ok {
				golog.Error("查询不到网格对冲单...")
			} else {
				analysis(hedgeRecord, first)
			}
		}
	case 4: // 已完成
		{
			if entrustOK {
				if entrust.Status == 2 {
					golog.Infof("该订单已完成, 委托已成交, 成交数量:[%s]", decimal.Decimal(entrust.CompleteNumber).String())
				} else if entrust.Status == 1 {
					golog.Infof("该订单已完成, 委托已撤销")
				} else if entrust.Status == 3 {
					golog.Infof("该订单已完成, 委托已部分成交[%s]", decimal.Decimal(entrust.CompleteNumber).String())
				}
			} else {
				golog.Info("订单已完成")
			}

			if record.IsOriginOrder {
				hedgeRecord, ok := getHedgeRecord(record.Id)
				if !ok {
					golog.Error("查询不到网格对冲单...")
				} else {
					analysis(hedgeRecord, first)
				}
			} else {
				originOrder, ok := getNextOriginOrder(record.Id, info.Price, record.GridIndex)
				if !ok {
					return
				}
				analysis(originOrder, first)
			}
		}
	case 5: // 已结束对冲
		{
			golog.Warn("已结束对冲")
		}
	case 6: // 已取消
		{
			golog.Warn("已取消")
		}
	case 7: // 已结束
		{
			golog.Warn("已结束")
		}
	}
}

func parseGrid() {
	params := []byte(robot.Params)
	for i := 0; i <= 100; i++ {
		if bytes.Contains(params, []byte(fmt.Sprintf("%d:", i))) {
			params = bytes.Replace(params, []byte(fmt.Sprintf("%d:", i)), []byte(fmt.Sprintf(`"%d":`, i)), 1)
		} else {
			break
		}
	}
	err := jsoniter.Unmarshal(params, &grid)
	if err != nil {
		golog.Error("解析网格参数失败.... ", err)
		os.Exit(1)
	}

	stype := "等比"
	coin := currency
	if grid.Type == 1 {
		stype = "等差"
	}
	if grid.Buy {
		stype += "买入"
	} else {
		stype += "卖出"
		coin = symbol
	}

	c := color.New(color.FgBlue, color.Bold, color.Concealed, color.Underline)
	golog.Debug(c.Sprintf("网格数量: %d", grid.GridAmount))
	golog.Debug(c.Sprintf("投入数量: %s%s", grid.TotalAmount.String(), coin))
	golog.Debug(c.Sprintf("最高价: %s", grid.UpperPrice.String()))
	golog.Debug(c.Sprintf("最低价: %s", grid.LowerPrice.String()))
	if grid.TriggerPrice.Cmp(decimal.Zero) > 0 {
		golog.Debug(c.Sprintf("触发价: %s", grid.TriggerPrice.String()))
	}
	golog.Debug(c.Sprintf("网格类型: %s\n", stype))
}

var (
	colorStatus  = noChange
	currentColor = color.FgBlue
	colors       = []color.Attribute{
		color.FgBlue,
		color.FgHiRed,
		color.FgHiGreen,
		color.FgYellow,
		color.FgHiMagenta,
		color.FgHiCyan,
	}
)

func randomColor(content int64) string {
	if colorStatus == change {
		cs := make([]color.Attribute, 0)
		for _, attribute := range colors {
			if currentColor != attribute {
				cs = append(cs, attribute)
			}
		}
		rn := rand.Intn(len(cs))
		currentColor = cs[rn]
	} else if colorStatus == reset {
		currentColor = color.FgBlue
	}

	rtn := fmt.Sprintf("%d", content)
	switch currentColor {
	case color.FgBlue:
		rtn = color.BlueString(rtn)
	case color.FgHiRed:
		rtn = color.RedString(rtn)
	case color.FgHiGreen:
		rtn = color.GreenString(rtn)
	case color.FgYellow:
		rtn = color.YellowString(rtn)
	case color.FgHiMagenta:
		rtn = color.MagentaString(rtn)
	case color.FgHiCyan:
		rtn = color.CyanString(rtn)
	}

	return rtn
}

const (
	reset = iota
	change
	noChange
)
