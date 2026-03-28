package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestVideoPlaybackOptionsBuilder(t *testing.T) {
	opts := NewAutoPlayVideoPlaybackOptions().
		WithLoop(true).
		WithMuted(true).
		WithVolume(120).
		WithStartTimeMS(1000).
		WithEndTimeMS(5000)

	if !opts.AutoPlay || !opts.LoopPlayback || !opts.Muted {
		t.Fatalf("unexpected video options flags: %+v", opts)
	}
	if opts.Volume != 100 {
		t.Fatalf("expected clamped volume 100, got %d", opts.Volume)
	}
	if opts.StartTimeMS == nil || *opts.StartTimeMS != 1000 {
		t.Fatalf("expected start time 1000ms, got %+v", opts.StartTimeMS)
	}
	if opts.EndTimeMS == nil || *opts.EndTimeMS != 5000 {
		t.Fatalf("expected end time 5000ms, got %+v", opts.EndTimeMS)
	}
}

func TestAudioPlaybackOptionsBuilder(t *testing.T) {
	opts := NewAutoPlayAudioPlaybackOptions().
		WithLoop(true).
		WithPlayAcrossSlides(true).
		WithVolume(110)

	if !opts.AutoPlay || !opts.LoopPlayback || !opts.PlayAcrossSlides {
		t.Fatalf("unexpected audio options flags: %+v", opts)
	}
	if opts.Volume != 100 {
		t.Fatalf("expected clamped volume 100, got %d", opts.Volume)
	}
}

func TestAddMediaWithPlaybackOptionsDelegatesToInsert(t *testing.T) {
	e := newMediaEditorFixture()
	if _, err := e.AddVideoWithPlaybackOptions(
		0,
		[]byte("video-bytes"),
		[]byte("poster-bytes"),
		"video/mp4",
		NewVideoPlaybackOptions().
			WithAltText(`Video "Alt" & Description`).
			WithMuted(true).
			WithLoop(true).
			WithVolume(120),
		10,
		20,
		300,
		200,
	); err != nil {
		t.Fatalf("AddVideoWithPlaybackOptions failed: %v", err)
	}

	if _, err := e.AddAudioWithPlaybackOptions(
		0,
		[]byte("audio-bytes"),
		"audio/mpeg",
		NewAutoPlayAudioPlaybackOptions().WithAltText("Audio alt text"),
		10,
		20,
		300,
		200,
	); err != nil {
		t.Fatalf("AddAudioWithPlaybackOptions failed: %v", err)
	}

	audioPath := filepath.Join(t.TempDir(), "playback-audio.wav")
	if err := os.WriteFile(audioPath, tinyWAVBytes(), 0o600); err != nil {
		t.Fatalf("write audio fixture: %v", err)
	}
	if _, err := e.AddAudioFromFileWithPlaybackOptions(
		0,
		audioPath,
		"audio/wav",
		NewAudioPlaybackOptions().WithLoop(true).WithVolume(85),
		10,
		20,
		300,
		200,
	); err != nil {
		t.Fatalf("AddAudioFromFileWithPlaybackOptions failed: %v", err)
	}

	slideXML := string(getFixturePart(t, e, "ppt/slides/slide1.xml"))
	if !strings.Contains(slideXML, `descr="Video &quot;Alt&quot; &amp; Description"`) {
		t.Fatalf("expected escaped video alt text descr attribute in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `descr="Audio alt text"`) {
		t.Fatalf("expected audio alt text descr attribute in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, "<p:timing>") {
		t.Fatalf("expected playback timing block in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, "<p:video>") || !strings.Contains(slideXML, "<p:audio>") {
		t.Fatalf("expected video+audio timing media nodes in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, "<p14:media") {
		t.Fatalf("expected p14 media extension tag in shape nvPr xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `repeatCount="indefinite"`) {
		t.Fatalf("expected loop timing repeatCount in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `delay="0"`) {
		t.Fatalf("expected autoplay delay=0 in slide xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `mute="1"`) {
		t.Fatalf("expected mute attribute in video timing node: %s", slideXML)
	}
	if !strings.Contains(slideXML, `vol="100000"`) {
		t.Fatalf("expected clamped video volume (100000) in timing xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `vol="85000"`) {
		t.Fatalf("expected audio volume (85000) in timing xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `numSld="1"`) {
		t.Fatalf("expected single-slide media timing node: %s", slideXML)
	}
	if !strings.Contains(slideXML, `r:embed="rId4"`) {
		t.Fatalf("expected media rel-id rId4 for video shape in slide xml: %s", slideXML)
	}
}

func TestAddAudioPlaybackOptions_PlayAcrossSlidesSetsNumSld(t *testing.T) {
	path := writeDeckFixture(t, "media-play-across.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
		elements.NewSlide("Slide 3"),
	})
	e, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	if _, err := e.AddAudioWithPlaybackOptions(
		1,
		[]byte("audio-bytes"),
		"audio/mpeg",
		NewAutoPlayAudioPlaybackOptions().WithPlayAcrossSlides(true).WithVolume(60),
		10,
		20,
		300,
		200,
	); err != nil {
		t.Fatalf("AddAudioWithPlaybackOptions failed: %v", err)
	}

	slidePart := e.Slides()[1].PartName
	slideXML := string(getFixturePart(t, e, slidePart))
	if !strings.Contains(slideXML, `numSld="2"`) {
		t.Fatalf("expected play-across-slides numSld=2 for slide index 1 of 3, xml: %s", slideXML)
	}
	if !strings.Contains(slideXML, `vol="60000"`) {
		t.Fatalf("expected audio volume 60000 in timing xml: %s", slideXML)
	}
}
