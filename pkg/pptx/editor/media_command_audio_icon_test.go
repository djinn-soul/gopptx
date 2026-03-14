package editor

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCommandAddAudioWithIconData(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-audio-icon-test.pptx", []elements.SlideContent{
		elements.NewSlide("Audio Icon Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	audioB64 := base64.StdEncoding.EncodeToString(tinyWAVBytes())
	iconB64 := base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	req := fmt.Sprintf(
		`{"api_version":1,"request_id":"a1","op":"add_audio","payload":{"slide_index":0,"x":120,"y":200,"w":900,"h":500,"mime_type":"audio/wav","data":"%s","icon_data":"%s"}}`,
		audioB64,
		iconB64,
	)
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_audio with icon_data failed: %s", resp)
	}

	slidePart := e.Slides()[0].PartName
	slideXML, ok := e.parts.Get(slidePart)
	if !ok {
		t.Fatalf("missing slide part %q", slidePart)
	}
	if !strings.Contains(string(slideXML), `<a:blip r:embed="`) {
		t.Fatalf("expected audio icon blip embed in slide xml: %s", string(slideXML))
	}
}

func tinyWAVBytes() []byte {
	return []byte{
		'R', 'I', 'F', 'F',
		0x25, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00,
		0x01, 0x00,
		0x01, 0x00,
		0x40, 0x1F, 0x00, 0x00,
		0x40, 0x1F, 0x00, 0x00,
		0x01, 0x00,
		0x08, 0x00,
		'd', 'a', 't', 'a',
		0x01, 0x00, 0x00, 0x00,
		0x80,
	}
}
