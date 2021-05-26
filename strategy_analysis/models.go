package main

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type GridRecord struct {
	Id              int64          `db:"id"`
	RobotId         int64          `db:"robotId"`
	UserId          int            `db:"userId"`
	MarkerId        int            `db:"marketId"`
	PlatformName    string         `db:"platformName"`
	IsBuy           Bool           `db:"isBuy"`
	OrderPrice      BigDecimal     `db:"orderPrice"`
	OrderAmount     BigDecimal     `db:"orderAmount"`
	ClientOrderId   string         `db:"clientOrderId"`
	OrderId         sql.NullString `db:"orderId"`
	OrderStatus     int            `db:"orderStatus"`
	OrderTime       Time           `db:"orderTime"`
	TradedAmount    BigDecimal     `db:"tradedAmount"`
	TradedMoney     BigDecimal     `db:"tradedMoney"`
	IsOriginOrder   Bool           `db:"isOrignOrder"`
	OriginRecordId  sql.NullInt64  `db:"orignRecordId"`
	GridIndex       int            `db:"gridIndex"`
	Income4Currency BigDecimal     `db:"income4Currency"`
	Status          int            `db:"status"`
	Income4Symbol   BigDecimal     `db:"income4Symbol"`
	TradeFee        BigDecimal     `db:"tradeFee"`
	TradeTime       Time           `db:"tradeTime"`
	IsIocOrder      Bool           `db:"isIocOrder"`
}

type Entrust struct {
	EntrustId           int64         `db:"entrustId"`
	UnitPrice           BigDecimal    `db:"unitPrice"`
	Numbers             BigDecimal    `db:"numbers"`
	TotalMoney          BigDecimal    `db:"totalMoney"`
	CompleteNumber      BigDecimal    `db:"completeNumber"`
	CompleteTotalNumber BigDecimal    `db:"completeTotalMoney"`
	WebId               int           `db:"webId"`
	SumToWeb            int           `db:"sumToWeb"`
	Types               int           `db:"types"`
	UserId              int           `db:"userId"`
	Status              int           `db:"status"`
	FreezeId            sql.NullInt64 `db:"freezeId"`
	SubmitTime          int64         `db:"submitTime"`
	AcctType            int           `db:"acctType"`
	FeeRate             BigDecimal    `db:"feeRate"`
	NeedRemoval         Bool          `db:"needRemoval"`
}

type Robot struct {
	Id               int64         `db:"id"`
	StrategyId       int           `db:"strategyId"`
	UserId           int           `db:"userId"`
	Name             string        `db:"name"`
	Frequency        sql.NullInt64 `db:"frequency"`
	Status           int           `db:"status"`
	CreateTime       Time          `db:"createTime"`
	UpdateTime       Time          `db:"updateTime"`
	InitialAsset     BigDecimal    `db:"initialAsset"`
	UserName         string        `db:"userName"`
	StrategyType     int           `db:"strategyType"`
	Params           string        `db:"params"`
	StartTime        Time          `db:"startTime"`
	StopTime         Time          `db:"stopTime"`
	GroupId          int           `db:"groupId"`
	Income           BigDecimal    `db:"income"`
	CurrentPercent   BigDecimal    `db:"currentPercent"`
	MarketName       string        `db:"marketName"`
	RunningTime      int64         `db:"runningTime"`
	CoinAmount       BigDecimal    `db:"coinAmount"`
	OriginPrice      BigDecimal    `db:"originPrice"`
	FaitAmount       BigDecimal    `db:"faitAmount"`
	StopPrice        BigDecimal    `db:"stopPrice"`
	ExtractedIncome  BigDecimal    `db:"extractedIncome"`
	CoinFee          BigDecimal    `db:"coinFee"`
	FaitFee          BigDecimal    `db:"faitFee"`
	RemainderCoinFee BigDecimal    `db:"remaindCoinFee"`
	RemainderFaitFee BigDecimal    `db:"remaindFaitFee"`
	BuyFeeCount      int           `db:"buyFeeCount"`
	TotalIncome      BigDecimal    `db:"totalIncome"`
	Fee              BigDecimal    `db:"fee"`
	FloatingIncome   BigDecimal    `db:"floatingIncome"`
	IsBuy            []byte        `db:"isBuy"`
}

type Grid struct {
	GridInfo     map[int]GridInfo `json:"gridInfo"`
	GridAmount   int              `json:"gridAmout"`
	Buy          bool             `json:"buy"`
	LowerPrice   decimal.Decimal  `json:"lowerPrice"`
	UpperPrice   decimal.Decimal  `json:"upperPrice"`
	TriggerPrice decimal.Decimal  `json:"triggerPrice"`
	TotalAmount  decimal.Decimal  `json:"totalAmount"`
	Type         int              `json:"type"`
}

type GridInfo struct {
	Amount decimal.Decimal `json:"amount"`
	Price  decimal.Decimal `json:"price"`
}
