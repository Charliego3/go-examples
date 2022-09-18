package freeze

import (
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
)

const (
	opPlus = "+"
	opSub  = "-"
)

var (
	WhitespaceReg = regexp.MustCompile("\\s")
	SymbolReg     = regexp.MustCompile("(qc|usdt|btc|usdc)$")
)

type Analyzer interface {
	UserID() string
	Execute(sql string)
	Numbers() decimal.Decimal
	Fund() string
	Result() (string, string)
}

func FindAnalyzer(table string, params []string) Analyzer {
	switch strings.ToLower(table) {
	case "currencyuser":
		return &CurrencyUser{params: params}
	case "loanasset":
		return &Loanasset{params: params}
	case "crossassets":
		return &Crossassets{params: params}
	case "pay_user":
		return &Payuser{params: params}
	}
	return nil
}

func execFields(sql string, fn func(field, op string)) {
	i1 := strings.Index(sql, "set ") + 4
	i2 := strings.Index(sql, " where ")
	subSQL := sql[i1:i2]
	subSQL = WhitespaceReg.ReplaceAllString(subSQL, "")
	fields := strings.Split(subSQL, ",")
	for _, field := range fields {
		arr := strings.Split(field, "=")
		length := len(arr[1])
		fn(arr[0], arr[1][length-2:length-1])
	}
}
