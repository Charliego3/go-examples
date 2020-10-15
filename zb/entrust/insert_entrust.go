package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", "root", "iPYDU0o3MRQOreEW", "tcp", "172.16.100.130", 3306, "zb_usdtqcentrust"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(-1)

	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 5000; i++ {
				entrustId := node.Generate().Int64()

				result, err := db.Exec(`INSERT INTO entrust (entrustId, unitPrice, numbers, totalMoney, completeNumber, completeTotalMoney, sumToWeb, webId, types, userId, status, freezeId, submitTime, feeRate, acctType, needRemoval) 
VALUES (?, 600.000000000000000, 1.000000000000000, 600.000000000000000, 1.000000000000000, 600.000000000000000, 8, 8, 1, 359797, 3, null, 1600843768269, 0.00100000, 0, false)`, entrustId)
				if err != nil {
					log.Println("EntrustId:", entrustId, err)
				}

				affected, _ := result.RowsAffected()

				log.Printf("EntrustId: %d, RowsAffected: %d", entrustId, affected)
				time.Sleep(time.Millisecond * 5)
			}
		}()
	}

	ints := make(chan int)
	<-ints
}
