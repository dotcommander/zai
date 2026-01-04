package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		commit         string
		expectedPrefix string
		expectedSuffix string
	}{
		{
			name:           "version without v prefix",
			version:        "0.1.0",
			commit:         "abc1234",
			expectedPrefix: "v0.1.0",
			expectedSuffix: "(abc1234)",
		},
		{
			name:           "version with v prefix",
			version:        "v0.2.0",
			commit:         "def5678",
			expectedPrefix: "v0.2.0",
			expectedSuffix: "(def5678)",
		},
		{
			name:           "unknown commit",
			version:        "0.1.0",
			commit:         "unknown",
			expectedPrefix: "v0.1.0",
			expectedSuffix: "",
		},
		{
			name:           "dirty version",
			version:        "v0.1.0+dirty",
			commit:         "a0e2dac",
			expectedPrefix: "v0.1.0+dirty",
			expectedSuffix: "(a0e2dac)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := Version
			origCommit := Commit
			defer func() {
				Version = origVersion
				Commit = origCommit
			}()

			// Set test values
			Version = tt.version
			Commit = tt.commit

			result := String()

			if !strings.HasPrefix(result, tt.expectedPrefix) {
				t.Errorf("String() = %q, want prefix %q", result, tt.expectedPrefix)
			}

			if tt.expectedSuffix != "" {
				if !strings.Contains(result, tt.expectedSuffix) {
					t.Errorf("String() = %q, want to contain %q", result, tt.expectedSuffix)
				}
			} else {
				if strings.Contains(result, "(") {
					t.Errorf("String() = %q, should not contain commit info", result)
				}
			}
		})
	}
}

func TestVersionPrefix(t *testing.T) {
	// Verify no double 'v' prefix
	tests := []struct {
		version string
		want    string
	}{
		{"0.1.0", "v0.1.0"},
		{"v0.1.0", "v0.1.0"},
		{"v1.2.3+dirty", "v1.2.3+dirty"},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			origVersion := Version
			origCommit := Commit
			defer func() {
				Version = origVersion
				Commit = origCommit
			}()

			Version = tt.version
			Commit = "unknown"

			result := String()
			if result != tt.want {
				t.Errorf("String() = %q, want %q", result, tt.want)
			}
		})
	}
}
