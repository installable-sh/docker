package version

import (
	"fmt"
	"runtime/debug"
)

// Get returns version information from build metadata.
// It uses VCS info embedded by Go at build time.
func Get() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) >= 7 {
				revision = setting.Value[:7]
			} else {
				revision = setting.Value
			}
		case "vcs.modified":
			if setting.Value == "true" {
				modified = "-dirty"
			}
		}
	}

	if revision == "" {
		return "dev"
	}

	return revision + modified
}

// Print outputs version information for a command.
func Print(command string) {
	fmt.Printf("%s version %s\n", command, Get())
}
