package main

//go:generate echo mm

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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
		cli.BoolFlag{
			Name:  "checknow",
			Usage: "Checks mirrors marked with checknow",
		},
		cli.StringFlag{
			Name:  "mirror",
			Usage: "If set, checks a specific mirror",
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

	locations, err := db.LocationsActive(c.Bool("checknow"))
	if err != nil {
		log.Fatal("db.LocationsActive:", err)
	}

	sentry := &Sentry{
		DB:        db,
		Locations: locations,
		StartTime: time.Now(),
	}

	mirrors, err := db.MirrorsActive(c.String("mirror"))
	if err != nil {
		log.Fatal("db.MirrorsActive:", err)
	}

	for _, m := range mirrors {
		err := sentry.CheckMirror(m)
		log.Printf("Error checking mirror: %s, err: %v", m.ID, err)
	}
}

type Sentry struct {
	DB        *bouncer.DB
	Locations []*bouncer.LocationsActiveResult
	StartTime time.Time
}

func (s *Sentry) CheckMirror(mirror *bouncer.MirrorsActiveResult) error {
	url, err := url.Parse(mirror.BaseURL)
	if err != nil {
		return fmt.Errorf("Could not parse url: %s, err: %v", mirror.BaseURL, err)
	}

	log.Println("Checking mirror", url.Host, "...")

	// Check DNS?

	req, err := http.NewRequest("HEAD", mirror.BaseURL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil || resp.StatusCode >= 500 {
		dberr := s.DB.MirrorSetHealth(mirror.ID, "0")
		if dberr != nil {
			log.Println(dberr)
		}
		if err != nil {
			dberr := s.DB.SentryLogInsert(s.StartTime, mirror.ID, "0", mirror.Rating, "No Response")
			if dberr != nil {
				log.Println(dberr)
			}
			return err
		}
		err = fmt.Errorf("%s returned %s", mirror.BaseURL, resp.Status)
		dberr = s.DB.SentryLogInsert(s.StartTime, mirror.ID, "0", mirror.Rating, "Bad Response")
		if dberr != nil {
			log.Println(dberr)
		}
		return err
	}

	for _, location := range s.Locations {
		err := s.CheckLocation(mirror, location)
		if err != nil {
			log.Printf("Error checking mirror: %s, location: %s, err: %v", mirror.ID, location.ID, err)
		}
	}

	return nil
}

func (s *Sentry) CheckLocation(mirror *bouncer.MirrorsActiveResult, location *bouncer.LocationsActiveResult) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 1 {
				return errors.New("Stopped after 1 redirect")
			}
			return nil
		},
	}

	path := strings.Replace(location.Path, ":lang", "en-US", -1)
	req, err := http.NewRequest("HEAD", mirror.BaseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	active, healthy := "1", "0"
	if resp.StatusCode == 200 && !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		active, healthy = "1", "1"
	} else if resp.StatusCode == 404 || resp.StatusCode == 403 {
		active, healthy = "0", "0"
	}
	err = s.DB.MirrorLocationUpdate(location.ID, mirror.ID, active, healthy)
	if err != nil {
		return err
	}

	return nil
}
