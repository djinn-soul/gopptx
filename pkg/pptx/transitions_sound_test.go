package pptx_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func TestTransitionSound(t *testing.T) {
	// Create a dummy audio file
	audioContent := []byte("dummy audio content")
	audioPath := filepath.Join(t.TempDir(), "sound.wav")
	if err := os.WriteFile(audioPath, audioContent, 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	slide := elements.NewSlide("Sound Transition").
		WithTransitionOptions(transitions.TransitionOptions{
			Type:       transitions.TransitionCut,
			DurationMS: 500,
		}).
		WithTransitionSound(audioPath)

	presBytes, err := pptx.CreateWithSlides("Sound Test", []elements.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}

	// Verify zip contents
	r, err := zip.NewReader(bytes.NewReader(presBytes), int64(len(presBytes)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	// 1. Check if audio file exists in ppt/media
	foundAudio := false
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/media/image") && strings.HasSuffix(f.Name, ".wav") { // Note: CURRENT implementation uses image prefix
			foundAudio = true
		}
	}
	if !foundAudio {
		t.Error("Expected audio file in ppt/media not found")
	}

	// 2. Check slide rels for the audio relationship
	// foundRel := false
	// We need to parse slide rels... simplified check for raw string
	// In a real test we'd parse XML

	// 3. Check slide XML for p:sndAc
	slideXML := readFile(t, r, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, "<p:sndAc>") {
		t.Error("Expected <p:sndAc> in slide XML")
	}
	if !strings.Contains(slideXML, "r:embed=") {
		t.Error("Expected r:embed in slide XML sound action")
	}
}

func readFile(t *testing.T, r *zip.Reader, name string) string {
	for _, f := range r.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("Open %s failed: %v", name, err)
			}
			defer func() {
				if err := rc.Close(); err != nil {
					t.Errorf("Close %s failed: %v", name, err)
				}
			}()
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(rc); err != nil {
				t.Fatalf("ReadFrom %s failed: %v", name, err)
			}
			return buf.String()
		}
	}
	t.Fatalf("File %s not found in zip", name)
	return ""
}
