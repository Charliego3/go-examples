package autoapi

type AccountType string
type OrderType string
type TradeType string

const (
	AccountTypeMain  AccountType = "0"
	AccountTypeLever AccountType = "1"
	AccountTypeCross AccountType = "2"
)

const (
	OrderTypeLimit    OrderType = "0"
	OrderTypePostOnly OrderType = "1"
	OrderTypeIoc      OrderType = "2"
)

const (
	TradeTypeSell         TradeType = "0"
	TradeTypeBuy          TradeType = "1"
	TradeTypePostOnlySell TradeType = "2"
	TradeTypePostOnlyBuy  TradeType = "3"
	TradeTypeIocSell      TradeType = "4"
	TradeTypeIocBuy       TradeType = "5"
)

func (tt TradeType) String() string {
	switch tt {
	case TradeTypeSell:
		return "限价卖出"
	case TradeTypeBuy:
		return "限价买入"
	case TradeTypePostOnlySell:
		return "PostOnly 卖出"
	case TradeTypePostOnlyBuy:
		return "PostOnly 买入"
	case TradeTypeIocSell:
		return "IOC 卖出"
	case TradeTypeIocBuy:
		return "IOC 买入"
	}
	return "--"
}

func TradeTypeByInt(types int) TradeType {
	switch types {
	case 0:
		return TradeTypeSell
	case 1:
		return TradeTypeBuy
	case 2:
		return TradeTypePostOnlySell
	case 3:
		return TradeTypePostOnlyBuy
	case 4:
		return TradeTypeIocSell
	case 5:
		return TradeTypeIocBuy
	}
	return TradeTypeBuy
}

func ReverseTradeType(types int) TradeType {
	switch types {
	case 0:
		return TradeTypeBuy
	case 1:
		return TradeTypeSell
	case 2:
		return TradeTypePostOnlyBuy
	case 3:
		return TradeTypePostOnlySell
	case 4:
		return TradeTypeIocBuy
	case 5:
		return TradeTypeIocSell
	}
	return TradeTypeBuy
}
