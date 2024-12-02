// go-bouncer is the web service behind https://download.mozilla.org/.
// Its purpose is to serve builds given a product, OS, and language.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	_ "github.com/mozilla-services/go-bouncer/mozlog"
)

const (
	// versionFilePath is the path to the `version.json` file in the Docker container.
	versionFilePath = "/app/version.json"
)

func main() {
	app := cli.NewApp()
	app.Name = "bouncer"
	app.Action = Main
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "cache-time",
			Value: 60,
			Usage: "Time, in seconds, for Cache-Control max-age",
		},
		cli.StringFlag{
			Name:   "addr",
			Value:  ":8888",
			Usage:  "Address on which to listen",
			EnvVar: "BOUNCER_ADDR",
		},
		cli.StringFlag{
			Name:   "db-dsn",
			Value:  "user:password@tcp(localhost:3306)/bouncer",
			Usage:  "Database DSN (https://github.com/go-sql-driver/mysql#dsn-data-source-name)",
			EnvVar: "BOUNCER_DB_DSN",
		},
		cli.StringFlag{
			Name:   "pin-https-header-name",
			Value:  "X-Forwarded-Proto",
			Usage:  "If this flag is set and the request header value equals https, an HTTPS redirect will always be returned",
			EnvVar: "BOUNCER_PIN_HTTPS_HEADER_NAME",
		},
		cli.StringFlag{
			Name:   "pinned-baseurl-http",
			Usage:  "The base URL for HTTP products. Scheme should be excluded, e.g. pinned-cdn.mozilla.com/pub",
			EnvVar: "BOUNCER_PINNED_BASEURL_HTTP",
		},
		cli.StringFlag{
			Name:   "pinned-baseurl-https",
			Usage:  "The base URL for HTTPS products. Scheme should be excluded, e.g. pinned-cdn.mozilla.com/pub",
			EnvVar: "BOUNCER_PINNED_BASEURL_HTTPS",
		},
		cli.StringFlag{
			Name:   "stub-root-url",
			Value:  "",
			Usage:  "Optional. Root URL of the stubattribution service, e.g. https://stubdownloader.services.mozilla.com/",
			EnvVar: "BOUNCER_STUB_ROOT_URL",
		},
	}
	app.RunAndExitOnError()
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	versionFile, err := os.ReadFile(versionFilePath)
	if err != nil {
		http.Error(w, "Could not read version file.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(versionFile)
}

// Main is the entrypoint of the application.
func Main(c *cli.Context) {
	db, err := NewDB(c.String("db-dsn"))
	if err != nil {
		log.Fatalf("Could not open DB: %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(300 * time.Second)

	if c.String("pinned-baseurl-http") == "" {
		log.Fatal("BOUNCER_PINNED_BASEURL_HTTP must be set")
	}

	if c.String("pinned-baseurl-https") == "" {
		log.Fatal("BOUNCER_PINNED_BASEURL_HTTPS must be set")
	}

	bouncerHandler := &BouncerHandler{
		db:                 db,
		CacheTime:          time.Duration(c.Int("cache-time")) * time.Second,
		PinHTTPSHeaderName: c.String("pin-https-header-name"),
		PinnedBaseURLHttp:  c.String("pinned-baseurl-http"),
		PinnedBaseURLHttps: c.String("pinned-baseurl-https"),
		StubRootURL:        c.String("stub-root-url"),
	}

	healthHandler := &HealthHandler{
		db:        db,
		CacheTime: 5 * time.Second,
	}

	mux := http.NewServeMux()

	mux.Handle("/__lbheartbeat__", healthHandler)
	mux.Handle("/__heartbeat__", healthHandler)
	mux.HandleFunc("/__version__", versionHandler)
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
