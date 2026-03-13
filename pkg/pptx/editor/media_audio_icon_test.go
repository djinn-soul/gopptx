package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestAddAudioWithIconEmbedsImageIcon(t *testing.T) {
	e := newMediaEditorFixture()
	shapeID, err := e.AddAudioWithIcon(
		0,
		[]byte("audio-bytes"),
		testutil.TinyPNG(),
		"audio/wav",
		100,
		120,
		240,
		120,
	)
	if err != nil {
		t.Fatalf("AddAudioWithIcon failed: %v", err)
	}
	if shapeID == 0 {
		t.Fatal("expected non-zero shape id")
	}

	slideXML := string(getFixturePart(t, e, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, `<a:blip r:embed="`) {
		t.Fatalf("expected embedded icon blip in audio shape: %s", slideXML)
	}

	relsXML := string(getFixturePart(t, e, "ppt/slides/_rels/slide1.xml.rels"))
	if !strings.Contains(relsXML, common.RelTypeAudio) ||
		!strings.Contains(relsXML, common.RelTypeMedia) ||
		!strings.Contains(relsXML, common.RelTypeImage) {
		t.Fatalf("expected audio/media/image relationships in rels xml: %s", relsXML)
	}
}
