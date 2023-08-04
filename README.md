# gurl
A simple http client cli written in Go.

## Install

```shell
go install github.com/pdrb/gurl@latest
```

## Usage

```text
Usage: gurl <command>

A simple http client cli written in Go.

Flags:
  -h, --help                     Show context-sensitive help.
  -a, --auth=STRING              Basic HTTP authentication in the format
                                 username:password.
  -b, --bearer-token=STRING      Set bearer auth token.
      --disable-redirect         Disable redirects (default: disabled).
  -H, --headers=KEY=VALUE;...    HTTP headers in the format:
                                 "header1=value1;header2=value2".
  -i, --insecure                 Allow insecure SSL connections (default:
                                 disabled).
  -t, --timeout=10000            Timeout in milliseconds.
      --tls-finger="go"          TLS Fingerprint: chrome, firefox, edge, safari,
                                 ios, android, random or go.
      --trace                    Show tracing/performance information (default:
                                 disabled).
  -u, --user-agent=STRING        Set User-Agent http header.
  -v, --verbose                  Enable verbose/debug mode (default: disabled).

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

## Example

```text
$ gurl post 'https://httpbin.org/anything?var1=value1&var2=value2' \
    -a user:pass -H 'header1=value1;header2=value2' \
    -d '{"user": "name"}' \
    -v --trace

2023/08/03 23:51:45.157111 DEBUG [req] HTTP/2 POST https://httpbin.org/anything?var1=value1&var2=value2
:authority: httpbin.org
:method: POST
:path: /anything?var1=value1&var2=value2
:scheme: https
content-type: application/json; charset=utf-8
authorization: Basic dXNlcjpwYXNz
header1: value1
header2: value2
user-agent: gurl 1.0.0
content-length: 16
accept-encoding: gzip

:status: 200
date: Fri, 04 Aug 2023 02:51:49 GMT
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
        "User-Agent": "gurl 1.0.0",
        "X-Amzn-Trace-Id": "Root=1-64cc67c2-6c372ba45440ccac33555a3d"
    },
    "json": {
        "user": "name"
    },
    "method": "POST",
    "origin": "187.0.35.180",
    "url": "https://httpbin.org/anything?var1=value1\u0026var2=value2"
}

------- TRACE INFO -------
TotalTime         : 4.0295301s
DNSLookupTime     : 3.661ms
TCPConnectTime    : 145.1952ms
TLSHandshakeTime  : 301.5689ms
FirstResponseTime : 3.5787528s
ResponseTime      : 0s
IsConnReused:     : false
RemoteAddr        : 100.26.90.23:443

the request total time is 4.0295301s, and costs 3.5787528s from connection ready to server respond first byte
```
