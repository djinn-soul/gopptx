package export

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

// mediaBoxRe matches the page box gopdf writes into the PDF catalog.
var mediaBoxRe = regexp.MustCompile(`/MediaBox\s*\[\s*0\s+0\s+([0-9.]+)\s+([0-9.]+)\s*\]`)

// readPDFPageSize returns the width/height in points of the first page.
func readPDFPageSize(t *testing.T, pdfPath string) (float64, float64) {
	t.Helper()
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		t.Fatalf("read pdf: %v", err)
	}
	m := mediaBoxRe.FindSubmatch(data)
	if m == nil {
		t.Fatalf("no /MediaBox found in %s", pdfPath)
	}
	w, err := strconv.ParseFloat(string(m[1]), 64)
	if err != nil {
		t.Fatalf("parse MediaBox width: %v", err)
	}
	h, err := strconv.ParseFloat(string(m[2]), 64)
	if err != nil {
		t.Fatalf("parse MediaBox height: %v", err)
	}
	return w, h
}

func slideWithBullets(title string, bullets ...string) elements.SlideContent {
	s := elements.NewSlide(title)
	s.Bullets = bullets
	return s
}

// writeDeck creates a PPTX on disk with the requested slide size and returns its path.
func writeDeck(t *testing.T, dir string, size pptx.SlideSize) string {
	t.Helper()
	data, err := pptx.CreateWithMetadata(
		pptx.Metadata{Metadata: common.Metadata{Title: "deck", SlideSize: size}},
		[]elements.SlideContent{slideWithBullets("Hello", "one", "two")},
	)
	if err != nil {
		t.Fatalf("create pptx: %v", err)
	}
	path := filepath.Join(dir, "deck.pptx")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write pptx: %v", err)
	}
	return path
}

// fakeSOffice puts a stub "soffice" on PATH for the duration of the test.
// It exits 0 and creates each path in createFiles as a zero-byte file, which
// mimics LibreOffice's habit of reporting success while producing nothing
// usable (e.g. when another instance holds the user-profile lock).
func fakeSOffice(t *testing.T, createFiles ...string) {
	t.Helper()

	binDir := t.TempDir()
	var name string
	var lines []string
	if runtime.GOOS == osWindows {
		name = "soffice.bat"
		lines = append(lines, "@echo off")
		for _, f := range createFiles {
			lines = append(lines, `type nul > "`+f+`"`)
		}
		lines = append(lines, "exit /b 0")
	} else {
		name = "soffice"
		lines = append(lines, "#!/bin/sh")
		for _, f := range createFiles {
			lines = append(lines, `: > "`+f+`"`)
		}
		lines = append(lines, "exit 0")
	}

	script := filepath.Join(binDir, name)
	sep := "\n"
	if runtime.GOOS == osWindows {
		sep = "\r\n"
	}
	// 0o700: the stub has to be executable for exec.LookPath to accept it.
	if err := os.WriteFile(script, []byte(strings.Join(lines, sep)+sep), 0o700); err != nil {
		t.Fatalf("write soffice stub: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// fakeSlowSOffice puts a stub "soffice" on PATH that hangs well past any
// timeout the test sets, standing in for LibreOffice blocking on a profile
// lock or a modal dialog.
func fakeSlowSOffice(t *testing.T) {
	t.Helper()

	binDir := t.TempDir()
	var name, body string
	if runtime.GOOS == osWindows {
		// ping is the portable way to sleep in a batch file.
		name, body = "soffice.bat", "@echo off\r\nping -n 30 127.0.0.1 > nul\r\nexit /b 0\r\n"
	} else {
		name, body = "soffice", "#!/bin/sh\nsleep 30\nexit 0\n"
	}

	script := filepath.Join(binDir, name)
	// 0o700: the stub has to be executable for exec.LookPath to accept it.
	if err := os.WriteFile(script, []byte(body), 0o700); err != nil {
		t.Fatalf("write soffice stub: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// TestLibreOfficeHonoursTimeout guards that a wedged converter fails the export
// instead of blocking the caller forever.
func TestLibreOfficeHonoursTimeout(t *testing.T) {
	dir := t.TempDir()
	pptxPath := writeDeck(t, dir, pptx.SlideSize4x3())
	pdfPath := filepath.Join(dir, "out.pdf")

	fakeSlowSOffice(t)

	start := time.Now()
	err := PDFFromFileWithOptions(pptxPath, pdfPath, PDFOptions{
		Driver:  PDFDriverLibreOffice,
		Timeout: 2 * time.Second,
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("error should mention the timeout, got: %v", err)
	}
	if elapsed > 20*time.Second {
		t.Errorf("export took %s; the timeout was not enforced", elapsed)
	}
}

// TestNativePDFPageSizeMatchesSlideSize guards that the native renderer sizes
// the PDF page from the deck's own <p:sldSz> rather than assuming 4:3.
func TestNativePDFPageSizeMatchesSlideSize(t *testing.T) {
	tests := []struct {
		name  string
		size  pptx.SlideSize
		wantW float64
		wantH float64
	}{
		{"4x3", pptx.SlideSize4x3(), 720, 540},
		{"16x9", pptx.SlideSize16x9(), 960, 540},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			pptxPath := writeDeck(t, dir, tc.size)
			pdfPath := filepath.Join(dir, "out.pdf")

			err := PDFFromFileWithOptions(pptxPath, pdfPath, PDFOptions{Driver: PDFDriverNative})
			if err != nil {
				t.Fatalf("native export: %v", err)
			}

			gotW, gotH := readPDFPageSize(t, pdfPath)
			if gotW != tc.wantW || gotH != tc.wantH {
				t.Errorf("page size = %.2fx%.2f pt, want %.2fx%.2f pt", gotW, gotH, tc.wantW, tc.wantH)
			}
		})
	}
}

// TestNativePDFReportsUnrenderableImage guards that a picture which cannot be
// drawn is reported rather than dropped, so callers do not ship a deck with
// silently missing content.
func TestNativePDFReportsUnrenderableImage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.pdf")

	slide := elements.NewSlide("Has a broken picture")
	slide.Images = []shapes.Image{{
		Data: []byte("this is not an image"),
		X:    styling.Inches(1),
		Y:    styling.Inches(1),
		CX:   styling.Inches(2),
		CY:   styling.Inches(2),
	}}

	err := PDFWithOptions("t", []elements.SlideContent{slide}, out, PDFOptions{Driver: PDFDriverNative})
	if err == nil {
		t.Fatal("expected an error for an undecodable image, got nil")
	}
	if !strings.Contains(err.Error(), "image") {
		t.Errorf("error should identify the failing image, got: %v", err)
	}
}

// TestNativePDFRendersEveryTable guards that tables beyond the first are drawn
// rather than overwritten.
func TestNativePDFRendersEveryTable(t *testing.T) {
	dir := t.TempDir()

	one := elements.NewSlide("One table")
	first := tables.NewTable([]styling.Length{styling.Inches(2)}).AddRow([]string{"alpha"})
	one.Table = &first

	two := elements.NewSlide("Two tables")
	firstOfTwo := tables.NewTable([]styling.Length{styling.Inches(2)}).AddRow([]string{"alpha"})
	two.Table = &firstOfTwo
	two.Tables = []tables.Table{
		tables.NewTable([]styling.Length{styling.Inches(2)}).AddRow([]string{"beta"}),
	}

	sizes := make(map[string]int64)
	for name, slide := range map[string]elements.SlideContent{"one": one, "two": two} {
		out := filepath.Join(dir, name+".pdf")
		err := PDFWithOptions("t", []elements.SlideContent{slide}, out, PDFOptions{Driver: PDFDriverNative})
		if err != nil {
			t.Fatalf("%s: native export: %v", name, err)
		}
		info, err := os.Stat(out)
		if err != nil {
			t.Fatalf("%s: stat: %v", name, err)
		}
		sizes[name] = info.Size()
	}

	if sizes["two"] <= sizes["one"] {
		t.Errorf("the second table was not drawn: one table = %d bytes, two tables = %d bytes",
			sizes["one"], sizes["two"])
	}
}

// TestLibreOfficeErrorsWhenNoPDFProduced covers the case where the output path
// equals the name LibreOffice derives from the input (deck.pptx -> deck.pdf in
// the same directory). The success path must still verify the file exists.
func TestLibreOfficeErrorsWhenNoPDFProduced(t *testing.T) {
	dir := t.TempDir()
	pptxPath := writeDeck(t, dir, pptx.SlideSize4x3())
	pdfPath := filepath.Join(dir, "deck.pdf") // same name soffice would generate

	fakeSOffice(t) // exits 0, writes nothing

	err := PDFFromFileWithOptions(pptxPath, pdfPath, PDFOptions{Driver: PDFDriverLibreOffice})
	if err == nil {
		t.Fatal("expected an error when LibreOffice produced no PDF, got nil")
	}
	if _, statErr := os.Stat(pdfPath); statErr == nil {
		t.Fatal("no PDF should exist, but one was found")
	}
}

// TestLibreOfficeErrorsWhenPDFIsEmpty guards against reporting success for a
// zero-byte PDF, which is not a readable document.
func TestLibreOfficeErrorsWhenPDFIsEmpty(t *testing.T) {
	dir := t.TempDir()
	pptxPath := writeDeck(t, dir, pptx.SlideSize4x3())
	pdfPath := filepath.Join(dir, "renamed.pdf")

	// soffice "succeeds" but leaves an empty deck.pdf next to the input.
	fakeSOffice(t, filepath.Join(dir, "deck.pdf"))

	err := PDFFromFileWithOptions(pptxPath, pdfPath, PDFOptions{Driver: PDFDriverLibreOffice})
	if err == nil {
		t.Fatal("expected an error for a zero-byte PDF, got nil")
	}
}

// TestLibreOfficeDoesNotClobberFilesInOutputDir guards that the intermediate
// PPTX is not written into the user's output directory under a name derived
// only from the deck title, which collides across concurrent exports and can
// destroy a pre-existing file.
func TestLibreOfficeDoesNotClobberFilesInOutputDir(t *testing.T) {
	dir := t.TempDir()
	pdfPath := filepath.Join(dir, "out.pdf")

	// A file the user already owns, whose name matches the temp name the
	// exporter derives from the title "deck".
	victim := filepath.Join(dir, "gopptx_deck_temp.pptx")
	const sentinel = "user data, must survive"
	if err := os.WriteFile(victim, []byte(sentinel), 0o600); err != nil {
		t.Fatalf("write victim: %v", err)
	}

	// Let the conversion "succeed" so the happy path runs end to end.
	fakeSOffice(t, filepath.Join(dir, "gopptx_deck_temp.pdf"), pdfPath)

	slides := []elements.SlideContent{elements.NewSlide("Hello")}
	_ = PDFWithOptions("deck", slides, pdfPath, PDFOptions{Driver: PDFDriverLibreOffice})

	got, err := os.ReadFile(victim)
	if err != nil {
		t.Fatalf("pre-existing file was deleted by the exporter: %v", err)
	}
	if string(got) != sentinel {
		t.Errorf("pre-existing file was overwritten: got %q, want %q", string(got), sentinel)
	}
}
