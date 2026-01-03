package version

// Version information, set via ldflags during build.
var (
	// Version is the current version of the application.
	Version = "0.1.0"

	// Commit is the git commit hash.
	Commit = "unknown"

	// Build is the build timestamp.
	Build = "unknown"
)
