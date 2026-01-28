package main

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/installable-sh/docker/internal/certs"
	"github.com/installable-sh/docker/internal/version"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type parsedArgs struct {
	showHelp    bool
	showVersion bool
	sendEnv     bool
	raw         bool
	noCache     bool
	url         string
	scriptArgs  []string
}

type fetchedScript struct {
	content string
	name    string
}

func parseArgs(args []string) parsedArgs {
	result := parsedArgs{}

	// Find the URL (first arg starting with http:// or https://)
	urlIndex := -1
	for i, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			urlIndex = i
			break
		}
	}

	var runArgs []string
	if urlIndex >= 0 {
		runArgs = args[:urlIndex]
		result.url = args[urlIndex]
		result.scriptArgs = args[urlIndex+1:]
	} else {
		runArgs = args
	}

	for _, arg := range runArgs {
		switch arg {
		case "--help", "-h":
			result.showHelp = true
		case "--version", "-v":
			result.showVersion = true
		case "+env":
			result.sendEnv = true
		case "+raw":
			result.raw = true
		case "+nocache":
			result.noCache = true
		}
	}

	return result
}

func main() {
	args := parseArgs(os.Args[1:])

	if args.showVersion {
		version.Print("RUN")
		os.Exit(0)
	}

	if args.showHelp || args.url == "" {
		fmt.Println("usage: RUN [+env] [+raw] [+nocache] <url> [args...]")
		fmt.Println("  +env      Send environment variables as X-Env-* headers")
		fmt.Println("  +raw      Print the script without executing")
		fmt.Println("  +nocache  Bypass CDN caches")
		os.Exit(0)
	}

	client, err := newHTTPClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create HTTP client: %v\n", err)
		os.Exit(1)
	}

	script, err := fetchScript(client, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch script: %v\n", err)
		os.Exit(1)
	}

	if args.raw {
		fmt.Print(script.content)
		return
	}

	if err := runScript(script, args.scriptArgs); err != nil {
		fmt.Fprintf(os.Stderr, "script error: %v\n", err)
		os.Exit(1)
	}
}

func isValidHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		// HTTP header names must be tokens (RFC 7230)
		// Allow: A-Z a-z 0-9 ! # $ % & ' * + - . ^ _ ` | ~
		if c <= ' ' || c >= 127 || strings.ContainsRune("\"(),/:;<=>?@[\\]{}", c) {
			return false
		}
	}
	return true
}

func fetchScript(client *retryablehttp.Client, args parsedArgs) (fetchedScript, error) {
	req, err := retryablehttp.NewRequest("GET", args.url, nil)
	if err != nil {
		return fetchedScript{}, err
	}

	userAgent := "run/1.0 (installable)"
	if ua := os.Getenv("USER_AGENT"); ua != "" {
		userAgent = ua
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/plain, text/x-shellscript, application/x-sh, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	if args.noCache {
		req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		req.Header.Set("Pragma", "no-cache")
	}

	if args.sendEnv {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && isValidHeaderName(parts[0]) {
				req.Header.Set("X-Env-"+parts[0], parts[1])
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fetchedScript{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fetchedScript{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Determine script name from Content-Disposition or URL path
	scriptName := ""
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		_, params, err := mime.ParseMediaType(cd)
		if err == nil && params["filename"] != "" {
			scriptName = params["filename"]
		}
	}
	if scriptName == "" {
		scriptName = path.Base(args.url)
		if scriptName == "" || scriptName == "/" || scriptName == "." {
			scriptName = "script.sh"
		}
	}

	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fetchedScript{}, fmt.Errorf("gzip error: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return fetchedScript{}, err
	}

	return fetchedScript{content: string(content), name: scriptName}, nil
}

func newHTTPClient() (*retryablehttp.Client, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certs.CACerts) {
		return nil, fmt.Errorf("failed to parse embedded CA certificates")
	}

	client := retryablehttp.NewClient()
	client.RetryMax = 0 // Unlimited retries
	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	return client, nil
}

func runScript(script fetchedScript, args []string) error {
	parser := syntax.NewParser()
	prog, err := parser.Parse(strings.NewReader(script.content), script.name)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	runner, err := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.Params(args...),
	)
	if err != nil {
		return fmt.Errorf("interpreter error: %w", err)
	}

	return runner.Run(context.Background(), prog)
}
