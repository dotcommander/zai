package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAPIURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		baseURL           string
		codingBaseURL     string
		wantBaseURL       string
		wantCodingBaseURL string
	}{
		{
			name:              "standard URL unchanged",
			baseURL:           "https://api.z.ai/api/paas/v4",
			codingBaseURL:     "https://api.z.ai/api/coding/paas/v4",
			wantBaseURL:       "https://api.z.ai/api/paas/v4",
			wantCodingBaseURL: "https://api.z.ai/api/coding/paas/v4",
		},
		{
			name:              "coding URL in base_url gets normalized",
			baseURL:           "https://api.z.ai/api/coding/paas/v4",
			codingBaseURL:     "https://api.z.ai/api/coding/paas/v4",
			wantBaseURL:       "https://api.z.ai/api/paas/v4",
			wantCodingBaseURL: "https://api.z.ai/api/coding/paas/v4",
		},
		{
			name:              "coding URL in base_url with empty coding_base_url",
			baseURL:           "https://api.z.ai/api/coding/paas/v4",
			codingBaseURL:     "",
			wantBaseURL:       "https://api.z.ai/api/paas/v4",
			wantCodingBaseURL: "https://api.z.ai/api/coding/paas/v4",
		},
		{
			name:              "both empty passthrough",
			baseURL:           "",
			codingBaseURL:     "",
			wantBaseURL:       "",
			wantCodingBaseURL: "",
		},
		{
			name:              "standard URL with empty coding preserves both",
			baseURL:           "https://api.z.ai/api/paas/v4",
			codingBaseURL:     "",
			wantBaseURL:       "https://api.z.ai/api/paas/v4",
			wantCodingBaseURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotBase, gotCoding := normalizeAPIURLs(tt.baseURL, tt.codingBaseURL)
			assert.Equal(t, tt.wantBaseURL, gotBase, "baseURL")
			assert.Equal(t, tt.wantCodingBaseURL, gotCoding, "codingBaseURL")
		})
	}
}
