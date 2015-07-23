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

// Location returns the path of the product/os combonation
func (d *DB) Location(productID, osID string) (id, path string, err error) {
	err = d.QueryRow(
		`SELECT id, path FROM mirror_locations
			WHERE product_id = ? AND os_id = ?`,
		productID, osID).Scan(&id, &path)

	return
}

type MirrorsResult struct {
	ID      string
	BaseURL string
	Rating  int
}

// Mirrors returns a list of valid mirrors
//
// If healthyonly is true and there are no mirrors, it will also search with
// healthyonly set to false
func (d *DB) Mirrors(sslOnly bool, lang, locationID string, healthyOnly bool) ([]MirrorsResult, error) {
	baseURLPrefix := "http://"
	if sslOnly {
		baseURLPrefix = "https://"
	}
	healthy := 1
	if !healthyOnly {
		healthy = 0
	}
	rows, err := d.Query(`
      SELECT
            mirror_mirrors.id,
            baseurl,
            rating
        FROM 
            mirror_mirrors
        JOIN
            mirror_location_mirror_map ON mirror_mirrors.id = mirror_location_mirror_map.mirror_id
        LEFT JOIN
            mirror_lmm_lang_exceptions AS lang_exc ON (mirror_location_mirror_map.id = lang_exc.location_mirror_map_id AND NOT lang_exc.language = ?)
        INNER JOIN
            geoip_mirror_region_map ON (geoip_mirror_region_map.mirror_id = mirror_mirrors.id)
        WHERE
            mirror_location_mirror_map.location_id = ? AND
            $cr_sql
            mirror_mirrors.active='1' AND 
            mirror_location_mirror_map.active ='1' AND
            mirror_location_mirror_map.healthy = ? AND
            mirror_mirrors.baseurl LIKE '`+baseURLPrefix+`%'
        ORDER BY rating
	`, lang, locationID, healthy)

	if err == nil {
		return nil, err
	}

	results := make([]MirrorsResult, 0)
	for rows.Next() {
		var tmp MirrorsResult
		err = rows.Scan(&tmp.ID, &tmp.BaseURL, &tmp.Rating)
		if err == nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if healthyOnly && len(results) == 0 {
		return d.Mirrors(sslOnly, lang, locationID, false)
	}

	return results, nil
}
