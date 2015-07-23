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
}
