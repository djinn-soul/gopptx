package editor

import (
	"bytes"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// Upstream #1049: an embedded movie could not be found or read out, because
// only images had a listing.
func TestListSlideMediaAndExtract(t *testing.T) {
	basePath := writeDeckFixture(t, "media-extraction.pptx", []elements.SlideContent{
		elements.NewSlide("Media"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	movie := []byte("fake-mp4-bytes")
	if _, err := ed.AddVideo(0, movie, testutil.TinyPNG(), "video/mp4", 0, 0, 100, 100); err != nil {
		t.Fatalf("add video: %v", err)
	}
	if _, err := ed.AddImageFromBytes(0, testutil.TinyPNG(), "png", 0, 0, 100, 100, nil); err != nil {
		t.Fatalf("add image: %v", err)
	}

	media, err := ed.ListSlideMedia(0)
	if err != nil {
		t.Fatalf("list slide media: %v", err)
	}

	var video, image int
	var videoPart string
	for _, ref := range media {
		switch ref.Kind {
		case MediaKindVideo:
			video++
			videoPart = ref.PartPath
			if ref.ContentType != "video/mp4" {
				t.Fatalf("expected video/mp4 content type, got %q", ref.ContentType)
			}
			if ref.SizeBytes != len(movie) {
				t.Fatalf("expected %d bytes, got %d", len(movie), ref.SizeBytes)
			}
		case MediaKindImage:
			image++
		}
	}
	if video == 0 {
		t.Fatalf("expected the video to be listed, got %+v", media)
	}
	if image == 0 {
		t.Fatalf("expected the image to be listed too, got %+v", media)
	}

	extracted, err := ed.ExtractMedia(videoPart)
	if err != nil {
		t.Fatalf("extract media: %v", err)
	}
	if !bytes.Equal(extracted, movie) {
		t.Fatalf("extracted bytes differ from the embedded movie")
	}

	// The returned slice is a copy: mutating it must not corrupt the package.
	extracted[0] = 'X'
	again, err := ed.ExtractMedia(videoPart)
	if err != nil {
		t.Fatalf("extract media again: %v", err)
	}
	if !bytes.Equal(again, movie) {
		t.Fatalf("ExtractMedia handed out the package's own buffer")
	}
}

func TestExtractMediaRejectsMissingPart(t *testing.T) {
	basePath := writeDeckFixture(t, "media-missing.pptx", []elements.SlideContent{
		elements.NewSlide("Media"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if _, err := ed.ExtractMedia("ppt/media/nothing.mp4"); err == nil {
		t.Fatalf("expected an error for a part that is not in the package")
	}
}
