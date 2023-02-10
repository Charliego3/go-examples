package main

import (
	json "github.com/json-iterator/go"
	"github.com/stretchr/objx"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestRepeatFundsType(t *testing.T) {
	const coinsJsonPath = "./bw6_new_coins.json"
	buf, err := os.ReadFile(coinsJsonPath)
	if err != nil {
		log.Fatal(err)
	}

	json, err := objx.FromJSON(string(buf))
	if err != nil {
		log.Fatal(err)
	}

	v := json.Get("coins")
	log.Println(v)
}

func TestUpdateFundsType(t *testing.T) {
	coins, err := objx.FromJSON(string(readJson("coins.json")))
	if err != nil {
		log.Fatal(err)
	}

	zbCoins := map[string]int{
		"btc": 2,
		"ltc": 3,
		"eth": 5,
		"etc": 7,
		"bts": 8,
		"eos": 9,
	}
	fundsTypes := make(map[int]string)

	value := coins.Get("coins").ReplaceObjxMap(func(i int, m objx.Map) objx.Map {
		var fundsType int
		if id, ok := zbCoins[m.Get("databaseKey").String()]; ok {
			fundsType = id
		} else {
			fundsType, err = strconv.Atoi(m.Get("fundsType").String())
			if err != nil {
				log.Fatal(err)
			}
		}

		m.Set("fundsType", fundsType)
		m.Set("isCanRecharge", false)
		m.Set("isCanWithdraw", false)

		if name, ok := fundsTypes[fundsType]; ok {
			log.Println("fundsType 已经存在", fundsType, name, m.Get("databaseKey"))
		} else {
			fundsTypes[fundsType] = m.Get("databaseKey").String()
		}
		return m
	})

	buf, err := json.MarshalIndent(objx.Map{"coins": value.Data()}, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("./coins_new.json", buf, os.ModePerm)
}

func readJson(path string) []byte {
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}
