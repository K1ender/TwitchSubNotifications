package database

import (
	"database/sql"
	"subalertor/config"
	"subalertor/logger"

	_ "github.com/mattn/go-sqlite3"
)

func MustInit(cfg *config.Config) *sql.DB {
	db, err := sql.Open("sqlite3", cfg.Database.File)
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Connected to database")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Database is ready")

	return db
}
