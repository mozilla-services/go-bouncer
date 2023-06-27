package main

import (
	"fmt"
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

func TestShouldAttribute(t *testing.T) {
	tests := []struct {
		In  *BouncerParams
		Out bool
	}{
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			true,
		},
		{
			&BouncerParams{
				OS:              "osx",
				Product:         "Firefox",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox",
				AttributionCode: "",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox-partial",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox-complete",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox-msi",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "Firefox-msix",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win64",
				Product:         "Firefox",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			true,
		},
		{
			&BouncerParams{
				OS:              "win64-aarch64",
				Product:         "Firefox",
				AttributionCode: "att-code",
				AttributionSig:  "att-sig",
			},
			true,
		},
		// https://github.com/mozilla-services/go-bouncer/issues/347
		{
			&BouncerParams{
				OS:              "win",
				Product:         "firefox-stub",
				AttributionCode: "c291cmNlPWFkZG9ucy5tb3ppbGxhLm9yZyZtZWRpdW09cmVmZXJyYWwmY2FtcGFpZ249bm9uLWZ4LWJ1dHRvbiZjb250ZW50PXJ0YTplMkk1WkdJeE5tRTBMVFpsWkdNdE5EZGxZeTFoTVdZMExXSTROakk1TW1Wa01qRXhaSDAmZXhwZXJpbWVudD0obm90IHNldCkmdmFyaWF0aW9uPShub3Qgc2V0KSZ1YT1lZGdlJnZpc2l0X2lkPShub3Qgc2V0KQ..",
				AttributionSig:  "att-sig",
				Referer:         "http://otherwebsite.com",
			},
			false,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "firefox-stub",
				AttributionCode: "c291cmNlPWFkZG9ucy5tb3ppbGxhLm9yZyZtZWRpdW09cmVmZXJyYWwmY2FtcGFpZ249bm9uLWZ4LWJ1dHRvbiZjb250ZW50PXJ0YTplMkk1WkdJeE5tRTBMVFpsWkdNdE5EZGxZeTFoTVdZMExXSTROakk1TW1Wa01qRXhaSDAmZXhwZXJpbWVudD0obm90IHNldCkmdmFyaWF0aW9uPShub3Qgc2V0KSZ1YT1lZGdlJnZpc2l0X2lkPShub3Qgc2V0KQ..",
				AttributionSig:  "att-sig",
				Referer:         "https://www.mozilla.org/",
			},
			true,
		},
		{
			&BouncerParams{
				OS:              "win",
				Product:         "firefox-stub",
				AttributionCode: "c291cmNlPWFkZG9ucy5tb3ppbGxhLm9yZyZtZWRpdW09cmVmZXJyYWwmY2FtcGFpZ249bm9uLWZ4LWJ1dHRvbiZjb250ZW50PXJ0YTplMkk1WkdJeE5tRTBMVFpsWkdNdE5EZGxZeTFoTVdZMExXSTROakk1TW1Wa01qRXhaSDAmZXhwZXJpbWVudD0obm90IHNldCkmdmFyaWF0aW9uPShub3Qgc2V0KSZ1YT1lZGdlJnZpc2l0X2lkPShub3Qgc2V0KQ..",
				AttributionSig:  "att-sig",
				Referer:         "https://www.mozilla.org/test/other/paths",
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("OS: %s, Product: %s, Code: %s, Sig: %s, Referer: %s", test.In.OS, test.In.Product, test.In.AttributionCode, test.In.AttributionSig, test.In.Referer), func(t *testing.T) {
			assert.Equal(t, test.Out, bouncerHandler.shouldAttribute(test.In))
		})
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

func TestIsWindows7UserAgent(t *testing.T) {
	uas := []struct {
		UA	string
		IsWin7	bool
	}{
		{"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0", false},                                           // firefox XP
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0", true},                                     // firefox win7
		{"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36", false},        // Chrome XP
		{"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36", true},  // Chrome win8.1
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", false},                                 // IE XP
		{"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP)", false},                                                              // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", false},                                                   // IE Vista
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.1; en-US)", true},                                                    // IE win7
		{"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.115", true},      // Edge win8
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58", false}, // Edge win10
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/601.7.8 (KHTML, like Gecko) Version/9.1.3 Safari/537.86.7", false},      // Safari OSX 10.9.5
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36", false},   // Chrome OSX 10.9.5
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9) AppleWebKit/537.71 (KHTML, like Gecko) Version/7.0 Safari/537.71", false},             // Safari OSX 10.9
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:47.0) Gecko/20100101 Firefox/47.0", false},                          // Firefox OSX 10.11
                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8", false},     // Safari OSX 10.12
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36", false}, // Chrome OSX 10.12
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:100.0) Gecko/20100101 Firefox/100.0", false},                         // Firefox OSX 10.13
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15", false},
	}
	for _, ua := range uas {
		assert.Equal(t, ua.IsWin7, isWindows7UserAgent(ua.UA), "ua: %v", ua.UA)
	}
}

func TestIsDeprecatedOSXAgent(t *testing.T) {
	uas := []struct {
	UA           string
		isDeprecated bool
	}{
		{"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0", false},                                           // firefox XP
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0", false},                                    // firefox non-XP
		{"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36", false},        // Chrome XP
		{"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36", false}, // Chrome non-XP
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", false},                                 // IE XP
		{"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP)", false},                                                              // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", false},                                                   // IE Vista
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.1; en-US)", false},                                                   // IE non-XP
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/601.7.8 (KHTML, like Gecko) Version/9.1.3 Safari/537.86.7", false},     // Safari OSX 10.9.5
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36", false},  // Chrome OSX 10.9.5
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9) AppleWebKit/537.71 (KHTML, like Gecko) Version/7.0 Safari/537.71", false},            // Safari OSX 10.9
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.11; rv:47.0) Gecko/20100101 Firefox/47.0", false},                          // Firefox OSX 10.11
                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8", true},     // Safari OSX 10.12
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36", true}, // Chrome OSX 10.12
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:100.0) Gecko/20100101 Firefox/100.0", true},                         // Firefox OSX 10.13
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15", false},
	}
	for _, ua := range uas {
		assert.Equal(t, ua.isDeprecated, isDeprecatedOSXAgent(ua.UA), "ua: %v", ua.UA)
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

func TestOsxEsrProduct(t *testing.T) {
	assert.Equal(t, "firefox-esr-next-pkg-latest-ssl", osxEsrProduct("firefox-pkg-latest-ssl"))
	assert.Equal(t, "firefox-esr-next-latest-ssl", osxEsrProduct("firefox-latest-ssl"))
}

func TestWin7EsrProduct(t *testing.T) {
	assert.Equal(t, "firefox-esr-next-latest-ssl", win7EsrProduct("firefox-latest-ssl"))
	assert.Equal(t, "firefox-esr-next-latest-ssl", win7EsrProduct("firefox-stub"))
	assert.Equal(t, "firefox-esr-next-msi-latest-ssl", win7EsrProduct("firefox-msi-latest-ssl"))
	assert.Equal(t, "firefox-msix-latest-ssl", win7EsrProduct("firefox-msix-latest-ssl"))
	assert.Equal(t, "firefox-pkg-latest-ssl", win7EsrProduct("firefox-pkg-latest-ssl"))
}

func TestWin7EsrOS(t *testing.T) {
	uas := []struct {
		UA string
		os string
	}{
		{"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0", "win"},                                           // firefox XP
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0", "win64"},                                     // firefox win7
		{"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2224.3 Safari/537.36", "win"},        // Chrome XP
		{"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2225.0 Safari/537.36", "win64"},  // Chrome win8.1
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", "win"},                                 // IE XP
		{"Mozilla/4.0 (compatible; MSIE 6.1; Windows XP)", "win"},                                                              // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", "win"},                                                   // IE Vista
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.1; en-US)", "win"},                                                    // IE win7
		{"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.115", "win"},      // Edge win8
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.58", "win64"}, // Edge win10
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/601.7.8 (KHTML, like Gecko) Version/9.1.3 Safari/537.86.7", "win"},      // Safari OSX 10.9.5
	}
	for _, ua := range uas {
		assert.Equal(t, ua.os, win7EsrOS("firefox-stub", ua.UA), "ua: %v", ua.UA)
	}
}

func BenchmarkSha1Product(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sha1Product("firefox-43.0.0")
		sha1Product("firefox-44.0b1")
	}
}
