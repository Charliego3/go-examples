package main

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"os"
	"regexp"
)

var regex = regexp.MustCompile(".*APW::获取到对头后处理:\\[], record:(.*)")

func main() {
	filePath := "/Users/charlie/Downloads/APW.log"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	dealCount := 0
	cancelCOunt := 0
	totalCount := 0

	first := ""
	last := ""

	reader := bufio.NewReader(file)
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}

		line := string(l)
		strs := regex.FindAllStringSubmatch(line, -1)
		if len(strs) == 0 {
			continue
		}

		m := make(map[string]interface{})
		err = jsoniter.Unmarshal([]byte(strs[0][1]), &m)
		if err != nil {
			log.Fatalln(line, err)
		}

		record := m["record"].(map[string]interface{})
		price, err := decimal.NewFromString(fmt.Sprintf("%v", record["unitPrice"]))
		if err != nil {
			log.Fatalln(line, err)
		}

		if first == "" {
			first = line
		}
		last = line

		if price.GreaterThan(decimal.Zero) {
			dealCount++
		} else {
			cancelCOunt++
		}
		totalCount++
	}

	log.Println(first)
	log.Println(last)
	log.Println("资金处理总条数:", totalCount, "成交处理条数:", dealCount, "撤单处理条数:", cancelCOunt)
}
