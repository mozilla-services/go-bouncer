package bouncer

import (
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

func TestMirrors(t *testing.T) {
	mirrors, err := testDB.Mirrors(false)
	assert.NoError(t, err)
	assert.Len(t, mirrors, 1)

	mirrors, err = testDB.Mirrors(true)
	assert.NoError(t, err)
	assert.Len(t, mirrors, 1)
	assert.Equal(t, "2", mirrors[0].ID)
}

func LocationsActive(t *testing.T) {
	locations, err := testDB.LocationsActive(false)
	assert.NoError(t, err)
	assert.Len(t, locations, 3)
}

func MirrorsActive(t *testing.T) {
	mirrors, err := testDB.MirrorsActive("")
	assert.NoError(t, err)
	assert.Len(t, mirrors, 2)
}
