package repos

type Area struct {
	ID       int    `db:"_id"`
	Name     string `db:"name"`
	PoiID    int    `db:"poi_id"`
	Rank     int    `db:"rank"`
	TenantID int    `db:"tenant_id"`
}

func FetchAreas() (areas []Area, err error) {
	rows, err := db.Queryx("SELECT * FROM tables_area ORDER BY rank")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var area Area
		err := rows.StructScan(&area)
		if err != nil {
			return nil, err
		}

		areas = append(areas, area)
	}
	return
}
