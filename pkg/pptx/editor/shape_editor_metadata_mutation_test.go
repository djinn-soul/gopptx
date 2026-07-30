package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestUpdateShapeClearsAltTextAndTitle(t *testing.T) {
	ed := newRectFixtureEditor(t)
	defer func() { _ = ed.Close() }()

	altText := "Accessible rectangle"
	title := "Rectangle title"
	if err := ed.UpdateShape(0, rectShapeID, common.ShapeUpdate{
		AltText: &altText,
		Title:   &title,
	}); err != nil {
		t.Fatalf("set shape metadata: %v", err)
	}
	empty := ""
	if err := ed.UpdateShape(0, rectShapeID, common.ShapeUpdate{
		AltText: &empty,
		Title:   &empty,
	}); err != nil {
		t.Fatalf("clear shape metadata: %v", err)
	}
	xml := rectSlideXML(t, ed)
	if strings.Contains(xml, `descr=`) || strings.Contains(xml, `title=`) {
		t.Fatalf("cleared metadata attributes remain in slide XML: %s", xml)
	}
}
