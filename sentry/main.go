package main

//go:generate echo mm

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/mozilla-services/go-bouncer/bouncer"
)

func main() {
	app := cli.NewApp()
	app.Name = "sentry"
	app.Action = Main
	app.Version = bouncer.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-dsn",
			Value:  "user:password@tcp(localhost:3306)/bouncer",
			Usage:  "database DSN (https://github.com/go-sql-driver/mysql#dsn-data-source-name)",
			EnvVar: "SENTRY_DB_DSN",
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
}