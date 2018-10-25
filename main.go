package main

//go:generate ./version.sh

import (
	"log"
	"net/http"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/go-bouncer/bouncer"
	_ "github.com/mozilla-services/go-bouncer/mozlog"
)

func main() {
	app := cli.NewApp()
	app.Name = "bouncer"
	app.Action = Main
	app.Version = bouncer.Version
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "cache-time",
			Value: 60,
			Usage: "Time, in seconds, for Cache-Control max-age",
		},
		cli.StringFlag{
			Name:   "addr",
			Value:  ":8888",
			Usage:  "address on which to listen",
			EnvVar: "BOUNCER_ADDR",
		},
		cli.StringFlag{
			Name:   "locations-file",
			Usage:  "Location file.",
			EnvVar: "BOUNCER_LOCATION_FILE",
		},
		cli.StringFlag{
			Name:   "pin-https-header-name",
			Value:  "X-Forwarded-Proto",
			Usage:  "If this flag is set and the request header value equals https, an https redirect will always be returned",
			EnvVar: "BOUNCER_PIN_HTTPS_HEADER_NAME",
		},
		cli.StringFlag{
			Name:   "pinned-baseurl-http",
			Usage:  "if this flag is set it will always be the base url for http products. Scheme should be excluded, e.g.,: pinned-cdn.mozilla.com/pub",
			EnvVar: "BOUNCER_PINNED_BASEURL_HTTP",
		},
		cli.StringFlag{
			Name:   "pinned-baseurl-https",
			Usage:  "if this flag is set it will always be the base url for https products. Scheme should be excluded, e.g.,: pinned-cdn.mozilla.com/pub",
			EnvVar: "BOUNCER_PINNED_BASEURL_HTTPS",
		},
		cli.StringFlag{
			Name:   "stub-root-url",
			Value:  "",
			Usage:  "Root url of service used to service modified stub installers e.g., https://stubdownloader.services.mozilla.com/",
			EnvVar: "BOUNCER_STUB_ROOT_URL",
		},
	}
	app.RunAndExitOnError()
}

func Main(c *cli.Context) {
	locations, err := ParseLocationsFile(c.String("locations-file"))
	if err != nil {
		log.Fatal("Error parsing locations-file: ", err)
	}
	bouncerHandler := &BouncerHandler{
		Locations:          locations,
		CacheTime:          time.Duration(c.Int("cache-time")) * time.Second,
		PinHttpsHeaderName: c.String("pin-https-header-name"),
		PinnedBaseURLHttp:  c.String("pinned-baseurl-http"),
		PinnedBaseURLHttps: c.String("pinned-baseurl-https"),
		StubRootURL:        c.String("stub-root-url"),
	}

	healthHandler := &HealthHandler{
		CacheTime: 5 * time.Second,
	}

	mux := http.NewServeMux()

	mux.Handle("/__heartbeat__", healthHandler)
	mux.Handle("/", bouncerHandler)

	server := &http.Server{
		Addr:    c.String("addr"),
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
