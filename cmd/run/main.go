package main

import (
	"fmt"
	"os"

	"github.com/installable-sh/docker/v1/internal/run"
)

func main() {
	cmd := run.New(os.Args[1:])
	if err := cmd.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
