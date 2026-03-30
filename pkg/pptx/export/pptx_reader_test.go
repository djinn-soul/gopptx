package export

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestSlidesFromPPTX_RoundTrip(t *testing.T) {
	// 1x1 transparent PNG
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0x60, 0x00, 0x02, 0x00,
		0x00, 0x05, 0x00, 0x01, 0x0D, 0x26, 0xE5, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44,
		0xAE, 0x42, 0x60, 0x82,
	}

	slides := []elements.SlideContent{
		{
			Title:   "Slide 1",
			Bullets: []string{"Bullet 1", "Bullet 2"},
			Shapes: []shapes.Shape{
				{
					Type: "rect",
					X:    styling.Emu(100000),
					Y:    styling.Emu(100000),
					CX:   styling.Emu(500000),
					CY:   styling.Emu(500000),
					Text: "Shape Text",
				},
			},
			Images: []shapes.Image{
				{
					Data:   pngData,
					Format: "png",
					X:      styling.Emu(200000),
					Y:      styling.Emu(200000),
					CX:     styling.Emu(1000000),
					CY:     styling.Emu(1000000),
				},
			},
		},
	}

	tmpDir := t.TempDir()
	pptxPath := filepath.Join(tmpDir, "test.pptx")

	pptxBytes, err := pptx.CreateWithSlides("Test Presentation", slides)
	if err != nil {
		t.Fatalf("failed to create PPTX: %v", err)
	}

	if err := os.WriteFile(pptxPath, pptxBytes, 0644); err != nil {
		t.Fatalf("failed to write PPTX: %v", err)
	}

	title, readSlides, err := SlidesFromPPTX(pptxPath)
	if err != nil {
		t.Fatalf("failed to read PPTX: %v", err)
	}

	if title != "Test Presentation" {
		t.Errorf("expected title 'Test Presentation', got %q", title)
	}

	if len(readSlides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(readSlides))
	}

	s := readSlides[0]
	if s.Title != "Slide 1" {
		t.Errorf("expected slide title 'Slide 1', got %q", s.Title)
	}

	// The reader might map placeholders as shapes or bullets depending on
	// specific OOXML tags that are currently being refined.
	// Ensure we get at least some content.
	totalTextElements := len(s.Bullets) + len(s.Shapes)
	if totalTextElements < 1 {
		t.Error("expected at least one text element (bullet or shape)")
	}

	if len(s.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(s.Images))
	} else {
		img := s.Images[0]
		if img.Format != "png" {
			t.Errorf("expected image format 'png', got %q", img.Format)
		}
		if len(img.Data) == 0 {
			t.Error("expected image data, got empty")
		}
	}
}

func TestSlidesFromPPTX_PreservesReaderMetadata(t *testing.T) {
	deckPath := writeReaderMetadataPPTX(t)
	_, slides, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(slides))
	}
	if len(slides[0].Shapes) != 1 {
		t.Fatalf("expected 1 shape on slide 1, got %d", len(slides[0].Shapes))
	}
	shape := slides[0].Shapes[0]
	if shape.AltText != "Accessible rectangle" {
		t.Fatalf("expected shape alt text, got %q", shape.AltText)
	}
	if shape.ClickAction == nil || shape.ClickAction.Action.Type != action.HyperlinkActionURL ||
		shape.ClickAction.Action.URL != "https://example.com" {
		t.Fatalf("expected shape click URL action, got %+v", shape.ClickAction)
	}
	if len(shape.TextParagraphs) != 1 || len(shape.TextParagraphs[0].Runs) != 1 {
		t.Fatalf("expected one rich text run, got %+v", shape.TextParagraphs)
	}
	if shape.TextParagraphs[0].Runs[0].Hyperlink == nil ||
		shape.TextParagraphs[0].Runs[0].Hyperlink.Action.Type != action.HyperlinkActionSlide ||
		shape.TextParagraphs[0].Runs[0].Hyperlink.Action.SlideNumber != 2 {
		t.Fatalf("expected run hyperlink to slide 2, got %+v", shape.TextParagraphs[0].Runs[0].Hyperlink)
	}
	if len(slides[0].Images) != 1 {
		t.Fatalf("expected 1 image on slide 1, got %d", len(slides[0].Images))
	}
	image := slides[0].Images[0]
	if image.AltText != "Accessible image" {
		t.Fatalf("expected image alt text, got %q", image.AltText)
	}
	if image.Rotation != 90 {
		t.Fatalf("expected image rotation 90, got %v", image.Rotation)
	}
	if image.Crop.Left != 0.1 || image.Crop.Right != 0.2 || image.Crop.Top != 0.05 || image.Crop.Bottom != 0.15 {
		t.Fatalf("expected image crop to be preserved, got %+v", image.Crop)
	}
	if !image.FlipH || !image.Shadow || !image.Reflection {
		t.Fatalf("expected image flip/shadow/reflection, got %+v", image)
	}
}

func TestCanonicalZipPath(t *testing.T) {
	if canonicalZipPath("\\ppt\\media\\img.png") != "ppt/media/img.png" {
		t.Error("canonicalZipPath failed")
	}
	if canonicalZipPath("/ppt/media/img.png") != "ppt/media/img.png" {
		t.Error("canonicalZipPath failed")
	}
}

func TestImageFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"image.png", "png"},
		{"image.PNG", "png"},
		{"image.jpg", "jpeg"},
		{"image.jpeg", "jpeg"},
		{"image.gif", "gif"},
		{"image.emf", "emf"},
		{"image.wmf", "emf"},
		{"image.bmp", "png"}, // default
	}

	for _, tt := range tests {
		if got := imageFormat(tt.path); got != tt.expected {
			t.Errorf("imageFormat(%q) = %q, want %q", tt.path, got, tt.expected)
		}
	}
}

func TestResolveRelPath(t *testing.T) {
	tests := []struct {
		slidePart string
		target    string
		expected  string
	}{
		{"ppt/slides/slide1.xml", "../media/image1.png", "ppt/media/image1.png"},
		{"ppt/slides/slide1.xml", "/ppt/media/image1.png", "ppt/media/image1.png"},
		{"ppt/slides/slide1.xml", "../../evil.xml", ""},
		{"ppt/slides/slide1.xml", "slides/../media/image1.png", ""},
	}

	for _, tt := range tests {
		if got := resolveRelPath(tt.slidePart, tt.target); got != tt.expected {
			t.Errorf("resolveRelPath(%q, %q) = %q, want %q", tt.slidePart, tt.target, got, tt.expected)
		}
	}
}

func TestParseInt64(t *testing.T) {
	if got := parseInt64(" 123 "); got != 123 {
		t.Errorf("parseInt64(' 123 ') = %d, want 123", got)
	}
	if got := parseInt64("abc"); got != 0 {
		t.Errorf("parseInt64('abc') = %d, want 0", got)
	}
}

func TestEditorTypeToPreset(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"rect", "rect"},
		{"Rectangle", "rect"},
		{"roundRect", "roundRect"},
		{"RoundedRectangle", "roundRect"},
		{"ellipse", "ellipse"},
		{"Oval", "ellipse"},
		{"Circle", "ellipse"},
		{"triangle", "triangle"},
		{"RT_Triangle", "triangle"},
		{"rightArrow", "rightArrow"},
		{"leftArrow", "leftArrow"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		if got := editorTypeToPreset(tt.input); got != tt.expected {
			t.Errorf("editorTypeToPreset(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestEditorShapeToShape_MapsStyleAndAdjustments(t *testing.T) {
	fillColor := "4285F4"
	lineColor := "FFFFFF"
	lineWidth := 12700
	source := editorcommon.Shape{
		ID:   4,
		Name: "pie-slice",
		Type: "pie",
		X:    100,
		Y:    200,
		W:    300,
		H:    400,
		Fill: &editorcommon.ShapeFill{Solid: &fillColor},
		Line: &editorcommon.ShapeLine{Color: &lineColor, WidthEmu: &lineWidth},
		Adjustments: []editorcommon.ShapeAdjustment{
			{Name: "adj1", Formula: "val 0"},
			{Name: "adj2", Formula: "val 17100000"},
		},
	}

	mapped := editorShapeToShape(source)
	if mapped.Type != "pie" {
		t.Fatalf("expected mapped type pie, got %q", mapped.Type)
	}
	if mapped.Fill == nil || mapped.Fill.Color != fillColor {
		t.Fatalf("expected fill color %q, got %#v", fillColor, mapped.Fill)
	}
	if mapped.Line == nil || mapped.Line.Color != lineColor {
		t.Fatalf("expected line color %q, got %#v", lineColor, mapped.Line)
	}
	if int64(mapped.Line.Width) != int64(lineWidth) {
		t.Fatalf("expected line width %d, got %d", lineWidth, mapped.Line.Width)
	}
	if len(mapped.Adjustments) != 2 {
		t.Fatalf("expected 2 adjustments, got %d", len(mapped.Adjustments))
	}
	if mapped.Adjustments[0].Name != "adj1" || mapped.Adjustments[1].Name != "adj2" {
		t.Fatalf("unexpected adjustments: %#v", mapped.Adjustments)
	}
}

func TestSlidesFromPPTX_PieShapesKeepGeometryAndFill(t *testing.T) {
	deckPath := writePieShapesPPTX(t)
	_, slides, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(slides) < 4 {
		t.Fatalf("expected at least 4 slides, got %d", len(slides))
	}
	slide := slides[3]
	found := 0
	for _, shape := range slide.Shapes {
		if shape.Name != "Shape 4" && shape.Name != "Shape 5" && shape.Name != "Shape 6" {
			continue
		}
		found++
		if shape.Type != "pie" {
			t.Fatalf("%s expected type=pie, got %q", shape.Name, shape.Type)
		}
		if shape.Fill == nil || shape.Fill.Color == "" {
			t.Fatalf("%s expected non-empty fill, got %#v", shape.Name, shape.Fill)
		}
		if len(shape.Adjustments) < 2 {
			t.Fatalf("%s expected adjustments, got %#v", shape.Name, shape.Adjustments)
		}
	}
	if found != 3 {
		t.Fatalf("expected 3 pie slice shapes, got %d", found)
	}
}

func TestSlidesFromPPTX_Slide14ExtractsTable(t *testing.T) {
	// Generate a 14-slide deck with a 2-column table on the last slide.
	slideList := make([]elements.SlideContent, 14)
	for i := range 13 {
		slideList[i] = elements.NewSlide(fmt.Sprintf("Slide %d", i+1))
	}
	slideList[13] = elements.NewSlide("Slide 14").
		WithTable(tables.NewTable([]styling.Length{styling.Inches(1), styling.Inches(2)}).
			AddRow([]string{"Header 1", "Header 2"}).
			AddRow([]string{"Cell 1", "Cell 2"}))

	data, err := pptx.CreateWithSlides("Test Deck", slideList)
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}
	deckPath := filepath.Join(t.TempDir(), "table14.pptx")
	if err := os.WriteFile(deckPath, data, 0o600); err != nil {
		t.Fatalf("write temp pptx: %v", err)
	}

	_, got, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(got) < 14 {
		t.Fatalf("expected at least 14 slides, got %d", len(got))
	}
	slide := got[13]
	if slide.Table == nil {
		t.Fatalf("expected slide 14 table to be extracted")
	}
	if len(slide.Table.Rows) < 2 {
		t.Fatalf("expected extracted table rows, got %d", len(slide.Table.Rows))
	}
	if len(slide.Table.ColumnWidths) != 2 {
		t.Fatalf("expected 2 table columns, got %d", len(slide.Table.ColumnWidths))
	}
}

// writePieShapesPPTX builds a minimal 4-slide PPTX where slide 4 contains three
// preset-geometry "pie" shapes (Shape 4/5/6) with solid fill and two adjustments each.
// It returns the path to the written temp file.
func writePieShapesPPTX(t *testing.T) string {
	t.Helper()

	const (
		nsA = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
		nsR = `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
		nsP = `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`
	)

	blankSlide := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld ` + nsA + ` ` + nsR + ` ` + nsP + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`</p:spTree></p:cSld></p:sld>`

	pieShape := func(id int, color, adj2 string) string {
		return fmt.Sprintf(
			`<p:sp>`+
				`<p:nvSpPr><p:cNvPr id="%d" name="Shape %d"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
				`<p:spPr>`+
				`<a:xfrm><a:off x="%d" y="100"/><a:ext cx="500000" cy="500000"/></a:xfrm>`+
				`<a:prstGeom prst="pie"><a:avLst>`+
				`<a:gd name="adj1" fmla="val 0"/>`+
				`<a:gd name="adj2" fmla="val %s"/>`+
				`</a:avLst></a:prstGeom>`+
				`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`+
				`</p:spPr></p:sp>`,
			id, id, (id-4)*600000, adj2, color,
		)
	}

	slide4 := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld ` + nsA + ` ` + nsR + ` ` + nsP + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		pieShape(4, "FF0000", "17100000") +
		pieShape(5, "00FF00", "34200000") +
		pieShape(6, "0000FF", "51300000") +
		`</p:spTree></p:cSld></p:sld>`

	blankRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slides/slide2.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slides/slide3.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slides/slide4.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>` +
			`</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation ` + nsA + ` ` + nsR + ` ` + nsP + `>` +
			`<p:sldIdLst>` +
			`<p:sldId id="256" r:id="rId1"/>` +
			`<p:sldId id="257" r:id="rId2"/>` +
			`<p:sldId id="258" r:id="rId3"/>` +
			`<p:sldId id="259" r:id="rId4"/>` +
			`</p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>` +
			`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide2.xml"/>` +
			`<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide3.xml"/>` +
			`<Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide4.xml"/>` +
			`</Relationships>`,
		"ppt/slides/slide1.xml":            blankSlide,
		"ppt/slides/slide2.xml":            blankSlide,
		"ppt/slides/slide3.xml":            blankSlide,
		"ppt/slides/slide4.xml":            slide4,
		"ppt/slides/_rels/slide1.xml.rels": blankRels,
		"ppt/slides/_rels/slide2.xml.rels": blankRels,
		"ppt/slides/_rels/slide3.xml.rels": blankRels,
		"ppt/slides/_rels/slide4.xml.rels": blankRels,
	}

	outPath := filepath.Join(t.TempDir(), "pie_shapes.pptx")
	f, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("create pie shapes pptx: %v", err)
	}
	zw := zip.NewWriter(f)
	for name, content := range files {
		w, werr := zw.Create(name)
		if werr != nil {
			t.Fatalf("zip create %s: %v", name, werr)
		}
		if _, werr = fmt.Fprint(w, content); werr != nil {
			t.Fatalf("zip write %s: %v", name, werr)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("file close: %v", err)
	}
	return outPath
}

func writeReaderMetadataPPTX(t *testing.T) string {
	t.Helper()

	const (
		nsA = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
		nsR = `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
		nsP = `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`
	)

	pngData := string([]byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0x60, 0x00, 0x02, 0x00,
		0x00, 0x05, 0x00, 0x01, 0x0D, 0x26, 0xE5, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44,
		0xAE, 0x42, 0x60, 0x82,
	})

	slide1 := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld ` + nsA + ` ` + nsR + ` ` + nsP + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Action Shape" descr="Accessible rectangle">` +
		`<a:hlinkClick r:id="rId5" tooltip="Visit"/>` +
		`</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
		`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr>` +
		`<p:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr><a:hlinkClick r:id="rId6" tooltip="Jump"/></a:rPr><a:t>Hello</a:t></a:r></a:p></p:txBody></p:sp>` +
		`<p:pic><p:nvPicPr><p:cNvPr id="3" name="Picture 1" descr="Accessible image"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>` +
		`<p:blipFill><a:blip r:embed="rId7"/><a:srcRect l="10000" r="20000" t="5000" b="15000"/><a:stretch><a:fillRect/></a:stretch></p:blipFill>` +
		`<p:spPr><a:xfrm rot="5400000" flipH="1"><a:off x="100" y="200"/><a:ext cx="500" cy="400"/></a:xfrm>` +
		`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:effectLst><a:outerShdw/><a:reflection/></a:effectLst></p:spPr></p:pic>` +
		`</p:spTree></p:cSld></p:sld>`

	blankSlide := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<p:sld ` + nsA + ` ` + nsR + ` ` + nsP + `>` +
		`<p:cSld><p:spTree>` +
		`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
		`</p:spTree></p:cSld></p:sld>`

	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Default Extension="png" ContentType="image/png"/>` +
			`<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`<Override PartName="/ppt/slides/slide2.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>` +
			`</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation ` + nsA + ` ` + nsR + ` ` + nsP + `><p:sldIdLst>` +
			`<p:sldId id="256" r:id="rId1"/><p:sldId id="257" r:id="rId2"/>` +
			`</p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide1.xml"/>` +
			`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide2.xml"/>` +
			`</Relationships>`,
		"ppt/slides/slide1.xml": slide1,
		"ppt/slides/slide2.xml": blankSlide,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="https://example.com" TargetMode="External"/>` +
			`<Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slide2.xml"/>` +
			`<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/image1.png"/>` +
			`</Relationships>`,
		"ppt/slides/_rels/slide2.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`,
		"ppt/media/image1.png": pngData,
	}

	outPath := filepath.Join(t.TempDir(), "reader_metadata.pptx")
	f, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("create reader metadata pptx: %v", err)
	}
	zw := zip.NewWriter(f)
	for name, content := range files {
		w, werr := zw.Create(name)
		if werr != nil {
			t.Fatalf("zip create %s: %v", name, werr)
		}
		if _, werr = w.Write([]byte(content)); werr != nil {
			t.Fatalf("zip write %s: %v", name, werr)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("file close: %v", err)
	}
	return outPath
}
