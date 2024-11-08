package bouncer

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// DB is a DB instance for running queries against the bouncer database
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

// AliasFor returns the alias for a product
//
// For example firefox-latest will resolve to the latest
// version of firefox.
func (d *DB) AliasFor(product string) (related string, err error) {
	err = d.QueryRow(
		"SELECT related_product FROM mirror_aliases WHERE alias = ?",
		product).Scan(&related)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		return "", err
	}
	return
}

// OSID returns the id of an operation system, by name
func (d *DB) OSID(name string) (id string, err error) {
	err = d.QueryRow(
		"SELECT id FROM mirror_os WHERE name = ?",
		name).Scan(&id)

	return
}

func (d *DB) ProductForLanguage(product, lang string) (productID string, sslOnly bool, err error) {
	sslInt := 0
	err = d.QueryRow(
		`SELECT prod.id, prod.ssl_only FROM mirror_products AS prod
		LEFT JOIN mirror_product_langs AS langs ON (prod.id = langs.product_id)
		WHERE prod.name LIKE ?
		AND (langs.language LIKE ? OR langs.language IS NULL)`,
		product, lang).Scan(&productID, &sslInt)

	if sslInt == 1 {
		sslOnly = true
	} else {
		sslOnly = false
	}
	return
}

// Location returns the path of the product/os combination
func (d *DB) Location(productID, osID string) (id, path string, err error) {
	err = d.QueryRow(
		`SELECT id, path FROM mirror_locations
			WHERE product_id = ? AND os_id = ?`,
		productID, osID).Scan(&id, &path)

	return
}
