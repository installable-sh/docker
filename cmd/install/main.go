package main

import (
	"fmt"
	"os"

	"github.com/installable-sh/docker/internal/version"
)

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			fmt.Println("usage: INSTALL [options] <url> [args...]")
			fmt.Println()
			fmt.Println("INSTALL is under development and will be available in a future release.")
			fmt.Println("It is intended for installation and setup tasks during Docker image builds.")
			os.Exit(0)
		case "--version", "-v":
			version.Print("INSTALL")
			os.Exit(0)
		}
	}

	fmt.Fprintln(os.Stderr, "INSTALL: coming soon")
	os.Exit(1)
}
