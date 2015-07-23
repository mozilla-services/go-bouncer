package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

func (d *DB) AliasFor(product string) (string, error) {
	related := ""
	err := d.QueryRow(
		"SELECT related_product FROM mirror_aliases WHERE alias = ?",
		product).Scan(&related)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		return "", err
	}
	return related, nil
}
