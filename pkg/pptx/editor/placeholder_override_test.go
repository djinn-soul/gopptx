package editor

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestUpdateSlidePlaceholderOverrideResolvesByName(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-name-resolve.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Title 1", phType: "title", phIndex: 0},
		{name: "Body 2", phType: "body", phIndex: 1},
	})

	size := 28
	x := styling.Inches(1)
	y := styling.Inches(2)
	cx := styling.Inches(4)
	cy := styling.Inches(2)
	slide := elements.NewSlide("Updated").WithPlaceholderOverride(
		shapes.PlaceholderTarget{Name: "Body 2"},
		shapes.PlaceholderOverrideOptions{
			X:  &x,
			Y:  &y,
			CX: &cx,
			CY: &cy,
			TextStyle: &shapes.PlaceholderTextStyle{
				SizePt: &size,
			},
		},
	)

	if err := editor.UpdateSlide(0, slide); err != nil {
		t.Fatalf("update slide: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "placeholder-name-resolve-out.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save: %v", err)
	}

	slideXML := string(readZipFileBytes(t, outPath, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, `<p:ph`) || !strings.Contains(slideXML, `type="body"`) ||
		!strings.Contains(slideXML, `idx="1"`) {
		t.Fatalf("expected body placeholder target in output xml")
	}
	if !strings.Contains(slideXML, `off x="914400" y="1828800"`) {
		t.Fatalf("expected geometry override x/y in output xml")
	}
	if !strings.Contains(slideXML, `ext cx="3657600" cy="1828800"`) {
		t.Fatalf("expected geometry override cx/cy in output xml")
	}
	if !strings.Contains(slideXML, `sz="2800"`) {
		t.Fatalf("expected text size override in output xml")
	}
}

func TestUpdateSlidePlaceholderOverrideNameNotFoundFails(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-name-missing.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Title 1", phType: "title", phIndex: 0},
	})

	slide := elements.NewSlide("Updated").WithPlaceholderOverride(
		shapes.PlaceholderTarget{Name: "Missing Placeholder"},
		shapes.PlaceholderOverrideOptions{},
	)
	err = editor.UpdateSlide(0, slide)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected name-not-found error, got %v", err)
	}
}

func TestUpdateSlidePlaceholderOverrideAmbiguousNameFails(t *testing.T) {
	basePath := writeDeckFixture(t, "placeholder-name-ambiguous.pptx", []elements.SlideContent{
		elements.NewSlide("Original"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	installPlaceholderSlideXML(editor, editor.slides[0].Part, []placeholderDef{
		{name: "Duplicate Name", phType: "body", phIndex: 1},
		{name: "Duplicate Name", phType: "body", phIndex: 2},
	})

	slide := elements.NewSlide("Updated").WithPlaceholderOverride(
		shapes.PlaceholderTarget{Name: "Duplicate Name"},
		shapes.PlaceholderOverrideOptions{},
	)
	err = editor.UpdateSlide(0, slide)
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous-name error, got %v", err)
	}
}

type placeholderDef struct {
	name    string
	phType  string
	phIndex int
}

func installPlaceholderSlideXML(editor *PresentationEditor, slidePart string, defs []placeholderDef) {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>
<p:spTree>
`)
	for i, def := range defs {
		id := i + 2
		b.WriteString(`<p:sp>
  <p:nvSpPr>
    <p:cNvPr id="`)
		b.WriteString(strconv.Itoa(id))
		b.WriteString(`" name="`)
		b.WriteString(def.name)
		b.WriteString(`"/>
    <p:cNvSpPr/>
    <p:nvPr>
      <p:ph type="`)
		b.WriteString(def.phType)
		if def.phIndex > 0 {
			b.WriteString(`" idx="`)
			b.WriteString(strconv.Itoa(def.phIndex))
		}
		b.WriteString(`"/>
    </p:nvPr>
  </p:nvSpPr>
</p:sp>
`)
	}
	b.WriteString(`</p:spTree>
</p:cSld>
</p:sld>`)
	editor.parts.Set(slidePart, []byte(b.String()))
}
