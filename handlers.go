package main

import (
	"net/http"
	"strings"
)

const DefaultLang = "en-US"
const DefaultOS = "win"

type BouncerHandler struct {
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
