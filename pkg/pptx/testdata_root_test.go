package pptx

import "path/filepath"

func rootTestdataPath(parts ...string) string {
	base := []string{"..", "..", "testdata"}
	base = append(base, parts...)
	return filepath.Join(base...)
}
