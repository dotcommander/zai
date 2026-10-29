//go:build windows

package app

import (
	"context"
	"errors"
	"testing"
)

func TestOpenWithCanceledContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := OpenWith(ctx, "https://example.com"); !errors.Is(err, context.Canceled) {
		t.Fatalf("OpenWith canceled context: got %v, want %v", err, context.Canceled)
	}
}

func TestOpenWithRejectsEmbeddedNUL(t *testing.T) {
	t.Parallel()
	if err := OpenWith(t.Context(), "https://example.com/\x00invalid"); err == nil {
		t.Fatal("OpenWith accepted a target containing NUL")
	}
}
