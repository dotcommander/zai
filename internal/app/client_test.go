package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClientChat tests the Chat method with mocked HTTP responses.
func TestClientChat(t *testing.T) {
	tests := []struct {
		name           string
		prompt         string
		mockResponse   ChatResponse
		mockStatusCode int
		expectError    bool
		expectedOutput string
	}{
		{
			name:   "successful chat response",
			prompt: "Hello",
			mockResponse: ChatResponse{
				ID:      "chat-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "glm-4.7",
				Choices: []Choice{
					{
						Message: Message{
							Role:    "assistant",
							Content: "Hi there! How can I help you today?",
						},
						FinishReason: "stop",
					},
				},
				Usage: Usage{
					PromptTokens:     10,
					CompletionTokens: 20,
					TotalTokens:      30,
				},
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
			expectedOutput: "Hi there! How can I help you today?",
		},
		{
			name:   "API error response",
			prompt: "Hello",
			mockResponse: ChatResponse{
				ID:      "error",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "glm-4.7",
				Choices: []Choice{},
			},
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:   "empty choices",
			prompt: "Hello",
			mockResponse: ChatResponse{
				ID:      "chat-123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "glm-4.7",
				Choices: []Choice{},
				Usage:   Usage{},
			},
			mockStatusCode: http.StatusOK,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/chat/completions", r.URL.Path)
				assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

				// Set status code
				w.WriteHeader(tt.mockStatusCode)

				// Write response
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create client with mock server URL
			config := ClientConfig{
				APIKey:  "test-api-key",
				BaseURL: server.URL,
				Model:   "glm-4.7",
				Timeout: 30 * time.Second,
				Verbose: false,
				RetryConfig: RetryConfig{
					MaxAttempts:    1,
					InitialBackoff: 1 * time.Second,
					MaxBackoff:     30 * time.Second,
				},
			}

			logger := DiscardLogger()
			client := NewClient(config, logger, nil, nil)

			// Execute chat
			ctx := context.Background()
			opts := DefaultChatOptions()

			response, err := client.Chat(ctx, tt.prompt, opts)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, response)
			}
		})
	}
}

// TestClientListModels tests the ListModels method.
func TestClientListModels(t *testing.T) {
	mockModels := []Model{
		{ID: "glm-4.7", Object: "model", Created: time.Now().Unix(), OwnedBy: "zai"},
		{ID: "glm-4.6", Object: "model", Created: time.Now().Unix(), OwnedBy: "zai"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/models", r.URL.Path)
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		response := ModelsResponse{
			Object: "list",
			Data:   mockModels,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
		Verbose: false,
		RetryConfig: RetryConfig{
			MaxAttempts:    1,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
		},
	}

	logger := DiscardLogger()
	client := NewClient(config, logger, nil, nil)

	ctx := context.Background()
	models, err := client.ListModels(ctx)

	require.NoError(t, err)
	assert.Len(t, models, 2)
	assert.Equal(t, "glm-4.7", models[0].ID)
	assert.Equal(t, "glm-4.6", models[1].ID)
}

// TestClientRetryLogic tests the retry logic with transient failures.
func TestClientRetryLogic(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 2 {
			// First attempt fails with 503
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Second attempt succeeds
		response := ChatResponse{
			ID:      "chat-123",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "glm-4.7",
			Choices: []Choice{
				{
					Message:      Message{Role: "assistant", Content: "Success after retry"},
					FinishReason: "stop",
				},
			},
			Usage: Usage{TotalTokens: 10},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Model:   "glm-4.7",
		Timeout: 30 * time.Second,
		Verbose: true,
		RetryConfig: RetryConfig{
			MaxAttempts:    3,
			InitialBackoff: 10 * time.Millisecond,
			MaxBackoff:     100 * time.Millisecond,
		},
	}

	logger := DiscardLogger()
	client := NewClient(config, logger, nil, nil)

	ctx := context.Background()
	opts := DefaultChatOptions()

	response, err := client.Chat(ctx, "test", opts)

	require.NoError(t, err)
	assert.Equal(t, "Success after retry", response)
	assert.Equal(t, 2, attemptCount)
}

// TestClientContextCancellation tests that context cancellation is respected.
func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay response
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ChatResponse{
			ID:      "chat-123",
			Choices: []Choice{{Message: Message{Content: "Response"}}},
		})
	}))
	defer server.Close()

	config := ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Model:   "glm-4.7",
		Timeout: 30 * time.Second,
		Verbose: false,
		RetryConfig: RetryConfig{
			MaxAttempts:    1,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
		},
	}

	logger := DiscardLogger()
	client := NewClient(config, logger, nil, nil)

	// Create context that cancels quickly
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	opts := DefaultChatOptions()
	_, err := client.Chat(ctx, "test", opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

// TestClientWithFileContent tests chat with file content included.
func TestClientWithFileContent(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	content := "This is test file content"
	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	tmpfile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read request body to verify file content is included
		var reqData ChatRequest
		json.NewDecoder(r.Body).Decode(&reqData)

		// The file content should be in the user message
		assert.Contains(t, reqData.Messages[len(reqData.Messages)-1].Content, content)

		response := ChatResponse{
			ID:      "chat-123",
			Choices: []Choice{{Message: Message{Content: "Response"}}},
			Usage:   Usage{TotalTokens: 100},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := ClientConfig{
		APIKey:  "test-api-key",
		BaseURL: server.URL,
		Model:   "glm-4.7",
		Timeout: 30 * time.Second,
		Verbose: false,
		RetryConfig: RetryConfig{
			MaxAttempts:    1,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
		},
	}

	logger := DiscardLogger()
	client := NewClient(config, logger, nil, nil)

	ctx := context.Background()
	opts := DefaultChatOptions()
	opts.FilePath = tmpfile.Name()

	response, err := client.Chat(ctx, "test", opts)

	require.NoError(t, err)
	assert.NotEmpty(t, response)
}

// TestIsRetryableError tests the isRetryableError function.
func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"timeout error", &testTimeoutError{true}, true},
		{"connection refused", fmt.Errorf("connection refused"), true},
		{"503 error", fmt.Errorf("API error: 503"), true},
		{"502 error", fmt.Errorf("API error: 502"), true},
		{"504 error", fmt.Errorf("API error: 504"), true},
		{"400 error", fmt.Errorf("API error: 400"), false},
		{"500 error", fmt.Errorf("API error: 500"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// testTimeoutError is a helper for testing timeout errors
type testTimeoutError struct {
	timeout bool
}

func (e *testTimeoutError) Error() string { return "timeout" }
func (e *testTimeoutError) Timeout() bool { return e.timeout }

// TestCalculateBackoff tests the exponential backoff calculation.
func TestCalculateBackoff(t *testing.T) {
	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	// Test that backoff increases with attempts
	backoff1 := calculateBackoff(1, initialBackoff, maxBackoff)
	backoff2 := calculateBackoff(2, initialBackoff, maxBackoff)
	backoff3 := calculateBackoff(3, initialBackoff, maxBackoff)

	assert.Greater(t, backoff2, backoff1/2) // Should generally increase
	assert.Greater(t, backoff3, backoff2/2)

	// Test that backoff is capped at maxBackoff
	backoff20 := calculateBackoff(20, initialBackoff, maxBackoff)
	assert.LessOrEqual(t, backoff20, maxBackoff+5*time.Second) // Allow some jitter above max
	assert.Greater(t, backoff20, maxBackoff-5*time.Second)     // But should be near max

	// Test that small backoffs work correctly
	smallBackoff := calculateBackoff(1, 100*time.Millisecond, 1*time.Second)
	assert.Greater(t, smallBackoff, 50*time.Millisecond)
	assert.Less(t, smallBackoff, 200*time.Millisecond)
}
