package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var gDB *sql.DB

func OpenDatabase() error {
	var err error
	gDB, err = sql.Open("sqlite", "store.db")
	if err != nil {
		return err
	}
	DBInitUsers()
	DBInitSessions()
	DBInitProducts()
	return nil
}

func CloseDatabase() error {
	return gDB.Close()
}
