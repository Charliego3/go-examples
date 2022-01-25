package repos

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	dbPath = "meituan.db"
)

var db *sql.DB

func init() {
	database, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_loc=auto", dbPath))
	if err != nil {
		log.Fatalln(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	db = database
}
