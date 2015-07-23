package main

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
	mirrors, err := testDB.Mirrors(false, "en-US", "1", true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mirrors))

	mirrors, err = testDB.Mirrors(true, "en-US", "1", true)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mirrors))
	assert.Equal(t, "2", mirrors[0].ID)
}
