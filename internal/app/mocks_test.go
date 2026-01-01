package app

import (
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/testify/mock"
)

// MockHTTPDoer is a mock implementation of HTTPDoer for testing.
type MockHTTPDoer struct {
	mock.Mock
}

func (m *MockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockFileReader is a mock implementation of FileReader for testing.
type MockFileReader struct {
	mock.Mock
}

func (m *MockFileReader) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// MockLogger is a mock implementation of Logger for testing.
type MockLogger struct {
	mock.Mock
	logs []string
}

func (m *MockLogger) Info(format string, args ...any) {
	msg := format
	for _, arg := range args {
		msg += " " + stringify(arg)
	}
	m.logs = append(m.logs, "INFO: "+msg)
	m.Called(format, args)
}

func (m *MockLogger) Warn(format string, args ...any) {
	msg := format
	for _, arg := range args {
		msg += " " + stringify(arg)
	}
	m.logs = append(m.logs, "WARN: "+msg)
	m.Called(format, args)
}

func (m *MockLogger) Error(format string, args ...any) {
	msg := format
	for _, arg := range args {
		msg += " " + stringify(arg)
	}
	m.logs = append(m.logs, "ERROR: "+msg)
	m.Called(format, args)
}

func (m *MockLogger) GetLogs() []string {
	return m.logs
}

func (m *MockLogger) ClearLogs() {
	m.logs = nil
}

// MockHistoryStore is a mock implementation of HistoryStore for testing.
type MockHistoryStore struct {
	mock.Mock
}

func (m *MockHistoryStore) Save(entry HistoryEntry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockHistoryStore) GetRecent(limit int) ([]HistoryEntry, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]HistoryEntry), args.Error(1)
}

func (m *MockHistoryStore) Path() string {
	args := m.Called()
	return args.String(0)
}

// MockSearchCache is a mock implementation of SearchCache for testing.
type MockSearchCache struct {
	mock.Mock
}

func (m *MockSearchCache) Get(query string, opts SearchOptions) ([]SearchResult, bool) {
	args := m.Called(query, opts)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]SearchResult), args.Bool(1)
}

func (m *MockSearchCache) Set(query string, opts SearchOptions, results []SearchResult, ttl int) error {
	args := m.Called(query, opts, results, ttl)
	return args.Error(0)
}

func (m *MockSearchCache) Clear() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSearchCache) Cleanup() error {
	args := m.Called()
	return args.Error(0)
}

// Helper function to stringify any value for logging
func stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return "..."
}

// NewMockHTTPResponse creates a mock HTTP response for testing.
func NewMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
