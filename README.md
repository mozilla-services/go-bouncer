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

### Setting up `bouncer-admin` in localdev

[bouncer-admin][] is the admin interface for go-bouncer. It can be optionally
set up by first cloning the repository (once):

```
git clone https://github.com/mozilla-services/bouncer-admin
```

Then, run `docker compose` as follows:

```
docker compose -f compose.yaml -f compose-admin.yaml up -d
```

Note: Every `docker compose` needs to specify both configuration files if the
intent is to interact with the `admin` container, e.g. `docker compose -f
compose.yaml -f compose-admin.yaml logs -f admin`.

The API is available at: http://127.0.0.1:9000/api/. The authenticated user is
`admin` with the traditional `admin` password. Here is an example to create a
new product:

```
curl -X POST 'http://admin:admin@127.0.0.1:9000/api/product_add/' -d 'product=A-Test-Product&ssl_only=true'
<?xml version="1.0" encoding="utf-8"?>
<products>
  <product id="4557" name="A-Test-Product"/>
</products>
```

## Environment variables

### `BOUNCER_ADDR`

Address on which to listen. The default value is: `:8888`

### `BOUNCER_DB_DSN`

Database DSN, see: https://github.com/go-sql-driver/mysql#dsn-data-source-name
for more details about the format.

### `BOUNCER_PIN_HTTPS_HEADER_NAME`

When this flag is set and the request header value equals https, an HTTPS
redirect will always be returned. The default value is: `X-Forwarded-Proto`,
which means this feature is enabled by default.

### `BOUNCER_PINNED_BASEURL_HTTP`

Configure the base URL for non-SSL only products.

### `BOUNCER_PINNED_BASEURL_HTTPS`

Same as `BOUNCER_PINNED_BASEURL_HTTP` but SSL-only products.

### `BOUNCER_STUB_ROOT_URL`

Optional. If set, bouncer will redirect requests with `attribution_sig` and
`attribution_code` parameters to the stubattribution service using this URL:

```
BOUNCER_STUB_ROOT_URL?product=PRODUCT&os=OS&lang=LANG&attribution_sig=ATTRIBUTION_SIG&attribution_code=ATTRIBUTION_CODE
```

[go-bouncer]: https://github.com/mozilla-services/go-bouncer/
[bouncer-admin]: https://github.com/mozilla-services/bouncer-admin/
