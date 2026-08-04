//nolint:mnd // Font table byte offsets and Calibri's fallback metrics are fixed values from the OpenType spec.
package export

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
	"sync"
)

// PowerPoint sizes its line box as a flat 1.2 x the point size, independent of
// the font. That was measured, not assumed: a probe deck rendered by PowerPoint
// itself puts 18pt Calibri, Segoe UI, Georgia, Verdana, Arial and Consolas on
// the same ~19.4pt pitch, even though those fonts' own hhea line heights span
// 1.14 to 1.33 em. See powerPointLineBoxFactor.
//
// Where the font does matter is the baseline inside that box. PowerPoint centres
// the font's ascent+descent in the line box (half-leading) and sets the baseline
// at the ascent. gopdf exposes glyph advances but not the hhea/OS-2 vertical
// metrics, so this file reads them straight out of the TTF and caches them per
// registered alias.
//
// baselineShiftFactor is expressed relative to gopdf's own placement: Cell()
// already drops the baseline by OS/2 typoAscender below the requested Y, so only
// the remainder is applied by the renderer.

// ttfLineMetrics holds the vertical metrics a renderer needs from one font file.
type ttfLineMetrics struct {
	UnitsPerEm   float64
	Ascender     float64
	Descender    float64 // negative, as stored in hhea
	LineGap      float64
	TypoAscender float64
}

// fallbackLineMetrics approximates Calibri and is used when a font file cannot
// be parsed, so layout degrades to a sane default instead of zero-height lines.
var fallbackLineMetrics = ttfLineMetrics{ //nolint:gochecknoglobals // Immutable default shared by all fallback lookups.
	UnitsPerEm:   2048,
	Ascender:     1536,
	Descender:    -512,
	LineGap:      452,
	TypoAscender: 1536,
}

// powerPointLineBoxFactor is the height of one line box at 100% line spacing,
// as a multiple of the point size. PowerPoint uses this same value for every
// font; paragraph-level lnSpc scales the pitch between lines but not this box,
// which is why the first baseline stays put when line spacing changes.
const powerPointLineBoxFactor = 1.2

// Sanity bound on the derived baseline. Fonts in the wild occasionally carry
// junk vertical metrics, and this renderer runs against whatever the host OS
// happens to ship. Clamping keeps a bad font looking slightly off instead of
// pushing its text clean off the slide.
const (
	minBaselineShiftFactor = -0.5
	maxBaselineShiftFactor = 1.0
)

// baselineShiftFactor is the extra multiple of the point size that must be added
// to gopdf's own baseline placement to match PowerPoint's. PowerPoint centres
// the font's ascent+descent inside the 1.2 em line box and puts the baseline at
// the ascent; gopdf has already applied typoAscender.
func (m ttfLineMetrics) baselineShiftFactor() float64 {
	if m.UnitsPerEm <= 0 {
		return fallbackLineMetrics.baselineShiftFactor()
	}
	ascent := m.Ascender / m.UnitsPerEm
	descent := -m.Descender / m.UnitsPerEm
	halfLeading := math.Max((powerPointLineBoxFactor-(ascent+descent))/2, 0)
	gopdfBaseline := m.TypoAscender / m.UnitsPerEm
	return clampFloat(halfLeading+ascent-gopdfBaseline, minBaselineShiftFactor, maxBaselineShiftFactor)
}

func clampFloat(v, lo, hi float64) float64 {
	return math.Min(math.Max(v, lo), hi)
}

//nolint:gochecknoglobals // Metrics are registered once per PDF font alias and read by every layout call.
var (
	pdfFontMetricsMu sync.RWMutex
	pdfFontMetrics   = map[string]ttfLineMetrics{}
)

// registerPDFFontMetrics parses fontPath and stores its metrics under alias.
// A parse failure is not fatal: the alias simply falls back to Calibri-like
// metrics, matching how font registration itself degrades.
func registerPDFFontMetrics(alias, fontPath string) {
	m, err := readTTFLineMetrics(fontPath)
	if err != nil {
		return
	}
	pdfFontMetricsMu.Lock()
	pdfFontMetrics[alias] = m
	pdfFontMetricsMu.Unlock()
}

// resetPDFFontMetrics drops every registered alias. Each export registers its
// own fonts, so stale metrics must not leak between documents.
func resetPDFFontMetrics() {
	pdfFontMetricsMu.Lock()
	pdfFontMetrics = map[string]ttfLineMetrics{}
	pdfFontMetricsMu.Unlock()
}

// lookupPDFFontMetrics returns the metrics registered for alias, or the
// Calibri-like fallback.
func lookupPDFFontMetrics(alias string) ttfLineMetrics {
	pdfFontMetricsMu.RLock()
	m, ok := pdfFontMetrics[alias]
	pdfFontMetricsMu.RUnlock()
	if !ok {
		return fallbackLineMetrics
	}
	return m
}

// metricsForFontHint resolves a PPTX typeface hint to the metrics of the font
// actually embedded for it.
func metricsForFontHint(fontHint string) ttfLineMetrics {
	return lookupPDFFontMetrics(resolvePDFFontAlias(fontHint))
}

// TTF/OTF table offsets used below, in bytes from the start of each table.
const (
	sfntHeaderSize      = 12
	sfntTableRecordSize = 16
	ttcFirstFontOffset  = 12

	headUnitsPerEmOffset  = 18
	hheaAscenderOffset    = 4
	hheaDescenderOffset   = 6
	hheaLineGapOffset     = 8
	os2TypoAscenderOffset = 68
)

var errFontTableMissing = errors.New("required font table missing")

// readTTFLineMetrics extracts the vertical metrics from a TrueType or OpenType
// file. TrueType Collections (.ttc) are read at their first font, which is the
// face gopdf embeds.
func readTTFLineMetrics(fontPath string) (ttfLineMetrics, error) {
	data, err := os.ReadFile(fontPath)
	if err != nil {
		return ttfLineMetrics{}, err
	}
	tables, err := parseSFNTTableDirectory(data)
	if err != nil {
		return ttfLineMetrics{}, err
	}

	head, okHead := tables["head"]
	hhea, okHhea := tables["hhea"]
	if !okHead || !okHhea {
		return ttfLineMetrics{}, errFontTableMissing
	}

	upem, err := readUint16At(data, head+headUnitsPerEmOffset)
	if err != nil {
		return ttfLineMetrics{}, err
	}
	ascender, err := readInt16At(data, hhea+hheaAscenderOffset)
	if err != nil {
		return ttfLineMetrics{}, err
	}
	descender, err := readInt16At(data, hhea+hheaDescenderOffset)
	if err != nil {
		return ttfLineMetrics{}, err
	}
	lineGap, err := readInt16At(data, hhea+hheaLineGapOffset)
	if err != nil {
		return ttfLineMetrics{}, err
	}
	if upem == 0 {
		return ttfLineMetrics{}, errFontTableMissing
	}

	// gopdf offsets its own baseline by OS/2 typoAscender. When OS/2 is absent
	// or truncated it uses zero, so mirror that rather than guessing.
	typoAscender := int16(0)
	if os2, ok := tables["OS/2"]; ok {
		if v, readErr := readInt16At(data, os2+os2TypoAscenderOffset); readErr == nil {
			typoAscender = v
		}
	}

	return ttfLineMetrics{
		UnitsPerEm:   float64(upem),
		Ascender:     float64(ascender),
		Descender:    float64(descender),
		LineGap:      float64(lineGap),
		TypoAscender: float64(typoAscender),
	}, nil
}

// parseSFNTTableDirectory maps table tags to their absolute byte offsets.
func parseSFNTTableDirectory(data []byte) (map[string]uint32, error) {
	if len(data) < sfntHeaderSize {
		return nil, errFontTableMissing
	}
	base := uint32(0)
	if string(data[0:4]) == "ttcf" {
		offset, err := readUint32At(data, ttcFirstFontOffset)
		if err != nil {
			return nil, err
		}
		base = offset
	}
	numTables, err := readUint16At(data, base+4)
	if err != nil {
		return nil, err
	}
	tables := make(map[string]uint32, numTables)
	for i := range uint32(numTables) {
		record := base + sfntHeaderSize + i*sfntTableRecordSize
		if uint64(record)+sfntTableRecordSize > uint64(len(data)) {
			return nil, errFontTableMissing
		}
		offset, err := readUint32At(data, record+8)
		if err != nil {
			return nil, err
		}
		tables[string(data[record:record+4])] = offset
	}
	return tables, nil
}

func readUint16At(data []byte, offset uint32) (uint16, error) {
	if uint64(offset)+2 > uint64(len(data)) {
		return 0, errFontTableMissing
	}
	return binary.BigEndian.Uint16(data[offset : offset+2]), nil
}

func readInt16At(data []byte, offset uint32) (int16, error) {
	v, err := readUint16At(data, offset)
	if err != nil {
		return 0, err
	}
	// hhea and OS/2 store these fields as signed 16-bit big-endian, and
	// descenders are genuinely negative. Reinterpreting the bits is the intended
	// conversion, not a lossy narrowing.
	return int16(v), nil //nolint:gosec // Deliberate two's-complement reinterpretation of a signed font value.
}

func readUint32At(data []byte, offset uint32) (uint32, error) {
	if uint64(offset)+4 > uint64(len(data)) {
		return 0, errFontTableMissing
	}
	return binary.BigEndian.Uint32(data[offset : offset+4]), nil
}
