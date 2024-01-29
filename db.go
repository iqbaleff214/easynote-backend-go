package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func database(uri string) (*sql.DB, error) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}