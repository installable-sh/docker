package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

//go:embed hack/ca-certificates.crt
var caCerts []byte

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("usage: RUN [+env] <url> [args...]")
		fmt.Println("  +env  Send environment variables as X-Env-* headers")
		os.Exit(0)
	}

	args := os.Args[1:]
	sendEnv := false

	if args[0] == "+env" {
		sendEnv = true
		args = args[1:]
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "error: +env requires a URL")
			os.Exit(1)
		}
	}

	url := args[0]
	scriptArgs := args[1:]

	client, err := newHTTPClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create HTTP client: %v\n", err)
		os.Exit(1)
	}

	script, err := fetchScript(client, url, sendEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch script: %v\n", err)
		os.Exit(1)
	}

	if err := runScript(script, scriptArgs); err != nil {
		fmt.Fprintf(os.Stderr, "script error: %v\n", err)
		os.Exit(1)
	}
}

func fetchScript(client *http.Client, url string, sendEnv bool) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if sendEnv {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				req.Header.Set("X-Env-"+parts[0], parts[1])
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	script, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(script), nil
}

func newHTTPClient() (*http.Client, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCerts) {
		return nil, fmt.Errorf("failed to parse embedded CA certificates")
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}, nil
}

func runScript(script string, args []string) error {
	parser := syntax.NewParser()
	prog, err := parser.Parse(strings.NewReader(script), "script.sh")
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
