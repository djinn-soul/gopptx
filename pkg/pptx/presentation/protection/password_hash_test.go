package protection

import (
	"encoding/base64"
	"testing"
)

func TestHashModifyPassword(t *testing.T) {
	// Test case based on a known-good Office hash.
	// Note: Finding a public test vector for this specific algorithm is hard,
	// so we will test the internal logic consistency and then verify with actual PowerPoint if possible.
	password := "test"
	saltBase64 := "uS6jM5k6pQ8="
	salt, _ := base64.StdEncoding.DecodeString(saltBase64)
	spinCount := 100000

	hash := HashModifyPassword(password, salt, spinCount)
	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	// Verify stability
	hash2 := HashModifyPassword(password, salt, spinCount)
	if hash != hash2 {
		t.Errorf("Hash changed between runs: %s != %s", hash, hash2)
	}

	// Verify password sensitivity
	hash3 := HashModifyPassword("wrong", salt, spinCount)
	if hash == hash3 {
		t.Error("Hash collision for different passwords")
	}
}
