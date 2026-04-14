package editor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
)

func readZipEntry(entry *zip.File) ([]byte, error) {
	rc, err := entry.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip entry %q: %w", entry.Name, err)
	}
	defer func() { _ = rc.Close() }()

	const maxInt = int(^uint(0) >> 1)
	size64 := entry.UncompressedSize64
	if size64 > 0 && size64 <= uint64(maxInt) {
		// Enforce the same upper bound as the unknown-size path to prevent a
		// crafted PPTX with a large declared size from triggering a huge allocation.
		if size64 > uint64(maxUnknownZipEntryBytes) {
			return nil, fmt.Errorf("zip entry %q declared size %d exceeds max %d bytes", entry.Name, size64, maxUnknownZipEntryBytes)
		}
		// Pre-allocate the exact size when available.
		data := make([]byte, int(size64))
		if _, err := io.ReadFull(rc, data); err != nil {
			return nil, fmt.Errorf("read zip entry %q: %w", entry.Name, err)
		}
		return data, nil
	}

	// Fallback for unknown sizes: use a pooled read buffer to reduce temporary allocations.
	chunkPtr, ok := zipReadChunkPool.Get().(*[]byte)
	if !ok || chunkPtr == nil {
		return nil, errors.New("zip read pool returned invalid buffer")
	}
	defer zipReadChunkPool.Put(chunkPtr)
	chunk := *chunkPtr
	var out bytes.Buffer
	limitedReader := io.LimitReader(rc, maxUnknownZipEntryBytes+1)
	if _, err := io.CopyBuffer(&out, limitedReader, chunk); err != nil {
		return nil, fmt.Errorf("read zip entry %q: %w", entry.Name, err)
	}
	if out.Len() > maxUnknownZipEntryBytes {
		return nil, fmt.Errorf("zip entry %q exceeds max size %d bytes", entry.Name, maxUnknownZipEntryBytes)
	}
	return out.Bytes(), nil
}
