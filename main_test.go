package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

// Http response
type HttpResponse struct {
	Body                string   `json:"body"`
	Method              string   `json:"method"`
	AuthorizationHeader []string `json:"auth_header"`
	ContentTypeHeader   []string `json:"content_header"`
	UserAgentHeader     []string `json:"agent_header"`
}

func Example() {
	// Create http server to run our tests against
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()

		// Create our response
		httpResp := HttpResponse{
			Body:                string(body),
			Method:              r.Method,
			AuthorizationHeader: r.Header["Authorization"],
			ContentTypeHeader:   r.Header["Content-Type"],
			UserAgentHeader:     r.Header["User-Agent"],
		}

		// Encode response to JSON
		jsonResp, _ := json.Marshal(httpResp)
		// Instead of writing the response back, print directly to stdout
		fmt.Println(string(jsonResp))
	}))
	defer testServer.Close()

	// Create JSON file for a test later
	jsonTestFile := "test.json"
	jsonData := []byte("{\"fromfile\": true}")
	err := os.WriteFile(jsonTestFile, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// Delete file when tests ends
	defer os.Remove(jsonTestFile)

	// Define tests
	tests := []struct {
		cmd []string
	}{
		{[]string{"./gurl", testServer.URL, "-X", "POST", "-u", "user:pass", "--impersonate", "chrome", "-d", "{\"name\": \"user\"}"}},
		{[]string{"./gurl", testServer.URL, "-X", "PUT", "-b", "Token", "-d", "{\"name\": \"user\"}", "-H", "h1=v1;h2=v2"}},
		{[]string{"./gurl", testServer.URL, "-X", "PATCH", "-c", "text", "--raw-response", "--impersonate", "firefox", "-d", "{\"name\": \"user\"}"}},
		{[]string{"./gurl", testServer.URL, "-X", "DELETE", "--impersonate", "safari", "--force-http-1"}},
		{[]string{"./gurl", testServer.URL, "-X", "DELETE", "-d", "{\"name\": \"user\"}", "--raw-response"}},
		{[]string{"./gurl", testServer.URL, "-X", "OPTIONS", "--disable-redirect"}},
		{[]string{"./gurl", testServer.URL, "-X", "GET", "--tls-finger", "android", "-k", "-r", "1"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "chrome", "-A", "MyAgent"}},
		{[]string{"./gurl", testServer.URL, "-X", "POST", "--tls-finger", "firefox", "-f", "test.json"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "edge"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "safari"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "ios"}},
		{[]string{"./gurl", testServer.URL, "--tls-finger", "random"}},
	}

	// Run tests
	for _, t := range tests {
		os.Args = t.cmd
		run()
	}

	// Output:
	// {"body":"{\"name\": \"user\"}","method":"POST","auth_header":["Basic dXNlcjpwYXNz"],"content_header":["application/json; charset=utf-8"],"agent_header":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"]}
	// {"body":"{\"name\": \"user\"}","method":"PUT","auth_header":["Bearer Token"],"content_header":["application/json; charset=utf-8"],"agent_header":["gurl 1.4.0"]}
	// {"body":"{\"name\": \"user\"}","method":"PATCH","auth_header":null,"content_header":["text"],"agent_header":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:105.0) Gecko/20100101 Firefox/105.0"]}
	// {"body":"","method":"DELETE","auth_header":null,"content_header":null,"agent_header":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15"]}
	// {"body":"{\"name\": \"user\"}","method":"DELETE","auth_header":null,"content_header":["text/plain; charset=utf-8"],"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"OPTIONS","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["MyAgent"]}
	// {"body":"{\"fromfile\": true}","method":"POST","auth_header":null,"content_header":["application/json; charset=utf-8"],"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
	// {"body":"","method":"GET","auth_header":null,"content_header":null,"agent_header":["gurl 1.4.0"]}
}
