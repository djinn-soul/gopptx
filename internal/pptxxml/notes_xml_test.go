package pptxxml

import (
	"strings"
	"testing"
)

func TestNotesMasterHidesDateTimeAndSlideNumberWhenDisabled(t *testing.T) {
	xml := NotesMaster(&NotesMasterSpec{
		ShowDateTime: false,
		ShowSlideNum: false,
	})
	if strings.Contains(xml, `type="dt"`) {
		t.Fatalf("did not expect date placeholder when ShowDateTime=false")
	}
	if strings.Contains(xml, `type="sldNum"`) {
		t.Fatalf("did not expect slide number placeholder when ShowSlideNum=false")
	}
	if strings.Contains(xml, `type="slidenum"`) {
		t.Fatalf("did not expect slidenum field when ShowSlideNum=false")
	}
}

func TestNotesMasterBackgroundSolidAndGradient(t *testing.T) {
	solid := NotesMaster(&NotesMasterSpec{
		Background: &SlideBackgroundSpec{
			Type:      "solid",
			SolidFill: &ShapeFillSpec{Color: "F0F0F0"},
		},
	})
	if !strings.Contains(solid, `<a:solidFill>`) || !strings.Contains(solid, `val="F0F0F0"`) {
		t.Fatalf("expected solid background fill in notes master xml")
	}

	gradient := NotesMaster(&NotesMasterSpec{
		Background: &SlideBackgroundSpec{
			Type: "gradient",
			GradientFill: &ShapeGradientFillSpec{
				Type: "linear",
				Stops: []ShapeGradientStopSpec{
					{PositionPct: 0, Color: "FFFFFF"},
					{PositionPct: 100, Color: "000000"},
				},
			},
		},
	})
	if !strings.Contains(gradient, `<a:gradFill`) {
		t.Fatalf("expected gradient background fill in notes master xml")
	}
}

func TestNotesMasterBackgroundPictureAndRelationships(t *testing.T) {
	xml := NotesMaster(&NotesMasterSpec{
		Background: &SlideBackgroundSpec{
			Type:        "picture",
			PictureFill: &ImageRef{RelID: "rId2"},
		},
	})
	if !strings.Contains(xml, `a:blip r:embed="rId2"`) {
		t.Fatalf("expected picture background embed rel in notes master xml")
	}

	rels := NotesMasterRelationships(1, []string{"../media/image1.png"})
	if !strings.Contains(rels, `Target="../theme/theme1.xml"`) {
		t.Fatalf("expected notes master rels to target theme1")
	}
	if !strings.Contains(rels, `Id="rId2"`) || !strings.Contains(rels, `/relationships/image"`) {
		t.Fatalf("expected image relationship in notes master rels")
	}
}
