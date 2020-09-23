package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/whimthen/kits/logger"
)

func main() {
	db, err := sql.Open("mysql", "root:rootroot@/test")
	if err != nil {
		logger.Fatal("Open mysql error: %+v", err)
	}

	for i := 0; i < 10000; i++ {
		result, err := db.Exec("insert into t(n) values (?)", i+1)
		if err != nil {
			logger.Fatal("Insert error: %+v", err)
		}
		affected, _ := result.RowsAffected()
		insertId, _ := result.LastInsertId()
		logger.Debug("RowsAffected: %d, LastInsertId: %d", affected, insertId)
	}
}
