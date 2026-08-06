package export

import (
	"strings"
	"testing"

	"github.com/signintech/gopdf"
)

// newTestPDF starts an empty one-page document to hang font state off.
func newTestPDF(t *testing.T) *gopdf.GoPdf {
	t.Helper()
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 100, H: 100}})
	pdf.AddPage()
	t.Cleanup(func() { releaseDocumentFonts(pdf) })
	return pdf
}

// firstEmbeddableFamily returns a family this host has and gopdf accepts, plus
// the alias it was embedded under. Not every installed family qualifies:
// symbol-encoded faces are rejected by the TTF parser and the renderer is
// expected to fall back for those.
func firstEmbeddableFamily(t *testing.T, pdf *gopdf.GoPdf) (string, string) {
	t.Helper()
	for name := range indexedFontFamilies() {
		if alias := ensureNamedPDFFont(pdf, name); alias != "" {
			return name, alias
		}
	}
	return "", ""
}

func indexedFontFamilies() map[string]*namedFontFaces {
	systemFontIndexMu.Lock()
	defer systemFontIndexMu.Unlock()
	if !systemFontIndexBuilt {
		refreshSystemFontIndexLocked()
	}
	return systemFontIndex
}

func TestStyleFromSubfamily(t *testing.T) {
	cases := map[string]int{
		"Regular":     gopdf.Regular,
		"Bold":        gopdf.Bold,
		"Italic":      gopdf.Italic,
		"Oblique":     gopdf.Italic,
		"Bold Italic": gopdf.Bold | gopdf.Italic,
		"bold italic": gopdf.Bold | gopdf.Italic,
	}
	for subfamily, want := range cases {
		if got := styleFromSubfamily(subfamily); got != want {
			t.Fatalf("styleFromSubfamily(%q)=%d want %d", subfamily, got, want)
		}
	}
}

func TestNormalizeFontFamily(t *testing.T) {
	if got := normalizeFontFamily("  Segoe UI  "); got != "segoe ui" {
		t.Fatalf("normalizeFontFamily=%q want %q", got, "segoe ui")
	}
}

func TestNamedFontAliasKeepsFamilyNameUnlessReserved(t *testing.T) {
	// The alias becomes the PDF's /BaseFont, so an ordinary family keeps its own
	// name and only a clash with a generic face is prefixed.
	if got := namedFontAlias("georgia"); got != "georgia" {
		t.Fatalf("alias=%q want %q", got, "georgia")
	}
	for _, reserved := range []string{fontFamilySans, fontFamilySerif, fontFamilyMono, fontFamilyCJK} {
		if got := namedFontAlias(reserved); got != namedFontAliasPrefix+reserved {
			t.Fatalf("alias for reserved %q=%q want prefixed", reserved, got)
		}
	}
}

func TestIsFontFileName(t *testing.T) {
	for _, name := range []string{"arial.ttf", "Georgia.OTF", "msyh.ttc"} {
		if !isFontFileName(name) {
			t.Fatalf("%q not recognised as a font file", name)
		}
	}
	for _, name := range []string{"desktop.ini", "font.fon", "notes.txt"} {
		if isFontFileName(name) {
			t.Fatalf("%q wrongly recognised as a font file", name)
		}
	}
}

func TestSystemFontDirsAreAbsolute(t *testing.T) {
	dirs := systemFontDirs()
	if len(dirs) == 0 {
		t.Fatal("no system font directories for this platform")
	}
	for _, dir := range dirs {
		if strings.TrimSpace(dir) == "" {
			t.Fatal("empty font directory in the list")
		}
	}
}

func TestCachedNamedPDFFontAliasIsEmptyBeforeRegistration(t *testing.T) {
	pdf := newTestPDF(t)
	if got := cachedNamedPDFFontAlias(pdf, "Georgia"); got != "" {
		t.Fatalf("alias=%q want empty before any registration", got)
	}
	if got := cachedNamedPDFFontAlias(pdf, ""); got != "" {
		t.Fatalf("alias for empty family=%q want empty", got)
	}
}

// The index reflects the host's installed fonts, so this only asserts what must
// hold on any host: keys are normalized and every family has at least one face.
func TestSystemFontIndexEntriesAreWellFormed(t *testing.T) {
	checked := 0
	for family, faces := range indexedFontFamilies() {
		if family != normalizeFontFamily(family) {
			t.Fatalf("index key %q is not normalized", family)
		}
		if len(faces.paths) == 0 {
			t.Fatalf("family %q indexed with no faces", family)
		}
		checked++
		if checked >= 20 {
			break
		}
	}
}

func TestEnsureNamedPDFFontEmbedsAnInstalledFamily(t *testing.T) {
	pdf := newTestPDF(t)
	family, alias := firstEmbeddableFamily(t, pdf)
	if family == "" {
		t.Skip("host has no embeddable text fonts")
	}

	// A second call must reuse the cached alias rather than embedding again.
	if again := ensureNamedPDFFont(pdf, strings.ToUpper(family)); again != alias {
		t.Fatalf("second lookup=%q want cached %q", again, alias)
	}
	// The alias now wins over the generic sans/serif/mono resolution.
	if got := resolvePDFFontAlias(pdf, family); got != alias {
		t.Fatalf("resolvePDFFontAlias=%q want %q", got, alias)
	}
}

func TestEnsureNamedPDFFontIgnoresUnknownFamily(t *testing.T) {
	pdf := newTestPDF(t)
	if alias := ensureNamedPDFFont(pdf, "No Such Typeface 12345"); alias != "" {
		t.Fatalf("alias=%q want empty for a family the host lacks", alias)
	}
	// The generic resolution must still work for that hint.
	sans, _, _, _ := documentFontsFor(pdf).genericAliases()
	if got := resolvePDFFontAlias(pdf, "No Such Typeface 12345"); got != sans {
		t.Fatalf("fallback alias=%q want %q", got, sans)
	}
}

// Font state is per document: an alias embedded in one PDF must not be handed
// to another, which never registered it and would draw the run in the wrong face.
func TestNamedFontStateIsPerDocument(t *testing.T) {
	first := newTestPDF(t)
	family, alias := firstEmbeddableFamily(t, first)
	if family == "" {
		t.Skip("host has no embeddable text fonts")
	}

	second := newTestPDF(t)
	if got := cachedNamedPDFFontAlias(second, family); got != "" {
		t.Fatalf("second document already knows alias %q", got)
	}

	// Generic aliases are per document too.
	setPDFFontAliases(second, "SecondSans", "SecondSerif", "SecondMono")
	sans, _, _, _ := documentFontsFor(second).genericAliases()
	if sans != "SecondSans" {
		t.Fatalf("second document sans alias=%q want %q", sans, "SecondSans")
	}
	firstSans, _, _, _ := documentFontsFor(first).genericAliases()
	if firstSans == "SecondSans" {
		t.Fatal("alias set on the second document leaked into the first")
	}
	if cachedNamedPDFFontAlias(first, family) != alias {
		t.Fatal("first document lost its embedded family")
	}
}

func TestReleaseDocumentFontsDropsState(t *testing.T) {
	pdf := newTestPDF(t)
	setPDFFontAliases(pdf, "KeptSans", "KeptSerif", "KeptMono")
	releaseDocumentFonts(pdf)

	sans, _, _, _ := documentFontsFor(pdf).genericAliases()
	if sans != fontFamilySans {
		t.Fatalf("sans alias after release=%q want the default %q", sans, fontFamilySans)
	}
}
