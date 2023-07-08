package zb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func TestBWCoins(t *testing.T) {
	f, err := excelize.OpenFile("/Users/charlie/Desktop/副本bw资金操作-换币120788.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := f.GetRows("原始表单")
	if err != nil {
		log.Fatal(err)
	}
	coins1 := make(map[string]struct{})
	coins2 := make(map[string]struct{})
	for _, row := range rows {
		name := row[7]
		if name == "" {
			fmt.Println("币种名称为空", row)
			continue
		}

		types := row[4]
		if types == "1" && strings.TrimSpace(row[8]) == "可提原封不动迁移钱包的" {
			if _, ok := coins1[name]; ok {
				fmt.Println("1 - 币种名称重复:", name)
				continue
			}
			coins1[name] = struct{}{}
		} else if types == "2" {
			if _, ok := coins2[name]; ok {
				fmt.Println("2 - 币种名称重复:", name)
				continue
			}
			coins2[name] = struct{}{}
		}
	}

	var names1 []string
	for k := range coins1 {
		names1 = append(names1, k)
	}
	var names2 []string
	for k := range coins2 {
		names2 = append(names2, k)
	}

	fmt.Printf("Coins1: %v\n Coins2: %v\n", names1, names2)
	coins, dbnames := generate1(t, 1, names1)
	c2, n2 := generate1(t, 2, names2)
	coins = append(coins, c2...)
	dbnames = append(dbnames, n2...)

	os.WriteFile("./temp.json", []byte(fmt.Sprintf("[\n\t%s\n]", strings.Join(coins, ",\n\t"))), os.ModePerm)
	fmt.Printf("Names: %v = %d\n", dbnames, len(dbnames))
}

var fundsType = 100

func generate1(t *testing.T, types int, names []string) ([]string, []string) {
	template := getTemplate(types)

	var coins, dbnames []string
	for _, name := range names {
		upperName := strings.ToUpper(name)
		dbName := "zb_" + name
		coins = append(coins, fmt.Sprintf(template, upperName, fundsType, name, upperName, upperName, upperName, dbName))
		dbnames = append(dbnames, dbName)
		fundsType++
	}
	return coins, dbnames
}

func getTemplate(types int) string {
	if types == 2 {
		return `{
		"txUrl": "#",
		"unitTag": "%s",
		"withdrawScale": 6,
		"fundsType": %d,
		"isLastCoin": false,
		"fatherFundsType": 0,
		"isDigtalCoin": true,
		"webUrl": "#",
		"coinArea": 0,
		"databaseKey": "%s",
		"isCoin": true,
		"isCurrencyUser": true,
		"isInnerTransfer": true,
		"isMain": true,
		"minFee": "5",
		"propEnName": "%s",
		"isCanRecharge": false,
		"isCanWithdraw": false,
		"unitDecimal": 8,
		"isNewDigtalCoin": true,
		"isShow": true,
		"minWithdrawAmount": "0.01",
		"propCnName": "%s",
		"propTag": "%s",
		"databasesName": "%s"
	}`
	}

	return `{
		"txUrl": "#",
		"withdrawScale": 6,
		"unitTag": "%s",
		"isLastCoin": false,
		"fatherFundsType": 0,
		"fundsType": %d,
		"isDigtalCoin": true,
		"webUrl": "#",
		"coinArea": 0,
		"databaseKey": "%s",
		"isCoin": true,
		"isCurrencyUser": true,
		"isInnerTransfer": true,
		"isMain": true,
		"minFee": "5",
		"isCanRecharge": false,
		"isCanWithdraw": false,
		"unitDecimal": 8,
		"isNewDigtalCoin": true,
		"isShow": true,
		"minWithdrawAmount": "0.01",
		"propCnName": "%s",
		"propEnName": "%s",
		"propTag": "%s",
		"databasesName": "%s",
		"chainShowName": "ERC20",
		"showLinkCoins": [
			{
				"fatherFundsType": 5,
				"useFatherUrl": true,
				"coinName": "eusdt",
				"isShow": true
			}
		],
		"withdrawCoins": [
			{
				"fatherFundsType": 5,
				"isDefault": true,
				"showName": "ERC20",
				"isOpen": true,
				"name": "eusdt"
			}
		]
	}`
}
