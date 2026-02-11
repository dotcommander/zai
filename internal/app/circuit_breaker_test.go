package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/dotcommander/zai/internal/config"
)

func TestCircuitBreakerIntegration(t *testing.T) {
	// Test circuit breaker configuration
	cbConfig := config.CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          1 * time.Second,
	}

	// Create a mock client with circuit breaker enabled
	cfg := ClientConfig{
		APIKey:         "test-key",
		BaseURL:        "https://api.test.com",
		Model:          "test-model",
		Timeout:        5 * time.Second,
		CircuitBreaker: cbConfig,
	}

	logger := DiscardLogger()
	history := &MockHistoryStore{}
	client := NewClient(cfg, logger, history, nil)

	// Test that circuit breakers are initialized
	cb := client.getCircuitBreaker("chat")
	if cb == nil {
		t.Fatal("Expected circuit breaker to be initialized")
	}

	// Test circuit breaker state transitions
	// First call should succeed (closed state)
	cb.Reset()

	// Success calls shouldn't open the circuit
	for i := 0; i < 2; i++ {
		err := cb.Execute(func() error {
			return nil // Success
		})
		if err != nil {
			t.Errorf("Expected success, got error: %v", err)
		}
	}

	// Circuit should still be closed after successes
	if cb.state != Closed {
		t.Errorf("Expected circuit to be closed, got %v", cb.state)
	}

	// Force failures to open the circuit
	for i := 0; i < 3; i++ {
		err := cb.Execute(func() error {
			return fmt.Errorf("API error")
		})
		if err == nil {
			t.Errorf("Expected failure, got success")
		}
	}

	// Circuit should now be open
	if cb.state != Open {
		t.Errorf("Expected circuit to be open, got %v", cb.state)
	}

	// Next call should fail immediately with circuit breaker error
	err := cb.Execute(func() error {
		return fmt.Errorf("API error")
	})
	if err == nil {
		t.Error("Expected circuit breaker error, got success")
	}
	if fmt.Sprintf("%v", err) != "circuit breaker 'chat' is open (timeout: 1s)" {
		t.Errorf("Expected specific circuit breaker error, got: %v", err)
	}

	// Wait for timeout
	time.Sleep(1100 * time.Millisecond)

	// Circuit should be in half-open state now
	err = cb.Execute(func() error {
		return nil // Success in half-open state
	})
	if err != nil {
		t.Errorf("Expected success in half-open state, got error: %v", err)
	}
	if cb.state != Closed {
		t.Errorf("Expected circuit to be reset to closed, got %v", cb.state)
	}
}
