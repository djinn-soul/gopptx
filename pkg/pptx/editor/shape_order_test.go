package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestMoveShape(t *testing.T) {
	xmlContent := `
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:cSld>
	<p:spTree>
	  <p:nvGrpSpPr/>
	  <p:grpSpPr/>
	  <p:sp>
		<p:nvSpPr><p:cNvPr id="2" name="Shape A"/></p:nvSpPr>
		<p:spPr>
			<a:xfrm><a:off x="100" y="100"/><a:ext cx="500" cy="50"/></a:xfrm>
		</p:spPr>
		<p:txBody><p:p><p:r><a:t>Hello</a:t></p:r></p:p></p:txBody>
	  </p:sp>
	  <p:sp>
		<p:nvSpPr><p:cNvPr id="3" name="Shape B"/></p:nvSpPr>
		<p:spPr>
			<a:xfrm><a:off x="100" y="100"/><a:ext cx="500" cy="50"/></a:xfrm>
		</p:spPr>
		<p:txBody><p:p><p:r><a:t>World</a:t></p:r></p:p></p:txBody>
	  </p:sp>
	  <p:sp>
		<p:nvSpPr><p:cNvPr id="4" name="Shape C"/></p:nvSpPr>
		<p:spPr>
			<a:xfrm><a:off x="100" y="100"/><a:ext cx="500" cy="50"/></a:xfrm>
		</p:spPr>
		<p:txBody><p:p><p:r><a:t>Test</a:t></p:r></p:p></p:txBody>
	  </p:sp>
	</p:spTree>
  </p:cSld>
</p:sld>`

	// Helper to recreate editor for each test
	setupEditor := func() *PresentationEditor {
		parts := NewPartStore()
		parts.Set("ppt/slides/slide1.xml", []byte(xmlContent))
		e := &PresentationEditor{
			parts: parts,
			slides: []common.EditorSlideRef{
				{Part: "ppt/slides/slide1.xml"},
			},
		}
		return e
	}

	t.Run("MoveToFront (Middle Shape)", func(t *testing.T) {
		e := setupEditor()
		err := e.MoveShapeToFront(0, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, _ := e.parts.Get("ppt/slides/slide1.xml")
		str := string(content)

		idx2 := strings.Index(str, `id="2"`)
		idx3 := strings.Index(str, `id="3"`)
		idx4 := strings.Index(str, `id="4"`)

		if idx2 >= idx4 || idx4 >= idx3 {
			t.Errorf("expected order: A, C, B. Got indices: A=%d, C=%d, B=%d", idx2, idx4, idx3)
		}
	})

	t.Run("MoveToBack (Middle Shape)", func(t *testing.T) {
		e := setupEditor()
		err := e.MoveShapeToBack(0, 3) // Move shape B (id 3) to back (beginning of drawing order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, _ := e.parts.Get("ppt/slides/slide1.xml")
		str := string(content)

		idx2 := strings.Index(str, `id="2"`)
		idx3 := strings.Index(str, `id="3"`)
		idx4 := strings.Index(str, `id="4"`)

		if idx3 >= idx2 || idx2 >= idx4 {
			t.Errorf("expected order: B, A, C. Got indices: B=%d, A=%d, C=%d", idx3, idx2, idx4)
		}
	})

	t.Run("MoveToFront (Already Front)", func(t *testing.T) {
		e := setupEditor()
		err := e.MoveShapeToFront(0, 4) // Shape C is already at the end
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, _ := e.parts.Get("ppt/slides/slide1.xml")
		str := string(content)

		idx2 := strings.Index(str, `id="2"`)
		idx3 := strings.Index(str, `id="3"`)
		idx4 := strings.Index(str, `id="4"`)

		if idx2 >= idx3 || idx3 >= idx4 {
			t.Errorf("expected order: A, B, C. Got indices A=%d, B=%d, C=%d", idx2, idx3, idx4)
		}
	})

	t.Run("MoveToBack (Already Back)", func(t *testing.T) {
		e := setupEditor()
		err := e.MoveShapeToBack(0, 2) // Shape A is already at the front/back (first in order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		content, _ := e.parts.Get("ppt/slides/slide1.xml")
		str := string(content)

		idx2 := strings.Index(str, `id="2"`)
		idx3 := strings.Index(str, `id="3"`)
		idx4 := strings.Index(str, `id="4"`)

		if idx2 >= idx3 || idx3 >= idx4 {
			t.Errorf("expected order: A, B, C. Got indices A=%d, B=%d, C=%d", idx2, idx3, idx4)
		}
	})
}

func TestMoveShapeRejectsMissingIDOnEmptyOrSingleShapeSlides(t *testing.T) {
	tests := []struct {
		name       string
		xmlContent string
	}{
		{
			name: "empty_shape_tree",
			xmlContent: `
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
	<p:spTree>
	  <p:nvGrpSpPr/>
	  <p:grpSpPr/>
	</p:spTree>
  </p:cSld>
</p:sld>`,
		},
		{
			name: "single_shape",
			xmlContent: `
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
	<p:spTree>
	  <p:nvGrpSpPr/>
	  <p:grpSpPr/>
	  <p:sp>
		<p:nvSpPr><p:cNvPr id="2" name="Only Shape"/></p:nvSpPr>
		<p:spPr/>
	  </p:sp>
	</p:spTree>
  </p:cSld>
</p:sld>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := NewPartStore()
			parts.Set("ppt/slides/slide1.xml", []byte(tc.xmlContent))
			e := &PresentationEditor{
				parts: parts,
				slides: []common.EditorSlideRef{
					{Part: "ppt/slides/slide1.xml"},
				},
			}

			errFront := e.MoveShapeToFront(0, 999)
			if errFront == nil || errFront.Error() != "shape with ID 999 not found" {
				t.Fatalf("MoveShapeToFront expected not-found error, got %v", errFront)
			}

			errBack := e.MoveShapeToBack(0, 999)
			if errBack == nil || errBack.Error() != "shape with ID 999 not found" {
				t.Fatalf("MoveShapeToBack expected not-found error, got %v", errBack)
			}
		})
	}
}
