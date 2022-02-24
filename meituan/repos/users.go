package repos

import (
	"github.com/sirupsen/logrus"
	"log"
)

func FetchUser() {
	row := db.QueryRowx("SELECT * FROM USER LIMIT 1")
	if row.Err() != nil {
		log.Fatalln(row.Err())
	}

	user := make(map[string]interface{})
	err := row.MapScan(user)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("User: %+v", user)
}
