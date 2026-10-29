package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

type bufferedDoer func(*http.Request) (*http.Response, error)

func (f bufferedDoer) Do(req *http.Request) (*http.Response, error) { return f(req) }

type bufferedBody struct {
	reader   io.Reader
	closeErr error
	closes   int
}

func (b *bufferedBody) Read(p []byte) (int, error) { return b.reader.Read(p) }
func (b *bufferedBody) Close() error {
	b.closes++
	return b.closeErr
}

type bufferedErrorReader struct{ err error }

func (r bufferedErrorReader) Read([]byte) (int, error) { return 0, r.err }

type bufferedRequestCase struct {
	name, endpoint, method, url string
	post                        bool
}

type bufferedOutcome struct {
	name            string
	status          int
	transport, read error
	cancel          bool
	wantErr         error
	wantPrefix      string
}

func TestBufferedRequest(t *testing.T) {
	t.Parallel()
	transportErr := errors.New("transport failure")
	readErr := errors.New("read failure")
	requests := []bufferedRequestCase{
		{"post base", "images/generations", http.MethodPost, "https://base.example/images/generations", true},
		{"post coding", "chat/completions", http.MethodPost, "https://coding.example/chat/completions", true},
		{"post embedding", "embeddings", http.MethodPost, "https://embedding.example/embeddings", true},
		{"get", "models", http.MethodGet, "https://base.example/models", false},
	}
	outcomes := []bufferedOutcome{
		{name: "success", status: http.StatusOK},
		{name: "provider", status: http.StatusCreated},
		{name: "transport", transport: transportErr, wantErr: transportErr, wantPrefix: "failed to send request: "},
		{name: "read", status: http.StatusOK, read: readErr, wantErr: readErr, wantPrefix: "failed to read response: "},
		{name: "provider read", status: http.StatusCreated, read: readErr, wantErr: readErr, wantPrefix: "read API error: "},
		{name: "canceled", cancel: true, wantErr: context.Canceled, wantPrefix: "failed to send request: "},
	}
	for _, request := range requests {
		t.Run(request.name, func(t *testing.T) {
			t.Parallel()
			for _, outcome := range outcomes {
				t.Run(outcome.name, func(t *testing.T) {
					t.Parallel()
					runBufferedCase(t, request, outcome)
				})
			}
		})
	}
}

func checkBufferedRequest(t *testing.T, req *http.Request, want bufferedRequestCase) {
	t.Helper()
	if req.Method != want.method || req.URL.String() != want.url {
		t.Errorf("request = %s %s, want %s %s", req.Method, req.URL, want.method, want.url)
	}
	if req.Header.Get("Authorization") != "Bearer test-key" {
		t.Errorf("authorization = %q", req.Header.Get("Authorization"))
	}
	contentType, language := "", ""
	if want.post {
		contentType, language = "application/json", "en-US,en"
	}
	if req.Header.Get("Content-Type") != contentType || req.Header.Get("Accept-Language") != language {
		t.Errorf("headers = %v", req.Header)
	}
}

func runBufferedCase(t *testing.T, request bufferedRequestCase, tc bufferedOutcome) {
	t.Helper()
	ctx := context.Background()
	if tc.cancel {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		cancel()
	}
	body := &bufferedBody{reader: strings.NewReader("response bytes")}
	if tc.read != nil {
		body.reader = bufferedErrorReader{tc.read}
	}
	calls := 0
	doer := bufferedDoer(func(req *http.Request) (*http.Response, error) {
		calls++
		checkBufferedRequest(t, req, request)
		if tc.cancel {
			return nil, req.Context().Err()
		}
		if tc.transport != nil {
			return nil, tc.transport
		}
		return &http.Response{StatusCode: tc.status, Body: body, Header: http.Header{"X-Request-Id": []string{"request-123"}}}, nil
	})
	client := NewClient(ClientConfig{APIKey: "test-key", BaseURL: "https://base.example", CodingBaseURL: "https://coding.example", EmbeddingBaseURL: "https://embedding.example", UseCoding: true}, DiscardLogger(), nil, doer)
	var got []byte
	var err error
	if request.post {
		got, err = client.executeJSONRequest(ctx, request.endpoint, map[string]string{"input": "value"})
	} else {
		got, err = client.executeGetRequest(ctx, request.endpoint)
	}
	if calls != 1 {
		t.Errorf("transport calls = %d, want 1", calls)
	}
	wantCloses := 1
	if tc.transport != nil || tc.cancel {
		wantCloses = 0
	}
	if body.closes != wantCloses {
		t.Errorf("body closes = %d, want %d", body.closes, wantCloses)
	}
	checkBufferedOutcome(t, tc, got, err)
}

func checkBufferedOutcome(t *testing.T, tc bufferedOutcome, got []byte, err error) {
	t.Helper()
	switch {
	case tc.wantErr != nil:
		if !errors.Is(err, tc.wantErr) || !strings.HasPrefix(err.Error(), tc.wantPrefix) {
			t.Fatalf("error = %v, want wrapped %v with prefix %q", err, tc.wantErr, tc.wantPrefix)
		}
	case tc.status == http.StatusCreated:
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusCreated || apiErr.Body != "response bytes" || apiErr.RequestID != "request-123" {
			t.Fatalf("provider error = %#v (%v)", apiErr, err)
		}
	default:
		if err != nil || string(got) != "response bytes" {
			t.Fatalf("response = %q, error = %v", got, err)
		}
	}
}

// This test replaces process-global stderr, so it must run serially.
func TestBufferedRequestCloseError(t *testing.T) {
	for _, post := range []bool{true, false} {
		name := "get"
		if post {
			name = "post"
		}
		t.Run(name, func(t *testing.T) {
			read, write, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			oldStderr := os.Stderr
			os.Stderr = write
			t.Cleanup(func() { os.Stderr = oldStderr; _ = read.Close(); _ = write.Close() })
			body := &bufferedBody{reader: strings.NewReader("ok"), closeErr: errors.New("close failure")}
			calls := 0
			client := NewClient(ClientConfig{BaseURL: "https://base.example"}, DiscardLogger(), nil, bufferedDoer(func(*http.Request) (*http.Response, error) {
				calls++
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
			}))
			var got []byte
			if post {
				got, err = client.executeJSONRequest(context.Background(), "test", nil)
			} else {
				got, err = client.executeGetRequest(context.Background(), "test")
			}
			if closeErr := write.Close(); closeErr != nil {
				t.Fatal(closeErr)
			}
			output, readErr := io.ReadAll(read)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if err != nil || string(got) != "ok" || calls != 1 || body.closes != 1 || string(output) != "failed to close response body: close failure\n" {
				t.Fatalf("response=%q err=%v calls=%d closes=%d stderr=%q", got, err, calls, body.closes, output)
			}
		})
	}
}
