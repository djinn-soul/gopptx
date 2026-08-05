//nolint:mnd // Font table byte offsets and Calibri's fallback metrics are fixed values from the OpenType spec.
package export

import (
	"encoding/binary"
	"errors"
	"math"
	"os"

	"github.com/signintech/gopdf"
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

	// Decoration and script metrics. Underline comes from the post table,
	// the rest from OS/2. All are in font units; zero means the font did not
	// state the value and the em-fraction accessors below substitute a default.
	UnderlinePosition  float64 // negative: below the baseline
	UnderlineThickness float64
	StrikeoutPosition  float64 // positive: above the baseline
	StrikeoutSize      float64
	SubscriptSize      float64
	SubscriptOffset    float64 // positive: below the baseline
	SuperscriptSize    float64
	SuperscriptOffset  float64 // positive: above the baseline
}

// fallbackLineMetrics approximates Calibri and is used when a font file cannot
// be parsed, so layout degrades to a sane default instead of zero-height lines.
var fallbackLineMetrics = ttfLineMetrics{ //nolint:gochecknoglobals // Immutable default shared by all fallback lookups.
	UnitsPerEm:   2048,
	Ascender:     1536,
	Descender:    -512,
	LineGap:      452,
	TypoAscender: 1536,

	UnderlinePosition:  -205, // -0.1 em
	UnderlineThickness: 102,  // 0.05 em
	StrikeoutPosition:  512,  // 0.25 em
	StrikeoutSize:      102,  // 0.05 em
	SubscriptSize:      1331, // 0.65 em
	SubscriptOffset:    287,  // 0.14 em
	SuperscriptSize:    1331, // 0.65 em
	SuperscriptOffset:  819,  // 0.4 em
}

// emFraction converts a value in font units to a fraction of the em, falling
// back to defaultEm when the font left the field at zero.
func (m ttfLineMetrics) emFraction(value, defaultEm float64) float64 {
	if m.UnitsPerEm <= 0 || value == 0 {
		return defaultEm
	}
	return value / m.UnitsPerEm
}

// underlineOffsetFactor is how far below the baseline the underline sits, as a
// fraction of the point size.
func (m ttfLineMetrics) underlineOffsetFactor() float64 {
	return -m.emFraction(m.UnderlinePosition, -0.1)
}

func (m ttfLineMetrics) underlineThicknessFactor() float64 {
	return m.emFraction(m.UnderlineThickness, 0.05)
}

// strikeoutOffsetFactor is how far above the baseline the strike line sits.
func (m ttfLineMetrics) strikeoutOffsetFactor() float64 {
	return m.emFraction(m.StrikeoutPosition, 0.25)
}

func (m ttfLineMetrics) strikeoutThicknessFactor() float64 {
	return m.emFraction(m.StrikeoutSize, 0.05)
}

func (m ttfLineMetrics) subscriptSizeFactor() float64 {
	return m.emFraction(m.SubscriptSize, 0.65)
}

func (m ttfLineMetrics) subscriptOffsetFactor() float64 {
	return m.emFraction(m.SubscriptOffset, 0.14)
}

func (m ttfLineMetrics) superscriptSizeFactor() float64 {
	return m.emFraction(m.SuperscriptSize, 0.65)
}

func (m ttfLineMetrics) superscriptOffsetFactor() float64 {
	return m.emFraction(m.SuperscriptOffset, 0.4)
}

// typoAscenderFactor is the baseline drop gopdf itself applies inside Cell().
func (m ttfLineMetrics) typoAscenderFactor() float64 {
	if m.UnitsPerEm <= 0 {
		return 0
	}
	return m.TypoAscender / m.UnitsPerEm
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

// baselineInternalLeadingFactor is the remaining gap between the baseline the
// hhea/OS2 model predicts and where PowerPoint actually draws: PowerPoint pads
// the top of the line box by part of the font's internal leading before the
// ascent. Measured against PowerPoint's own PDF export at 8, 11 and 24pt, in
// both top- and centre-anchored boxes, the offset is a constant 0.09 em.
const baselineInternalLeadingFactor = 0.09

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
	shift := halfLeading + ascent - gopdfBaseline + baselineInternalLeadingFactor
	return clampFloat(shift, minBaselineShiftFactor, maxBaselineShiftFactor)
}

func clampFloat(v, lo, hi float64) float64 {
	return math.Min(math.Max(v, lo), hi)
}

// registerPDFFontMetrics parses fontPath and stores its metrics under alias for
// the document being written. A parse failure is not fatal: the alias simply
// falls back to Calibri-like metrics, matching how font registration itself
// degrades.
func registerPDFFontMetrics(pdf *gopdf.GoPdf, alias, fontPath string) {
	m, err := readTTFLineMetrics(fontPath)
	if err != nil {
		return
	}
	documentFontsFor(pdf).putMetrics(alias, m)
}

// lookupPDFFontMetrics returns the metrics registered for alias in this
// document, or the Calibri-like fallback.
func lookupPDFFontMetrics(pdf *gopdf.GoPdf, alias string) ttfLineMetrics {
	m, ok := documentFontsFor(pdf).lookupMetrics(alias)
	if !ok {
		return fallbackLineMetrics
	}
	return m
}

// metricsForFontHint resolves a PPTX typeface hint to the metrics of the font
// actually embedded for it in this document.
func metricsForFontHint(pdf *gopdf.GoPdf, fontHint string) ttfLineMetrics {
	return lookupPDFFontMetrics(pdf, resolvePDFFontAlias(pdf, fontHint))
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

	// post table: underline placement, in font units.
	postUnderlinePositionOffset  = 8
	postUnderlineThicknessOffset = 10

	// OS/2 table: strike-through placement and the script metrics PowerPoint
	// uses to draw subscript and superscript runs.
	os2SubscriptYSizeOffset     = 12
	os2SubscriptYOffsetOffset   = 16
	os2SuperscriptYSizeOffset   = 20
	os2SuperscriptYOffsetOffset = 24
	os2StrikeoutSizeOffset      = 26
	os2StrikeoutPositionOffset  = 28
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

	metrics := ttfLineMetrics{
		UnitsPerEm: float64(upem),
		Ascender:   float64(ascender),
		Descender:  float64(descender),
		LineGap:    float64(lineGap),
	}
	readTTFDecorationMetrics(data, tables, &metrics)
	return metrics, nil
}

// readTTFDecorationMetrics fills in the optional post and OS/2 fields. Every one
// of them is allowed to be missing: the em-fraction accessors substitute a
// default, so a font without a post table still gets an underline.
func readTTFDecorationMetrics(data []byte, tables map[string]uint32, metrics *ttfLineMetrics) {
	if post, ok := tables["post"]; ok {
		metrics.UnderlinePosition = readInt16OrZero(data, post+postUnderlinePositionOffset)
		metrics.UnderlineThickness = readInt16OrZero(data, post+postUnderlineThicknessOffset)
	}
	os2, ok := tables["OS/2"]
	if !ok {
		return
	}
	// gopdf offsets its own baseline by OS/2 typoAscender. When OS/2 is absent
	// or truncated it uses zero, so mirror that rather than guessing.
	metrics.TypoAscender = readInt16OrZero(data, os2+os2TypoAscenderOffset)
	metrics.SubscriptSize = readInt16OrZero(data, os2+os2SubscriptYSizeOffset)
	metrics.SubscriptOffset = readInt16OrZero(data, os2+os2SubscriptYOffsetOffset)
	metrics.SuperscriptSize = readInt16OrZero(data, os2+os2SuperscriptYSizeOffset)
	metrics.SuperscriptOffset = readInt16OrZero(data, os2+os2SuperscriptYOffsetOffset)
	metrics.StrikeoutSize = readInt16OrZero(data, os2+os2StrikeoutSizeOffset)
	metrics.StrikeoutPosition = readInt16OrZero(data, os2+os2StrikeoutPositionOffset)
}

func readInt16OrZero(data []byte, offset uint32) float64 {
	v, err := readInt16At(data, offset)
	if err != nil {
		return 0
	}
	return float64(v)
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
