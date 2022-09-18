package freeze

import (
	"github.com/shopspring/decimal"
	"strings"
)

type Payuser struct {
	params []string
	fund   string

	succ bool
}

func (p *Payuser) UserID() string {
	return p.params[len(p.params)-2]
}

func (p *Payuser) Execute(sql string) {
	if strings.Contains(sql, "set freezeths=freezeths-? where") ||
		strings.Contains(sql, "set eths=eths+?,freezeths=freezeths-? where") {
		p.succ = true
		p.fund = "ETH"
	} else if strings.Contains(sql, "set freez_money=freez_money-? where") ||
		strings.Contains(sql, "set balance_money=balance_money+?,freez_money=freez_money-? where") {
		p.succ = true
		p.fund = "USDT"
	}
}

func (p *Payuser) Numbers() decimal.Decimal {
	return decimal.Decimal{}
}

func (p *Payuser) Fund() string {
	return p.fund
}

func (p *Payuser) Result() (string, string) {
	if !p.succ {
		return "", ""
	}
	return p.UserID() + "_" + p.Fund(), "现货"
}
