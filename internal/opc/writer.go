package opc

import (
	"archive/zip"
	"io"

	"github.com/djinn-soul/gopptx/internal/zipfast"
)

// Writer handles the creation of the PPTX (ZIP) package.
type Writer struct {
	zipWriter *zip.Writer
}

// NewWriter creates a new OPC writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		zipWriter: zipfast.NewWriter(w),
	}
}

// AddFile adds a file to the package.
func (w *Writer) AddFile(name string, content []byte) error {
	f, err := w.zipWriter.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}

// Close closes the underlying zip writer.
func (w *Writer) Close() error {
	return w.zipWriter.Close()
}
