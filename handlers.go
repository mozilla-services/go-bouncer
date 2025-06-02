package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	defaultLang          = "en-US"
	defaultOS            = "win"
	fxPre2024LastNightly = "firefox-nightly-pre2024"
	fxPre2024LastBeta    = "127.0b9"
	fxPre2024LastRelease = "127.0"
	esr115Product        = "firefox-esr115-latest-ssl"
)

var (
	// detects windows 7/8/8.1 clients
	windowsRegexForESR115 = regexp.MustCompile(`Windows (?:NT 6\.(1|2|3))`)
	// this is used to verify the referer header
	allowedReferrerRegexp = regexp.MustCompile(`^https://www\.(mozilla\.org|firefox\.com)/`)
	// detects partner aliases
	fxPartnerAlias = regexp.MustCompile(`^partner-firefox-release-([^-]*)-(.*)-latest$`)
	// detects x64 clients
	win64Regex = regexp.MustCompile(`Win64|WOW64`)
)

func isUserAgentOnlyCompatibleWithESR115(userAgent string) bool {
	return windowsRegexForESR115.MatchString(userAgent)
}

func hasAllowedReferrer(referrer string) bool {
	return allowedReferrerRegexp.MatchString(referrer)
}

func isWin64UserAgent(userAgent string) bool {
	return win64Regex.MatchString(userAgent)
}

// isPre2024StubUserAgent is used to detect stub installers that pin the
// "DigiCert SHA2 Assured ID Code Signing CA" intermediate.
func isPre2024StubUserAgent(userAgent string) bool {
	return userAgent == "NSIS InetBgDL (Mozilla)"
}

func pre2024Product(product string) string {
	productParts := strings.SplitN(product, "-", 2)
	if len(productParts) < 2 {
		return product
	}

	partnerMatch := fxPartnerAlias.FindStringSubmatch(product)
	if partnerMatch != nil && partnerMatch[1] == "unitedinternet" {
		return "firefox-" + fxPre2024LastRelease + "-unitedinternet-" + partnerMatch[2]
	}

	if productParts[0] != "firefox" {
		return product
	}

	switch productParts[1] {
	case "nightly-latest", "nightly-latest-l10n":
		return fxPre2024LastNightly
	case "nightly-latest-ssl", "nightly-latest-l10n-ssl":
		return fxPre2024LastNightly + "-ssl"
	case "beta-latest-ssl":
		return "firefox-" + fxPre2024LastBeta + "-ssl"
	case "beta-latest":
		return "firefox-" + fxPre2024LastBeta
	case "devedition-latest-ssl":
		return "devedition-" + fxPre2024LastBeta + "-ssl"
	case "devedition-latest":
		return "devedition-" + fxPre2024LastBeta
	case "latest-ssl":
		return "firefox-" + fxPre2024LastRelease + "-ssl"
	case "latest":
		return "firefox-" + fxPre2024LastRelease
	}

	return product

}

// HealthResult represents service health
type HealthResult struct {
	DB      bool `json:"db"`
	Healthy bool `json:"healthy"`
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
	db *DB

	CacheTime time.Duration
}

func (h *HealthHandler) check() *HealthResult {
	result := &HealthResult{
		DB:      true,
		Healthy: true,
	}

	err := h.db.Ping()
	if err != nil {
		result.DB = false
		result.Healthy = false
		log.Printf("HealthHandler err: %v", err)
	}
	return result
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
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
	db *DB

	CacheTime          time.Duration
	PinHTTPSHeaderName string
	PinnedBaseURLHttp  string
	PinnedBaseURLHttps string
	StubRootURL        string
}

// URL returns the final redirect URL given a lang, os and product
// if the string is == "", no mirror or location was found
func (b *BouncerHandler) URL(pinHTTPS bool, lang, os, product string) (string, error) {
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
	locationPath = strings.Replace(locationPath, ":lang", lang, -1)

	mirrorBaseURL := "http://" + b.PinnedBaseURLHttp
	if pinHTTPS || sslOnly {
		mirrorBaseURL = "https://" + b.PinnedBaseURLHttps
	}

	return mirrorBaseURL + locationPath, nil
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

func (b *BouncerHandler) shouldPinHTTPS(req *http.Request) bool {
	if b.PinHTTPSHeaderName == "" {
		return false
	}

	return req.Header.Get(b.PinHTTPSHeaderName) == "https"
}

func fromRTAMO(attributionCode string) bool {
	// base64 decode the attribution_code value to see if it matches the RTAMO regex
	// This uses '.' as padding because Bedrock is using this library to encode the values:
	// https://pypi.org/project/querystringsafe-base64/
	var base64Decoder = base64.URLEncoding.WithPadding('.')
	sDec, err := base64Decoder.DecodeString(attributionCode)
	if err != nil {
		log.Printf("Error decoding %s: %s ", attributionCode, err.Error())
		return false
	}
	q, err := url.ParseQuery(string(sDec))
	if err != nil {
		log.Printf("Error parsing the attribution_code query parameter: %s", err.Error())
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
		for _, s := range []string{"win", "win64", "win64-aarch64", "osx"} {
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
	// if there is a referer header from a known allowed site.
	// https://github.com/mozilla-services/go-bouncer/issues/347
	if fromRTAMO(reqParams.AttributionCode) {
		return hasAllowedReferrer(reqParams.Referer)
	}

	return true
}

func (b *BouncerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqParams := BouncerParamsFromValues(req.URL.Query(), req.Header)

	if reqParams.Product == "" {
		http.Redirect(w, req, "https://www.mozilla.org/", http.StatusFound)
		return
	}

	if reqParams.OS == "" {
		reqParams.OS = defaultOS
	}

	if reqParams.Lang == "" {
		reqParams.Lang = defaultLang
	}

	// If attribution_code is set, redirect to the stub service.
	if b.shouldAttribute(reqParams) {
		stubURL := b.stubAttributionURL(reqParams)
		http.Redirect(w, req, stubURL, http.StatusFound)
		return
	}

	// We want to return ESR115 when... the product is for Firefox
	shouldReturnESR115 := strings.HasPrefix(reqParams.Product, "firefox-") &&
		// and the product is _not_ an MSI build
		!strings.Contains(reqParams.Product, "-msi") &&
		// and the product is _not_ a partial or complete update (MAR files)
		!strings.Contains(reqParams.Product, "-partial") &&
		!strings.Contains(reqParams.Product, "-complete") &&
		// and the OS param specifies windows
		strings.HasPrefix(reqParams.OS, "win") &&
		// and the User-Agent says it's a Windows 7/8/8.1 client
		isUserAgentOnlyCompatibleWithESR115(req.UserAgent()) &&
		// and the request doesn't come from an allowed site
		!hasAllowedReferrer(reqParams.Referer)

	// Send the latest compatible ESR product if we detect that this is the best option for the client.
	if shouldReturnESR115 {
		// Override the OS if we detect a x64 client that attempts to get a stub installer.
		if strings.Contains(reqParams.Product, "-stub") && isWin64UserAgent(req.UserAgent()) {
			reqParams.OS = "win64"
		}
		reqParams.Product = esr115Product
	}

	// If the user is an "old" stub installer, send a pre-2024-cert-rotation product.
	if isPre2024StubUserAgent(req.UserAgent()) {
		reqParams.Product = pre2024Product(reqParams.Product)
	}

	url, err := b.URL(b.shouldPinHTTPS(req), reqParams.Lang, reqParams.OS, reqParams.Product)
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

	http.Redirect(w, req, url, http.StatusFound)
}
