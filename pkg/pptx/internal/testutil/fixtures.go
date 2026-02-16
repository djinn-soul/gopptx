package testutil

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// RootTestdataDir searches upward for the testdata directory.
func RootTestdataDir() string {
	base := "../../testdata"
	for range 5 {
		if _, err := os.Stat(base); err == nil {
			return base
		}
		base = "../" + base
	}
	return base
}

// FixtureZipReader reads a fixture file from testdata/ppt_rs and returns a [zip.Reader].
func FixtureZipReader(t *testing.T, fixtureName string) *zip.Reader {
	t.Helper()
	fixturePath := filepath.Join(RootTestdataDir(), "ppt_rs", fixtureName)
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixturePath, err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read fixture %s: %v", fixturePath, err)
	}
	return zr
}

// ReadAllSlidesXML reads and concatenates all slide XMLs from a [zip.Reader].
func ReadAllSlidesXML(t *testing.T, zr *zip.Reader) string {
	t.Helper()
	names := make([]string, 0)
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, "ppt/slides/slide") || !strings.HasSuffix(f.Name, ".xml") {
			continue
		}
		names = append(names, f.Name)
	}
	sort.Strings(names)

	var b strings.Builder
	for _, name := range names {
		b.WriteString(ReadZipFile(t, zr, name))
	}
	return b.String()
}

// AssertContainsTokens checks that xml contains all specified token substrings.
func AssertContainsTokens(t *testing.T, label string, xml string, tokens []string) {
	t.Helper()
	for _, token := range tokens {
		if !strings.Contains(xml, token) {
			t.Fatalf("%s missing token %q", label, token)
		}
	}
}
