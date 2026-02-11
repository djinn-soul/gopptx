package pptx

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesProducesPackage(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Intro").AddBullet("First").AddBullet("Second"),
		NewSlide("Details"),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	required := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"ppt/presentation.xml",
		"ppt/_rels/presentation.xml.rels",
		"ppt/slides/slide1.xml",
		"ppt/slides/slide2.xml",
		"ppt/slides/_rels/slide1.xml.rels",
		"ppt/slides/_rels/slide2.xml.rels",
		"ppt/slideLayouts/slideLayout1.xml",
		"ppt/slideMasters/slideMaster1.xml",
		"ppt/theme/theme1.xml",
		"docProps/core.xml",
		"docProps/app.xml",
	}

	present := map[string]bool{}
	for _, f := range zr.File {
		present[f.Name] = true
	}
	for _, name := range required {
		if !present[name] {
			t.Fatalf("missing %s", name)
		}
	}
}

func TestCreateWithSlidesEscapesText(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Title & More").AddBullet("Use <tag>").AddBullet("5 > 3"),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	var slideXML string
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			r, err := f.Open()
			if err != nil {
				t.Fatalf("open slide: %v", err)
			}
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(r); err != nil {
				_ = r.Close()
				t.Fatalf("read slide: %v", err)
			}
			_ = r.Close()
			slideXML = buf.String()
			break
		}
	}

	if slideXML == "" {
		t.Fatalf("slide XML not found")
	}

	checks := []string{
		"Title &amp; More",
		"Use &lt;tag&gt;",
		"5 &gt; 3",
	}

	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesValidation(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{{}})
	if err == nil {
		t.Fatalf("expected error for empty title")
	}
}

func TestCreateWithSlidesEmbedsImage(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "sample.png")
	if err := os.WriteFile(imgPath, testutil.TinyPNG, 0o600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	slides := []SlideContent{
		NewSlide("Image Slide").AddImage(NewImage(imgPath, 1200000, 1700000, 2400000, 1800000)),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	if !zipHasFile(zr, "ppt/media/image1.png") {
		t.Fatalf("missing embedded media file")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `relationships/image"`) {
		t.Fatalf("expected image relationship in slide rels")
	}
	if !strings.Contains(relsXML, `Target="../media/image1.png"`) {
		t.Fatalf("expected image media target in slide rels")
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `a:blip r:embed="rId2"`) {
		t.Fatalf("expected image embed reference in slide xml")
	}
}

func TestCreateWithSlidesFailsForMissingImage(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Image Slide").AddImage(NewImage("does-not-exist.png", 1, 1, 1, 1)),
	}
	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected error for missing image file")
	}
}

func TestCreateWithSlidesEmbedsTable(t *testing.T) {
	table := NewTable([]int64{2743400, 2743400, 2743400}).
		AddRow([]string{"Name", "Status", "Owner"}).
		AddRow([]string{"Parser", "Done", "Core Team"})

	slides := []SlideContent{
		NewSlide("Table Slide").WithTable(table),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, "<a:tbl>") {
		t.Fatalf("expected table XML in slide")
	}
	if !strings.Contains(slideXML, "<a:t>Name</a:t>") {
		t.Fatalf("expected table cell text in slide")
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("unexpected bullet content shape when table is present")
	}
}

func TestCreateWithSlidesRejectsInvalidTable(t *testing.T) {
	table := NewTable([]int64{2000000, 2000000}).
		AddRow([]string{"A"})

	slides := []SlideContent{
		NewSlide("Broken Table").WithTable(table),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected table validation error")
	}
}

func TestCreateWithSlidesEmbedsStyledTableCell(t *testing.T) {
	table := NewTable([]int64{2743400, 2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").WithBold(true).WithBackgroundColor("1F497D"),
			NewTableCell("Value"),
		}).
		AddRow([]string{"Row 1", "Plain"})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Styled Table").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`<a:rPr lang="en-US" dirty="0" b="1"/>`,
		`<a:tcPr><a:solidFill><a:srgbClr val="1F497D"/></a:solidFill></a:tcPr>`,
		`<a:t>Header</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in table XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsStyledTableInvalidColor(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").WithBackgroundColor("NOTHEX"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Styled Table").WithTable(table)})
	if err == nil {
		t.Fatalf("expected styled table color validation error")
	}
}

func zipHasFile(zr *zip.Reader, name string) bool {
	for _, f := range zr.File {
		if f.Name == name {
			return true
		}
	}
	return false
}

func readZipFile(t *testing.T, zr *zip.Reader, name string) string {
	t.Helper()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		r, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer func() { _ = r.Close() }()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		return buf.String()
	}
	t.Fatalf("file %s not found in zip", name)
	return ""
}
