package main

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/transerver/commons/logger"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	maxPriceLimitRate decimal.Decimal
	minPriceLimitRate decimal.Decimal
)

type Trader struct {
	user User

	ctx    context.Context
	market Market
	dialer *WebsocketDialer

	uif *Receiver
	uf  *Receiver
	uir *Receiver

	funds map[string]*Fund
	mutex sync.RWMutex

	fc     chan Fund
	logger *logger.Logger

	random *rand.Rand
	td     time.Duration
}

func NewTrader(u User, market Market, td time.Duration) *Trader {
	return &Trader{
		user:   u,
		market: market,
		logger: logger.NewLogger(logger.WithPrefix("%d:%s", u.ID, u.Username)),
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		td:     td,
		funds:  make(map[string]*Fund),
	}
}

func (t *Trader) Run(ctx context.Context) {
	t.uif = NewReceiver(UserIncrAssetType)
	t.uf = NewReceiver(UserAssetType)
	t.uir = NewReceiver(UserIncrRecordType)
	t.ctx = ctx
	go t.do()
}

func (t *Trader) Funds() map[string]*Fund {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.funds
}

func (t *Trader) Fund(coin string) Fund {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	f, ok := t.funds[coin]
	if !ok {
		return Fund{}
	}

	return *f
}

func (t *Trader) UpdateFund(coin string, f *Fund) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.funds[coin] = f
}

func (t *Trader) SetFunds(f map[string]*Fund) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.funds = f
}

func (t *Trader) prepare() error {
	t.dialer = NewWebsocketDialerWithSuffix(t.ctx, Settings.Websocket.Address, t.user.Username, t.uif, t.uf, t.uir)

	err := t.dialer.Connect(t.user)
	if err != nil {
		return err
	}

	err = t.dialer.SubscribeUserAsset()
	if err != nil {
		return err
	}

	err = t.dialer.SubscribeUserIncrAsset()
	if err != nil {
		return err
	}

	err = t.dialer.SubscribeIncrRecord(t.market)
	if err != nil {
		return err
	}

	return nil
}

func (t *Trader) do() {
	err := t.prepare()
	if err != nil {
		return
	}

	ticker := time.NewTicker(t.td)
	defer ticker.Stop()

	for {
		select {
		case buf, ok := <-t.uf.C:
			if !ok {
				t.uf.C = nil
				t.logger.Debugf("关闭了全量资金channel")
				break
			}
			t.acceptAssets(buf)

		case buf, ok := <-t.uif.C:
			if !ok {
				t.uif.C = nil
				break
			}

			t.acceptIncrAsset(buf)
			break

		case f, ok := <-t.fc:
			if !ok {
				break
			}
			t.fundCheck(f)

		case <-ticker.C:
			isBuy := t.random.Intn(2)
			price := t.getPrice(isBuy == 1)
			if price.IsZero() {
				break
			}

			number := t.getNumber(price)

			total := price.Mul(number).Truncate(int32(t.market.CurrencyBix + t.market.SymbolBix))
			fund := t.Fund(t.market.Currency)
			if isBuy == 1 && fund.Available.Cmp(total) < 0 {
				// logger.Warnf("准备买入下单 -> 价格: %s, 数量: %s, Available: %s, Total: %s, 买入: %t", price, number, fund.Available, total, isBuy == 1)
				_ = t.recharge(t.market.Currency, total.Mul(oneHundred))
				break
			} else {
				if t.Fund(t.market.Symbol).Available.Cmp(number) < 0 {
					// logger.Warnf("准备卖出下单 -> 价格: %s, 数量: %s, Available: %s, Total: %s, 买入: %t", price, number, fund.Available, total, isBuy == 1)
					_ = t.recharge(t.market.Symbol, total.Mul(oneHundred))
					break
				}
			}

			Order(t.user, t.market, t.logger, price, number, isBuy)

		case buf, ok := <-t.uir.C:
			if !ok {
				t.uir.C = nil
				break
			}
			_ = buf

			// t.logger.Infof("收到增量委托: %s", buf)
			break

		case <-t.ctx.Done():
			t.logger.Debug("退出Trader")
			t.SetFunds(nil)
			return
		}
	}
}

func (t *Trader) getNumber(price decimal.Decimal) decimal.Decimal {
	amount := t.market.MinAmount
	total := amount.Mul(price)
	if total.LessThanOrEqual(t.market.MinExchange) {
		amount = t.market.MinExchange.Div(price).RoundUp(int32(t.market.SymbolBix))
	}

	return amount
}

func (t *Trader) getPrice(isBuy bool) decimal.Decimal {
	depth := Depth()
	if depth.NotValid() {
		t.logger.Warnf("Depth not valid")
		return decimal.Decimal{}
	}

	// min := depth.MinBuy()
	// max := depth.MaxSell()

	min := minPriceLimitRate.Mul(depth.CurrentPrice)
	max := maxPriceLimitRate.Mul(depth.CurrentPrice)

	// lp := min.IntPart()
	// le := min.Mod(one).Coefficient().Int64()
	//
	// hp := max.IntPart()
	// he := max.Mod(one).Coefficient().Int64()

	lp, le := splitDecimal(min)
	hp, he := splitDecimal(max)

	if he < le {
		he, le = le, he
	}

	var number string
	if he == le {
		number = strconv.FormatInt(t.random.Int63n(hp-lp)+lp, 10)
	} else {
		number = fmt.Sprintf("%d.%d", t.random.Int63n(hp-lp)+lp, int32(t.random.Int63n(he-le)+le))
	}

	price, err := decimal.NewFromString(number)
	if err != nil {
		t.logger.Errorf("随机生成价格失败: %v", err)
		return decimal.Decimal{}
	}

	// logger.Debugf("Price: %s, Min: %s, Max: %s, lp: %d, le: %d, hp: %d, he: %d", price, min, max, lp, le, hp, he)

	// buyHigh := maxPriceLimitRate.Mul(depth.CurrentPrice)
	// if isBuy && price.Cmp(buyHigh) > 0 {
	// 	// t.logger.Errorf("买入委托价格异常: %s", price)
	// 	price = buyHigh.RoundDown(int32(t.market.CurrencyBix))
	// } else {
	// 	sellLow := minPriceLimitRate.Mul(depth.CurrentPrice)
	// 	if price.Cmp(sellLow) < 0 {
	// 		// t.logger.Errorf("卖出委托价格异常: %s", price)
	// 		// return decimal.Decimal{}
	// 		price = sellLow.RoundUp(int32(t.market.SymbolBix))
	// 	}
	// }
	return price.Truncate(int32(t.market.CurrencyBix))
}

func splitDecimal(d decimal.Decimal) (int64, int64) {
	ds := strings.SplitN(d.String(), ".", 2)
	d1, err := strconv.ParseInt(ds[0], 0, 64)
	if err != nil {
		return 0, 0
	}

	if len(ds) == 1 {
		return d1, 0
	}

	ds2 := strings.TrimLeftFunc(ds[1], func(r rune) bool {
		return r == '0'
	})
	d2, err := strconv.ParseInt(ds2, 0, 64)
	if err != nil {
		return 0, 0
	}

	return d1, d2
}

func (t *Trader) fundCheck(f Fund) {
	if f.Available.GreaterThan(t.market.MinAmount) {
		return
	}

	amount := t.market.MinAmount.Mul(oneHundred)
	if err := t.recharge(f.Name, amount); err != nil {
		return
	}

	f.Available = f.Available.Add(amount)
	t.UpdateFund(f.Name, &f)
}

func (t *Trader) recharge(coin string, amount decimal.Decimal) error {
	err := Recharge(t.user.ID, coin, amount)
	if err != nil {
		return err
	}

	t.logger.Infof("充值成功 -> %s: %s", coin, amount.String())
	return nil
}

func (t *Trader) acceptIncrAsset(buf []byte) {
	asset, err := t.unmarshalAsset(buf)
	if err != nil {
		t.logger.Errorf("解析增量资金出错: %v", err)
		return
	}

	if len(asset.Coins) <= 0 {
		t.logger.Warnf("收到增量资产变更, 但变更币种为空")
		return
	}

	fund := asset.Coins[0].ToFund()
	t.UpdateFund(fund.Name, fund)
	t.fc <- *fund
	// t.logger.Infof("收到增量资产: %+v", fund)
}

func (t *Trader) acceptAssets(buf []byte) {
	assets, err := t.unmarshalAsset(buf)
	if err != nil {
		t.logger.Errorf("解析全量资金出错: %v", err)
		return
	}

	funds := make(map[string]*Fund)
	for _, coin := range assets.Coins {
		fund := coin.ToFund()
		funds[fund.Name] = fund
		t.logger.Debugf("币种: %s, 余额: %s, 冻结: %s", fund.Name, fund.Available.String(), fund.Freeze.String())
	}

	t.SetFunds(funds)
	err = t.dialer.UnSubscribeUserAsset()
	if err != nil {
		t.logger.Errorf("取消订阅全量资金失败: %v", err)
	}

	t.fc = make(chan Fund, 1000)
	for _, fund := range funds {
		t.fc <- *fund
	}
}

func (t *Trader) unmarshalAsset(buf []byte) (UserAsset, error) {
	var asset UserAsset
	err := json.Unmarshal(buf, &asset)
	if err != nil {
		return UserAsset{}, err
	}

	return asset, nil
}
