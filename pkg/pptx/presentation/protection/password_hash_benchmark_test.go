package protection

import (
	"encoding/base64"
	"testing"
)

func BenchmarkHashModifyPassword(b *testing.B) {
	salt, err := base64.StdEncoding.DecodeString("uS6jM5k6pQ8=")
	if err != nil {
		b.Fatalf("decode salt: %v", err)
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = HashModifyPassword("test-password", salt, 100000)
	}
}
