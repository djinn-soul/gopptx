package common

import (
	"crypto/rand"
	"fmt"
)

// NewGUID generates a fresh GUID string in {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX} format.
func NewGUID() (string, error) {
	const (
		guidByteLen        = 16
		versionByteIndex   = 6
		variantByteIndex   = 8
		versionMask        = 0x0f
		version4Bits       = 0x40
		variantMask        = 0x3f
		rfc4122VariantBits = 0x80
	)
	b := make([]byte, guidByteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read random bytes for GUID: %w", err)
	}
	// GUID version 4 (random)
	b[versionByteIndex] = (b[versionByteIndex] & versionMask) | version4Bits
	b[variantByteIndex] = (b[variantByteIndex] & variantMask) | rfc4122VariantBits
	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%012X}",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
