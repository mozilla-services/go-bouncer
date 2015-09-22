package sentry

import (
	"bufio"
	"net/http"
	"os"
	"testing"

	"github.com/mozilla-services/go-bouncer/bouncer"
	"github.com/stretchr/testify/assert"
)

type fileRoundTripper string

func (f fileRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	file, err := os.Open(string(f))
	if err != nil {
		return nil, err
	}
	defer file.Close()
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
