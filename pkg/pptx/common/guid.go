package common

import (
	"crypto/rand"
	"fmt"
)

// NewGUID generates a fresh GUID string in {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX} format.
func NewGUID() (string, error) {
	const guidByteLen = 16
	b := make([]byte, guidByteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes for GUID: %w", err)
	}
	// GUID version 4 (random)
	//nolint:mnd // GUID bit manipulation
	b[6] = (b[6] & 0x0f) | 0x40
	//nolint:mnd // GUID bit manipulation
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%012X}",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
