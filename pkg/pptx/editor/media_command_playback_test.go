package editor

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCommandAddVideoWithPlaybackOptions(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-video-playback-test.pptx", []elements.SlideContent{
		elements.NewSlide("Video Playback Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	videoB64 := base64.StdEncoding.EncodeToString([]byte("video-bytes"))
	posterB64 := base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	req := fmt.Sprintf(
		`{"api_version":1,"request_id":"v1","op":"add_video","payload":{"slide_index":0,"x":120,"y":200,"w":900,"h":500,"mime_type":"video/mp4","data":"%s","poster_data":"%s","auto_play":true,"loop":true,"muted":true,"volume":90}}`,
		videoB64,
		posterB64,
	)
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_video with playback options failed: %s", resp)
	}

	slidePart := e.Slides()[0].PartName
	slideXML, ok := e.parts.Get(slidePart)
	if !ok {
		t.Fatalf("missing slide part %q", slidePart)
	}
	xmlText := string(slideXML)
	if !strings.Contains(xmlText, "<p:video>") {
		t.Fatalf("expected video timing node: %s", xmlText)
	}
	if !strings.Contains(xmlText, "<p14:media") {
		t.Fatalf("expected p14 media extension tag in shape nvPr: %s", xmlText)
	}
	if !strings.Contains(xmlText, `delay="0"`) {
		t.Fatalf("expected autoplay delay=0 in timing xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `repeatCount="indefinite"`) {
		t.Fatalf("expected loop repeatCount in timing xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `mute="1"`) {
		t.Fatalf("expected mute attr in timing xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, `vol="90000"`) {
		t.Fatalf("expected volume=90000 in timing xml: %s", xmlText)
	}
}

func TestCommandAddAudioWithIconAndPlaybackOptions(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-audio-playback-test.pptx", []elements.SlideContent{
		elements.NewSlide("S1").AddBullet("b1"),
		elements.NewSlide("S2").AddBullet("b2"),
		elements.NewSlide("S3").AddBullet("b3"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	audioB64 := base64.StdEncoding.EncodeToString(tinyWAVBytes())
	iconB64 := base64.StdEncoding.EncodeToString(testutil.TinyPNG())
	req := fmt.Sprintf(
		`{"api_version":1,"request_id":"a2","op":"add_audio","payload":{"slide_index":1,"x":120,"y":200,"w":900,"h":500,"mime_type":"audio/wav","data":"%s","icon_data":"%s","auto_play":true,"loop":true,"play_across_slides":true,"hide_during_show":true,"volume":80}}`,
		audioB64,
		iconB64,
	)
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_audio with icon+playback options failed: %s", resp)
	}

	slidePart := e.Slides()[1].PartName
	slideXML, ok := e.parts.Get(slidePart)
	if !ok {
		t.Fatalf("missing slide part %q", slidePart)
	}
	xmlText := string(slideXML)
	if !strings.Contains(xmlText, `<a:blip r:embed="`) {
		t.Fatalf("expected icon blip in slide xml: %s", xmlText)
	}
	if !strings.Contains(xmlText, "<p:audio>") {
		t.Fatalf("expected audio timing node: %s", xmlText)
	}
	if !strings.Contains(xmlText, "<p14:media") {
		t.Fatalf("expected p14 media extension tag in shape nvPr: %s", xmlText)
	}
	if !strings.Contains(xmlText, `numSld="2"`) {
		t.Fatalf("expected across-slides numSld=2 from slide index 1 of 3: %s", xmlText)
	}
	if !strings.Contains(xmlText, `showWhenStopped="0"`) {
		t.Fatalf("expected hide_during_show => showWhenStopped=0: %s", xmlText)
	}
	if !strings.Contains(xmlText, `vol="80000"`) {
		t.Fatalf("expected volume=80000 in timing xml: %s", xmlText)
	}
}
