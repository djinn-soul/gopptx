package protection

import (
	"bytes"
	"testing"
)

func TestCFBWriter_Internal(t *testing.T) {
	// ensureMinCFBStreamSize
	small := []byte{1, 2, 3}
	padded := ensureMinCFBStreamSize(small)
	if len(padded) != cfbMinRegularBytes {
		t.Errorf("Expected padding to %d, got %d", cfbMinRegularBytes, len(padded))
	}
	large := bytes.Repeat([]byte{1}, cfbMinRegularBytes+10)
	notPadded := ensureMinCFBStreamSize(large)
	if len(notPadded) != len(large) {
		t.Error("Should not pad if already large enough")
	}

	// buildCompoundFile
	info := []byte("info")
	pkg := []byte("pkg")
	cfb, err := buildCompoundFile(info, pkg, "Info", "Pkg")
	if err != nil {
		t.Fatalf("buildCompoundFile failed: %v", err)
	}
	if !bytes.HasPrefix(cfb, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
		t.Error("Missing CFB header")
	}

	// newDirEntry errors
	_, err = newDirEntry("", 1, 0, 0, 0, 0, 0)
	if err == nil {
		t.Error("Expected error for empty name")
	}
	_, err = newDirEntry("very long name that exceeds the sixty four byte limit for compound file binary format directory entries", 1, 0, 0, 0, 0, 0)
	if err == nil {
		t.Error("Expected error for long name")
	}

	// sectorsNeeded
	if sectorsNeeded(0) != 1 {
		t.Error("sectorsNeeded(0) should be 1")
	}
	if sectorsNeeded(512) != 1 {
		t.Error("sectorsNeeded(512) should be 1")
	}
	if sectorsNeeded(513) != 2 {
		t.Error("sectorsNeeded(513) should be 2")
	}
}

func TestCFBWriter_LargeFileError(t *testing.T) {
	// Trigger fatSectors > 109
	// Each FAT sector holds 128 entries. 110 FAT sectors * 128 = 14080 sectors.
	// 14080 * 512 bytes = ~7.2MB
	// Let's use a huge size to trigger the limit.
	huge := make([]byte, 8*1024*1024)
	_, err := buildCompoundFile(huge, huge, "I", "P")
	if err == nil {
		t.Error("Expected error for too large file")
	}
}
