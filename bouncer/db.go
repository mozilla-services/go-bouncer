package bouncer

import (
	"database/sql"
	"strconv"
	"time"

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
            mirror_mirrors.active='1' AND 
            mirror_location_mirror_map.active ='1' AND
            mirror_location_mirror_map.healthy = ? AND
            mirror_mirrors.baseurl LIKE '`+baseURLPrefix+`%'
        ORDER BY rating
	`, lang, locationID, healthy)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]MirrorsResult, 0)
	for rows.Next() {
		var tmp MirrorsResult
		err = rows.Scan(&tmp.ID, &tmp.BaseURL, &tmp.Rating)
		if err != nil {
			return nil, err
		}
		results = append(results, tmp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

type LocationsActiveResult struct {
	ID   string
	Path string
}

// LocationsActive returns all active locations
func (d *DB) LocationsActive(checkNow bool) ([]*LocationsActiveResult, error) {
	sql := `SELECT mirror_locations.id, mirror_locations.path
		FROM mirror_locations
		INNER JOIN mirror_products ON mirror_locations.product_id = mirror_products.id
		WHERE mirror_products.active='1'`

	if checkNow {
		sql += ` AND mirror_products.checknow='1'`
	}

	rows, err := d.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*LocationsActiveResult, 0)
	for rows.Next() {
		tmp := new(LocationsActiveResult)
		err = rows.Scan(&tmp.ID, &tmp.Path)
		if err != nil {
			return nil, err
		}
		results = append(results, tmp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

type MirrorsActiveResult struct {
	ID      string
	BaseURL string
	Rating  string
	Name    string
	IP      string
	Host    string
}

// MirrorsActive returns all active mirrors
func (d *DB) MirrorsActive(checkMirror string) ([]*MirrorsActiveResult, error) {
	params := []interface{}{}
	sql := `SELECT id, baseurl, rating, name,
				FROM mirror_mirrors WHERE active='1'`
	if checkMirror != "" {
		if _, err := strconv.Atoi(checkMirror); err == nil {
			params = []interface{}{checkMirror}
			sql += ` AND id = ?`
		} else {
			params = []interface{}{"%" + checkMirror + "%", "%" + checkMirror + "%"}
			sql += ` AND (baseurl LIKE %?% OR name LIKE %?%`
		}
	} else {
		sql += ` ORDER BY name`
	}

	rows, err := d.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*MirrorsActiveResult, 0)
	for rows.Next() {
		tmp := new(MirrorsActiveResult)
		err = rows.Scan(&tmp.ID, &tmp.BaseURL, &tmp.Rating, &tmp.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, tmp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// MirrorLocationUpdate updates or inserts status for a mirror location
func (d *DB) MirrorLocationUpdate(locationID, mirrorID, active, healthy string) error {
	sql := `INSERT INTO mirror_location_mirror_map
			(location_id, mirror_id, active, healthy) VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE active=VALUES(active), healthy=VALUES(healthy)`

	_, err := d.Exec(sql, locationID, mirrorID, active, healthy)
	return err
}

// MirrorSetHealth updates health for a mirror
func (d *DB) MirrorSetHealth(mirrorID, healthy string) error {
	sql := `UPDATE mirror_location_mirror_map SET healthy = ? WHERE mirror_id = ?`

	_, err := d.Exec(sql, healthy, mirrorID)
	return err
}

// SentryLogInsert insert in to the sentry log
func (d *DB) SentryLogInsert(logDate time.Time, mirrorID, active, rating, reason string) error {
	sql := `INSERT INTO sentry_log (log_date, mirror_id, mirror_active, mirror_rating, reason) VALUES (FROM_UNIXTIME(?), ?, ?, ?, ?)`
	_, err := d.Exec(sql, logDate.Unix(), mirrorID, active, rating, reason)
	return err
}

// MirrorUpdateRating updates mirror rating
func (d *DB) MirrorUpdateRating(mirrorID, rating string) error {
	sql := `UPDATE mirror_mirrors SET rating = ? WHERE id = ?`
	_, err := d.Exec(sql, rating, mirrorID)
	return err
}

// SentryLogUpdateReason updates sentry_log reason
func (d *DB) SentryLogUpdateReason(mirrorID, reason string, logUnixTime int64) error {
	sql := `UPDATE sentry_log SET reason=? WHERE log_date=FROM_UNIXTIME(?) AND mirror_id=?`
	_, err := d.Exec(sql, reason, logUnixTime, mirrorID)
	return err
}
