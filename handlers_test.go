package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mozilla-services/go-bouncer/bouncer"
	"github.com/stretchr/testify/assert"
)

const testDSN = "root@tcp(127.0.0.1:3306)/bouncer_test"

var bouncerHandler *BouncerHandler
var bouncerHandlerPinned *BouncerHandler

func init() {
	testDB, err := bouncer.NewDB(testDSN)
	if err != nil {
		log.Fatal(err)
	}

	bouncerHandler = &BouncerHandler{
		db:                 testDB,
		StubRootURL:        "https://stub/",
		PinHttpsHeaderName: "X-Forwarded-Proto",
		PinnedBaseURLHttp:  "download.cdn.mozilla.net/pub",
		PinnedBaseURLHttps: "download-installer.cdn.mozilla.net/pub",
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
			true,
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
		{
			&BouncerParams{
				OS:              "win",
				Product:         "firefox-stub",
				AttributionCode: "c291cmNlPWFkZG9ucy5tb3ppbGxhLm9yZyZtZWRpdW09cmVmZXJyYWwmY2FtcGFpZ249bm9uLWZ4LWJ1dHRvbiZjb250ZW50PXJ0YTplMkk1WkdJeE5tRTBMVFpsWkdNdE5EZGxZeTFoTVdZMExXSTROakk1TW1Wa01qRXhaSDAmZXhwZXJpbWVudD0obm90IHNldCkmdmFyaWF0aW9uPShub3Qgc2V0KSZ1YT1lZGdlJnZpc2l0X2lkPShub3Qgc2V0KQ..",
				AttributionSig:  "att-sig",
				// Bogus referer
				Referer: "https://www-mozilla.org/",
			},
			false,
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
			`https://stub/?attribution_code=source%3Dgoogle.com%26medium%3Dorganic%26campaign%3D%28not+set%29%26content%3D%28not+set%29&attribution_sig=anhmacsig&lang=en-US&os=osx&product=firefox`,
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
	assert.Equal(t, "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg", w.Body.String())
}

func TestBouncerHandlerValid(t *testing.T) {
	defaultUA := "Mozilla/5.0 (Windows NT 7.0; rv:10.0) Gecko/20100101 Firefox/43.0"
	testRequests := []struct {
		URL              string
		ExpectedLocation string
		UserAgent        string
	}{
		{"http://test/?product=firefox-latest&os=osx&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/mac/en-US/Firefox%2039.0.dmg", defaultUA},
		{"http://test/?product=firefox-latest&os=win64&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-SSL&os=win64&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win64/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-nightly-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-nightly-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central/firefox-128.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla 2024)"},                      // new stub
		{"http://test/?product=Firefox-nightly-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central/firefox-128.0a1.en-US.win32.installer.exe", defaultUA},
		{"http://test/?product=Firefox-nightly-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-nightly-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla 2024)"},                 // new stub
		{"http://test/?product=Firefox-nightly-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", defaultUA},
		{"http://test/?product=Firefox-nightly-latest-l10n-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-nightly-latest-l10n-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla 2024)"},                 // new stub
		{"http://test/?product=Firefox-nightly-latest-l10n-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", defaultUA},
		{"http://test/?product=Firefox-nightly-latest-l10n&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/2024/05/2024-05-06-09-48-55-mozilla-central-l10n/firefox-127.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-nightly-latest-l10n&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", "NSIS InetBgDL (Mozilla 2024)"},                 // new stub
		{"http://test/?product=Firefox-nightly-latest-l10n&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/nightly/latest-mozilla-central-l10n/firefox-128.0a1.en-US.win32.installer.exe", defaultUA},
		{"http://test/?product=Firefox-beta-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/127.0b9/win32/en-US/Firefox%20Setup%20127.0b9.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-beta-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/win32/en-US/Firefox%20Setup%2039.0.exe", "NSIS InetBgDL (Mozilla 2024)"},  // new stub
		{"http://test/?product=Firefox-beta-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/win32/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-devedition-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/devedition/releases/127.0b9/win32/en-US/Firefox%20Setup%20127.0b9.exe", "NSIS InetBgDL (Mozilla)"},      // old stub
		{"http://test/?product=Firefox-devedition-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/devedition/releases/128.0b1/win32/en-US/Firefox%20Setup%20128.0b1.exe", "NSIS InetBgDL (Mozilla 2024)"}, // new stub
		{"http://test/?product=Firefox-devedition-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/devedition/releases/128.0b1/win32/en-US/Firefox%20Setup%20128.0b1.exe", defaultUA},
		{"http://test/?product=Firefox-beta-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/127.0b9/win32/en-US/Firefox%20Setup%20127.0b9.exe", "NSIS InetBgDL (Mozilla)"}, // old stub
		{"http://test/?product=Firefox-beta-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win32/en-US/Firefox%20Setup%2039.0.exe", "NSIS InetBgDL (Mozilla 2024)"},  // new stub
		{"http://test/?product=Firefox-beta-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/39.0/win32/en-US/Firefox%20Setup%2039.0.exe", defaultUA},
		{"http://test/?product=Firefox-devedition-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/devedition/releases/127.0b9/win32/en-US/Firefox%20Setup%20127.0b9.exe", "NSIS InetBgDL (Mozilla)"},      // old stub
		{"http://test/?product=Firefox-devedition-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/devedition/releases/128.0b1/win32/en-US/Firefox%20Setup%20128.0b1.exe", "NSIS InetBgDL (Mozilla 2024)"}, // new stub
		{"http://test/?product=Firefox-devedition-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/devedition/releases/128.0b1/win32/en-US/Firefox%20Setup%20128.0b1.exe", defaultUA},
		{"http://test/?product=Firefox-latest&os=win&lang=en-US", "http://download.cdn.mozilla.net/pub/firefox/releases/127.0/win32/en-US/Firefox%20Setup%20127.0.exe", "NSIS InetBgDL (Mozilla)"},                                                  // old stub
		{"http://test/?product=Firefox-latest-ssl&os=win&lang=en-US", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/127.0/win32/en-US/Firefox%20Setup%20127.0.exe", "NSIS InetBgDL (Mozilla)"},                                   // old stub
		{"http://test/?product=partner-firefox-release-unitedinternet-foo-latest&os=win&lang=de", "http://download.cdn.mozilla.net/pub/firefox/releases/partners/foo/bar/127.0/win32/de/Firefox%20Setup%20127.0.exe", "NSIS InetBgDL (Mozilla)"},    // old stub
		{"http://test/?product=partner-firefox-release-unitedinternet-foo-latest&os=win&lang=de", "http://download.cdn.mozilla.net/pub/firefox/releases/partners/foo/bar/39.0/win32/de/Firefox%20Setup%2039.0.exe", "NSIS InetBgDL (Mozilla 2024)"}, // new stub
		{"http://test/?product=partner-firefox-release-unitedinternet-foo-latest&os=win&lang=de", "http://download.cdn.mozilla.net/pub/firefox/releases/partners/foo/bar/39.0/win32/de/Firefox%20Setup%2039.0.exe", defaultUA},
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

func TestBouncerHandlerPre2024(t *testing.T) {
	testRequests := []struct {
		URL string
	}{
		{"http://test/?product=unknown&os=win&lang=en-US"},
		{"http://test/?product=notfirefox-nightly-latest-ssl&os=win&lang=en-US"},
		{"http://test/?product=firefox-unknown&os=win&lang=en-US"},
	}

	for _, testRequest := range testRequests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", testRequest.URL, nil)
		req.Header.Set("User-Agent", "NSIS InetBgDL (Mozilla)")

		bouncerHandler.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code, "url: %v", testRequest.URL)
	}
}

func TestIsUserAgentOnlyCompatibleWithESR115(t *testing.T) {
	uas := []struct {
		UA     string
		IsWin7 bool
	}{
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", true},                                                                 // IE 64bits Win7
		{"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.7.62 Version/11.01", true},                                                                       // Opera 11 Win7
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", false},                                                         // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", false},                                                                           // IE Vista
		{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36", true}, // Edge Win7
		{"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.71 Safari/537.36", true},                    // Chrome Win8
		{"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:124.0) Gecko/20100101 Firefox/124.0", true},                                                           // Firefox Win8.1
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36", false},                     // Safari Win10
		{"Mozilla/5.0 (Windows NT 611; WOW64; Trident/7.0; rv:11.0) like Gecko", false},                                                                // Bogus
	}
	for _, ua := range uas {
		assert.Equal(t, ua.IsWin7, isUserAgentOnlyCompatibleWithESR115(ua.UA), "ua: %v", ua.UA)
	}
}

func TestIsWin64UserAgent(t *testing.T) {
	uas := []struct {
		UA      string
		IsWin64 bool
	}{
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", true},                                                                  // IE 64bits Win7
		{"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.7.62 Version/11.01", false},                                                                       // Opera 11 Win7 32bits
		{"Mozilla/5.0 (Windows; U; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727)", false},                                                          // IE XP
		{"Mozilla/5.0 (Windows; U; MSIE 7.0; Windows NT 6.0; en-US)", false},                                                                            // IE Vista
		{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.36", true},  // Edge Win7
		{"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.71 Safari/537.36", true},                     // Chrome Win8
		{"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:124.0) Gecko/20100101 Firefox/124.0", true},                                                            // Firefox Win8.1
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36", true},                       // Safari Win10
		{"Mozilla/5.0 (Windows NT 611; WOW64; Trident/7.0; rv:11.0) like Gecko", true},                                                                  // Bogus
		{"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44", true}, // Edge 100 64bits (Windows 7 SP1)
	}
	for _, ua := range uas {
		assert.Equal(t, ua.IsWin64, isWin64UserAgent(ua.UA), "ua: %v", ua.UA)
	}
}

func TestBouncerHandlerForWindowsOnlyCompatibleWithESR115(t *testing.T) {
	for _, tc := range []struct {
		userAgent string
		platform  string
	}{
		// Win 7
		{"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36", "win32"},
		{"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", "win64"},
		{"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)", "win64"},
		// Win 8
		{"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36", "win32"},
		{"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.71 Safari/537.36", "win64"},
		{"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36", "win64"},
		// Win 8.1
		{"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36", "win32"},
		{"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:124.0) Gecko/20100101 Firefox/124.0", "win64"},
		{"Mozilla/5.0 (Windows NT 6.3; Win64; x64; Trident/7.0; Touch; LCJB; rv:11.0) like Gecko", "win64"},
	} {
		// This is for stub installers.
		for _, url := range []string{
			"http://test/?product=firefox-stub&os=win&lang=en-US",
			"http://test/?product=firefox-beta-stub&os=win&lang=en-US",
			"http://test/?product=firefox-esr-stub&os=win&lang=en-US",
		} {
			// For stub installers, we need to adjust the `os` param for x64 builds.
			expectedLocation := fmt.Sprintf("https://download-installer.cdn.mozilla.net/pub/firefox/releases/115.16.1esr/%s/en-US/Firefox%%20Setup%%20115.16.1esr.exe", tc.platform)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", tc.userAgent)

			bouncerHandler.ServeHTTP(w, req)

			assert.Equal(t, 302, w.Code, "userAgent: %v, url: %v", tc.userAgent, url)
			assert.Equal(t, expectedLocation, w.HeaderMap.Get("Location"), "userAgent: %v, url: %v", tc.userAgent, url)
		}

		// This is for other firefox 32-bit products.
		for _, url := range []string{
			"http://test/?product=firefox-beta&os=win&lang=en-US",
			"http://test/?product=firefox-devedition&os=win&lang=en-US",
			"http://test/?product=firefox-nightly-latest-ssl&os=win&lang=en-US",
			"http://test/?product=firefox-ssl-latest&os=win&lang=en-US",
			"http://test/?product=firefox-unknown&os=win&lang=en-US",
		} {
			expectedLocation := "//download-installer.cdn.mozilla.net/pub/firefox/releases/115.16.1esr/win32/en-US/Firefox%20Setup%20115.16.1esr.exe"

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", tc.userAgent)

			bouncerHandler.ServeHTTP(w, req)

			assert.Equal(t, 302, w.Code, "userAgent: %v, url: %v", tc.userAgent, url)
			// We don't need to assert the scheme.
			assert.True(t, strings.HasSuffix(w.HeaderMap.Get("Location"), expectedLocation), "userAgent: %v, url: %v", tc.userAgent, url)
		}

		// This is for other firefox 64-bit products.
		for _, url := range []string{
			"http://test/?product=firefox-beta&os=win64&lang=en-US",
			"http://test/?product=firefox-devedition&os=win64&lang=en-US",
			"http://test/?product=firefox-nightly-latest-ssl&os=win64&lang=en-US",
			"http://test/?product=firefox-ssl-latest&os=win64&lang=en-US",
			"http://test/?product=firefox-unknown&os=win64&lang=en-US",
		} {
			expectedLocation := "//download-installer.cdn.mozilla.net/pub/firefox/releases/115.16.1esr/win64/en-US/Firefox%20Setup%20115.16.1esr.exe"

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", tc.userAgent)

			bouncerHandler.ServeHTTP(w, req)

			assert.Equal(t, 302, w.Code, "userAgent: %v, url: %v", tc.userAgent, url)
			// We don't need to assert the scheme.
			assert.True(t, strings.HasSuffix(w.HeaderMap.Get("Location"), expectedLocation), "userAgent: %v, url: %v", tc.userAgent, url)
		}

		// This is for MSI builds.
		expectedLocation := "https://download-installer.cdn.mozilla.net/pub/firefox/releases/131.0.3/win64/en-US/Firefox%20Setup%20131.0.3.msi"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "http://test/?product=firefox-msi-latest-ssl&os=win64&lang=en-US", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:124.0) Gecko/20100101 Firefox/124.0")

		bouncerHandler.ServeHTTP(w, req)

		assert.Equal(t, 302, w.Code)
		assert.Equal(t, expectedLocation, w.HeaderMap.Get("Location"))

		// This is for unrelated products.
		for _, url := range []string{
			"http://test/?product=unknown&os=win&lang=en-US",
			"http://test/?product=notfirefox-nightly-latest-ssl&os=win&lang=en-US",
			"http://test/?product=thunderbird-something-latest-ssl&os=win&lang=en-US",
			"http://test/?product=firefox-115.17.0esr-complete&os=win&lang=en-US",
			"http://test/?product=firefox-115.17.0esr-partial-115.16.1esr&os=win&lang=en-US",
		} {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", tc.userAgent)

			bouncerHandler.ServeHTTP(w, req)

			assert.Equal(t, 404, w.Code, "userAgent: %v, url: %v", tc.userAgent, url)
		}
	}
}

func TestBouncerHandlerForWindowsOnlyCompatibleWithESR115WithMozorgReferrer(t *testing.T) {
	expectedLocation := "http://download.cdn.mozilla.net/pub/firefox/releases/39.0/win32/en-US/Firefox%20Setup%2039.0.exe"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://test/?product=firefox-latest&os=win&lang=en-US", nil)
	req.Header.Set("Referer", "https://www.mozilla.org/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:124.0) Gecko/20100101 Firefox/124.0")

	bouncerHandler.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
	assert.Equal(t, expectedLocation, w.HeaderMap.Get("Location"))
}

func TestHealthHandler(t *testing.T) {
	testDB, err := bouncer.NewDB(testDSN)
	if err != nil {
		log.Fatal(err)
	}

	h := &HealthHandler{
		db: testDB,
	}
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	h.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, fmt.Sprintf(`{"db":true,"healthy":true,"version":"%s"}`, bouncer.Version), w.Body.String())
}
