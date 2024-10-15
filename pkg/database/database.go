package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(connstr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connstr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
