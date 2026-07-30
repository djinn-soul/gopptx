package editor

import (
	"path/filepath"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// Upstream #1020: <a:custGeom> was emitted by AddFreeformShape and never parsed,
// so a freeform's points could not be read back through the public API.
func TestGetShapesReadsFreeformGeometry(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-freeform-read.pptx", []elements.SlideContent{
		elements.NewSlide("Freeform"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	points := []freeformPoint{{X: 100, Y: 200}, {X: 900, Y: 200}, {X: 500, Y: 800}}
	shapeID, err := ed.AddFreeformShape(0, points, true)
	if err != nil {
		t.Fatalf("add freeform: %v", err)
	}

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}

	var freeform *common.FreeformGeometry
	for _, shape := range shapes {
		if shape.ID == shapeID {
			freeform = shape.Freeform
		}
	}
	if freeform == nil {
		t.Fatalf("expected freeform geometry on shape %d, got shapes %+v", shapeID, shapes)
	}
	if len(freeform.Paths) != 1 {
		t.Fatalf("expected one path, got %d", len(freeform.Paths))
	}

	path := freeform.Paths[0]
	if path.W != 800 || path.H != 600 {
		t.Fatalf("expected path space 800x600, got %dx%d", path.W, path.H)
	}
	wantKinds := []string{
		common.FreeformSegmentMoveTo,
		common.FreeformSegmentLineTo,
		common.FreeformSegmentLineTo,
		common.FreeformSegmentClose,
	}
	if len(path.Segments) != len(wantKinds) {
		t.Fatalf("expected %d segments, got %+v", len(wantKinds), path.Segments)
	}
	for i, want := range wantKinds {
		if path.Segments[i].Type != want {
			t.Fatalf("segment %d: got %q, want %q", i, path.Segments[i].Type, want)
		}
	}
	// The writer stores points local to the shape origin, so the first vertex
	// (100,200) becomes (0,0) and the last (500,800) becomes (400,600).
	first := path.Segments[0].Points
	if len(first) != 1 || first[0].X != 0 || first[0].Y != 0 {
		t.Fatalf("unexpected moveTo point: %+v", first)
	}
	last := path.Segments[2].Points
	if len(last) != 1 || last[0].X != 400 || last[0].Y != 600 {
		t.Fatalf("unexpected final lnTo point: %+v", last)
	}
	if len(path.Segments[3].Points) != 0 {
		t.Fatalf("close segment should carry no points, got %+v", path.Segments[3].Points)
	}
}

func TestParseFreeformKeepsGuideFormulasAndArcs(t *testing.T) {
	shapeXML := []byte(`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:nvSpPr><p:cNvPr id="7" name="Curve"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
		`<p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="100" cy="100"/></a:xfrm>` +
		`<a:custGeom><a:pathLst><a:path w="200" h="100" fill="none" stroke="0">` +
		`<a:moveTo><a:pt x="wd2" y="0"/></a:moveTo>` +
		`<a:cubicBezTo><a:pt x="10" y="10"/><a:pt x="20" y="20"/><a:pt x="30" y="30"/></a:cubicBezTo>` +
		`<a:arcTo wR="50" hR="25" stAng="0" swAng="10800000"/>` +
		`</a:path></a:pathLst></a:custGeom></p:spPr></p:sp>`)

	parsed, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parse shape: %v", err)
	}
	if parsed.Freeform == nil || len(parsed.Freeform.Paths) != 1 {
		t.Fatalf("expected one freeform path, got %+v", parsed.Freeform)
	}

	path := parsed.Freeform.Paths[0]
	if path.Fill != "none" || path.Stroke == nil || *path.Stroke {
		t.Fatalf("expected fill=none and stroke=false, got fill=%q stroke=%v", path.Fill, path.Stroke)
	}
	if len(path.Segments) != 3 {
		t.Fatalf("expected three segments, got %+v", path.Segments)
	}

	// A guide formula is preserved verbatim rather than reported as coordinate 0.
	moveTo := path.Segments[0].Points[0]
	if moveTo.XExpr != "wd2" || moveTo.X != 0 {
		t.Fatalf("expected preserved x formula, got %+v", moveTo)
	}
	if len(path.Segments[1].Points) != 3 {
		t.Fatalf("cubicBezTo should carry three points, got %+v", path.Segments[1].Points)
	}

	arc := path.Segments[2]
	if arc.WidthRadius == nil || *arc.WidthRadius != 50 || arc.HeightRadius == nil || *arc.HeightRadius != 25 {
		t.Fatalf("unexpected arc radii: %+v", arc)
	}
	if arc.SwingAngle == nil || *arc.SwingAngle != 180 {
		t.Fatalf("expected 180 degree swing angle, got %+v", arc.SwingAngle)
	}
}

// Upstream #1020 also asks for transparency on fills other than a solid fill.
func TestParseGradientStopTransparency(t *testing.T) {
	shapeXML := []byte(`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:nvSpPr><p:cNvPr id="3" name="Faded"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
		`<p:spPr><a:gradFill><a:gsLst>` +
		`<a:gs pos="0"><a:srgbClr val="FF0000"><a:alpha val="25000"/></a:srgbClr></a:gs>` +
		`<a:gs pos="100000"><a:srgbClr val="0000FF"/></a:gs>` +
		`</a:gsLst></a:gradFill></p:spPr></p:sp>`)

	parsed, err := parseShapeProperties(shapeXML)
	if err != nil {
		t.Fatalf("parse shape: %v", err)
	}
	if parsed.Fill == nil || parsed.Fill.Gradient == nil || len(parsed.Fill.Gradient.Stops) != 2 {
		t.Fatalf("expected two gradient stops, got %+v", parsed.Fill)
	}
	first := parsed.Fill.Gradient.Stops[0]
	if first.Transparency == nil || *first.Transparency != 0.75 {
		t.Fatalf("expected 0.75 transparency on the first stop, got %+v", first.Transparency)
	}
	if parsed.Fill.Gradient.Stops[1].Transparency != nil {
		t.Fatalf("a stop without <a:alpha> should report no transparency")
	}
}

// Upstream #1020: an autoshape filled with an image exposed neither the fill nor
// the image, because only <p:pic> shapes were treated as carrying an image.
func TestGetShapesReadsPictureFill(t *testing.T) {
	fixture := filepath.Join(t.TempDir(), "shape-picture-fill.pptx")
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Default Extension="png" ContentType="image/png"/>` +
			`<Override PartName="/ppt/presentation.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>` +
			`<Override PartName="/ppt/slides/slide1.xml" ` +
			`ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>` +
			`</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" ` +
			`Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" ` +
			`Target="ppt/presentation.xml"/>` +
			`</Relationships>`,
		"ppt/presentation.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:sldIdLst><p:sldId id="256" r:id="rId1"/></p:sldIdLst></p:presentation>`,
		"ppt/_rels/presentation.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" ` +
			`Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" ` +
			`Target="slides/slide1.xml"/></Relationships>`,
		"ppt/slides/slide1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
			`<p:cSld><p:spTree>` +
			`<p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` +
			`<p:sp><p:nvSpPr><p:cNvPr id="2" name="Filled Shape"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>` +
			`<p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="100" cy="100"/></a:xfrm>` +
			`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>` +
			`<a:blipFill><a:blip r:embed="rId9"/><a:srcRect l="10000" t="0" r="20000" b="0"/>` +
			`<a:stretch><a:fillRect/></a:stretch></a:blipFill>` +
			`</p:spPr></p:sp>` +
			`</p:spTree></p:cSld></p:sld>`,
		"ppt/slides/_rels/slide1.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId9" ` +
			`Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" ` +
			`Target="../media/image1.png"/></Relationships>`,
		"ppt/media/image1.png": string(testutil.TinyPNG()),
	}
	parts := make(map[string][]byte, len(files))
	for name, body := range files {
		parts[name] = []byte(body)
	}
	if err := writeZipFixtureBytes(fixture, parts); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	ed, err := OpenPresentationEditor(fixture)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	if len(shapes) != 1 {
		t.Fatalf("expected one shape, got %d", len(shapes))
	}

	fill := shapes[0].Fill
	if fill == nil || fill.Picture == nil {
		t.Fatalf("expected picture fill on shape, got %+v", fill)
	}
	if fill.Picture.RelID != "rId9" {
		t.Fatalf("expected rel id rId9, got %q", fill.Picture.RelID)
	}
	if fill.Picture.ImagePart != "ppt/media/image1.png" {
		t.Fatalf("expected resolved image part, got %q", fill.Picture.ImagePart)
	}
	if fill.Picture.Mode != "stretch" {
		t.Fatalf("expected stretch mode, got %q", fill.Picture.Mode)
	}
	if fill.Picture.Crop == nil || fill.Picture.Crop.Left != 0.1 || fill.Picture.Crop.Right != 0.2 {
		t.Fatalf("unexpected crop: %+v", fill.Picture.Crop)
	}

	// The same shape now answers image metadata queries, which previously
	// required a <p:pic>.
	meta, err := ed.GetImageMetadata(0, 2)
	if err != nil {
		t.Fatalf("image metadata for picture-filled shape: %v", err)
	}
	if meta == nil {
		t.Fatalf("expected image metadata for picture-filled shape")
	}
}
