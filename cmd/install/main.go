package main

import (
	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			fmt.Println("usage: INSTALL [options] <url> [args...]")
			fmt.Println()
			fmt.Println("INSTALL is under development and will be available in a future release.")
			fmt.Println("It is intended for installation and setup tasks during Docker image builds.")
			os.Exit(0)
		}
	}

	fmt.Fprintln(os.Stderr, "INSTALL: coming soon")
	os.Exit(1)
}
