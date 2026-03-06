package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestTextParityFixture_AdvancedFrameParagraphAndRunMappings(t *testing.T) {
	basePath := writeDeckFixture(t, "text-parity-fixture.pptx", []elements.SlideContent{
		elements.NewSlide("Text Parity"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeID, err := editor.AddShape(0, "rect", 120, 120, 2200, 1200)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}

	addr := "https://example.com/docs"
	allCaps := true
	smallCaps := true
	textFrameRotation := 45.0
	textFrameColumns := 2
	indent := 228600
	hanging := 114300
	runs := []common.TextRun{
		{
			Text:      "HELLO",
			AllCaps:   &allCaps,
			Hyperlink: &common.Hyperlink{Address: &addr},
		},
		{
			Text:      "world",
			SmallCaps: &smallCaps,
		},
	}
	updates := common.ShapeUpdate{
		Runs: &runs,
		TextFrame: &common.TextFrame{
			MarginLeft:    intPtr(12700),
			MarginRight:   intPtr(25400),
			MarginTop:     intPtr(38100),
			MarginBottom:  intPtr(50800),
			WordWrap:      boolPtr(false),
			AutoFitType:   strPtr("normal"),
			VerticalAlign: strPtr("bottom"),
			Orientation:   strPtr("vert270"),
			Columns:       &textFrameColumns,
			Rotation:      &textFrameRotation,
		},
		Paragraph: &common.Paragraph{
			Indent:  &indent,
			Hanging: &hanging,
		},
	}
	if err := editor.UpdateShape(0, shapeID, updates); err != nil {
		t.Fatalf("update shape: %v", err)
	}

	partPath := editor.slides[0].Part
	slideData, ok := editor.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide part %q", partPath)
	}
	slideXML := string(slideData)
	expectedTokens := []string{
		`lIns="12700"`,
		`rIns="25400"`,
		`tIns="38100"`,
		`bIns="50800"`,
		`wrap="none"`,
		`anchor="b"`,
		`vert="vert270"`,
		`numCol="2"`,
		`rot="2700000"`,
		`<a:normAutoFit`,
		`marL="228600"`,
		`indent="-114300"`,
		`cap="all"`,
		`cap="small"`,
		`hlinkClick`,
	}
	for _, token := range expectedTokens {
		if !strings.Contains(slideXML, token) {
			t.Fatalf("expected token %q in slide xml: %s", token, slideXML)
		}
	}
}

func intPtr(v int) *int {
	return &v
}

func boolPtr(v bool) *bool {
	return &v
}

func strPtr(v string) *string {
	return &v
}
