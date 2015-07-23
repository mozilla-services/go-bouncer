package main

import (
	"net/http"
	"strings"
)

const DefaultLang = "en-US"
const DefaultOS = "win"

// BouncerHandler is the primary handler for this application
type BouncerHandler struct {
	db *DB
}

// Url returns the final redirect URL given a lang, os and product
// if the string is == "", no mirror or location was found
func (b *BouncerHandler) URL(lang, os, product string) (string, error) {
	product, err := b.db.AliasFor(product)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (b *BouncerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	queryVals := req.URL.Query()

	os := queryVals.Get("os")
	product := queryVals.Get("product")
	lang := queryVals.Get("lang")

	if product == "" {
		http.Redirect(w, req, "http://www.mozilla.org/", 302)
		return
	}
	if os == "" {
		os = DefaultOS
	}
	if lang == "" {
		lang = DefaultLang
	}

	product = strings.TrimSpace(strings.ToLower(product))
	os = strings.TrimSpace(strings.ToLower(os))
}
