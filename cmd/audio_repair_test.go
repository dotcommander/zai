package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetAudioFlags(t *testing.T) {
	t.Helper()
	file, video, stream, jsonOutput, preprocess, vad := audioFile, audioVideo, audioStream, audioJSON, audioPreprocess, audioVAD
	t.Cleanup(func() {
		audioFile, audioVideo, audioStream, audioJSON, audioPreprocess, audioVAD = file, video, stream, jsonOutput, preprocess, vad
	})
	audioFile, audioVideo, audioStream, audioJSON, audioPreprocess, audioVAD = "", "", false, false, false, false
}

func TestAudioLocalSourceOwnership(t *testing.T) {
	// Serial: audio flags and stdin are process globals.
	resetAudioFlags(t)
	path := filepath.Join(t.TempDir(), "input.wav")
	if err := os.WriteFile(path, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	audioFile = path
	for _, canceled := range []bool{false, true} {
		ctx, cancel := context.WithCancel(t.Context())
		if canceled {
			cancel()
		}
		source, err := determineAudioSource(ctx)
		cancel()
		if err != nil || source.Path != path || source.Owned {
			t.Fatalf("source=%+v error=%v", source, err)
		}
		mgr := &TempFileManager{}
		if source.Owned {
			mgr.Add(source.Path)
		}
		mgr.Cleanup()
		if data, err := os.ReadFile(path); err != nil || string(data) != "keep" {
			t.Fatalf("local source changed: %q, %v", data, err)
		}
	}
	// A preprocessing failure also executes the production deferred cleanup.
	audioPreprocess = true
	t.Setenv("PATH", t.TempDir())
	if err := runAudioTranscription(audioCmd, nil); err == nil {
		t.Fatal("expected missing ffmpeg")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestAudioOwnedStdinCleanup(t *testing.T) {
	// Serial: audio flags and stdin are process globals.
	resetAudioFlags(t)
	path := filepath.Join(t.TempDir(), "stdin.wav")
	if err := os.WriteFile(path, []byte("audio"), 0600); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	stdin := os.Stdin
	os.Stdin = f
	t.Cleanup(func() { os.Stdin = stdin })
	audioFile = "-"
	source, err := determineAudioSource(t.Context())
	if err != nil || !source.Owned {
		t.Fatalf("source=%+v err=%v", source, err)
	}
	mgr := &TempFileManager{}
	mgr.Add(source.Path)
	mgr.Cleanup()
	if _, err := os.Stat(source.Path); !os.IsNotExist(err) {
		t.Fatalf("temp survived cleanup: %v", err)
	}
}

func TestAudioStreamGuards(t *testing.T) {
	// Serial: command flags are process globals.
	resetAudioFlags(t)
	audioStream, audioJSON = true, true
	if err := runAudioTranscription(audioCmd, nil); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("error=%v", err)
	}
	audioJSON = false
	file, err := os.Create(filepath.Join(t.TempDir(), "large.wav"))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(25*1024*1024 + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	audioFile = file.Name()
	audioPreprocess = true
	if err := runAudioTranscription(audioCmd, nil); err == nil || !strings.Contains(err.Error(), "without --stream") {
		t.Fatalf("error=%v", err)
	}
	if _, err := os.Stat(audioFile); err != nil {
		t.Fatal(err)
	}
}
