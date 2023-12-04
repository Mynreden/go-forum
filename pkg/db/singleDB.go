package db

import (
	"database/sql"
	"sync"
)

type DB struct {
	Db *sql.DB
}

var (
	instance *DB
	once     sync.Once
)

func GetSingleDBInstance() *DB {
	once.Do(func() {
		instance = &DB{}
	})
	return instance
}

func (db *DB) OpenDB(dsn string) error {
	var err error
	db.Db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	return db.Db.Ping()
}
