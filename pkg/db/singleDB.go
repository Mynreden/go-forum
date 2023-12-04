package db

import "database/sql"

type DB struct {
	db *sql.DB
}

func (db *DB) OpenDB(dsn string) (*DB, error) {
	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}
	db.db = database
	return db, nil
}
