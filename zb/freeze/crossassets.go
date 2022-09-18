package freeze

import "github.com/shopspring/decimal"

type Crossassets struct {
	params []string

	accountAdd bool
	accountSub bool
	freezeAdd  bool
	freezeSub  bool

	succ bool
}

func (c *Crossassets) UserID() string {
	return c.params[len(c.params)-3]
}

func (c *Crossassets) Execute(sql string) {
	execFields(sql, func(field, op string) {
		if field == "account" && op == opPlus {
			c.accountAdd = true
		} else if field == "account" && op == opSub {
			c.accountSub = true
		} else if field == "freeze" && op == opPlus {
			c.freezeAdd = true
		} else if field == "freeze" && op == opSub {
			c.freezeSub = true
		}
	})

	c.succ = (c.accountAdd && c.freezeSub) ||
		(c.accountAdd && !c.freezeSub) ||
		(c.accountSub && !c.freezeAdd) ||
		(c.freezeAdd && !c.accountSub) ||
		(c.freezeSub && !c.accountAdd)
}

func (c *Crossassets) Numbers() decimal.Decimal {
	return decimal.Decimal{}
}

func (c *Crossassets) Fund() string {
	return c.params[len(c.params)-2]
}

func (c *Crossassets) Result() (string, string) {
	if !c.succ {
		return "", ""
	}
	return c.UserID() + "_" + c.Fund(), "全仓杠杆"
}
