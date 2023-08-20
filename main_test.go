package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func echoResp(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	fmt.Printf("%v,%v,%v,%v\n", string(body), r.Method, r.Header["Authorization"], r.Header["Content-Type"])
}

func Example() {
	// Create http server to run our tests against
	testServer := httptest.NewServer(http.HandlerFunc(echoResp))
	defer testServer.Close()

	// Define tests
	tests := []struct {
		cmd []string
	}{
		{[]string{"./gurl", testServer.URL, "-X", "post", "-u", "user:pass", "--impersonate", "chrome", "-d", "{\"name\": \"user\"}"}},
		{[]string{"./gurl", testServer.URL, "-X", "put", "-b", "Token", "-d", "{\"name\": \"user\"}", "-H", "h1=v1;h2=v2"}},
		{[]string{"./gurl", testServer.URL, "-X", "patch", "-c", "text", "--raw-response", "--impersonate", "firefox"}},
		{[]string{"./gurl", testServer.URL, "-X", "delete", "--impersonate", "safari", "--force-http-1"}},
		{[]string{"./gurl", testServer.URL, "-X", "options", "--disable-redirect"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "android", "-k", "-r", "1"}},
	}

	// Run tests
	for _, t := range tests {
		os.Args = t.cmd
		run()
	}

	// Output:
	// {"name": "user"},POST,[Basic dXNlcjpwYXNz],[application/json; charset=utf-8]
	// {"name": "user"},PUT,[Bearer Token],[application/json; charset=utf-8]
	// ,PATCH,[],[text]
	// ,DELETE,[],[]
	// ,OPTIONS,[],[]
	// ,GET,[],[]
}
