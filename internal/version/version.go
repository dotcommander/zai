package version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Version information, extracted automatically from debug.BuildInfo.
var (
	// Version is the current version of the application.
	Version = "0.1.0"

	// Commit is the git commit hash (first 7 chars).
	Commit = "unknown"

	// Build is the build timestamp.
	Build = "unknown"
)

func init() {
	// Extract version info from debug.BuildInfo automatically
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// Get version from module (if available)
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		Version = info.Main.Version
	}

	// Extract VCS info from build settings
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			// Take first 7 chars of commit hash (like git short sha)
			if len(setting.Value) >= 7 {
				Commit = setting.Value[:7]
			} else {
				Commit = setting.Value
			}
		case "vcs.time":
			Build = setting.Value
		}
	}
}

// String returns a formatted version string.
func String() string {
	var b strings.Builder
	// Version might already have 'v' prefix from module version
	if strings.HasPrefix(Version, "v") {
		b.WriteString(Version)
	} else {
		fmt.Fprintf(&b, "v%s", Version)
	}
	if Commit != "unknown" {
		fmt.Fprintf(&b, " (%s)", Commit)
	}
	return b.String()
}
