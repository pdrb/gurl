# gurl

[![Go Report Card](https://goreportcard.com/badge/github.com/pdrb/gurl)](https://goreportcard.com/report/github.com/pdrb/gurl)
[![CI](https://github.com/pdrb/gurl/actions/workflows/ci.yml/badge.svg)](https://github.com/pdrb/gurl/actions/workflows/ci.yml)
[![Release](https://github.com/pdrb/gurl/actions/workflows/release.yml/badge.svg)](https://github.com/pdrb/gurl/actions/workflows/release.yml)
[![Releases](https://img.shields.io/github/v/release/pdrb/gurl.svg)](https://github.com/pdrb/gurl/releases)
[![LICENSE](https://img.shields.io/github/license/pdrb/gurl)](https://github.com/pdrb/gurl/blob/main/LICENSE)

A simple http client cli written in Go.

Thanks to [req](https://github.com/imroc/req) (and [utls](https://github.com/refraction-networking/utls)) there are some cool features like:

* TLS Fingerprinting
* HTTP Fingerprinting
* Basic/Bearer Authentication
* Custom Certificates
* Proxy Support
* Retries
* Tracing

Check usage below for full list of supported options.

Whenever possible, the options are similar to `curl` with some sensible defaults applied like `application/json` content type for post, put and patch methods.

Also, if the response is valid json, it will be automatically prettified (this can be disabled using `--raw-response`).

## Install

Install compiling from source using Go:

```shell
go install github.com/pdrb/gurl@latest
```

Or download the appropriate pre-built binary from [Releases](https://github.com/pdrb/gurl/releases).

## Usage

```text
Usage: gurl <url>

A simple http client cli written in Go.

Arguments:
  <url>    Url to access.

Flags:
  -h, --help                    Show context-sensitive help.
  -u, --auth=auth               Basic HTTP authentication in the format
                                username:password.
  -b, --bearer-token=token      Set bearer auth token.
      --ca-cert=file            CA certificate file.
      --client-cert=cert-file,key-file,...
                                Client cert and key files separated by comma:
                                "cert.pem,key.pem".
  -c, --content-type=content    Content-Type http header, default is
                                application/json for POST, PUT and PATCH
                                methods.
  -d, --data=payload            Data payload (request body).
  -f, --data-file=file          Read data payload from file.
      --disable-redirect        Disable redirects.
      --force-http-1            Force HTTP/1.1 to be used.
  -H, --headers=h1=v1;h2=v2     HTTP headers in the format:
                                "header1=value1;header2=value2".
      --impersonate="none"      Fully impersonate chrome, firefox or safari
                                browser (this will automatically set headers,
                                headers order and tls fingerprint).
  -k, --insecure                Allow insecure SSL connections.
  -X, --method="GET"            Http method: GET, HEAD, POST, PUT, PATCH,
                                DELETE or OPTIONS.
  -o, --output-file=file        Save response to file.
      --proxy=proxy             Proxy to use, e.g.:
                                "http://user:pass@myproxy:8080".
      --raw-response            Print raw response string (disable json
                                prettify).
  -r, --retries=0               Number of retries in case of errors and http
                                status code >= 500.
  -t, --timeout=10000           Timeout in milliseconds.
      --tls-finger="go"         TLS Fingerprint: chrome, firefox, edge, safari,
                                ios, android, random or go.
      --trace                   Show tracing/performance information.
  -A, --user-agent=agent        Set User-Agent http header.
  -v, --verbose                 Enable verbose/debug mode.
  -V, --version                 Show version and exit.
```

## Example

```text
$ gurl -X POST 'https://httpbin.org/anything?var1=value1&var2=value2' \
    -u user:pass \
    -H 'header1=value1;header2=value2' \
    -d '{"user": "name"}' \
    -v --trace

2023/09/03 00:21:10.926979 DEBUG [req] HTTP/2 POST https://httpbin.org/anything?var1=value1&var2=value2
:authority: httpbin.org
:method: POST
:path: /anything?var1=value1&var2=value2
:scheme: https
content-type: application/json; charset=utf-8
authorization: Basic dXNlcjpwYXNz
header1: value1
header2: value2
user-agent: gurl 1.5.0
content-length: 16
accept-encoding: gzip

:status: 200
date: Sun, 03 Sep 2023 03:21:12 GMT
content-type: application/json
content-length: 644
server: gunicorn/19.9.0
access-control-allow-origin: *
access-control-allow-credentials: true

{
  "args": {
    "var1": "value1",
    "var2": "value2"
  },
  "data": "{\"user\": \"name\"}",
  "files": {},
  "form": {},
  "headers": {
    "Accept-Encoding": "gzip",
    "Authorization": "Basic dXNlcjpwYXNz",
    "Content-Length": "16",
    "Content-Type": "application/json; charset=utf-8",
    "Header1": "value1",
    "Header2": "value2",
    "Host": "httpbin.org",
    "User-Agent": "gurl 1.5.0",
    "X-Amzn-Trace-Id": "Root=1-64f3fba8-52bd3b8337869ce91c0b44d1"
  },
  "json": {
    "user": "name"
  },
  "method": "POST",
  "origin": "187.0.34.127",
  "url": "https://httpbin.org/anything?var1=value1&var2=value2"
}

------- TRACE INFO -------
TotalTime         : 1.0029421s
DNSLookupTime     : 6.4776ms
TCPConnectTime    : 198.6111ms
TLSHandshakeTime  : 423.5118ms
FirstResponseTime : 368.5909ms
ResponseTime      : 549.8µs
IsConnReused:     : false
RemoteAddr        : 54.175.87.239:443

the request total time is 1.0029421s, and costs 423.5118ms on tls handshake
```
