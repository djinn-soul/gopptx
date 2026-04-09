package editor

import (
	"io"
	"testing"
)

func BenchmarkEditorSaveToWriterUnencrypted(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	b.ResetTimer()
	for b.Loop() {
		if err := editor.SaveToWriter(io.Discard); err != nil {
			b.Fatalf("SaveToWriter: %v", err)
		}
	}
}

func BenchmarkEditorSaveToWriterEncrypted(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()
	editor.metadata.Protection.EncryptPassword = "bench-password"
	if err := editor.SaveToWriter(io.Discard); err != nil {
		b.Skipf("encrypted SaveToWriter unavailable in this environment: %v", err)
	}

	b.ResetTimer()
	for b.Loop() {
		if err := editor.SaveToWriter(io.Discard); err != nil {
			b.Fatalf("SaveToWriter encrypted: %v", err)
		}
	}
}

func BenchmarkEditorSaveToBytesUnencrypted(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()

	b.ResetTimer()
	for b.Loop() {
		if _, err := editor.SaveToBytes(); err != nil {
			b.Fatalf("SaveToBytes: %v", err)
		}
	}
}

func BenchmarkEditorSaveToBytesEncrypted(b *testing.B) {
	editor := openBenchEditor(b)
	defer func() { _ = editor.Close() }()
	editor.metadata.Protection.EncryptPassword = "bench-password"
	if _, err := editor.SaveToBytes(); err != nil {
		b.Skipf("encrypted SaveToBytes unavailable in this environment: %v", err)
	}

	b.ResetTimer()
	for b.Loop() {
		if _, err := editor.SaveToBytes(); err != nil {
			b.Fatalf("SaveToBytes encrypted: %v", err)
		}
	}
}
