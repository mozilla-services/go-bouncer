package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/cli"
	_ "github.com/mozilla-services/go-bouncer/mozlog"
)

func main() {
	app := cli.NewApp()
	app.Name = "bouncer"
	app.Action = Main
	app.Version = Version
	app.Flags = []cli.Flag{
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
	}
	app.RunAndExitOnError()
}

func Main(c *cli.Context) {
	db, err := NewDB(c.String("db-dsn"))
	if err != nil {
		log.Fatalf("Could not open DB: %v", err)
	}
	defer db.Close()

	bouncerHandler := &BouncerHandler{
		db: db,
	}

	server := &http.Server{
		Addr:    c.String("addr"),
		Handler: bouncerHandler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
