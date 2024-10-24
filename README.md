# go-bouncer

[![CircleCI](https://circleci.com/gh/mozilla-services/go-bouncer/tree/master.svg?style=svg)](https://circleci.com/gh/mozilla-services/go-bouncer/?branch=master)

The project behind https://download.mozilla.org/ :fire:

## Getting started (for development)

```
docker compose up -d
```

You can then call the service using `127.0.0.1:8000` directly or via the Nginx
proxy at `127.0.0.1:18000`:

```
$ curl -H "User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)" -I 'http://127.0.0.1:8000/?product=firefox-ssl&os=win'
HTTP/1.1 302 Found
Cache-Control: max-age=60
Content-Type: text/html; charset=utf-8
Location:
https://download-installer.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe
Date: Wed, 23 Oct 2024 08:24:42 GMT

$ curl -H "User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)" -I 'http://127.0.0.1:18000/?product=firefox-ssl&os=win'
HTTP/1.1 302 Found
Server: nginx/1.25.5
Date: Wed, 23 Oct 2024 08:24:43 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 134
Connection: keep-alive
Cache-Control: max-age=60
Location: https://download-installer.cdn.mozilla.net/pub/firefox/releases/43.0.1/win32/en-US/Firefox%20Setup%2043.0.1.exe
x-debug-cache-key: upstream_bouncer/?product=firefox-ssl&os=winwinxpother
```

This nginx config looks similar to the one we have on production but it isn't
exactly the same. In addition to that, it adds some debugging capabilities like
the following headers:

- `x-debug-cache-key`: the computed cache key
- `x-debug-referer`: the referer value, if any

## Environment Variables

Note: This section is not exhaustive.

### `BOUNCER_PINNED_BASEURL_HTTP`

If this is a unset, bouncer will randomly pick a healthy mirror from the
database and return its base url. If this option is set, the mirror table is
completely ignored and `BOUNCER_PINNED_BASEURL_HTTP` will be returned instead.

This option acts on non ssl only products.

Example: `BOUNCER_PINNED_BASEURL=download-sha1.cdn.mozilla.net/pub`

### `BOUNCER_PINNED_BASEURL_HTTPS`

This option is exactly the same as `BOUNCER_PINNED_BASEURL_HTTP` but acts on ssl
only products.

### `BOUNCER_STUB_ROOT_URL`

If set, bouncer will redirect requests with `attribution_sig` and
`attribution_code` parameters to
`BOUNCER_STUB_ROOT_URL?product=PRODUCT&os=OS&lang=LANG&attribution_sig=ATTRIBUTION_SIG&attribution_code=ATTRIBUTION_CODE`.

Example: `BOUNCER_STUB_ROOT_URL=https://stubdownloader.services.mozilla.com/`

[go-bouncer]: https://github.com/mozilla-services/go-bouncer/
