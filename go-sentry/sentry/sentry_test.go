package sentry

import (
	"bufio"
	"net/http"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/mozilla-services/go-bouncer/bouncer"
	"github.com/stretchr/testify/assert"
)

type roundTripRecorder struct {
	FileName    string
	LastRequest *http.Request
}

func fileRoundTripper(f string) *roundTripRecorder {
	return &roundTripRecorder{FileName: f}
}

func (r *roundTripRecorder) RoundTrip(req *http.Request) (*http.Response, error) {
	file, err := os.Open(r.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r.LastRequest = req
	return http.ReadResponse(bufio.NewReader(file), req)
}

func TestHeadMirror(t *testing.T) {
	s := &Sentry{}

	mirror := &bouncer.MirrorsActiveResult{}

	s.roundTripper = fileRoundTripper("./fixtures/200.txt")
	err := s.HeadMirror(mirror)
	assert.NoError(t, err)

	s.roundTripper = fileRoundTripper("./fixtures/404.txt")
	err = s.HeadMirror(mirror)
	assert.NoError(t, err)

	s.roundTripper = fileRoundTripper("./fixtures/500.txt")
	err = s.HeadMirror(mirror)
	assert.EqualError(t, err, "Bad Response: 500 Internal Service Error")
}

func TestCheckLocation(t *testing.T) {
	s := &Sentry{client: DefaultClient}

	mirror := &bouncer.MirrorsActiveResult{
		BaseURL: "http://download-installer.cdn.mozilla.net/pub",
	}

	locations := []struct {
		Path string
		Lang string
	}{
		{"/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe", "en-US"},
		{"/thunderbird/releases/3.1a1/win32/:lang/Firefox%20Setup%2039.0.exe", "en-US"},
		{"/thunderbird/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe", "en-US"},
		{"/seamonkey/releases/39.0/win32/:lang/Firefox%20Setup%2039.0.exe", "en-US"},
		{"/firefox/releases/39.0/win32/:lang/Firefox%20Setup%2039.0-euBallot.exe", "en-US"},
	}

	for _, loc := range locations {
		location := &bouncer.LocationsActiveResult{
			Path: loc.Path,
		}

		runLog := logrus.WithField("testing", true)

		rt := fileRoundTripper("./fixtures/200.binary.txt")
		s.client.Transport = rt
		res := s.checkLocation(mirror, location, runLog)
		assert.True(t, res.Active)
		assert.True(t, res.Healthy)
		assert.Contains(t, rt.LastRequest.URL.String(), "/"+loc.Lang+"/")

		rt = fileRoundTripper("./fixtures/404.txt")
		s.client.Transport = rt
		res = s.checkLocation(mirror, location, runLog)
		assert.False(t, res.Active)
		assert.False(t, res.Healthy)

		rt = fileRoundTripper("./fixtures/200.txt")
		s.client.Transport = rt
		res = s.checkLocation(mirror, location, runLog)
		assert.True(t, res.Active)
		assert.False(t, res.Healthy)

		rt = fileRoundTripper("./fixtures/500.txt")
		s.client.Transport = rt
		res = s.checkLocation(mirror, location, runLog)
		assert.True(t, res.Active)
		assert.False(t, res.Healthy)

		// Infinite redirect
		rt = fileRoundTripper("./fixtures/302.txt")
		s.client.Transport = rt
		res = s.checkLocation(mirror, location, runLog)
		assert.True(t, res.Active)
		assert.False(t, res.Healthy)
	}
}
