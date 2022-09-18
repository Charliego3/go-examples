package freeze

import (
	"github.com/shopspring/decimal"
)

type CurrencyUser struct {
	params []string

	freezePlusExecuted  bool
	freezeSubExecuted   bool
	accountPlusExecuted bool
	accountSubExecuted  bool

	succ bool
}

func (c *CurrencyUser) UserID() string {
	return c.params[len(c.params)-3]
}

func (c *CurrencyUser) Numbers() decimal.Decimal {
	return decimal.RequireFromString(c.params[0])
}

func (c *CurrencyUser) Fund() string {
	return c.params[len(c.params)-2]
}

func (c *CurrencyUser) Execute(sql string) {
	defer func() {
		c.freezePlusExecuted = false
		c.accountPlusExecuted = false
	}()

	execFields(sql, func(field, op string) {
		if field == "freeze" && op == opPlus {
			c.freezePlusExecuted = true
		} else if field == "freeze" && op == opSub {
			c.freezeSubExecuted = true
		} else if field == "account" && op == opPlus {
			c.accountPlusExecuted = true
		} else if field == "account" && op == opSub {
			c.accountSubExecuted = true
		}
	})

	// if c.accountPlusExecuted && c.freezeSubExecuted { // 撤单
	// 	fmt.Println(c.UserID(), "撤单失败")
	// } else if (c.freezeSubExecuted && !c.accountPlusExecuted) ||
	// 	(c.accountPlusExecuted && !c.freezeSubExecuted) { // 资产处理失败
	// fmt.Println(c.UserID(), c.Fund(), "资金处理失败")
	// c.succ = true
	// }
	c.succ = (c.accountPlusExecuted && c.freezeSubExecuted) ||
		(c.freezeSubExecuted && !c.accountPlusExecuted) ||
		(c.accountPlusExecuted && !c.freezeSubExecuted)
}

func (c *CurrencyUser) Result() (string, string) {
	if !c.succ {
		return "", ""
	}
	return c.UserID() + "_" + c.Fund(), "现货"
}
