package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/cli"
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
	}
	app.RunAndExitOnError()
}

func Main(c *cli.Context) {

	server := &http.Server{
		Addr: c.String("addr"),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
