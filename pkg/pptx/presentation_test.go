package pptx_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesProducesPackage(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Intro").AddBullet("First").AddBullet("Second"),
		pptx.NewSlide("Details"),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
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
	slides := []pptx.SlideContent{
		pptx.NewSlide("Title & More").AddBullet("Use <tag>").AddBullet("5 > 3"),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

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
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{{}})
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

	slides := []pptx.SlideContent{
		pptx.NewSlide("Image Slide").AddImage(pptx.NewImage(imgPath, 1200000, 1700000, 2400000, 1800000)),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	if !testutil.ZipHasFile(zr, "ppt/media/image1.png") {
		t.Fatalf("missing embedded media file")
	}

	relsXML := testutil.ReadZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `relationships/image"`) {
		t.Fatalf("expected image relationship in slide rels")
	}
	if !strings.Contains(relsXML, `Target="../media/image1.png"`) {
		t.Fatalf("expected image media target in slide rels")
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `a:blip r:embed="rId2"`) {
		t.Fatalf("expected image embed reference in slide xml")
	}
}

func TestCreateWithSlidesFailsForMissingImage(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Image Slide").AddImage(pptx.NewImage("does-not-exist.png", 1, 1, 1, 1)),
	}
	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected error for missing image file")
	}
}

func TestCreateWithSlidesEmbedsTable(t *testing.T) {
	table := pptx.NewTable([]int64{2743400, 2743400, 2743400}).
		AddRow([]string{"Name", "Status", "Owner"}).
		AddRow([]string{"Parser", "Done", "Core Team"})

	slides := []pptx.SlideContent{
		pptx.NewSlide("Table Slide").WithTable(table),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
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
	table := pptx.NewTable([]int64{2000000, 2000000}).
		AddRow([]string{"A"})

	slides := []pptx.SlideContent{
		pptx.NewSlide("Broken Table").WithTable(table),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected table validation error")
	}
}

func TestCreateWithSlidesEmbedsStyledTableCell(t *testing.T) {
	table := pptx.NewTable([]int64{2743400, 2743400}).
		AddStyledRow([]pptx.TableCell{
			pptx.NewTableCell("Header").WithBold(true).WithBackgroundColor("1F497D"),
			pptx.NewTableCell("Value"),
		}).
		AddRow([]string{"Row 1", "Plain"})

	data, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("Styled Table").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
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
	table := pptx.NewTable([]int64{2743400}).
		AddStyledRow([]pptx.TableCell{
			pptx.NewTableCell("Header").WithBackgroundColor("NOTHEX"),
		})

	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("Broken Styled Table").WithTable(table)})
	if err == nil {
		t.Fatalf("expected styled table color validation error")
	}
}

func TestCreateWithSlidesEmbedsBackgroundColor(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Red Background").WithBackgroundColor("#FF0000"),
	}

	data, err := pptx.CreateWithSlides("Background Test", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	needle := `<a:srgbClr val="FF0000"/>`
	if !strings.Contains(slideXML, needle) {
		t.Fatalf("expected background color %q in slide XML", needle)
	}
}

func TestCreateWithSlidesTitleAlignment(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Right Title").WithTitleAlign("r"),
	}

	data, err := pptx.CreateWithSlides("Title Alignment Test", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	needle := `<a:pPr algn="r"/>`
	if !strings.Contains(slideXML, needle) {
		t.Fatalf("expected title alignment %q in slide XML", needle)
	}
}

func TestCreateWithSlidesContentVerticalAlignment(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Centered Content").
			AddBullet("Bullet 1").
			WithContentVAlign("ctr"),
	}

	data, err := pptx.CreateWithSlides("Vertical Alignment Test", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
	needle := `anchor="ctr"`
	if !strings.Contains(slideXML, needle) {
		t.Fatalf("expected content vertical alignment %q in slide XML", needle)
	}
}

func TestCreateWithSlidesTitleFontAndNumber(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Styled Slide").
			WithTitleFont("Consolas").
			WithSlideNumber(true),
	}

	data, err := pptx.CreateWithSlides("Font and Number Test", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	// Check Title Font
	fontNeedle := `typeface="Consolas"`
	if !strings.Contains(slideXML, fontNeedle) {
		t.Fatalf("expected title font %q in slide XML", fontNeedle)
	}

	// Check Slide Number
	sldNumNeedle := `type="sldNum"`
	if !strings.Contains(slideXML, sldNumNeedle) {
		t.Fatalf("expected slide number placeholder %q in slide XML", sldNumNeedle)
	}
	fldNeedle := `type="slnum"`
	if !strings.Contains(slideXML, fldNeedle) {
		t.Fatalf("expected slide number field %q in slide XML", fldNeedle)
	}
}
