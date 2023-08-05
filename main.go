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
const gurlVersion = "1.0.0"

// Cli arguments
var cli struct {
	Auth            string            `help:"Basic HTTP authentication in the format username:password." short:"a"`
	BearerToken     string            `help:"Set bearer auth token." short:"b"`
	DisableRedirect bool              `help:"Disable redirects." default:"false"`
	Headers         map[string]string `help:"HTTP headers in the format: \"header1=value1;header2=value2\"." short:"H"`
	Insecure        bool              `help:"Allow insecure SSL connections." short:"i" default:"false"`
	RawResponse     bool              `help:"Print raw response string (disable json prettify)." default:"false"`
	Timeout         int               `help:"Timeout in milliseconds." short:"t" default:"10000"`
	TlsFinger       string            `help:"TLS Fingerprint: chrome, firefox, edge, safari, ios, android, random or go." enum:"chrome, firefox, edge, safari, ios, android, random, go" default:"go"`
	Trace           bool              `help:"Show tracing/performance information." default:"false"`
	UserAgent       string            `help:"Set User-Agent http header." short:"u"`
	Verbose         bool              `help:"Enable verbose/debug mode." short:"v" default:"false"`

	Get struct {
		Url string `arg:"" help:"Url to access."`
	} `cmd:"" help:"GET HTTP method."`
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
func setContentHeader(httpMethod string) {
	methods := []string{"post", "put", "patch"}
	if slices.Contains(methods, httpMethod) {
		req.SetCommonContentType("application/json; charset=utf-8")
	}
}

// Configure our http request
func configRequest(ctx *kong.Context) {
	// Set default http scheme if no scheme is provided
	req.SetScheme("http")
	// Set application/json content-type for post, put and patch methods
	setContentHeader(ctx.Args[0])
	if cli.Auth != "" {
		splitAuth := strings.Split(cli.Auth, ":")
		user, pass := splitAuth[0], splitAuth[1]
		req.SetCommonBasicAuth(user, pass)
	}
	if cli.BearerToken != "" {
		req.SetCommonBearerAuthToken(cli.BearerToken)
	}
	if cli.DisableRedirect {
		req.SetRedirectPolicy(req.NoRedirectPolicy())
	}
	if len(cli.Headers) > 0 {
		req.SetCommonHeaders(cli.Headers)
	}
	if cli.Insecure {
		req.EnableInsecureSkipVerify()
	}
	// Sites to check finger hash:
	// - https://tls.peet.ws/api/clean
	// - https://tools.scrapfly.io/api/fp/ja3
	if cli.TlsFinger != "go" {
		switch cli.TlsFinger {
		case "chrome":
			req.SetTLSFingerprintChrome()
		case "firefox":
			req.SetTLSFingerprintFirefox()
		case "edge":
			req.SetTLSFingerprintEdge()
		case "safari":
			req.SetTLSFingerprintSafari()
		case "ios":
			req.SetTLSFingerprintIOS()
		case "android":
			req.SetTLSFingerprintAndroid()
		case "random":
			req.SetTLSFingerprintRandomized()
		}
	}
	if cli.Trace {
		req.EnableTraceAll()
	}
	if cli.UserAgent != "" {
		req.SetUserAgent(cli.UserAgent)
	} else {
		req.SetUserAgent(fmt.Sprintf("gurl %v", gurlVersion))
	}
	if cli.Verbose {
		req.EnableDumpAllWithoutBody().EnableDebugLog()
	}
	req.SetTimeout(time.Duration(cli.Timeout) * time.Millisecond)
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

// Log error and quit if an error ocurred
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Do HEAD request
func doHeadRequest() *req.Response {
	// Always dump headers for HEAD requests
	req.EnableDumpAllWithoutBody()
	resp, err := req.Head(cli.Head.Url)
	checkErr(err)
	return resp
}

// Do GET request
func doGetRequest() *req.Response {
	resp, err := req.Get(cli.Get.Url)
	checkErr(err)
	return resp
}

// Do POST request
func doPostRequest() *req.Response {
	if cli.Post.ContentType != "" {
		req.SetCommonContentType(cli.Post.ContentType)
	}
	resp, err := req.SetBody(cli.Post.Data).Post(cli.Post.Url)
	checkErr(err)
	return resp
}

// Do PUT request
func doPutRequest() *req.Response {
	if cli.Put.ContentType != "" {
		req.SetCommonContentType(cli.Put.ContentType)
	}
	resp, err := req.SetBody(cli.Put.Data).Put(cli.Put.Url)
	checkErr(err)
	return resp
}

// Do PATCH request
func doPatchRequest() *req.Response {
	if cli.Patch.ContentType != "" {
		req.SetCommonContentType(cli.Patch.ContentType)
	}
	resp, err := req.SetBody(cli.Patch.Data).Patch(cli.Patch.Url)
	checkErr(err)
	return resp
}

// Do DELETE request
func doDeleteRequest() *req.Response {
	if cli.Delete.ContentType != "" {
		req.SetCommonContentType(cli.Delete.ContentType)
	}
	// We need to declare the vars outside if/else scope to avoid unused/undeclared vars errors
	var resp *req.Response
	var err error
	// According to https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/DELETE
	// DELETE method may have a body
	if cli.Delete.Data != "" {
		resp, err = req.SetBody(cli.Delete.Data).Delete(cli.Delete.Url)
	} else {
		resp, err = req.Delete(cli.Delete.Url)
	}
	checkErr(err)
	return resp
}

// Do OPTIONS request
func doOptionsRequest() *req.Response {
	resp, err := req.Options(cli.Options.Url)
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
	// Configure http request based on cli arguments
	configRequest(ctx)
	// Store response pointer
	var resp *req.Response
	// Execute cli command accordingly
	switch ctx.Command() {
	case "version":
		fmt.Println(gurlVersion)
		os.Exit(0)
	case "head <url>":
		resp = doHeadRequest()
	case "get <url>":
		resp = doGetRequest()
	case "post <url>":
		resp = doPostRequest()
	case "put <url>":
		resp = doPutRequest()
	case "patch <url>":
		resp = doPatchRequest()
	case "delete <url>":
		resp = doDeleteRequest()
	case "options <url>":
		resp = doOptionsRequest()
	}
	// Print raw response or a prettified json
	printResponse(resp.String())
	// Show trace info if needed
	showTraceInfo(resp)
}
