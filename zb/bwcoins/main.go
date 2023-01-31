package main

import (
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/stretchr/objx"
	"log"
	"os"
	"strings"
)

var (
	bcoins    []BwCoin
	resources map[string]objx.Map
)

func main() {
	readyXlsx()
	readyJson()

	buildNewCoins()

	log.Println("done....")
}

func buildNewCoins() {
	var coins []objx.Map
	for _, coin := range bcoins {
		if r, ok := resources[coin.ID]; !ok {
			log.Printf("%q, %q, %q 找不到币种...\n", coin.Name, coin.NewName, coin.ID)
			continue
		} else {
			upperName := strings.ToUpper(coin.NewName)
			mark := r.Get("mark").String()
			if mark == "" {
				mark = coin.NewName
			}
			m := objx.MSI(
				"isCanWithdraw", coin.Withdraw,
				"isCurrencyUser", true,
				"isNewDigtalCoin", true,
				"databaseKey", coin.NewName,
				"txUrl", "",
				"fundsType", coin.ID,
				"unitTag", upperName,
				"minWithdrawAmount", r.Get("onceDrawLimit").String(),
				"isLastCoin", false,
				"chainShowName", mark,
				"isCanRecharge", false,
				"propTag", upperName,
				"propCnName", upperName,
				"propEnName", upperName,
				"isMain", false,
				"isDigtalCoin", true,
				"withdrawScale", r.Get("defaultDecimal").String(),
				"isInnerTransfer", true,
				"minFee", r.Get("minFee").String(),
				"isShow", true,
				"fatherFundsType", 0,
				"databasesName", "zb_"+coin.NewName,
				"webUrl", "",
				"isCoin", true,
				"coinArea", 0,
				"unitDecimal", r.Get("defaultDecimal").String(),
			)

			if coin.Withdraw {
				m.Set("withdrawCoins", []objx.Map{
					objx.MSI(
						"fatherFundsType", 5,
						"isDefault", false,
						"showName", "ERC20",
						"isOpen", true,
						"name", coin.Name,
					),
				})
			}

			coins = append(coins, m)
		}
	}

	c := map[string][]objx.Map{
		"coins": coins,
	}
	buf, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("/Users/charlie/dev/go/temp/zb/bwcoins/coins.json", buf, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func readyJson() {
	const filepath = "/Users/charlie/dev/go/temp/zb/bwcoins/currecies.json"
	buf, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	slice := objx.MustFromJSONSlice(string(buf))
	resources = make(map[string]objx.Map, len(slice))
	for _, m := range slice {
		resources[m.Get("currencyId").String()] = m
	}
}

func readyXlsx() {
	const sheetName = "原始表单"
	const filepath = "/Users/charlie/Downloads/副本bw资金操作统计1202V3.xlsx"
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range rows {
		if i < 2 {
			continue
		}

		id := row[1]
		name := strings.ReplaceAll(row[3], "无", "")
		newname := row[7]
		if id == "" && name == "" && newname == "" {
			continue
		}

		withdraw := "可提原封不动迁移钱包的" == strings.TrimSpace(row[8])
		bcoins = append(bcoins, BwCoin{
			ID:       id,
			Name:     name,
			NewName:  newname,
			Withdraw: withdraw,
		})
	}
}

type BwCoin struct {
	ID       string
	Name     string
	NewName  string
	Withdraw bool
}
