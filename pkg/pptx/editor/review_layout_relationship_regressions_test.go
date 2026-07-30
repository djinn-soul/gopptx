package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestLayoutAndMasterShapeAddRejectInvalidExtents(t *testing.T) {
	tests := []struct {
		name string
		add  func(*PresentationEditor) (int, error)
	}{
		{
			name: "layout shape",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddLayoutShape("ppt/slideLayouts/slideLayout1.xml", "rect", 0, 0, 0, 10)
			},
		},
		{
			name: "layout text box",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddLayoutTextBox("ppt/slideLayouts/slideLayout1.xml", "text", 0, 0, 10, -1)
			},
		},
		{
			name: "master shape",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddMasterShape("ppt/slideMasters/slideMaster1.xml", "rect", 0, 0, -1, 10)
			},
		},
		{
			name: "master text box",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddMasterTextBox("ppt/slideMasters/slideMaster1.xml", "text", 0, 0, 10, 0)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editor := newLayoutFixtureEditor(t)
			if _, err := tc.add(editor); err == nil {
				t.Fatal("expected invalid extents to be rejected")
			}
		})
	}
}

func TestLayoutAndMasterTextBoxesCarryTextBoxMarker(t *testing.T) {
	tests := []struct {
		name string
		part string
		add  func(*PresentationEditor) (int, error)
	}{
		{
			name: "layout",
			part: "ppt/slideLayouts/slideLayout1.xml",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddLayoutTextBox("ppt/slideLayouts/slideLayout1.xml", "text", 0, 0, 10, 10)
			},
		},
		{
			name: "master",
			part: "ppt/slideMasters/slideMaster1.xml",
			add: func(editor *PresentationEditor) (int, error) {
				return editor.AddMasterTextBox("ppt/slideMasters/slideMaster1.xml", "text", 0, 0, 10, 10)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editor := newLayoutFixtureEditor(t)
			if _, err := tc.add(editor); err != nil {
				t.Fatalf("add text box: %v", err)
			}
			part, _ := editor.parts.Get(tc.part)
			if !strings.Contains(string(part), `<p:cNvSpPr txBox="1"/>`) {
				t.Fatalf("text box marker missing: %s", part)
			}
		})
	}
}

func TestResolvePictureFillUsesExternalRelationshipID(t *testing.T) {
	shape := parsedShape{Fill: &common.ShapeFill{Picture: &common.PictureFill{
		External: "rId9",
	}}}
	resolvePictureFillPart(&shape, "ppt/slideLayouts/slideLayout1.xml", map[string]common.EditorRelationship{
		"rId9": {
			ID:         "rId9",
			Target:     "https://example.test/image.png",
			TargetMode: "External",
		},
	})

	if got := shape.Fill.Picture.External; got != "https://example.test/image.png" {
		t.Fatalf("external target = %q", got)
	}
}

func TestResolvePictureFillUsesOwningPartForInternalTarget(t *testing.T) {
	shape := parsedShape{Fill: &common.ShapeFill{Picture: &common.PictureFill{
		RelID: "rId7",
	}}}
	resolvePictureFillPart(&shape, "ppt/slideLayouts/slideLayout1.xml", map[string]common.EditorRelationship{
		"rId7": {ID: "rId7", Target: "../media/image1.png"},
	})

	if got := shape.Fill.Picture.ImagePart; got != "ppt/media/image1.png" {
		t.Fatalf("image part = %q", got)
	}
}
