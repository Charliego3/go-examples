package repos

type Table struct {
	ID       int    `db:"_id"`
	AreaID   int    `db:"area_id"`
	Name     string `db:"name"`
	PoiID    int    `db:"poi_id"`
	Rank     int    `db:"rank"`
	Seats    int    `db:"seats"`
	TableID  int    `db:"table_id"`
	TenantID int    `db:"tenant_id"`
}

func FetchTables(areaID int) (tables []Table, err error) {
	rows, err := db.Queryx("SELECT * FROM tables WHERE area_id = ? ORDER BY rank", areaID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table Table
		err := rows.StructScan(&table)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}
	return
}
