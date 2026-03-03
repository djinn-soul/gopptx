package elements

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestMaster_Methods(t *testing.T) {
	m := NewMaster().
		WithBackground(NewSolidBackground("FF0000")).
		AddShape(shapes.NewRectangle(0, 0, 1, 1)).
		AddImage(shapes.Image{Data: []byte("fake"), Format: "png"}).
		WithFooter("Footer").
		WithColorMapping("lt1", "dk1").
		WithTitleStyle([]TextLevelStyle{{Level: 0, Bold: true}}).
		WithBodyStyle([]TextLevelStyle{{Level: 0, Italic: true}}).
		WithOtherStyle([]TextLevelStyle{{Level: 0, Color: "0000FF"}})

	if m.FooterText != "Footer" { t.Error("WithFooter failed") }
	if m.Background == nil { t.Error("WithBackground failed") }
	if len(m.Shapes) != 1 { t.Error("AddShape failed") }
	if len(m.Images) != 1 { t.Error("AddImage failed") }
	if m.ColorMapping.BG1 != "lt1" { t.Error("WithColorMapping failed") }
	if len(m.TxStyles.TitleStyle) != 1 { t.Error("TxStyles failed") }
}

func TestNotesMaster_Methods(t *testing.T) {
	nm := NewNotesMaster().
		WithHeader("Header").
		WithFooter("Footer").
		WithDateTime(true).
		WithSlideNumber(true).
		WithBodyStyle([]TextLevelStyle{{Level: 0, Bold: true}})

	if nm.HeaderText != "Header" || nm.FooterText != "Footer" || !nm.ShowDateTime || !nm.ShowSlideNum {
		t.Error("NotesMaster options failed")
	}
}
