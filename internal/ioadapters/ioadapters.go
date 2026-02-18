package ioadapters

import (
	"errors"
	"io"
	"slices"
)

type readerAtAdapter struct {
	r         io.Reader
	readBytes []byte
}

// ToReaderAt returns an [io.ReaderAt] from an [io.Reader].
// If the reader already implements [io.ReaderAt], it is returned directly.
// Otherwise, it returns an adapter that buffers the read data.
func ToReaderAt(r io.Reader) io.ReaderAt {
	ra, ok := r.(io.ReaderAt)
	if ok {
		return ra
	}
	return &readerAtAdapter{
		r: r,
	}
}

func (r *readerAtAdapter) ReadAt(p []byte, off int64) (int, error) {
	if int(off)+len(p) > len(r.readBytes) {
		err := r.expandBuffer(int(off) + len(p))
		if err != nil {
			return 0, err
		}
	}
	return bytesReaderAt(r.readBytes).ReadAt(p, off)
}

func (r *readerAtAdapter) expandBuffer(newSize int) error {
	if cap(r.readBytes) < newSize {
		r.readBytes = slices.Grow(r.readBytes, newSize-cap(r.readBytes))
	}

	newPart := r.readBytes[len(r.readBytes):newSize]
	n, err := r.r.Read(newPart)
	switch {
	case err == nil:
		r.readBytes = r.readBytes[:newSize]
	case errors.Is(err, io.EOF):
		r.readBytes = r.readBytes[:len(r.readBytes)+n]
	default:
		return err
	}
	return nil
}

// BytesReadAt returns an [io.ReaderAt] view over a byte slice.
func BytesReadAt(src []byte, dst []byte, off int64) (int, error) {
	return bytesReaderAt(src).ReadAt(dst, off)
}

type bytesReaderAt []byte

func (bra bytesReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, errors.New("ioadapters: negative offset")
	}
	if off >= int64(len(bra)) {
		return 0, io.EOF
	}
	n := copy(p, bra[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}
