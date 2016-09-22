package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mozilla-services/go-bouncer/bouncer"
)

const DefaultLang = "en-US"
const DefaultOS = "win"
const firefoxSHA1ReleaseAliasSuffix = "sha1"
const firefoxSHA1BetaAliasSuffix = "beta-sha1"
const firefoxSHA1ESRAliasSuffix = "esr-sha1"

type xpRelease struct {
	Version string
}

// detects Windows XP and Vista clients
var windowsXPRegex = regexp.MustCompile(`Windows (?:NT 5.1|XP|NT 5.2|NT 6.0)`)

var tBirdWinXPLastRelease = xpRelease{"38.5.0"}
var tBirdWinXPLastBeta = xpRelease{"43.0b1"}

func isWindowsXPUserAgent(userAgent string) bool {
	return windowsXPRegex.MatchString(userAgent)
}

func isNotNumber(r rune) bool {
	return !unicode.IsNumber(r)
}

// a < b = -1
// a == b = 0
// a > b = 1
func compareVersions(a, b string) int {
	if a == b {
		return 0
	}
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i, verA := range aParts {
		if len(bParts) <= i {
			return 1
		}
		verB := bParts[i]

		aInt, err := strconv.Atoi(strings.TrimRightFunc(verA, isNotNumber))
		if err != nil {
			aInt = 0
		}
		bInt, err := strconv.Atoi(strings.TrimRightFunc(verB, isNotNumber))
		if err != nil {
			bInt = 0
		}

		if aInt > bInt {
			return 1
		}
		if aInt < bInt {
			return -1
		}
	}
	return 0
}

func tBirdSha1Product(productSuffix string) string {
	switch productSuffix {
	case "beta", "beta-latest":
		return tBirdWinXPLastBeta.Version
	case "ssl":
		return tBirdWinXPLastRelease.Version + "-ssl"
	case "latest":
		return tBirdWinXPLastRelease.Version
	}

	productSuffixParts := strings.SplitN(productSuffix, "-", 2)
	ver := productSuffixParts[0]

	possibleVersion := tBirdWinXPLastRelease
	if strings.Contains(ver, ".0b") {
		possibleVersion = tBirdWinXPLastBeta
	}

	if compareVersions(ver, possibleVersion.Version) == -1 {
		return productSuffix
	}

	if len(productSuffixParts) == 1 {
		return possibleVersion.Version
	}

	if productSuffixParts[1] == "ssl" {
		return possibleVersion.Version + "-ssl"
	}

	return productSuffix
}

func firefoxSha1Product(productSuffix string) string {
	// Example list of products:
	// Firefox-48.0-Complete
	// Firefox-48.0build1-Complete
	// Firefox-48.0
	// Firefox-48.0-SSL
	// Firefox-48.0-stub
	// Firefox-48.0build1-Partial-47.0build3
	// Firefox-48.0build1-Partial-47.0.1build1
	// Firefox-48.0build1-Partial-48.0b10build1
	// Firefox-48.0-Partial-47.0
	// Firefox-48.0-Partial-47.0.1
	// Firefox-48.0-Partial-48.0b10

	// Example list of aliases:
	// firefox-beta-latest
	// firefox-beta-sha1
	// Firefox-beta-stub
	// firefox-esr-latest
	// firefox-esr-sha1
	// firefox-latest
	// firefox-sha1
	// Firefox-stub

	// Do not touch products ending with "sha1"
	if strings.HasSuffix(productSuffix, "-sha1") {
		return productSuffix
	}

	// Do not touch completes and partials
	if strings.HasSuffix(productSuffix, "-complete") || strings.Contains(productSuffix, "-partial-") {
		return productSuffix
	}
	switch productSuffix {
	// special product manually created for aurora
	case "aurora", "aurora-stub":
		return "aurora-sha1"
	// Bouncer aliases, no version specified
	case "stub", "latest":
		return firefoxSHA1ReleaseAliasSuffix
	case "beta-latest", "beta-stub":
		return firefoxSHA1BetaAliasSuffix
	case "esr-latest":
		return firefoxSHA1ESRAliasSuffix
	}

	productSuffixParts := strings.SplitN(productSuffix, "-", 2)
	ver := productSuffixParts[0]

	if strings.Contains(ver, "esr") {
		return firefoxSHA1ESRAliasSuffix
	} else if strings.Contains(ver, ".0b") {
		return firefoxSHA1BetaAliasSuffix
	}
	// Fallback to release
	return firefoxSHA1ReleaseAliasSuffix
}

func sha1Product(product string) string {
	productParts := strings.SplitN(product, "-", 2)
	if len(productParts) == 1 {
		return product
	}

	if productParts[0] == "firefox" {
		return "firefox-" + firefoxSha1Product(productParts[1])
	}

	if productParts[0] == "thunderbird" {
		return "thunderbird-" + tBirdSha1Product(productParts[1])
	}

	return product
}

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

	CacheTime     time.Duration
	PinnedBaseURL string
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

	mirrorBaseURL, err := b.mirrorBaseURL(sslOnly, lang, locationID)
	if err != nil || mirrorBaseURL == "" {
		return "", err
	}

	locationPath = strings.Replace(locationPath, ":lang", lang, -1)

	return mirrorBaseURL + locationPath, nil
}

func (b *BouncerHandler) mirrorBaseURL(sslOnly bool, lang, locationID string) (string, error) {
	if b.PinnedBaseURL != "" {
		if sslOnly {
			return "https://" + b.PinnedBaseURL, nil
		}
		return "http://" + b.PinnedBaseURL, nil
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

	return mirror.BaseURL, nil
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

	// HACKS
	// If the user is coming from windows xp or vista, send a sha1
	// signed product
	// HACKS
	if (os == "win" || os == "win64") && isWindowsXPUserAgent(req.UserAgent()) {
		product = sha1Product(product)
	}

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

	http.Redirect(w, req, url, 302)
}
