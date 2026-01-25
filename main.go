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
		fmt.Println("usage: RUN <url> [args...]")
		os.Exit(0)
	}

	url := os.Args[1]

	client, err := newHTTPClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create HTTP client: %v\n", err)
		os.Exit(1)
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch script: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "failed to fetch script: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	script, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read script: %v\n", err)
		os.Exit(1)
	}

	if err := runScript(string(script), os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "script error: %v\n", err)
		os.Exit(1)
	}
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
