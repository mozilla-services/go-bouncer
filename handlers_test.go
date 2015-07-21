package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBouncerHandlerParams(t *testing.T) {
	bouncer := &BouncerHandler{}
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://test/?os=mac&lang=en-US", nil)
	assert.NoError(t, err)

	bouncer.ServeHTTP(w, req)
	assert.Equal(t, 302, w.Code)
	assert.Equal(t, "http://www.mozilla.org/", w.HeaderMap.Get("Location"))
}
