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
	product, err := testDB.AliasFor("firefox-latest")
	assert.NoError(t, err)
	assert.Equal(t, "Firefox", product)
}
