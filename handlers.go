package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PagerDuty/godspeed"
	"github.com/mozilla-services/go-bouncer/bouncer"
)

const DefaultLang = "en-US"
const DefaultOS = "win"

// HealthResult represents service health
type HealthResult struct {
	DB      bool   `json:"db"`
	Healthy bool   `json:"healthy"`
	Version string `json:"version"`
}

// JSON returns json string
func (h *HealthResult) JSON() []byte {
	res, err := json.Marshal(h)
	if err != nil {
		log.Printf("HealthResult.JSON err: %v", err)
		return []byte{}
	}
	return res
}

// HealthHandler returns 200 if the app looks okay
type HealthHandler struct {
	db *bouncer.DB

	CacheTime time.Duration
}

func (h *HealthHandler) check() *HealthResult {
	result := &HealthResult{
		DB:      true,
		Healthy: true,
		Version: bouncer.Version,
	}

	err := h.db.Ping()
	if err != nil {
		result.DB = false
		result.Healthy = false
		log.Printf("HealthHandler err: %v", err)
	}
	return result
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.CacheTime > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", h.CacheTime/time.Second))
	}

	w.Header().Set("Content-Type", "application/json")

	result := h.check()
	if !result.Healthy {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result.JSON())
}

// BouncerHandler is the primary handler for this application
type BouncerHandler struct {
	db *bouncer.DB

	CacheTime time.Duration

	godspeedChan chan []string
}

// SetGodspeed sets the handlers godspeed client
func (h *BouncerHandler) SetGodspeed(gs *godspeed.Godspeed, buffer int) {
	h.godspeedChan = make(chan []string, buffer)
	go func(gs *godspeed.Godspeed, c chan []string) {
		for tags := range h.godspeedChan {
			err := gs.Incr("bouncer.download", tags)
			if err != nil {
				log.Println(err)
			}

		}
	}(gs, h.godspeedChan)
}

func (h *BouncerHandler) downloadMetric(product, lang, os string) {
	if h.godspeedChan == nil {
		return
	}

	tags := []string{
		"product:" + product,
		"language:" + lang,
		"os:" + os,
	}
	select {
	case h.godspeedChan <- tags:
	default:
		log.Println("downloadMetric: godspeed buffer full sending: %v", tags)
	}
}

func randomMirror(mirrors []bouncer.MirrorsResult) *bouncer.MirrorsResult {
	totalRatings := 0
	for _, m := range mirrors {
		totalRatings += m.Rating
	}
	for _, m := range mirrors {
		// Intn(x) returns from [0,x) and we need [1,x], so adding 1
		rand := rand.Intn(totalRatings) + 1
		if rand <= m.Rating {
			return &m
		}
		totalRatings -= m.Rating
	}

	// This shouldn't happen
	if len(mirrors) == 0 {
		return nil
	}
	return &mirrors[0]
}

// URL returns the final redirect URL given a lang, os and product
// if the string is == "", no mirror or location was found
func (b *BouncerHandler) URL(lang, os, product string) (string, error) {
	product, err := b.db.AliasFor(product)
	if err != nil {
		return "", err
	}

	osID, err := b.db.OSID(os)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	}

	productID, sslOnly, err := b.db.ProductForLanguage(product, lang)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	}

	locationID, locationPath, err := b.db.Location(productID, osID)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	}

	mirrors, err := b.db.Mirrors(sslOnly, lang, locationID, true)
	if err != nil {
		return "", err
	}

	if len(mirrors) == 0 {
		// try again, looking for unhealthy mirrors
		mirrors, err = b.db.Mirrors(sslOnly, lang, locationID, false)
		if err != nil {
			return "", err
		}
	}

	if len(mirrors) == 0 {
		return "", nil
	}

	mirror := randomMirror(mirrors)
	if mirror == nil {
		return "", nil
	}

	locationPath = strings.Replace(locationPath, ":lang", lang, -1)

	return mirror.BaseURL + locationPath, nil
}

func (b *BouncerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	queryVals := req.URL.Query()

	printOnly := queryVals.Get("print")
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

	url, err := b.URL(lang, os, product)
	if err != nil {
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if url == "" {
		http.NotFound(w, req)
		return
	}

	if b.CacheTime > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", b.CacheTime/time.Second))
	}

	// If ?print=yes, print the resulting URL instead of 302ing
	if printOnly == "yes" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(url))
		return
	}

	b.downloadMetric(product, lang, os)

	http.Redirect(w, req, url, 302)
}
