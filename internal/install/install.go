package install

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/installable-sh/docker/v1/internal/version"
)

// Install represents the INSTALL command with parsed arguments.
type Install struct {
	ShowHelp    bool
	ShowVersion bool

	// IO streams (defaults to os.Stdin/Stdout/Stderr)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New parses command-line arguments and returns an Install command.
func New(args []string) *Install {
	i := &Install{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			i.ShowHelp = true
		case "--version", "-v":
			i.ShowVersion = true
		}
	}

	return i
}

// Exec executes the INSTALL command.
func (i *Install) Exec(_ context.Context) error {
	if i.ShowVersion {
		version.Print("INSTALL")
		return nil
	}

	if i.ShowHelp {
		_, _ = fmt.Fprintln(i.Stdout, "usage: INSTALL [options] <url> [args...]")
		_, _ = fmt.Fprintln(i.Stdout)
		_, _ = fmt.Fprintln(i.Stdout, "INSTALL is under development and will be available in a future release.")
		_, _ = fmt.Fprintln(i.Stdout, "It is intended for installation and setup tasks during Docker image builds.")
		return nil
	}

	_, _ = fmt.Fprintln(i.Stderr, "INSTALL: coming soon")
	return fmt.Errorf("not implemented")
}
