package tplx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func readZipParts(zr *zip.Reader) (map[string][]byte, error) {
	parts := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.Contains(f.Name, "..") || strings.HasPrefix(f.Name, "/") {
			return nil, fmt.Errorf("tplx: unsafe zip entry path %q", f.Name)
		}
		partData, err := readZipFileBytes(f)
		if err != nil {
			return nil, fmt.Errorf("tplx: open %s: %w", f.Name, err)
		}
		parts[f.Name] = partData
	}
	return parts, nil
}

func readZipFileBytes(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err //nolint:wrapcheck // Caller wraps with part-name context.
	}
	defer func() {
		_ = rc.Close()
	}()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(rc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// repackZip rebuilds the ZIP archive, replacing parts with modified data.
func repackZip(original *zip.Reader, parts map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	written := make(map[string]bool)
	for _, f := range original.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if _, ok := parts[f.Name]; !ok {
			continue
		}
		if err := writeArchiveEntry(zw, f, parts); err != nil {
			return nil, err
		}
		written[f.Name] = true
	}

	if err := writeNewArchiveParts(zw, parts, written); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err //nolint:wrapcheck // Preserve ZIP finalization errors.
	}
	return buf.Bytes(), nil
}

func writeArchiveEntry(zw *zip.Writer, f *zip.File, parts map[string][]byte) error {
	hdr := f.FileHeader
	w, err := zw.CreateHeader(&hdr)
	if err != nil {
		return err //nolint:wrapcheck // ZIP writer header failures are surfaced directly.
	}
	if data, ok := parts[f.Name]; ok {
		return writeArchiveBytes(w, data)
	}
	data, err := readZipFileBytes(f)
	if err != nil {
		return err //nolint:wrapcheck // Preserve ZIP entry open failures for passthrough files.
	}
	return writeArchiveBytes(w, data)
}

func writeArchiveBytes(w io.Writer, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return err //nolint:wrapcheck // Preserve write failures for generated and passthrough payloads.
	}
	return nil
}

func writeNewArchiveParts(zw *zip.Writer, parts map[string][]byte, written map[string]bool) error {
	for name, data := range parts {
		if written[name] {
			continue
		}
		w, err := zw.Create(name)
		if err != nil {
			return err //nolint:wrapcheck // ZIP writer create failures are surfaced directly.
		}
		if _, err = w.Write(data); err != nil {
			return err //nolint:wrapcheck // Preserve write failures for newly added parts.
		}
	}
	return nil
}
