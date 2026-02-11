package common

import (
	"archive/zip"
)

// WriteFile is a helper to write a string content to a zip file.
func WriteFile(zw *zip.Writer, path string, content string) error {
	w, err := zw.Create(path)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

// WriteBinaryFile is a helper to write binary content to a zip file.
func WriteBinaryFile(zw *zip.Writer, path string, data []byte) error {
	w, err := zw.Create(path)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
