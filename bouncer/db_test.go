package bouncer

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDB *DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = NewDB("root@tcp(127.0.0.1:3306)/bouncer_test")
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	os.Exit(m.Run())
}

func TestAliasFor(t *testing.T) {
	res, err := testDB.AliasFor("firefox-latest")
	assert.NoError(t, err)
	assert.Equal(t, "Firefox", res)

	res, err = testDB.AliasFor("firefox-123")
	assert.NoError(t, err)
	// When the alias is not found, we return the product.
	assert.Equal(t, "firefox-123", res)
}

func TestOSID(t *testing.T) {
	res, err := testDB.OSID("win64")
	assert.NoError(t, err)
	assert.Equal(t, "1", res)
}

func TestProductForLanguage(t *testing.T) {
	res, sslOnly, err := testDB.ProductForLanguage("Firefox", "en-US")
	assert.NoError(t, err)
	assert.False(t, sslOnly)
	assert.Equal(t, "1", res)

	res, sslOnly, err = testDB.ProductForLanguage("Firefox-SSL", "en-US")
	assert.NoError(t, err)
	assert.True(t, sslOnly)
	assert.Equal(t, "2", res)
}

func TestLocation(t *testing.T) {
	// We need some IDs before we can invoke `Location()`.
	productID, _, _ := testDB.ProductForLanguage("Firefox", "en-US")
	osID, _ := testDB.OSID("win64")

	id, path, err := testDB.Location(productID, osID)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	assert.Equal(t, "/firefox/releases/39.0/win64/:lang/Firefox%20Setup%2039.0.exe", path)

	// Unknown location
	_, _, err = testDB.Location("some-product-id", osID)
	assert.Error(t, err, sql.ErrNoRows)
}
