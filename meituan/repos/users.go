package repos

import "log"

func FindUser() {
	row := db.QueryRow("SELECT * FROM USER LIMIT 1")
	if row.Err() != nil {
		log.Fatalln(row.Err())
	}

	log.Println(row)
}
