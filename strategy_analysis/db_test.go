package main

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestGetGridRecordByRobotId(t *testing.T) {
	initDB(t)

	records, ok := getGridRecordByRobotId(10, 10)

	if !ok {
		t.Fatal("获取不到数据")
	}

	for _, record := range records {
		t.Logf("R: %d, GridIndex: %d", record.Id, record.GridIndex)
	}
}

func initDB(t *testing.T) {
	var err error
	strategyDB, err = sqlx.Open("mysql", "root:iPYDU0o3MRQOreEW@tcp(172.16.100.130:3306)/strategy")

	if err != nil {
		t.Fatal("数据库连接异常")
	}

	err = strategyDB.Ping()
	if err != nil {
		t.Fatal("Ping DB error:", err)
		return
	}
}
