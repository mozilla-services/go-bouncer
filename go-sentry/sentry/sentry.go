package sentry

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mozilla-services/go-bouncer/bouncer"
)

// The default http.Client for HeadLocation
var DefaultClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 1 {
			return errors.New("Stopped after 1 redirect")
		}
		return nil
	},
}

// Sentry contains sentry operations
type Sentry struct {
	DB      *bouncer.DB
	Verbose bool

	locations   []*bouncer.LocationsActiveResult
	mirrors     []*bouncer.MirrorsActiveResult
	startTime   time.Time
	runLck      sync.Mutex
	locationSem chan bool
	mirrorSem   chan bool

	client       *http.Client
	roundTripper http.RoundTripper
}

// New returns a new Sentry
func New(db *bouncer.DB, checknow bool, mirror string, mirrorRoutines, locRoutines int) (*Sentry, error) {
	locations, err := db.LocationsActive(checknow)
	if err != nil {
		return nil, fmt.Errorf("db.LocationsActive: %v", err)
	}

	mirrors, err := db.MirrorsActive(mirror)
	if err != nil {
		return nil, fmt.Errorf("db.MirrorsActive: %v", err)
	}

	return &Sentry{
		DB:           db,
		locations:    locations,
		mirrors:      mirrors,
		locationSem:  make(chan bool, locRoutines),
		mirrorSem:    make(chan bool, mirrorRoutines),
		client:       DefaultClient,
		roundTripper: http.DefaultTransport,
	}, nil
}

// Run starts a full sentry run
func (s *Sentry) Run() error {
	s.runLck.Lock()
	defer s.runLck.Unlock()

	wg := sync.WaitGroup{}

	s.startTime = time.Now()
	for _, mirror := range s.mirrors {
		s.mirrorSem <- true
		wg.Add(1)
		go func(mirror *bouncer.MirrorsActiveResult) {
			defer func() {
				<-s.mirrorSem
				wg.Done()
			}()
			if err := s.checkMirror(mirror); err != nil {
				log.Printf("Error checking mirror: %s err: %s", mirror.BaseURL, err)
			}
		}(mirror)
	}

	wg.Wait()
	return nil
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

type checkLocationResult struct {
	Active  bool
	Healthy bool
}

func (s *Sentry) checkLocation(mirror *bouncer.MirrorsActiveResult, location *bouncer.LocationsActiveResult, runLog *lockedWriter) *checkLocationResult {
	lang := "en-US"

	if strings.Contains(location.Path, "/firefox/") &&
		!strings.Contains(location.Path, "/namoroka/") &&
		!strings.Contains(location.Path, "/devpreview/") &&
		!strings.Contains(location.Path, "3.6b1") &&
		!strings.Contains(location.Path, "wince-arm") &&
		!strings.Contains(strings.ToLower(location.Path), "euballot") {

		lang = "zh-TW"
	} else if strings.Contains(location.Path, "/thunderbird/") {
		if strings.Contains(location.Path, "3.1a1") {
			lang = "tr"
		} else {
			lang = "zh-TW"
		}
	} else if strings.Contains(location.Path, "/seamonkey/") {
		if strings.Contains(location.Path, "2.0.5") || strings.Contains(location.Path, "2.0.6") {
			lang = "zh-CN"
		} else {
			lang = "tr"
		}
	} else if strings.Contains(strings.ToLower(location.Path), "-euballot") {
		lang = "sv-SE"
	}

	path := strings.Replace(location.Path, ":lang", lang, -1)
	url := mirror.BaseURL + path

	start := time.Now()
	active, healthy := true, false

	resp, err := s.HeadLocation(url)
	elapsed := time.Now().Sub(start)
	if err != nil {
		active, healthy = true, false
		runLog.Printf("%s TOOK=%v ERR=%v\n", url, elapsed, err)
	} else {
		if resp.StatusCode == 200 && !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			active, healthy = true, true
		} else if resp.StatusCode == 404 || resp.StatusCode == 403 {
			active, healthy = false, false
		}

		runLog.Printf("%s TOOK=%v RC=%d\n", url, elapsed, resp.StatusCode)
	}
	return &checkLocationResult{Active: active, Healthy: healthy}
}

func (s *Sentry) checkMirror(mirror *bouncer.MirrorsActiveResult) error {
	runLog := newLockedWriter()
	runLog.Printf("Checking mirror %s ...\n", mirror.BaseURL)

	// Check overall mirror health
	err := s.HeadMirror(mirror)
	if err != nil {
		if dberr := s.DB.MirrorSetHealth(mirror.ID, "0"); dberr != nil {
			return fmt.Errorf("MirrorSetHealth: %v", dberr)
		}
		if dberr := s.DB.SentryLogInsert(s.startTime, mirror.ID, "0", mirror.Rating, err.Error()); dberr != nil {
			return fmt.Errorf("SentryLogInsert: %v", dberr)
		}
		return fmt.Errorf("HeadMirror: %v", err)
	}

	// Check locations
	wg := sync.WaitGroup{}
	for _, location := range s.locations {
		s.locationSem <- true
		wg.Add(1)
		go func(location *bouncer.LocationsActiveResult) {
			defer func() {
				<-s.locationSem
				wg.Done()
			}()

			res := s.checkLocation(mirror, location, runLog)
			if err := s.DB.MirrorLocationUpdate(location.ID, mirror.ID, boolToString(res.Active), boolToString(res.Healthy)); err != nil {
				runLog.Printf("MirrorLocationUpdate err: %v", err)
				return
			}

		}(location)
	}

	wg.Wait()

	if err := s.DB.SentryLogInsert(s.startTime, mirror.ID, "1", mirror.Rating, runLog.String()); err != nil {
		log.Println(err)
	}

	if s.Verbose {
		log.Println(runLog.String())
	}

	return nil
}

// HeadMirror returns error if mirror is not healthy
func (s *Sentry) HeadMirror(mirror *bouncer.MirrorsActiveResult) error {
	// Check DNS?

	req, err := http.NewRequest("HEAD", mirror.BaseURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.roundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 500 {
		return fmt.Errorf("Bad Response: %s", resp.Status)
	}
	return nil

}

// HeadLocation makes a HEAD request to url and returns the response
func (s *Sentry) HeadLocation(url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req)
}
