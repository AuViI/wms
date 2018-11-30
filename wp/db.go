package wp

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const Driver = "sqlite3"

func createDB(file string) (*sql.DB, error) {
	db, err := sql.Open(Driver, file)
	if err != nil {
		return nil, err
	}

	db.Exec(`create table wp (id integer not null primary key autoincrement,
		location text not null,
		year integer not null,
		month integer not null,
		day integer not null,
		content text not null)`)

	return db, nil
}

func closeDB(db *sql.DB) {
	db.Close()
}

func getDB(db *sql.DB, r *Request) (*Response, error) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`SELECT content FROM wp
		WHERE location = ? AND year = ? AND month = ? AND day = ?`)
	defer stmt.Close()
	var content string
	var err error
	if r.isDate() {
		err = stmt.QueryRow(r.location, r.year, r.month, r.day).Scan(&content)
	} else {
		t := time.Now().Local()
		err = stmt.QueryRow(r.location, t.Year(), t.Month(), t.Day()).Scan(&content)
	}
	tx.Commit()
	return &Response{content}, err
}

/// TODO error checking
func setDB(db *sql.DB, r *Request, c *Response) error {
	var stmt *sql.Stmt
	var err error

	co, gerr := getDB(db, r)
	tx, _ := db.Begin()
	if gerr == nil {
		fmt.Printf("change '%s' to '%s' at '%s'\n", co.content, c.content, r.String())
		stmt, err = tx.Prepare(`UPDATE wp
			SET content = ?
			WHERE location = ? AND year = ? AND month = ? AND day = ?`)
		stmt.Exec(c.content, r.location, r.year, r.month, r.day)
	} else {
		stmt, err = tx.Prepare(`INSERT INTO wp
			(location, year, month, day, content) VALUES
			(?, ?, ?, ?, ?)`)
		stmt.Exec(r.location, r.year, r.month, r.day, c.content)
	}
	defer stmt.Close()
	tx.Commit()
	return err
}
