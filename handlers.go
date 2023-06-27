package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mozilla-services/go-bouncer/bouncer"
)

const DefaultLang = "en-US"
const DefaultOS = "win"
const firefoxSHA1ESRAliasSuffix = "sha1"

type xpRelease struct {
	Version string
}

// detects Windows XP and Vista clients
var windowsXPRegex = regexp.MustCompile(`Windows (?:NT 5.1|XP|NT 5.2|NT 6.0)`)

// detects OSX 10.12, 10.11 and 10.12 clients
// Examples:
// Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8
// Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36
// Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:100.0) Gecko/20100101 Firefox/100.0
var deprecatedOSXRegex = regexp.MustCompile(`^Mozilla/5.0 \(Macintosh; Intel Mac OS X 10[\._](12|13|14)`)
var deprecatedOSXPkgProduct = "firefox-esr-next-pkg-latest-ssl"
var deprecatedOSXDmgProduct = "firefox-esr-next-latest-ssl"

// detects Windows 7 / 8 / 8.1 clients
var windows7Regex = regexp.MustCompile(`Windows (?:NT 6.[123])`)
var deprecatedWin7ExeProduct = "firefox-esr-next-latest-ssl"
var deprecatedWin7MsiProduct = "firefox-esr-next-msi-latest-ssl"

// detects 64-bit windows
var win64Regex = regexp.MustCompile(`; (?:Win|WOW)64`)

var tBirdWinXPLastRelease = xpRelease{"38.5.0"}
var tBirdWinXPLastBeta = xpRelease{"43.0b1"}

func isDeprecatedOSXAgent(userAgent string) bool {
	return deprecatedOSXRegex.MatchString(userAgent)
}

func isWindowsXPUserAgent(userAgent string) bool {
	return windowsXPRegex.MatchString(userAgent)
}

func isWindows7UserAgent(userAgent string) bool {
	return windows7Regex.MatchString(userAgent)
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
	return firefoxSHA1ESRAliasSuffix
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

func osxEsrProduct(product string) string {
	if product == "firefox-pkg-latest-ssl" {
		return deprecatedOSXPkgProduct
	}
	if product == "firefox-latest-ssl" {
		return deprecatedOSXDmgProduct
	}
	return product
}

func win7EsrProduct(product string) string {
	if product == "firefox-latest-ssl" {
		return deprecatedWin7ExeProduct
	}
	if product == "firefox-msi-latest-ssl" {
		return deprecatedWin7MsiProduct
	}
	if product == "firefox-stub" {
		return deprecatedWin7ExeProduct
	}
	return product
}

func win7EsrOS(product string, ua string) string {
	if product != "firefox-stub" {
		return ""
	}
	if win64Regex.MatchString(ua) {
		return "win64"
	}
	return "win"
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

	CacheTime          time.Duration
	PinHttpsHeaderName string
	PinnedBaseURLHttp  string
	PinnedBaseURLHttps string
	StubRootURL        string
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
func (b *BouncerHandler) URL(pinHttps bool, lang, os, product string) (string, error) {
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

	_, locationPath, err := b.db.Location(productID, osID)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	}

	mirrorBaseURL, err := b.mirrorBaseURL(pinHttps || sslOnly)
	if err != nil || mirrorBaseURL == "" {
		return "", err
	}

	locationPath = strings.Replace(locationPath, ":lang", lang, -1)

	return mirrorBaseURL + locationPath, nil
}

func (b *BouncerHandler) mirrorBaseURL(sslOnly bool) (string, error) {
	if b.PinnedBaseURLHttps != "" && sslOnly {
		return "https://" + b.PinnedBaseURLHttps, nil
	}

	if b.PinnedBaseURLHttp != "" && !sslOnly {
		return "http://" + b.PinnedBaseURLHttp, nil
	}

	mirrors, err := b.db.Mirrors(sslOnly)
	if err != nil {
		return "", err
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

func (b *BouncerHandler) stubAttributionURL(reqParams *BouncerParams) string {
	query := url.Values{}
	query.Set("lang", reqParams.Lang)
	query.Set("os", reqParams.OS)
	query.Set("product", reqParams.Product)
	query.Set("attribution_code", reqParams.AttributionCode)
	query.Set("attribution_sig", reqParams.AttributionSig)

	return b.StubRootURL + "?" + query.Encode()
}

func (b *BouncerHandler) shouldPinHttps(req *http.Request) bool {
	if b.PinHttpsHeaderName == "" {
		return false
	}

	return req.Header.Get(b.PinHttpsHeaderName) == "https"
}

func fromRTAMO(attribution_code string) bool {
	// base64 decode the attribution_code value to see if it matches the RTAMO regex
	// This uses '.' as padding because Bedrock is using this library to encode the values:
	// https://pypi.org/project/querystringsafe-base64/
	var base64Decoder = base64.URLEncoding.WithPadding('.')
	sDec, err := base64Decoder.DecodeString(attribution_code)
	if err != nil {
		log.Printf("Error decoding %s: %s ", attribution_code, err.Error())
		return false
	}
	q, err := url.ParseQuery(string(sDec))
	if err != nil {
		log.Printf("Error parsing the attribution_code query parameters: %s", err.Error())
		return false
	}

	content := q.Get("content")
	matched, err := regexp.MatchString(`^rta:`, content)
	if err != nil {
		log.Printf("Error matching RTAMO regex: %s", err.Error())
		return false
	}
	if matched {
		return true
	}
	return false
}

func (b *BouncerHandler) shouldAttribute(reqParams *BouncerParams) bool {
	validOs := func() bool {
		// Only include windows.
		for _, s := range []string{"win", "win64", "win64-aarch64"} {
			if reqParams.OS == s {
				return true
			}
		}
		return false
	}

	if b.StubRootURL == "" {
		return false
	}

	if reqParams.AttributionCode == "" {
		return false
	}
	if reqParams.AttributionSig == "" {
		return false
	}

	if !validOs() {
		return false
	}

	// Exclude updates, MSI, and MSIX installers
	// Technically, -msi covers -msix as well, but both are here to
	// prevent a future footgun where -msi is removed, but we still
	// need -msix covered.
	for _, s := range []string{"-partial", "-complete", "-msi", "-msix"} {
		if strings.Contains(reqParams.Product, s) {
			return false
		}
	}

	// Check if the request is coming from RTAMO, and if so, only attribute
	// if there is a referer header from https://www.mozilla.org/
	// https://github.com/mozilla-services/go-bouncer/issues/347
	if fromRTAMO(reqParams.AttributionCode) {
		refererMatch, err := regexp.MatchString(`^https://www.mozilla.org/`, reqParams.Referer)
		if err != nil {
			log.Printf("Error matching www.mozilla.org regex: %s", err.Error())
			return false
		}

		if !refererMatch {
			return false
		} else {
			return true
		}
	}

	return true
}

func (b *BouncerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqParams := BouncerParamsFromValues(req.URL.Query(), req.Header)

	if reqParams.Product == "" {
		http.Redirect(w, req, "https://www.mozilla.org/", 302)
		return
	}

	if reqParams.OS == "" {
		reqParams.OS = DefaultOS
	}
	if reqParams.Lang == "" {
		reqParams.Lang = DefaultLang
	}

	isWinXpClient := isWindowsXPUserAgent(req.UserAgent())

	// HACKS
	// If the user is coming from 32-bit windows xp or vista, send a sha1 signed product.
	// If the user is coming from an old version of OSX, change their product to ESR
	// If the user is coming from windows 7/8/8.1, change their product to ESR
	// HACKS
	if reqParams.OS == "win" && isWinXpClient {
		reqParams.Product = sha1Product(reqParams.Product)
	} else if reqParams.OS == "osx" && isDeprecatedOSXAgent(req.UserAgent()) {
		reqParams.Product = osxEsrProduct(reqParams.Product)
	} else if strings.HasPrefix(reqParams.OS, "win") && isWindows7UserAgent(req.UserAgent()) {
		reqParams.Product = win7EsrProduct(reqParams.Product)
		os := win7EsrOS(reqParams.Product, req.UserAgent())
		if os != "" {
			reqParams.OS = os
		}
	}

	// If the client is not WinXP and attribution_code is set, redirect to the stub service
	if b.shouldAttribute(reqParams) && !isWinXpClient {
		stubURL := b.stubAttributionURL(reqParams)
		http.Redirect(w, req, stubURL, 302)
		return
	}

	url, err := b.URL(b.shouldPinHttps(req), reqParams.Lang, reqParams.OS, reqParams.Product)
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
	if reqParams.PrintOnly {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(url))
		return
	}

	http.Redirect(w, req, url, 302)
}
