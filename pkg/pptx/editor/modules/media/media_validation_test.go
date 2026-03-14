package media

import (
	"strings"
	"testing"
)

func TestResolveVideoExtension(t *testing.T) {
	t.Parallel()

	ext, err := ResolveVideoExtension("video/quicktime", "", true)
	if err != nil {
		t.Fatalf("ResolveVideoExtension failed: %v", err)
	}
	if ext != "mov" {
		t.Fatalf("expected mov extension, got %q", ext)
	}

	ext, err = ResolveVideoExtension("", "demo.webm", false)
	if err != nil {
		t.Fatalf("ResolveVideoExtension path fallback failed: %v", err)
	}
	if ext != "webm" {
		t.Fatalf("expected webm extension from path, got %q", ext)
	}

	ext, err = ResolveVideoExtension("video/x-ms-wmv", "", true)
	if err != nil {
		t.Fatalf("ResolveVideoExtension wmv mime failed: %v", err)
	}
	if ext != "wmv" {
		t.Fatalf("expected wmv extension, got %q", ext)
	}

	ext, err = ResolveVideoExtension("video/x-matroska", "", true)
	if err != nil {
		t.Fatalf("ResolveVideoExtension mkv mime failed: %v", err)
	}
	if ext != "mkv" {
		t.Fatalf("expected mkv extension, got %q", ext)
	}
}

func TestResolveVideoExtensionRejectsUnsupportedMime(t *testing.T) {
	t.Parallel()

	_, err := ResolveVideoExtension("video/unknown", "", true)
	if err == nil || !strings.Contains(err.Error(), "unsupported video mime type") {
		t.Fatalf("expected unsupported video mime error, got: %v", err)
	}
}

func TestResolveAudioExtensionRequiresMimeForInlineData(t *testing.T) {
	t.Parallel()

	_, err := ResolveAudioExtension("", "", true)
	if err == nil || !strings.Contains(err.Error(), "audio mime type is required") {
		t.Fatalf("expected audio mime required error, got: %v", err)
	}
}

func TestResolveAudioExtensionFromPath(t *testing.T) {
	t.Parallel()

	ext, err := ResolveAudioExtension("", "track.m4a", false)
	if err != nil {
		t.Fatalf("ResolveAudioExtension path fallback failed: %v", err)
	}
	if ext != "m4a" {
		t.Fatalf("expected m4a extension from path, got %q", ext)
	}

	ext, err = ResolveAudioExtension("audio/ogg", "", true)
	if err != nil {
		t.Fatalf("ResolveAudioExtension ogg mime failed: %v", err)
	}
	if ext != "ogg" {
		t.Fatalf("expected ogg extension, got %q", ext)
	}

	ext, err = ResolveAudioExtension("audio/flac", "", true)
	if err != nil {
		t.Fatalf("ResolveAudioExtension flac mime failed: %v", err)
	}
	if ext != "flac" {
		t.Fatalf("expected flac extension, got %q", ext)
	}
}
