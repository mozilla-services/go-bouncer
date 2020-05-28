package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mozilla-services/go-bouncer/bouncer"
	"github.com/stretchr/testify/assert"
)

var bouncerHandler *BouncerHandler
var bouncerHandlerPinned *BouncerHandler

func init() {
	testDB, err := bouncer.NewDB("root@tcp(127.0.0.1:3306)/bouncer_test")
	if err != nil {
		log.Fatal(err)
	}

	bouncerHandler = &BouncerHandler{
		db:                 testDB,
		StubRootURL:        "https://stub/",
		PinHttpsHeaderName: "X-Forwarded-Proto",
	}
	bouncerHandlerPinned = &BouncerHandler{
		db:                 testDB,
		PinnedBaseURLHttp:  "download-sha1.cdn.mozilla.net/pub",
		PinnedBaseURLHttps: "download-sha1.cdn.mozilla.net/pub",
		PinHttpsHeaderName: "X-Forwarded-Proto",
	}
}

func TestBouncerHandlerAttributionCode(t *testing.T) {
	tests := []struct {
		In  string
		Out string
	}{
		{
			`http://test/?product=Firefox&os=osx&lang=en-US&attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D(not%20set)%26content%3D(not%20set)&attribution_sig=anhmacsig`,
			`http://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg`,
		},
		{
			`http://test/?product=Firefox&os=win&lang=en-US&attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D(not%20set)%26content%3D(not%20set)&attribution_sig=anhmacsig`,
			`https://stub/?attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D%28not+set%29%26content%3D%28not+set%29&attribution_sig=anhmacsig&lang=en-US&os=win&product=firefox`,
		},
		{
			`http://test/?product=Firefox-stub&os=win&lang=en-US&attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D(not%20set)%26content%3D(not%20set)&attribution_sig=anhmacsig`,
			`https://stub/?attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D%28not+set%29%26content%3D%28not+set%29&attribution_sig=anhmacsig&lang=en-US&os=win&product=firefox-stub`,
		},
	}
	for _, test := range tests {
		w := httptest.NewRecorder()

		req, err := http.NewRequest("GET", test.In, nil)
		assert.NoError(t, err)

		bouncerHandler.ServeHTTP(w, req)
		assert.Equal(t, 302, w.Code)
		assert.Equal(t, test.Out, w.HeaderMap.Get("Location"))
	}
}

func TestBouncerHandlerParams(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://test/?os=mac&lang=en-US", nil)
	assert.NoError(t, err)

	bouncerHandler.ServeHTTP(w, req)
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "https://www.mozilla.org/", w.HeaderMap.Get("Location"))
}

func TestBouncerShouldPinHttps(t *testing.T) {
	bouncerHandler.PinHttpsHeaderName = ""
	req, err := http.NewRequest("GET", "http://test/?product=firefox-latest&os=osx&lang=en-US", nil)
	assert.NoError(t, err)
	assert.Equal(t, false, bouncerHandler.shouldPinHttps(req))

	req.Header.Set("X-Forwarded-Proto", "https")
	assert.Equal(t, false, bouncerHandler.shouldPinHttps(req))

	bouncerHandler.PinHttpsHeaderName = "X-Forwarded-Proto"

	assert.Equal(t, true, bouncerHandler.shouldPinHttps(req))

	req.Header.Set("X-Forwarded-Proto", "http")
	assert.Equal(t, false, bouncerHandler.shouldPinHttps(req))

	req.Header.Del("X-Forwarded-Proto")
	assert.Equal(t, false, bouncerHandler.shouldPinHttps(req))
}

func TestBouncerHandlerPrintQuery(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://test/?product=firefox-latest&os=osx&lang=en-US&print=yes", nil)
	assert.NoError(t, err)

	bouncerHandler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "http://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg", w.Body.String())
}

func TestBouncerHandlerPinnedValid(t *testing.T) {
	defaultUA := "Mozilla/5.0 (Windows NT 7.0; rv:10.0) Gecko/20100101 Firefox/43.0"
	testRequests := []struct {
		URL              string
		ExpectedLocation string
		UserAgent        string
		XForwardedProto  string
	}{
		{"http://test/?product=firefox-latest&os=osx&lang=en-US", "http://download-sha1.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg", defaultUA, "http"},
		{"http://test/?product=firefox-latest&os=win64&lang=en-US", "http://download-sha1.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA, "http"},
		{"http://test/?product=firefox-latest&os=win64&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA, "https"},
		{"http://test/?product=Firefox-SSL&os=win64&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA, "http"},
		{"http://test/?product=Firefox-SSL&os=win&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", "http"},  // Windows XP
		{"http://test/?product=Firefox-SSL&os=win&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 6.0; SV1; .NET CLR 2.0.50727)", "http"},  // Windows Vista
		{"http://test/?product=Firefox-SSL&os=win64&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", "http"},    // Windows XP 64 bit - should get normal win64 build
		{"http://test/?product=Firefox-stub&os=win&lang=en-US", "https://download-sha1.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", "http"}, // Windows XP no stub
	}

	for _, testRequest := range testRequests {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", testRequest.URL, nil)
		req.Header.Set("X-Forwarded-Proto", testRequest.XForwardedProto)
		assert.NoError(t, err, "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)

		req.Header.Set("User-Agent", testRequest.UserAgent)

		bouncerHandlerPinned.ServeHTTP(w, req)
		assert.Equal(t, 302, w.Code, "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)
		assert.Equal(t, testRequest.ExpectedLocation, w.HeaderMap.Get("Location"), "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)
	}
}

func TestBouncerHandlerValid(t *testing.T) {
	defaultUA := "Mozilla/5.0 (Windows NT 7.0; rv:10.0) Gecko/20100101 Firefox/43.0"
	testRequests := []struct {
		URL              string
		ExpectedLocation string
		UserAgent        string
	}{
		{"http://test/?product=firefox-latest&os=osx&lang=en-US", "http://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg", defaultUA},
		{"http://test/?product=firefox-latest&os=win64&lang=en-US", "http://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-SSL&os=win64&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-SSL&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)"},  // Windows XP
		{"http://test/?product=Firefox-SSL&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 6.0; SV1; .NET CLR 2.0.50727)"},  // Windows Vista
		{"http://test/?product=Firefox-SSL&os=win64&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)"},    // Windows XP 64 bit - should get normal win64 build
		{"http://test/?product=Firefox-stub&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe", "Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)"}, // Windows XP no stub
	}

	for _, testRequest := range testRequests {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", testRequest.URL, nil)
		assert.NoError(t, err, "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)

		req.Header.Set("User-Agent", testRequest.UserAgent)

		bouncerHandler.ServeHTTP(w, req)
		assert.Equal(t, 302, w.Code, "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)
		assert.Equal(t, testRequest.ExpectedLocation, w.HeaderMap.Get("Location"), "url: %v ua: %v", testRequest.URL, testRequest.UserAgent)
	}
}

func TestIsWindowsXPUserAgent(t *testing.T) {
	uas := []struct {
		UA   string
		IsXP bool
	}{
		{"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0", true},                                            // firefox XP
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0", false},                                    // firefox non-XP
		{"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36", true},         // Chrome XP
		{"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36", false}, // Chrome non-XP
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", true},                                  // IE XP
		{"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP)", true},                                                               // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", true},                                                    // IE Vista
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.1; en-US)", false},                                                   // IE non-XP
	}
	for _, ua := range uas {
		assert.Equal(t, ua.IsXP, isWindowsXPUserAgent(ua.UA), "ua: %v", ua.UA)
	}
}

func TestSha1Product(t *testing.T) {
	// Ignore products ending with sha1
	assert.Equal(t, "firefox-something-sha1", sha1Product("firefox-something-sha1"))
	assert.Equal(t, "firefox-45.0-sha1", sha1Product("firefox-45.0-sha1"))
	assert.Equal(t, "firefox-45.0.2-sha1", sha1Product("firefox-45.0.2-sha1"))
	assert.Equal(t, "firefox-49.0b1-sha1", sha1Product("firefox-49.0b1-sha1"))
	assert.Equal(t, "firefox-49.0b2-sha1", sha1Product("firefox-49.0b2-sha1"))
	assert.Equal(t, "firefox-45.0esr-sha1", sha1Product("firefox-45.0esr-sha1"))
	assert.Equal(t, "firefox-45.0.2esr-sha1", sha1Product("firefox-45.0.2esr-sha1"))
	assert.Equal(t, "firefox-45.1.0esr-sha1", sha1Product("firefox-45.1.0esr-sha1"))
	assert.Equal(t, "firefox-45.1.2esr-sha1", sha1Product("firefox-45.1.2esr-sha1"))

	// Ignore partials and completes
	assert.Equal(t, "firefox-42.0.0-complete", sha1Product("firefox-42.0.0-complete"))
	assert.Equal(t, "firefox-48.0-partial-41.0.2build1", sha1Product("firefox-48.0-partial-41.0.2build1"))
	assert.Equal(t, "firefox-43.0.2-complete", sha1Product("firefox-43.0.2-complete"))
	assert.Equal(t, "firefox-44.0-complete", sha1Product("firefox-44.0-complete"))
	assert.Equal(t, "firefox-45.0b1-complete", sha1Product("firefox-45.0b1-complete"))
	assert.Equal(t, "firefox-48.0-partial-42.0b1", sha1Product("firefox-48.0-partial-42.0b1"))
	assert.Equal(t, "firefox-48.0b9-partial-48.0b1", sha1Product("firefox-48.0b9-partial-48.0b1"))

	// ignore product wihtout dashes
	assert.Equal(t, "firefox", sha1Product("firefox"))

	// Aliases with no version specified
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-latest"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-beta-latest"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-beta-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-esr-latest"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-esr-stub"))

	// Aurora is special a bit
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-aurora"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-aurora-stub"))

	// Beta versions
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0b1"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-49.0b8"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0b1-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-49.0b8-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0b1-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-49.0b8-ssl"))

	// ESR
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0esr"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0.1esr"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.0esr"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.1esr"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0esr-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0.1esr-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.0esr-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.1esr-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0esr-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.0.1esr-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.0esr-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-45.3.1esr-ssl"))

	// Everything else starting with firefox should go to firefox-sha1
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-42.0.0"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0.1"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0.4"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-42.0.0-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0.1-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0.4-stub"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-42.0.0-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0.1-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0.4-ssl"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-42.0.0-something-new"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0-ssl-something-new"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-48.0.1-ssl-something-new"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0-something-old"))
	assert.Equal(t, "firefox-sha1", sha1Product("firefox-8.0.4-ssl-something-old"))

	assert.Equal(t, "thunderbird-38.5.0", sha1Product("thunderbird-38.6.0"))
	assert.Equal(t, "thunderbird-38.5.0", sha1Product("thunderbird-39.0.0"))

	assert.Equal(t, "thunderbird-38.4.0", sha1Product("thunderbird-38.4.0"))

	assert.Equal(t, "thunderbird-43.0b1", sha1Product("thunderbird-43.0b2"))
	assert.Equal(t, "thunderbird-43.0b1", sha1Product("thunderbird-44.0b1"))

	assert.Equal(t, "thunderbird-42.0b1", sha1Product("thunderbird-42.0b1"))
}

func BenchmarkSha1Product(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sha1Product("firefox-43.0.0")
		sha1Product("firefox-44.0b1")
	}
}
