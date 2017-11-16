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
			Name:   "db-dsn",
			Value:  "user:password@tcp(localhost:3306)/bouncer",
			Usage:  "database DSN (https://github.com/go-sql-driver/mysql#dsn-data-source-name)",
			EnvVar: "BOUNCER_DB_DSN",
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
	db, err := bouncer.NewDB(c.String("db-dsn"))
	if err != nil {
		log.Fatalf("Could not open DB: %v", err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(300 * time.Second)

	bouncerHandler := &BouncerHandler{
		db:                 db,
		CacheTime:          time.Duration(c.Int("cache-time")) * time.Second,
		PinnedBaseURLHttp:  c.String("pinned-baseurl-http"),
		PinnedBaseURLHttps: c.String("pinned-baseurl-https"),
		StubRootURL:        c.String("stub-root-url"),
	}

	healthHandler := &HealthHandler{
		db:        db,
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
