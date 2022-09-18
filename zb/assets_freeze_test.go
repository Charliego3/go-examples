package zb

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/whimthen/temp/zb/freeze"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Funds struct {
	ShowName  string `json:"showName"`
	FundsType int    `json:"fundsType"`
}

var (
	marketMap = make(map[string]map[string]string)
	fundsMap  = make(map[string]string)
)

func init() {
	var funds []Funds
	_, err := resty.New().SetTimeout(time.Minute).R().SetResult(&funds).Get("https://api.zb.com/data/v1/coins")
	if err != nil || len(funds) == 0 {
		panic(err)
	}

	for _, f := range funds {
		fundsMap[strconv.Itoa(f.FundsType)] = f.ShowName
	}
}

func TestStatisticsAssets(t *testing.T) {
	paths := []string{
		"/Users/charlie/Downloads/13-pan_qc.log",
		"/Users/charlie/Downloads/13-pan_usdc (1).log",
		"/Users/charlie/Downloads/13-pan_usdt.log",
		"/Users/charlie/Downloads/13-pan_btc.log",
	}
	for _, path := range paths {
		analyze(path)
	}

	blue := color.New(color.FgHiBlue, color.Bold)
	red := color.New(color.FgHiRed, color.Bold, color.Underline)
	for market, users := range marketMap {
		if len(users) == 1 {
			for userKey, accountType := range users {
				userId, fund := getUserFund(userKey)
				fmt.Printf("市场:[%s], 用户: [%s], 币种: [%s]\n",
					blue.Sprint(market), red.Sprint(userId),
					red.Sprint(fund+"("+accountType+")"))
			}
			continue
		}

		fmt.Printf("市场:[%s]\n", blue.Sprint(market))
		for userKey, accountType := range users {
			userId, fund := getUserFund(userKey)
			fmt.Printf("\t 用户: [%s], 币种: [%s]\n", red.Sprint(userId), red.Sprint(fund+"("+accountType+")"))
		}
	}
}

func getUserFund(key string) (string, string) {
	uf := strings.Split(key, "_")
	fund, ok := fundsMap[uf[1]]
	if !ok {
		fund = uf[1]
	}
	return uf[0], fund
}

func analyze(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var market string
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		if len(buf) == 0 {
			continue
		}

		line := string(buf)
		_ = line

		name := getName(line)
		if name != "" {
			market = name
			continue
		}

		table, sql, params := getSqls(line)
		// fmt.Printf("\tTable: %q, SQL: %q, Params: \"%+v\"\n", table, sql, params)

		analyzer := freeze.FindAnalyzer(table, params)
		if analyzer == nil {
			continue
		}

		analyzer.Execute(sql)
		user, accountType := analyzer.Result()
		if user == "" {
			continue
		}

		users, ok := marketMap[market]
		if !ok {
			users = make(map[string]string)
			marketMap[market] = users
		}
		if value, ok := users[user]; !ok {
			users[user] = accountType
		} else if !strings.Contains(value, accountType) {
			users[user] = value + " | " + accountType
		}
	}
}

func getName(line string) string {
	if !strings.Contains(line, " export ") {
		return ""
	}

	s := freeze.WhitespaceReg.Split(line, 2)[0]
	return strings.Split(s, "_")[1]
}

var unique = make(map[string]struct{})

// getSqls returns tableName, sql, params slice
func getSqls(line string) (string, string, []string) {
	i1 := strings.Index(line, "：") + 3
	i2 := strings.Index(line, ",参数：:")
	i3 := strings.Index(line, "的执行与预期行数不一")
	fullsql := line[i1:i3]
	if _, ok := unique[fullsql]; ok {
		return "", "", nil // repeat
	}

	unique[fullsql] = struct{}{}
	sql := strings.ToLower(line[i1:i2])
	params := line[i2+11 : i3]
	return freeze.WhitespaceReg.Split(sql, 3)[1], sql, strings.Split(params, ":")
}
