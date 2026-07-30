package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// rectShapeID is the id of the plain autoshape appended by newRectFixtureEditor.
const rectShapeID = 3

// newRectFixtureEditor opens a deck whose slide holds a picture plus a plain
// rectangle, and returns an editor over it. The rectangle is what these tests
// target: it exercises the non-picture path.
func newRectFixtureEditor(t *testing.T) *PresentationEditor {
	t.Helper()
	content := strings.Replace(
		pictureFixtureSlideXML(),
		`</p:pic></p:spTree>`,
		`</p:pic><p:sp><p:nvSpPr><p:cNvPr id="3" name="Rect 3"/><p:cNvSpPr/><p:nvPr/>`+
			`</p:nvSpPr><p:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="1000" cy="1000"/>`+
			`</a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:sp></p:spTree>`,
		1,
	)
	pptxPath := createPictureFixtureWithSlideXML(t, []byte(content), testutil.TinyPNG())

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	return ed
}

func rectSlideXML(t *testing.T, ed *PresentationEditor) string {
	t.Helper()
	slideBytes, ok := ed.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide part not found after update")
	}
	return string(slideBytes)
}

// Rotation and flips are attributes of <a:xfrm>, which every shape kind
// carries, so they must apply to ordinary autoshapes and not only to pictures.
// They used to be lumped in with cropping and rejected on anything but a
// picture, which made rotating a rectangle impossible.
func TestUpdateShapeRotatesNonPictureShape(t *testing.T) {
	rotation45 := 45.0
	rotation90 := 90.0
	enabled := true
	disabled := false

	tests := []struct {
		name    string
		updates common.ShapeUpdate
		want    []string
		absent  []string
	}{
		{
			name:    "rotation",
			updates: common.ShapeUpdate{Rotation: &rotation45},
			want:    []string{`rot="2700000"`},
		},
		{
			name:    "flip horizontal",
			updates: common.ShapeUpdate{FlipH: &enabled},
			want:    []string{`flipH="1"`},
		},
		{
			name: "rotation with both flips",
			updates: common.ShapeUpdate{
				Rotation: &rotation90,
				FlipH:    &enabled,
				FlipV:    &enabled,
			},
			want: []string{`rot="5400000"`, `flipH="1"`, `flipV="1"`},
		},
		{
			name:    "flip disabled leaves no attribute",
			updates: common.ShapeUpdate{FlipH: &disabled},
			absent:  []string{`flipH="1"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := newRectFixtureEditor(t)
			defer func() { _ = ed.Close() }()

			if err := ed.UpdateShape(0, rectShapeID, tt.updates); err != nil {
				t.Fatalf("UpdateShape on a non-picture shape: %v", err)
			}

			xmlData := rectSlideXML(t, ed)
			for _, want := range tt.want {
				if !strings.Contains(xmlData, want) {
					t.Errorf("slide XML missing %s\n%s", want, xmlData)
				}
			}
			for _, absent := range tt.absent {
				if strings.Contains(xmlData, absent) {
					t.Errorf("slide XML unexpectedly contains %s", absent)
				}
			}
		})
	}
}

func TestPictureOnlyVersusTransformUpdateFields(t *testing.T) {
	rotation := 10.0
	flip := true
	text := "hi"

	cases := []struct {
		name        string
		updates     common.ShapeUpdate
		pictureOnly bool
		transform   bool
	}{
		{"crop", common.ShapeUpdate{Crop: &common.ImageCrop{}}, true, false},
		{"rotation", common.ShapeUpdate{Rotation: &rotation}, false, true},
		{"flipH", common.ShapeUpdate{FlipH: &flip}, false, true},
		{"flipV", common.ShapeUpdate{FlipV: &flip}, false, true},
		{"text only", common.ShapeUpdate{Text: &text}, false, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPictureOnlyUpdateFields(tt.updates); got != tt.pictureOnly {
				t.Errorf("hasPictureOnlyUpdateFields = %v, want %v", got, tt.pictureOnly)
			}
			if got := hasTransformAttrUpdateFields(tt.updates); got != tt.transform {
				t.Errorf("hasTransformAttrUpdateFields = %v, want %v", got, tt.transform)
			}
		})
	}
}
