package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/TylerBrock/colorjson"
	"github.com/alecthomas/kong"
	"github.com/imroc/req/v3"
)

// Program version
const gurlVersion = "1.1.0"

// Cli arguments
var cli struct {
	Auth            string            `help:"Basic HTTP authentication in the format username:password." short:"a"`
	BearerToken     string            `help:"Set bearer auth token." short:"b"`
	DisableRedirect bool              `help:"Disable redirects." default:"false"`
	Headers         map[string]string `help:"HTTP headers in the format: \"header1=value1;header2=value2\"." short:"H"`
	Impersonate     string            `help:"Fully impersonate chrome, firefox or safari browser (this will automatically set headers, headers order and tls fingerprint)." enum:"chrome, firefox, safari, none" default:"none"`
	Insecure        bool              `help:"Allow insecure SSL connections." short:"i" default:"false"`
	RawResponse     bool              `help:"Print raw response string (disable json prettify)." default:"false"`
	Retries         int               `help:"Number of retries in case of errors and http status code >= 500." short:"r" default:"0"`
	Timeout         int               `help:"Timeout in milliseconds." short:"t" default:"10000"`
	TlsFinger       string            `help:"TLS Fingerprint: chrome, firefox, edge, safari, ios, android, random or go." enum:"chrome, firefox, edge, safari, ios, android, random, go" default:"go"`
	Trace           bool              `help:"Show tracing/performance information." default:"false"`
	UserAgent       string            `help:"Set User-Agent http header." short:"u"`
	Verbose         bool              `help:"Enable verbose/debug mode." short:"v" default:"false"`

	Get struct {
		Url string `arg:"" help:"Url to access."`
	} `cmd:"" help:"GET HTTP method." default:"withargs"`
	Head struct {
		Url string `arg:"" help:"Url to access."`
	} `cmd:"" help:"HEAD HTTP method."`
	Post struct {
		Url         string `arg:"" help:"Url to access."`
		Data        string `help:"Data payload (request body)." short:"d"`
		ContentType string `help:"Content-Type http header, default is application/json." short:"c"`
	} `cmd:"" help:"POST HTTP method."`
	Put struct {
		Url         string `arg:"" help:"Url to access."`
		Data        string `help:"Data payload (request body)." short:"d"`
		ContentType string `help:"Content-Type http header, default is application/json." short:"c"`
	} `cmd:"" help:"PUT HTTP method."`
	Patch struct {
		Url         string `arg:"" help:"Url to access."`
		Data        string `help:"Data payload (request body)." short:"d"`
		ContentType string `help:"Content-Type http header, default is application/json." short:"c"`
	} `cmd:"" help:"PATCH HTTP method."`
	Delete struct {
		Url         string `arg:"" help:"Url to access."`
		Data        string `help:"Data payload (request body)." short:"d"`
		ContentType string `help:"Content-Type http header, default is text/plain." short:"c"`
	} `cmd:"" help:"DELETE HTTP method."`
	Options struct {
		Url string `arg:"" help:"Url to access."`
	} `cmd:"" help:"OPTIONS HTTP method."`
	Version struct{} `cmd:"" help:"Show version and exit."`
}

// Set application/json content-type http header for post, put and patch methods
func setContentHeader(httpMethod string, request *req.Request) {
	methods := []string{"post", "put", "patch"}
	if slices.Contains(methods, httpMethod) {
		request.SetContentType("application/json; charset=utf-8")
	}
}

// Configure our http request
func configRequest(ctx *kong.Context, request *req.Request) {
	// Set default http scheme if no scheme is provided
	request.GetClient().SetScheme("http")
	// Set client timeout
	request.GetClient().SetTimeout(time.Duration(cli.Timeout) * time.Millisecond)
	// Set application/json content-type for post, put and patch methods
	setContentHeader(ctx.Args[0], request)
	if cli.Auth != "" {
		splitAuth := strings.Split(cli.Auth, ":")
		user, pass := splitAuth[0], splitAuth[1]
		request.SetBasicAuth(user, pass)
	}
	if cli.BearerToken != "" {
		request.SetBearerAuthToken(cli.BearerToken)
	}
	if cli.DisableRedirect {
		request.GetClient().SetRedirectPolicy(req.NoRedirectPolicy())
	}
	if len(cli.Headers) > 0 {
		request.SetHeaders(cli.Headers)
	}
	if cli.Insecure {
		request.GetClient().EnableInsecureSkipVerify()
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
	var jsonObj map[string]interface{}
	if cli.RawResponse {
		fmt.Print(rawStr)
		return
	}
	err := json.Unmarshal([]byte(rawStr), &jsonObj)
	if err != nil {
		fmt.Print(rawStr)
	} else {
		fomatter := colorjson.NewFormatter()
		fomatter.Indent = 2
		prettyJson, _ := fomatter.Marshal(jsonObj)
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

// Do HEAD request
func doHeadRequest(request *req.Request) *req.Response {
	// Always dump headers for HEAD requests
	request.GetClient().EnableDumpAllWithoutBody()
	resp, err := request.Head(cli.Head.Url)
	checkErr(err)
	return resp
}

// Do GET request
func doGetRequest(request *req.Request) *req.Response {
	resp, err := request.Get(cli.Get.Url)
	checkErr(err)
	return resp
}

// Do POST request
func doPostRequest(request *req.Request) *req.Response {
	if cli.Post.ContentType != "" {
		request.SetContentType(cli.Post.ContentType)
	}
	resp, err := request.SetBody(cli.Post.Data).Post(cli.Post.Url)
	checkErr(err)
	return resp
}

// Do PUT request
func doPutRequest(request *req.Request) *req.Response {
	if cli.Put.ContentType != "" {
		request.SetContentType(cli.Put.ContentType)
	}
	resp, err := request.SetBody(cli.Put.Data).Put(cli.Put.Url)
	checkErr(err)
	return resp
}

// Do PATCH request
func doPatchRequest(request *req.Request) *req.Response {
	if cli.Patch.ContentType != "" {
		request.SetContentType(cli.Patch.ContentType)
	}
	resp, err := request.SetBody(cli.Patch.Data).Patch(cli.Patch.Url)
	checkErr(err)
	return resp
}

// Do DELETE request
func doDeleteRequest(request *req.Request) *req.Response {
	if cli.Delete.ContentType != "" {
		request.SetContentType(cli.Delete.ContentType)
	}
	// We need to declare the vars outside if/else scope to avoid unused/undeclared vars errors
	var resp *req.Response
	var err error
	// According to https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/DELETE
	// DELETE method may have a body
	if cli.Delete.Data != "" {
		resp, err = request.SetBody(cli.Delete.Data).Delete(cli.Delete.Url)
	} else {
		resp, err = request.Delete(cli.Delete.Url)
	}
	checkErr(err)
	return resp
}

// Do OPTIONS request
func doOptionsRequest(request *req.Request) *req.Response {
	resp, err := request.Options(cli.Options.Url)
	checkErr(err)
	return resp
}

// Main function
func main() {
	// Parse cli arguments
	ctx := kong.Parse(&cli,
		kong.Name("gurl"),
		kong.Description("A simple http client cli written in Go."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	// Create a new request object from client
	request := req.C().R()
	// Configure http request based on cli arguments
	configRequest(ctx, request)
	// Store response pointer
	var resp *req.Response
	// Execute cli command accordingly
	switch ctx.Command() {
	case "version":
		fmt.Println(gurlVersion)
		os.Exit(0)
	case "head <url>":
		resp = doHeadRequest(request)
	case "get <url>":
		resp = doGetRequest(request)
	case "post <url>":
		resp = doPostRequest(request)
	case "put <url>":
		resp = doPutRequest(request)
	case "patch <url>":
		resp = doPatchRequest(request)
	case "delete <url>":
		resp = doDeleteRequest(request)
	case "options <url>":
		resp = doOptionsRequest(request)
	}
	// Print raw response or a prettified json
	printResponse(resp.String())
	// Show trace info if needed
	showTraceInfo(resp)
}
