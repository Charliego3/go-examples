package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kataras/golog"
	"github.com/shopspring/decimal"
	"strings"
)

var (
	entrustDB     *sqlx.DB
	entrustBackDB *sqlx.DB
	strategyDB    *sqlx.DB
)

func connect(url string, types int) bool {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		golog.Error("数据库连接异常", err)
		return false
	}

	err = db.Ping()
	if err != nil {
		golog.Error("数据库检测异常", err)
		return false
	}

	if types == 0 {
		strategyDB = db
	} else if types == 1 {
		entrustDB = db
	} else {
		entrustBackDB = db
	}

	return true
}

func loadRobot(robotId int64) (robot *Robot, ok bool) {
	robot = &Robot{}
	row := strategyDB.QueryRowx("SELECT * FROM robot WHERE id = ?", robotId)
	err := row.StructScan(robot)
	if err != nil {
		golog.Errorf("查询机器人[%d]失败, %v", robotId, err)
		return nil, false
	}
	return robot, true
}

func getGridRecordById(id int64) (record *GridRecord, ok bool) {
	record = &GridRecord{}
	row := strategyDB.QueryRowx("SELECT * FROM gridRecordV2 WHERE id = ?", id)
	err := row.StructScan(record)
	if err != nil {
		golog.Error("Scan GridRecord Err: ", err)
		return nil, false
	}
	return record, true
}

func getHedgeRecord(originId int64) (record *GridRecord, ok bool) {
	record = &GridRecord{}
	row := strategyDB.QueryRowx("SELECT * FROM gridRecordV2 WHERE orignRecordId = ? LIMIT 1", originId)
	err := row.StructScan(record)
	if err != nil {
		golog.Error("Scan GridRecord Err: ", err)
		return nil, false
	}
	return record, true
}

func getNextOriginOrder(id int64, price decimal.Decimal, gridIndex int) (record *GridRecord, ok bool) {
	record = &GridRecord{}
	row := strategyDB.QueryRowx("SELECT * FROM gridRecordV2 WHERE isOrignOrder = TRUE AND orderPrice = ? AND gridIndex = ? AND id > ? LIMIT 1", price.String(), gridIndex, id)
	err := row.StructScan(record)
	if err != nil {
		golog.Error("Scan GridRecord Err: ", err)
		return nil, false
	}
	return record, true
}

func getGridRecordByRobotId(robotId int64) (records []*GridRecord, ok bool) {
	err := strategyDB.Select(&records, "SELECT * FROM gridRecordV2 WHERE robotId = ? AND isOrignOrder = TRUE ORDER BY id LIMIT 100", robotId)
	if err != nil {
		golog.Errorf("查询机器人[%d]的网格记录失败, Err: %v", robotId, err)
		return
	}
	if len(records) > 0 {
		ok = true
	}
	return
}

func getEntrustById(entrustId string) (entrust *Entrust, ok bool) {
	entrust, ok = getEntrustByIdWithDB(entrustDB, entrustId)
	if ok {
		return
	}

	if entrustBackDB == nil {
		backURL := env + fmt.Sprintf(entrustDBName+"_backup", strings.Replace(market, "_", "", -1))
		ok := connect(backURL, 2)
		if !ok {
			return nil, false
		}
	}

	return getEntrustByIdWithDB(entrustBackDB, entrustId)
}

func getEntrustByIdWithDB(db *sqlx.DB, entrustId string) (entrust *Entrust, ok bool) {
	entrust = &Entrust{}
	row := db.QueryRowx("SELECT * FROM entrust WHERE entrustId = ?", entrustId)
	err := row.StructScan(entrust)
	if err != nil {
		return nil, false
	}
	return entrust, true
}
