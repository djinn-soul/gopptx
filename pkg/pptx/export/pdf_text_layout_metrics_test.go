package export

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

// calibriMetrics are Calibri's real hhea/OS-2 values, used so the derived
// factors can be asserted without depending on a font file being installed.
func calibriMetrics() ttfLineMetrics {
	return ttfLineMetrics{
		UnitsPerEm:   2048,
		Ascender:     1536,
		Descender:    -512,
		LineGap:      452,
		TypoAscender: 1536,
	}
}

func TestLineHeightIsFontIndependent(t *testing.T) {
	t.Parallel()

	// Measured from PowerPoint's own render: 18pt Calibri, Segoe UI, Georgia,
	// Verdana, Arial and Consolas all land on the same ~19.4pt pitch even though
	// their hhea line heights span 1.14 to 1.33 em. So the line box is 1.2 x the
	// point size for every font.
	if got := pdfLineHeight(18); math.Abs(got-21.6) > 0.001 {
		t.Fatalf("18pt line height=%v want 21.6", got)
	}
	if got := pdfLineHeight(44); math.Abs(got-52.8) > 0.001 {
		t.Fatalf("44pt line height=%v want 52.8", got)
	}
}

func TestBaselineShiftFactorIsRelativeToGopdfPlacement(t *testing.T) {
	t.Parallel()

	// Calibri: ascent 0.75 em, descent 0.25 em, so half-leading inside the
	// 1.2 em line box is (1.2 - 1.0) / 2 = 0.1 em and the baseline sits at
	// 0.85 em. gopdf's Cell() has already applied typoAscender (0.75 em), so the
	// renderer adds the remaining 0.1 em.
	got := calibriMetrics().baselineShiftFactor()
	if math.Abs(got-0.1) > 0.0005 {
		t.Fatalf("calibri baseline shift factor=%v want ~0.1", got)
	}
}

func TestLineMetricsFallBackWhenUnitsPerEmMissing(t *testing.T) {
	t.Parallel()

	var empty ttfLineMetrics
	if empty.baselineShiftFactor() != fallbackLineMetrics.baselineShiftFactor() {
		t.Fatalf("zero metrics baseline shift=%v want fallback", empty.baselineShiftFactor())
	}
}

func TestReadTTFLineMetricsParsesRealFont(t *testing.T) {
	t.Parallel()

	path := findTestFontPath()
	if path == "" {
		t.Skip("no system font available to parse")
	}
	m, err := readTTFLineMetrics(path)
	if err != nil {
		t.Fatalf("readTTFLineMetrics(%q) error: %v", path, err)
	}
	if m.UnitsPerEm <= 0 {
		t.Fatalf("unitsPerEm=%v want > 0", m.UnitsPerEm)
	}
	if m.Ascender <= 0 {
		t.Fatalf("ascender=%v want > 0", m.Ascender)
	}
	if m.Descender > 0 {
		t.Fatalf("descender=%v want <= 0 (hhea stores it negative)", m.Descender)
	}
	// Every mainstream UI font sits well inside this band; anything outside it
	// means the table offsets were misread.
	if f := (m.Ascender - m.Descender) / m.UnitsPerEm; f < 0.9 || f > 2.0 {
		t.Fatalf("ascent+descent=%v em out of plausible range", f)
	}
}

func TestLineMetricsClampJunkVerticalMetrics(t *testing.T) {
	t.Parallel()

	// A font claiming a 40 em ascender must not push its text off the slide.
	junk := ttfLineMetrics{UnitsPerEm: 2048, Ascender: 81920, Descender: -81920, LineGap: 0}
	if got := junk.baselineShiftFactor(); got != maxBaselineShiftFactor {
		t.Fatalf("junk baseline shift factor=%v want clamp to %v", got, maxBaselineShiftFactor)
	}

	// ...and one claiming a huge typoAscender must not drag it upward off it.
	inverted := ttfLineMetrics{UnitsPerEm: 2048, Ascender: 1536, Descender: -512, TypoAscender: 81920}
	if got := inverted.baselineShiftFactor(); got != minBaselineShiftFactor {
		t.Fatalf("inverted baseline shift factor=%v want clamp to %v", got, minBaselineShiftFactor)
	}
}

func TestReadTTFLineMetricsRejectsNonFont(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "not-a-font.ttf")
	if err := os.WriteFile(path, []byte("nope"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if _, err := readTTFLineMetrics(path); err == nil {
		t.Fatal("readTTFLineMetrics on a non-font returned no error")
	}
}

func TestMetricsForFontHintFallsBackForUnregisteredAlias(t *testing.T) {
	resetPDFFontMetrics()
	got := metricsForFontHint("Nonexistent Font")
	if got != fallbackLineMetrics {
		t.Fatalf("unregistered hint metrics=%+v want fallback", got)
	}
}

// findTestFontPath reuses the renderer's own font search rather than repeating a
// hard-coded path list, so the test cannot drift from production and skips
// cleanly on a machine (or CI image) that ships none of these families.
func findTestFontPath() string {
	for _, family := range []string{fontFamilySans, fontFamilySerif, fontFamilyMono} {
		if paths := systemFontPathsForFamily(family); len(paths) > 0 {
			return paths[0]
		}
	}
	return ""
}
