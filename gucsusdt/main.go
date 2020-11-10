package main

import (
	"bytes"
	"encoding/csv"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("/Users/nzlong/Downloads/sqlresult_5389199.csv")
	if err != nil {
		panic(err)
	}
	newb := bytes.ReplaceAll(file, []byte("b'0'"), []byte("b0"))
	reader := csv.NewReader(bytes.NewReader(newb[3:]))

	entrustIdBuyIndex, entrustIdSellIndex, _, index := 4, 6, 8, 0
	var entrustIdBuyArr, entrustIdSellArr []string
	const queryEntrust = "SELECT * FROM zb_gucsusdtentrust.entrust WHERE entrustId IN (%s)"

	for {
		fields, err := reader.Read()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if index > 0 {
			entrustIdSellArr = append(entrustIdSellArr, fields[entrustIdSellIndex])
			entrustIdBuyArr = append(entrustIdBuyArr, fields[entrustIdBuyIndex])
		}
		index++
	}

	color.Red(queryEntrust, strings.Join(entrustIdBuyArr, ", "))
	color.Red(queryEntrust, strings.Join(entrustIdSellArr, ", "))

	min, _ := decimal.NewFromString("3")

	// 买单用户ID
	file, err = ioutil.ReadFile("/Users/nzlong/Downloads/master-sql1.csv")
	if err != nil {
		panic(err)
	}
	newb = bytes.ReplaceAll(file, []byte("b'0'"), []byte("b0"))
	reader = csv.NewReader(bytes.NewReader(newb[3:]))

	index = 0
	numbersIndex, userIdIndex := 2, 9
	buyUserIdMap := map[string]struct{}{}
	var buyUserIds []string

	for {
		fields, err := reader.Read()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if index > 0 {
			numbers, _ := decimal.NewFromString(fields[numbersIndex])
			if numbers.LessThan(min) {
				buyUserIdMap[fields[userIdIndex]] = struct{}{}
			}
		}
		index++
	}

	if len(buyUserIdMap) > 0 {
		for userId := range buyUserIdMap {
			buyUserIds = append(buyUserIds, userId)
		}
	}

	color.Blue("买单下单数量小于%s的用户ID: [%+v]", min.String(), strings.Join(buyUserIds, ", "))

	// 卖单用户ID
	file, err = ioutil.ReadFile("/Users/nzlong/Downloads/master-sql2.csv")
	if err != nil {
		panic(err)
	}
	newb = bytes.ReplaceAll(file, []byte("b'0'"), []byte("b0"))
	reader = csv.NewReader(bytes.NewReader(newb[3:]))

	index = 0
	numbersIndex, userIdIndex = 2, 9
	sellUserIdMap := map[string]struct{}{}
	var sellUserIds []string

	for {
		fields, err := reader.Read()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if index > 0 {
			numbers, _ := decimal.NewFromString(fields[numbersIndex])
			if numbers.LessThan(min) {
				sellUserIdMap[fields[userIdIndex]] = struct{}{}
			}
		}
		index++
	}

	if len(sellUserIdMap) > 0 {
		for userId := range sellUserIdMap {
			sellUserIds = append(sellUserIds, userId)
		}
	}

	color.Blue("卖单下单数量小于%s的用户ID: [%+v]", min.String(), strings.Join(sellUserIds, ", "))
}
