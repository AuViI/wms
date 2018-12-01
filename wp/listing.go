package wp

import (
	"bytes"
	"encoding/json"
	"io"
)

type DBRow struct {
	Location string `json:"loc"`
	Year     int64  `json:"y"`
	Month    int64  `json:"m"`
	Day      int64  `json:"d"`
	Content  string `json:"c"`
}

func listJson(w io.Writer) error {
	db, err := createDB("wp.db")
	if err != nil {
		return err
	}
	defer closeDB(db)

	rows, err := db.Query(`SELECT
		location, year, month, day, content
		FROM wp
		ORDER BY year desc, month desc, day desc, location asc
		LIMIT 20;`)

	if err != nil {
		return err
	}

	dbrows := make([]DBRow, 0, 20)

	for rows.Next() {
		var r DBRow
		if err := rows.Scan(&(r.Location), &(r.Year), &(r.Month), &(r.Day), &(r.Content)); err != nil {
			return err
		}
		dbrows = append(dbrows, r)
	}

	b, err := json.Marshal(dbrows)
	if err != nil {
		return err
	}

	io.Copy(w, bytes.NewBuffer(b))

	return nil
}
