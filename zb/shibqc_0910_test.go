package zb

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRecord(t *testing.T) {
	file, err := os.Open("/Users/charlie/Downloads/shibqc/tomcat/logs/catalina.out")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	userMap := make(map[string]map[string]struct{})
	unique := make(map[string]struct{})

	reader := bufio.NewReader(file)
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		line := string(buf)

		if !strings.Contains(line, "的执行与预期行数不一，预期影响：1行，实际影响：0行,导致事务回滚") {
			continue
		}

		i1 := strings.Index(line, "：")
		i2 := strings.Index(line, ",参数：:")
		i3 := strings.Index(line, "的执行与预期行数不一")
		sql := line[i1+3 : i2]
		// fullsql := line[i1+3 : i3]
		x := line[i2+11 : i3]
		params := strings.Split(x, ":")
		database := sql[7:strings.Index(sql, " set")]

		if _, ok := unique[sql]; !ok {
			unique[sql] = struct{}{}
		}

		var userId string
		if strings.HasPrefix(sql, "update loanasset set fiatFreeze=fiatFreeze-? where UserId=? and marketName=? and fiatFreeze>=?") {
			userId = params[1] // + ": 买入"
		} else if strings.HasPrefix(sql, "update loanasset set fiatAmount=fiatAmount+?,fiatFreeze=fiatFreeze-? where UserId=? and marketName=? and fiatFreeze>=?") {
			userId = params[2] // + ": 撤销买入"
		} else if strings.HasPrefix(sql, "update currencyUser set account=account-?,freeze=freeze+? where UserId=? and currency=? and account>=?") {
			continue
		}

		users, ok := userMap[database]
		if !ok {
			users = make(map[string]struct{})
			userMap[database] = users
		}
		users[userId] = struct{}{}
	}

	for sql := range unique {
		fmt.Printf("SQL: %s\n", sql)
	}

	for _, m := range userMap {
		// fmt.Println(database)
		for userId, _ := range m {
			fmt.Printf("%s,", userId)
		}
	}
}
