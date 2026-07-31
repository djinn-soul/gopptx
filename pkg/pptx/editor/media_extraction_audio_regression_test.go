package editor

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestListSlideMediaClassifiesGenericAudioRelationship(t *testing.T) {
	basePath := writeDeckFixture(t, "media-audio-extraction.pptx", []elements.SlideContent{
		elements.NewSlide("Audio"),
	})
	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if _, err := editor.AddAudio(0, []byte("audio-bytes"), "audio/mpeg", 0, 0, 100, 100); err != nil {
		t.Fatalf("add audio: %v", err)
	}
	media, err := editor.ListSlideMedia(0)
	if err != nil {
		t.Fatalf("list media: %v", err)
	}

	audioCount := 0
	for _, ref := range media {
		if ref.Kind == MediaKindVideo {
			t.Fatalf("audio relationship misclassified as video: %+v", ref)
		}
		if ref.Kind == MediaKindAudio {
			audioCount++
		}
	}
	if audioCount < 2 {
		t.Fatalf("expected audio and generic-media relationships as audio, got %+v", media)
	}
}
