package storage

import "database/sql"

type EventSubStore interface{}

type SQLiteEventSubStore struct {
	db *sql.DB
}

func NewSQLiteEventSubStore(db *sql.DB) *SQLiteEventSubStore {
	return &SQLiteEventSubStore{
		db: db,
	}
}
