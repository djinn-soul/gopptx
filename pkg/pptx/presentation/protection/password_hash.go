package protection

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"unicode/utf16"
)

// HashModifyPassword implements the SHA-512 hashing algorithm for p:modifyVerifier.
// It follows the algorithm:
// 1. Convert password to UTF-16LE bytes (truncated to 255 chars).
// 2. Initial hash = SHA-512(salt + password_bytes).
// 3. Iterative hash = SHA-512(prev_hash + uint32LE(iteration_index)).
func HashModifyPassword(password string, salt []byte, spinCount int) string {
	// 1. Password Encoding (UTF-16LE)
	runes := []rune(password)
	if len(runes) > 255 {
		runes = runes[:255]
	}
	u16 := utf16.Encode(runes)
	pwdBytes := make([]byte, len(u16)*2)
	for i, v := range u16 {
		binary.LittleEndian.PutUint16(pwdBytes[i*2:], v)
	}

	// 2. Initial Hashing
	initial := make([]byte, len(salt)+len(pwdBytes))
	copy(initial, salt)
	copy(initial[len(salt):], pwdBytes)
	hash := sha512.Sum512(initial)

	// 3. Iterative Hashing (Spin Count)
	var spinInput [sha512.Size + 4]byte
	for i := range spinCount {
		copy(spinInput[:sha512.Size], hash[:])
		binary.LittleEndian.PutUint32(spinInput[sha512.Size:], uint32(i))
		hash = sha512.Sum512(spinInput[:])
	}

	return base64.StdEncoding.EncodeToString(hash[:])
}
