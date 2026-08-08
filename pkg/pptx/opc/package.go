// Package opc exposes the OPC container a .pptx file is: a zip of parts, each
// addressed by path, plus the relationship files that bind them together.
//
// gopptx models decks through typed APIs, and everything they do goes through
// an internal package writer with no read, remove or list — so a caller who
// needed a part the typed API does not cover had nowhere to go. This package is
// that escape hatch, and nothing more: it does not interpret parts, it hands
// them over.
package opc

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

// Package is an open OPC container held in memory.
type Package struct {
	parts map[string][]byte
}

// ErrPartNotFound is returned when a named part is not in the package.
var ErrPartNotFound = errors.New("part not found")

// New returns an empty package. A package with no [Content_Types].xml is not a
// valid deck; callers building one from scratch must add the parts themselves.
func New() *Package {
	return &Package{parts: map[string][]byte{}}
}

// Open reads a .pptx (or any OPC container) from disk.
func Open(filePath string) (*Package, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read package %s: %w", filePath, err)
	}
	return OpenBytes(data)
}

// OpenBytes reads a package from memory.
func OpenBytes(data []byte) (*Package, error) {
	return OpenReader(bytes.NewReader(data), int64(len(data)))
}

// OpenReader reads a package from any random-access source.
func OpenReader(reader io.ReaderAt, size int64) (*Package, error) {
	zr, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}

	pkg := New()
	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}
		rc, openErr := file.Open()
		if openErr != nil {
			return nil, fmt.Errorf("open part %s: %w", file.Name, openErr)
		}
		content, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read part %s: %w", file.Name, readErr)
		}
		pkg.parts[normalizePartPath(file.Name)] = content
	}
	return pkg, nil
}

// Has reports whether the package contains a part.
func (p *Package) Has(partPath string) bool {
	_, ok := p.parts[normalizePartPath(partPath)]
	return ok
}

// Part returns a part's bytes. The returned slice is a copy, so a caller
// cannot mutate the package by writing through it.
func (p *Package) Part(partPath string) ([]byte, error) {
	content, ok := p.parts[normalizePartPath(partPath)]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPartNotFound, partPath)
	}
	return append([]byte(nil), content...), nil
}

// PartString returns a part's bytes as a string, for the XML parts.
func (p *Package) PartString(partPath string) (string, error) {
	content, err := p.Part(partPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SetPart adds or replaces a part.
func (p *Package) SetPart(partPath string, content []byte) {
	p.parts[normalizePartPath(partPath)] = append([]byte(nil), content...)
}

// SetPartString adds or replaces a part from a string.
func (p *Package) SetPartString(partPath, content string) {
	p.SetPart(partPath, []byte(content))
}

// RemovePart deletes a part, reporting whether it was there. Removing a part
// does not remove the relationships that point at it — see Relationships.
func (p *Package) RemovePart(partPath string) bool {
	key := normalizePartPath(partPath)
	if _, ok := p.parts[key]; !ok {
		return false
	}
	delete(p.parts, key)
	return true
}

// PartPaths lists every part, sorted, so output is deterministic.
func (p *Package) PartPaths() []string {
	out := make([]string, 0, len(p.parts))
	for partPath := range p.parts {
		out = append(out, partPath)
	}
	sort.Strings(out)
	return out
}

// PartPathsWithPrefix lists the parts under a directory, sorted.
func (p *Package) PartPathsWithPrefix(prefix string) []string {
	normalized := normalizePartPath(prefix)
	var out []string
	for _, partPath := range p.PartPaths() {
		if strings.HasPrefix(partPath, normalized) {
			out = append(out, partPath)
		}
	}
	return out
}

// PartCount is how many parts the package holds.
func (p *Package) PartCount() int {
	return len(p.parts)
}

// Save writes the package to disk.
func (p *Package) Save(filePath string) error {
	var buf bytes.Buffer
	if err := p.SaveTo(&buf); err != nil {
		return err
	}
	if err := os.WriteFile(filePath, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("write package %s: %w", filePath, err)
	}
	return nil
}

// SaveTo writes the package to a writer. Parts are written in sorted order so
// the same package always produces the same bytes.
func (p *Package) SaveTo(w io.Writer) error {
	zw := zip.NewWriter(w)
	for _, partPath := range p.PartPaths() {
		writer, err := zw.Create(partPath)
		if err != nil {
			return fmt.Errorf("create part %s: %w", partPath, err)
		}
		if _, err = writer.Write(p.parts[partPath]); err != nil {
			return fmt.Errorf("write part %s: %w", partPath, err)
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("finish package: %w", err)
	}
	return nil
}

// Bytes renders the package to a byte slice.
func (p *Package) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := p.SaveTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// RelationshipsPath is the .rels part that belongs to a part: the package
// relationships for "", and _rels/<name>.rels beside anything else.
func RelationshipsPath(partPath string) string {
	normalized := normalizePartPath(partPath)
	if normalized == "" {
		return "_rels/.rels"
	}
	return path.Join(path.Dir(normalized), "_rels", path.Base(normalized)+".rels")
}

// normalizePartPath puts a path in the form the zip uses: forward slashes, no
// leading one.
func normalizePartPath(partPath string) string {
	cleaned := strings.ReplaceAll(strings.TrimSpace(partPath), "\\", "/")
	return strings.TrimPrefix(cleaned, "/")
}
