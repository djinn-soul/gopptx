package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Video and audio existed only on the editor path: gopptx could insert media
// into an existing deck but could not build one that contained any.
func TestGeneratedDeckCarriesVideo(t *testing.T) {
	video := shapes.NewVideoFromBytes(
		[]byte("fake-mp4-bytes"), "mp4",
		styling.Inches(1), styling.Inches(1), styling.Inches(6), styling.Inches(3.4),
	).WithAutoPlay(true).WithLoop(true).WithAltText("Product demo")

	slide := elements.NewSlide("Demo")
	slide.Media = []shapes.Media{video}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})

	if _, ok := parts["ppt/media/media1.mp4"]; !ok {
		t.Fatal("the media part was not written")
	}
	slideXML := parts["ppt/slides/slide1.xml"]
	if !strings.Contains(slideXML, "<a:videoFile") {
		t.Fatalf("no videoFile element: %s", slideXML)
	}
	if !strings.Contains(slideXML, "p14:media") {
		t.Fatalf("no p14:media extension: %s", slideXML)
	}
	if !strings.Contains(slideXML, "<p:timing>") {
		t.Fatalf("autoplay needs a timing tree: %s", slideXML)
	}

	rels := parts["ppt/slides/_rels/slide1.xml.rels"]
	for _, want := range []string{"relationships/video", "office/2007/relationships/media"} {
		if !strings.Contains(rels, want) {
			t.Fatalf("slide rels missing %s: %s", want, rels)
		}
	}
	if !strings.Contains(parts["[Content_Types].xml"], `Extension="mp4"`) {
		t.Fatal("no content type declared for the media extension")
	}
}

func TestGeneratedDeckCarriesAudio(t *testing.T) {
	audio := shapes.NewAudioFromBytes([]byte("fake-mp3"), "mp3", styling.Inches(1), styling.Inches(1)).
		WithPlayAcrossSlides(true)

	slide := elements.NewSlide("Soundtrack")
	slide.Media = []shapes.Media{audio}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide, elements.NewSlide("Next")})

	if _, ok := parts["ppt/media/media1.mp3"]; !ok {
		t.Fatal("the audio part was not written")
	}
	if !strings.Contains(parts["ppt/slides/slide1.xml"], "<a:audioFile") {
		t.Fatalf("no audioFile element: %s", parts["ppt/slides/slide1.xml"])
	}
	// Playing across slides is a numSld greater than one on the timing node.
	if !strings.Contains(parts["ppt/slides/slide1.xml"], `numSld="2"`) {
		t.Fatalf("play-across-slides not applied: %s", parts["ppt/slides/slide1.xml"])
	}
}

// An online video links out instead of embedding, so it has no media part.
func TestOnlineVideoIsLinkedNotEmbedded(t *testing.T) {
	online := shapes.NewOnlineVideo(
		"https://example.com/watch?v=1",
		styling.Inches(1), styling.Inches(1), styling.Inches(6), styling.Inches(3.4),
	)
	slide := elements.NewSlide("Hosted")
	slide.Media = []shapes.Media{online}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})

	for name := range parts {
		if strings.HasPrefix(name, "ppt/media/media") {
			t.Fatalf("an online video should embed nothing, found %s", name)
		}
	}
	rels := parts["ppt/slides/_rels/slide1.xml.rels"]
	if !strings.Contains(rels, `TargetMode="External"`) {
		t.Fatalf("the link should be an external relationship: %s", rels)
	}
}

// A poster is the still shown before playback.
func TestVideoPosterIsWrittenAndReferenced(t *testing.T) {
	video := shapes.NewVideoFromBytes(
		[]byte("fake-mp4"), "mp4",
		styling.Inches(1), styling.Inches(1), styling.Inches(4), styling.Inches(2),
	).WithPoster([]byte("fake-png"), "png")

	slide := elements.NewSlide("Demo")
	slide.Media = []shapes.Media{video}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	if _, ok := parts["ppt/media/mediaPoster1.png"]; !ok {
		t.Fatal("the poster image was not written")
	}
	if !strings.Contains(parts["ppt/slides/slide1.xml"], "<a:blip r:embed=") {
		t.Fatalf("the poster is not referenced: %s", parts["ppt/slides/slide1.xml"])
	}
}
