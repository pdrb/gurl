# gurl
A simple http client cli written in Go.

## Install

Install compiling from source using Go:

```shell
go install github.com/pdrb/gurl@latest
```

Or download the appropriate pre-built binary from [Releases](https://github.com/pdrb/gurl/releases).

## Usage

```text
Usage: gurl <command>

A simple http client cli written in Go.

Flags:
  -h, --help                   Show context-sensitive help.
  -u, --auth=auth              Basic HTTP authentication in the format
                               username:password.
  -b, --bearer-token=token     Set bearer auth token.
      --ca-cert=file           CA certificate file.
      --client-cert=cert-file,key-file,...
                               Client cert and key files separated by comma:
                               "cert.pem,key.pem".
      --disable-redirect       Disable redirects.
      --force-http-1           Force HTTP/1.1 to be used.
  -H, --headers=h1=v1;h2=v2    HTTP headers in the format:
                               "header1=value1;header2=value2".
      --impersonate="none"     Fully impersonate chrome, firefox or safari
                               browser (this will automatically set headers,
                               headers order and tls fingerprint).
  -k, --insecure               Allow insecure SSL connections.
  -o, --output-file=file       Save response to file.
      --proxy=proxy            Proxy to use, e.g.:
                               "http://user:pass@myproxy:8080".
      --raw-response           Print raw response string (disable json
                               prettify).
  -r, --retries=0              Number of retries in case of errors and http
                               status code >= 500.
  -t, --timeout=10000          Timeout in milliseconds.
      --tls-finger="go"        TLS Fingerprint: chrome, firefox, edge, safari,
                               ios, android, random or go.
      --trace                  Show tracing/performance information.
  -A, --user-agent=agent       Set User-Agent http header.
  -v, --verbose                Enable verbose/debug mode.

Commands:
  get        GET HTTP method.
  head       HEAD HTTP method.
  post       POST HTTP method.
  put        PUT HTTP method.
  patch      PATCH HTTP method.
  delete     DELETE HTTP method.
  options    OPTIONS HTTP method.
  version    Show version and exit.

Run "gurl <command> --help" for more information on a command.
```

Each command has it's own help, for example:

```text
$ gurl post --help

Usage: gurl post <url>

POST HTTP method.

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
      --disable-redirect        Disable redirects.
      --force-http-1            Force HTTP/1.1 to be used.
  -H, --headers=h1=v1;h2=v2     HTTP headers in the format:
                                "header1=value1;header2=value2".
      --impersonate="none"      Fully impersonate chrome, firefox or safari
                                browser (this will automatically set headers,
                                headers order and tls fingerprint).
  -k, --insecure                Allow insecure SSL connections.
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

  -d, --data=payload            Data payload (request body).
  -f, --data-file=file          Read data payload from file.
  -c, --content-type=content    Content-Type http header, default is
                                application/json.
```

## Example

```text
$ gurl post 'https://httpbin.org/anything?var1=value1&var2=value2' \
    -u user:pass -H 'header1=value1;header2=value2' \
    -d '{"user": "name"}' \
    -v --trace

2023/08/10 01:10:18.662247 DEBUG [req] HTTP/2 POST https://httpbin.org/anything?var1=value1&var2=value2
:authority: httpbin.org
:method: POST
:path: /anything?var1=value1&var2=value2
:scheme: https
header2: value2
user-agent: gurl 1.2.0
content-type: application/json; charset=utf-8
authorization: Basic dXNlcjpwYXNz
header1: value1
content-length: 16
accept-encoding: gzip

:status: 200
date: Thu, 10 Aug 2023 04:10:20 GMT
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
    "User-Agent": "gurl 1.2.0",
    "X-Amzn-Trace-Id": "Root=1-64d4632b-6cbed42b1c30160e29d66ed0"
  },
  "json": {
    "user": "name"
  },
  "method": "POST",
  "origin": "187.0.35.180",
  "url": "https://httpbin.org/anything?var1=value1&var2=value2"
}

------- TRACE INFO -------
TotalTime         : 1.4490617s
DNSLookupTime     : 3.7486ms
TCPConnectTime    : 149.5495ms
TLSHandshakeTime  : 304.5126ms
FirstResponseTime : 985.7845ms
ResponseTime      : 544.1µs
IsConnReused:     : false
RemoteAddr        : 54.236.190.246:443

the request total time is 1.4490617s, and costs 985.7845ms from connection ready to server respond first byte
```
