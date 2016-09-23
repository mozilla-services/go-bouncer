# go-bouncer [![Build Status](https://travis-ci.org/mozilla-services/go-bouncer.svg?branch=master)](https://travis-ci.org/mozilla-services/go-bouncer) [![GoDoc](https://godoc.org/github.com/mozilla-services/go-bouncer?status.svg)](https://godoc.org/github.com/mozilla-services/go-bouncer)

A Go port of the [user facing portion](https://github.com/mozilla/tuxedo/tree/master/bouncer) as part of the [Bouncer project](https://wiki.mozilla.org/Bouncer).

## Environment Variables
### `BOUNCER_PINNED_BASEURL_HTTP`
If this is a unset, bouncer will randomly pick a healthy mirror from the database and return its base url. If this option is set, the mirror table is completely ignored and `BOUNCER_PINNED_BASEURL_HTTP` will be returned instead.

This option acts on non ssl only products.

Example: `BOUNCER_PINNED_BASEURL=download-sha1.cdn.mozilla.net/pub`

### `BOUNCER_PINNED_BASEURL_HTTPS`
This option is exactly the same as `BOUNCER_PINNED_BASEURL_HTTP` but acts on ssl only products.
