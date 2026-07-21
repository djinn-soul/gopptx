// Package zipfast provides zip writers tuned for OPC packages.
//
// The stdlib archive/zip defaults to [flate.DefaultCompression] (level 6).
// Level 6 uses the full flate compressor, whose Reset zeroes a ~192KB hash
// table on every zip entry. An OPC package has hundreds of entries, many of
// them tiny (_rels files, small XML parts), so that per-entry reset dominates
// the actual compression work: CPU profiles of the package writer attribute
// ~62% of samples to flate.(*compressor).reset alone.
//
// Levels 1 and 2 use flate's deflatefast path, which keeps far less state and
// resets cheaply. Measured on a 400-entry package, level 1 writes ~8x faster
// than level 6 while producing output ~22% larger. Levels 2 and above showed
// no meaningful gain over the default.
package zipfast

import (
	"archive/zip"
	"compress/flate"
	"io"
	"sync"
)

//nolint:gochecknoglobals // Pooling flate writers is the entire point of this package.
var flateWriterPool sync.Pool

// pooledFlateWriter returns its flate.Writer to the pool on Close.
type pooledFlateWriter struct {
	writer *flate.Writer
	closed bool
}

func (p *pooledFlateWriter) Write(b []byte) (int, error) {
	return p.writer.Write(b)
}

// Close flushes the entry and releases the underlying writer back to the pool.
// It is idempotent: a second call is a no-op, so a double Close cannot hand the
// same flate.Writer to two concurrent entries.
func (p *pooledFlateWriter) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	err := p.writer.Close()
	// Drop the reference to the surrounding zip entry writer so the pooled
	// object does not pin it until the next Reset.
	p.writer.Reset(io.Discard)
	flateWriterPool.Put(p.writer)
	p.writer = nil
	return err
}

// newCompressor builds the zip.Compressor installed by [Register].
func newCompressor(out io.Writer) (io.WriteCloser, error) {
	if pooled, ok := flateWriterPool.Get().(*flate.Writer); ok && pooled != nil {
		pooled.Reset(out)
		return &pooledFlateWriter{writer: pooled}, nil
	}
	fresh, err := flate.NewWriter(out, flate.BestSpeed)
	if err != nil {
		return nil, err
	}
	return &pooledFlateWriter{writer: fresh}, nil
}

// Register installs the pooled [flate.BestSpeed] compressor on zw, replacing
// the stdlib default for [zip.Deflate] entries. Entries written with
// [zip.Store], and entries copied verbatim via CreateRaw, are unaffected.
func Register(zw *zip.Writer) {
	zw.RegisterCompressor(zip.Deflate, newCompressor)
}

// NewWriter returns a [zip.Writer] wrapping w with the pooled BestSpeed
// compressor already registered. It is a drop-in replacement for
// [zip.NewWriter] at OPC package write sites.
func NewWriter(w io.Writer) *zip.Writer {
	zw := zip.NewWriter(w)
	Register(zw)
	return zw
}
