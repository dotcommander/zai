package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/dotcommander/zai/internal/config"
)

// QuotaClient uses a separate, no-redirect monitoring transport.
type QuotaClient struct {
	APIKey        string
	BaseURL       string
	CodingBaseURL string
	Settings      config.QuotaConfig
	Doer          HTTPDoer
}

type QuotaUnavailable struct{ Reason string }

func (e *QuotaUnavailable) Error() string { return "quota unavailable: " + e.Reason }

func quotaReason(err error) string {
	var unavailable *QuotaUnavailable
	if errors.As(err, &unavailable) {
		return unavailable.Reason
	}
	return "request_failed"
}

func quotaOrigin(value string, hosts []string) (string, bool) {
	u, err := url.Parse(value)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" || u.Opaque != "" || u.RawQuery != "" || u.Fragment != "" || u.Host != u.Hostname() || !slices.Contains(hosts, u.Hostname()) {
		return "", false
	}
	return "https://" + u.Hostname(), true
}

// Usage retrieves one bounded monitoring response without generation, retries or redirects.
func (c *QuotaClient) Usage(ctx context.Context) ([]QuotaWindow, error) {
	route, hosts, err := config.QuotaSafety()
	if err != nil {
		return nil, &QuotaUnavailable{Reason: "unsupported_endpoint"}
	}
	standard, valid := quotaOrigin(c.BaseURL, hosts)
	coding, codingValid := quotaOrigin(c.CodingBaseURL, hosts)
	if !valid || !codingValid || standard != coding {
		return nil, &QuotaUnavailable{Reason: "unsupported_endpoint"}
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return nil, &QuotaUnavailable{Reason: "missing_credentials"}
	}
	if c.Settings.Timeout <= 0 || c.Settings.MaxResponseBytes <= 0 || c.Settings.MaxResponseBytes > 10485760 {
		return nil, &QuotaUnavailable{Reason: "invalid_response"}
	}
	requestCtx, cancel := context.WithTimeout(ctx, c.Settings.Timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, standard+route, nil)
	if err != nil {
		return nil, &QuotaUnavailable{Reason: "unsupported_endpoint"}
	}
	request.Header.Set("Authorization", c.APIKey)
	request.Header.Set("Accept", "application/json")
	doer := c.Doer
	if doer == nil {
		doer = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	}
	response, err := doer.Do(request)
	if err != nil {
		return nil, &QuotaUnavailable{Reason: "request_failed"}
	}
	if response == nil || response.Body == nil {
		return nil, &QuotaUnavailable{Reason: "invalid_response"}
	}
	defer func() { _ = response.Body.Close() }() // Best-effort cleanup must not expose raw transport errors.
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return nil, &QuotaUnavailable{Reason: "authentication_failed"}
	}
	if response.StatusCode >= 300 && response.StatusCode < 400 || response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusMethodNotAllowed || response.StatusCode == http.StatusNotImplemented {
		return nil, &QuotaUnavailable{Reason: "unsupported_endpoint"}
	}
	if response.StatusCode != http.StatusOK {
		return nil, &QuotaUnavailable{Reason: "request_failed"}
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, c.Settings.MaxResponseBytes+1))
	if err != nil || int64(len(body)) > c.Settings.MaxResponseBytes {
		return nil, &QuotaUnavailable{Reason: "invalid_response"}
	}
	if code := quotaCode(body); code == "401" || code == "403" {
		return nil, &QuotaUnavailable{Reason: "authentication_failed"}
	}
	windows, err := parseQuotaWindows(body)
	if err != nil {
		return nil, &QuotaUnavailable{Reason: "invalid_response"}
	}
	return windows, nil
}
