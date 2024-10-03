package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/alecthomas/kong"
	"github.com/hokaccha/go-prettyjson"
	"github.com/imroc/req/v3"
)

// Program version
const gurlVersion = "1.7.0"

// Cli arguments
var cli struct {
	// Args
	Url string `arg:"" help:"Url to access."`

	// Flags
	Auth            string            `help:"Basic HTTP authentication in the format username:password." placeholder:"auth" short:"u"`
	BearerToken     string            `help:"Set bearer auth token." placeholder:"token" short:"b"`
	CACert          string            `help:"CA certificate file." placeholder:"file" type:"path"`
	ClientCert      []string          `help:"Client cert and key files separated by comma: \"cert.pem,key.pem\"." placeholder:"cert-file,key-file" type:"path"`
	ContentType     string            `help:"Content-Type http header, default is application/json for POST, PUT and PATCH methods." placeholder:"content" short:"c"`
	Data            string            `help:"Data payload (request body)." xor:"data" placeholder:"payload" short:"d"`
	DataFile        string            `help:"Read data payload from file." xor:"data" placeholder:"file" short:"f" type:"path"`
	DisableRedirect bool              `help:"Disable redirects." default:"false"`
	ForceHttp1      bool              `help:"Force HTTP/1.1 to be used." default:"false"`
	Headers         map[string]string `help:"HTTP headers in the format: \"header1=value1;header2=value2\"." placeholder:"h1=v1;h2=v2" short:"H"`
	Impersonate     string            `help:"Fully impersonate chrome, firefox or safari browser (this will automatically set headers, headers order and tls fingerprint)." enum:"chrome, firefox, safari, none" default:"none"`
	Insecure        bool              `help:"Allow insecure SSL connections." short:"k" default:"false"`
	Method          string            `help:"Http method: GET, HEAD, POST, PUT, PATCH, DELETE or OPTIONS." enum:"GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS" short:"X" default:"GET"`
	OutputFile      string            `help:"Save response to file." short:"o" placeholder:"file" type:"path"`
	Proxy           string            `help:"Proxy to use, e.g.: \"http://user:pass@myproxy:8080\"." placeholder:"proxy"`
	RawResponse     bool              `help:"Print raw response string (disable json prettify)." default:"false"`
	Retries         int               `help:"Number of retries in case of errors and http status code >= 500." short:"r" default:"0"`
	Timeout         int               `help:"Timeout in milliseconds." short:"t" default:"10000"`
	TlsFinger       string            `help:"TLS Fingerprint: chrome, firefox, edge, safari, ios, android, random or go." enum:"chrome, firefox, edge, safari, ios, android, random, go" default:"go"`
	Trace           bool              `help:"Show tracing/performance information." default:"false"`
	UserAgent       string            `help:"Set User-Agent http header." placeholder:"agent" short:"A"`
	Verbose         bool              `help:"Enable verbose/debug mode." short:"v" default:"false"`
	Version         kong.VersionFlag  `help:"Show version and exit." short:"V"`
}

// Set application/json content-type http header for POST, PUT and PATCH methods
func setContentHeader(httpMethod string, request *req.Request) {
	methods := []string{"POST", "PUT", "PATCH"}
	if slices.Contains(methods, httpMethod) {
		request.SetContentType("application/json; charset=utf-8")
	}
}

// Configure our http request
func configRequest(request *req.Request) {
	// Set default http scheme if no scheme is provided
	request.GetClient().SetScheme("http")
	// Set client timeout
	request.GetClient().SetTimeout(time.Duration(cli.Timeout) * time.Millisecond)
	// Set application/json content-type for POST, PUT and PATCH methods
	setContentHeader(cli.Method, request)
	if cli.Auth != "" {
		splitAuth := strings.Split(cli.Auth, ":")
		user, pass := splitAuth[0], splitAuth[1]
		request.SetBasicAuth(user, pass)
	}
	if cli.BearerToken != "" {
		request.SetBearerAuthToken(cli.BearerToken)
	}
	if cli.CACert != "" {
		request.GetClient().SetRootCertsFromFile(cli.CACert)
	}
	if len(cli.ClientCert) > 1 {
		request.GetClient().SetCertFromFile(cli.ClientCert[0], cli.ClientCert[1])
	}
	if cli.DisableRedirect {
		request.GetClient().SetRedirectPolicy(req.NoRedirectPolicy())
	}
	if cli.ForceHttp1 {
		request.GetClient().EnableForceHTTP1()
	}
	if len(cli.Headers) > 0 {
		request.SetHeaders(cli.Headers)
	}
	if cli.Insecure {
		request.GetClient().EnableInsecureSkipVerify()
	}
	if cli.OutputFile != "" {
		request.SetOutputFile(cli.OutputFile)
		// Register callback to show download progress
		ProgressCallback := func(info req.DownloadInfo) {
			if info.Response.Response != nil {
				log.Printf("Downloading %.2f%%", float64(info.DownloadedSize)/float64(info.Response.ContentLength)*100.0)
			}
		}
		request.SetDownloadCallback(ProgressCallback)
	}
	if cli.Proxy != "" {
		request.GetClient().SetProxyURL(cli.Proxy)
	}
	if cli.Retries > 0 {
		request.SetRetryCount(cli.Retries)
		// Exponential backoff
		request.SetRetryBackoffInterval(1*time.Second, 5*time.Second)
		// Retry in case of errors or http status >= 500
		request.AddRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.StatusCode >= 500
		})
		// Log to stderr if a retry occurs
		request.AddRetryHook(func(resp *req.Response, err error) {
			req := resp.Request.RawRequest
			log.Printf("Retrying %v request to %v", req.Method, req.URL)
		})
	}
	// Sites to check finger hash:
	// - https://tls.peet.ws/api/clean
	// - https://tools.scrapfly.io/api/fp/ja3
	if cli.TlsFinger != "go" {
		switch cli.TlsFinger {
		case "chrome":
			request.GetClient().SetTLSFingerprintChrome()
		case "firefox":
			request.GetClient().SetTLSFingerprintFirefox()
		case "edge":
			request.GetClient().SetTLSFingerprintEdge()
		case "safari":
			request.GetClient().SetTLSFingerprintSafari()
		case "ios":
			request.GetClient().SetTLSFingerprintIOS()
		case "android":
			request.GetClient().SetTLSFingerprintAndroid()
		case "random":
			request.GetClient().SetTLSFingerprintRandomized()
		}
	}
	if cli.Trace {
		request.GetClient().EnableTraceAll()
	}
	if cli.UserAgent != "" {
		request.GetClient().SetUserAgent(cli.UserAgent)
	} else {
		request.GetClient().SetUserAgent(fmt.Sprintf("gurl %v", gurlVersion))
	}
	if cli.Verbose {
		request.GetClient().EnableDumpAllWithoutBody().EnableDebugLog()
	}
	// Set impersonate as the last step to override possible earlier configurations
	if cli.Impersonate != "none" {
		switch cli.Impersonate {
		case "chrome":
			request.GetClient().ImpersonateChrome()
		case "firefox":
			request.GetClient().ImpersonateFirefox()
		case "safari":
			request.GetClient().ImpersonateSafari()
		}
	}
}

// Print raw string response or a prettified json if possible
func printResponse(rawStr string) {
	if cli.RawResponse {
		fmt.Print(rawStr)
		return
	}
	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(rawStr), &jsonObj)
	if err != nil {
		fmt.Print(rawStr)
	} else {
		prettyJson, _ := prettyjson.Marshal(jsonObj)
		fmt.Print(string(prettyJson))
	}
}

// Print trace information
func showTraceInfo(resp *req.Response) {
	if cli.Trace {
		trace := resp.TraceInfo()
		fmt.Print("\n\n------- TRACE INFO -------\n")
		fmt.Println(trace)
		fmt.Printf("\n%v\n", trace.Blame())
	}
}

// Log error and quit if an error occurred
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Read file content
func readFile(filename string) string {
	content, err := os.ReadFile(filename)
	checkErr(err)
	return string(content)
}

// Do HEAD request
func doHeadRequest(request *req.Request) *req.Response {
	// Always dump headers for HEAD requests
	request.GetClient().EnableDumpAllWithoutBody()
	resp, err := request.Head(cli.Url)
	checkErr(err)
	return resp
}

// Do GET request
func doGetRequest(request *req.Request) *req.Response {
	resp, err := request.Get(cli.Url)
	checkErr(err)
	return resp
}

// Configure payload
func configPayload(request *req.Request) string {
	if cli.ContentType != "" {
		request.SetContentType(cli.ContentType)
	}
	var payload string
	if cli.Data != "" {
		payload = cli.Data
	} else if cli.DataFile != "" {
		payload = readFile(cli.DataFile)
	}
	return payload
}

// Do POST request
func doPostRequest(request *req.Request) *req.Response {
	payload := configPayload(request)
	if payload != "" {
		request.SetBody(payload)
	}
	resp, err := request.Post(cli.Url)
	checkErr(err)
	return resp
}

// Do PUT request
func doPutRequest(request *req.Request) *req.Response {
	payload := configPayload(request)
	if payload != "" {
		request.SetBody(payload)
	}
	resp, err := request.Put(cli.Url)
	checkErr(err)
	return resp
}

// Do PATCH request
func doPatchRequest(request *req.Request) *req.Response {
	payload := configPayload(request)
	if payload != "" {
		request.SetBody(payload)
	}
	resp, err := request.Patch(cli.Url)
	checkErr(err)
	return resp
}

// Do DELETE request
func doDeleteRequest(request *req.Request) *req.Response {
	// According to https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/DELETE
	// DELETE method may have a body
	payload := configPayload(request)
	if payload != "" {
		request.SetBody(payload)
	}
	resp, err := request.Delete(cli.Url)
	checkErr(err)
	return resp
}

// Do OPTIONS request
func doOptionsRequest(request *req.Request) *req.Response {
	resp, err := request.Options(cli.Url)
	checkErr(err)
	return resp
}

// Run our cli
func run() {
	// Parse cli arguments
	kong.Parse(&cli,
		kong.Name("gurl"),
		kong.Description("A simple http client cli written in Go."),
		kong.UsageOnError(),
		kong.Vars{"version": gurlVersion},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	// Create a new request object from client
	request := req.C().R()
	// Configure http request based on cli arguments
	configRequest(request)
	// Store response pointer
	var resp *req.Response
	// Execute cli command accordingly
	switch cli.Method {
	case "HEAD":
		resp = doHeadRequest(request)
	case "GET":
		resp = doGetRequest(request)
	case "POST":
		resp = doPostRequest(request)
	case "PUT":
		resp = doPutRequest(request)
	case "PATCH":
		resp = doPatchRequest(request)
	case "DELETE":
		resp = doDeleteRequest(request)
	case "OPTIONS":
		resp = doOptionsRequest(request)
	}
	// Print raw response or a prettified json
	printResponse(resp.String())
	// Show trace info if needed
	showTraceInfo(resp)
}

// Main function
func main() {
	run()
}
