package vba

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/richardlehane/mscfb"
)

// CFBEntry describes one directory/stream entry in a Compound File Binary blob.
type CFBEntry struct {
	Name string
	Path string
	Size int64
}

// CFBInfo contains parsed metadata for a Compound File Binary blob.
type CFBInfo struct {
	ID       string
	Created  time.Time
	Modified time.Time
	Entries  []CFBEntry
}

// InspectCFB parses a CFB blob using mscfb and returns enumerated entries.
func InspectCFB(data []byte) (*CFBInfo, error) {
	if len(data) == 0 {
		return nil, errors.New("empty CFB data")
	}

	reader, err := mscfb.New(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse CFB: %w", err)
	}

	info := &CFBInfo{
		ID:       reader.ID(),
		Created:  reader.Created(),
		Modified: reader.Modified(),
	}
	for {
		entry, nextErr := reader.Next()
		if nextErr == nil {
			info.Entries = append(info.Entries, CFBEntry{
				Name: entry.Name,
				Path: strings.Join(entry.Path, "/"),
				Size: entry.Size,
			})
			continue
		}
		if errors.Is(nextErr, io.EOF) {
			break
		}
		return nil, fmt.Errorf("walk CFB entries: %w", nextErr)
	}
	return info, nil
}

// InspectCFB parses the project's binary blob as CFB and returns metadata.
func (p *VBAProject) InspectCFB() (*CFBInfo, error) {
	if p == nil {
		return nil, errors.New("nil VBA project")
	}
	return InspectCFB(p.Data)
}
