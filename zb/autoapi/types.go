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
