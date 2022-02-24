package repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

const (
	dbPath = "./hty.db"
)

var db *sqlx.DB

func init() {
	database, err := sqlx.Open("sqlite3", fmt.Sprintf("file:%s?_loc=auto", dbPath))
	if err != nil {
		log.Fatalln(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	db = database
}
