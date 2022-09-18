package freeze

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

type Loanasset struct {
	params []string

	fiatFreezeSub bool
	fiatFreezeAdd bool
	fiatAmountSub bool
	fiatAmountAdd bool

	fund string
	succ bool
}

func (l *Loanasset) UserID() string {
	return l.params[len(l.params)-3]
}

func (l *Loanasset) Execute(sql string) {
	execFields(sql, func(field, op string) {
		if field == "fiatamount" && op == opPlus {
			l.fiatAmountAdd = true
		} else if field == "fiatamount" && op == opSub {
			l.fiatAmountSub = true
		} else if field == "fiatfreeze" && op == opPlus {
			l.fiatFreezeAdd = true
		} else if field == "fiatfreeze" && op == opSub {
			l.fiatFreezeSub = true
		}

		if strings.HasPrefix(field, "fiat") {
			l.fund = "QC"
		} else if l.fund != "" {
			fmt.Println("逐仓杠杆存在多个币种!!!", l.UserID(), l.Fund())
		} else {
			l.fund = SymbolReg.ReplaceAllString(l.params[len(l.params)-2], "")
		}
	})

	l.succ = (l.fiatAmountAdd && l.fiatFreezeSub) ||
		(l.fiatFreezeSub && !l.fiatAmountAdd) ||
		(l.fiatFreezeAdd && !l.fiatAmountSub) ||
		(l.fiatAmountAdd && !l.fiatFreezeAdd) ||
		(l.fiatAmountSub && !l.fiatFreezeAdd)
}

func (l *Loanasset) Numbers() decimal.Decimal {
	return decimal.Decimal{}
}

func (l *Loanasset) Fund() string {
	return l.fund
}

func (l *Loanasset) market() string {
	return l.params[len(l.params)-2]
}

func (l *Loanasset) Result() (string, string) {
	if !l.succ {
		return "", ""
	}
	return l.UserID() + "_" + l.Fund(), "逐仓杠杆"
}
