package main

import (
	"net/http"
	"net/url"
	"strings"
)

// BouncerParams holds/parses params for incoming bouncer requests
type BouncerParams struct {
	PrintOnly       bool
	OS              string
	Product         string
	Lang            string
	AttributionCode string
	AttributionSig  string
	Referer         string
}

// BouncerParamsFromValues constructs parameter list from incoming request Values
func BouncerParamsFromValues(vals url.Values, headers http.Header) *BouncerParams {
	return &BouncerParams{
		PrintOnly:       vals.Get("print") == "yes",
		OS:              strings.TrimSpace(strings.ToLower(vals.Get("os"))),
		Product:         strings.TrimSpace(strings.ToLower(vals.Get("product"))),
		Lang:            vals.Get("lang"),
		AttributionCode: vals.Get("attribution_code"),
		AttributionSig:  vals.Get("attribution_sig"),
		Referer:         headers.Get("Referer"),
	}
}
