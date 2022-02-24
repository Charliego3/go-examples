package repos

type User struct {
	ID         int64  `db:"_id"`
	IsLast     bool   `db:"is_last"`
	IsMaster   bool   `db:"is_master"`
	Login      string `db:"login"`
	PoiId      int    `db:"poi_id"`
	PoiName    string `db:"poi_name"`
	PosId      string `db:"pos_id"`
	TenantID   int    `db:"tenant_id"`
	Token      string `db:"token"`
	Validation string `db:"validation"`
}

func FetchUser() (user User, err error) {
	row := db.QueryRowx("SELECT * FROM USER LIMIT 1")
	if err = row.Err(); err != nil {
		return
	}

	err = row.StructScan(&user)
	return
}
