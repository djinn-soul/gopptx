// examples/75-media-ole-api demonstrates media (audio/video) and OLE object embedding.
//
// Uses the PresentationEditor to insert:
//   - Audio from bytes (tiny WAV) with an icon
//   - Video from file if a fixture exists, otherwise logs a skip notice
//   - OLE object from bytes if a fixture exists, otherwise logs a skip notice
//
// Real fixture files can be placed in examples/assets/75/:
//   - sample.mp4  (video)
//   - poster.png  (video poster / OLE icon)
//   - sample_ole.bin (OLE binary)
//
// Run with: go run ./examples/75-media-ole-api/main.go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

const (
	outputDir  = "examples/output"
	outputFile = "75_media_ole_api.pptx"
	assetsDir  = "examples/assets/75"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-media-ole-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	basePath := filepath.Join(tmpDir, "base.pptx")

	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Media & OLE Embedding").
			AddBullet("Audio from bytes – embedded using AddAudioWithIcon.").
			AddBullet("Video from file – requires sample.mp4 + poster.png fixture.").
			AddBullet("OLE object – requires sample_ole.bin + poster.png fixture."),
	}
	if err := pptx.WriteFile(basePath, "Media & OLE API Demo", baseSlides); err != nil {
		return fmt.Errorf("write base deck: %w", err)
	}

	ed, err := pptx.OpenPresentationEditor(basePath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	// --- 1. Audio from bytes (always works – no fixture needed) ---
	audioOpts := editor.NewAutoPlayAudioPlaybackOptions().
		WithLoop(false).
		WithPlayAcrossSlides(false).
		WithVolume(80).
		WithAltText("Embedded WAV audio clip")

	if _, err := ed.AddAudioWithPlaybackOptions(
		0,
		tinyWAV(),
		"audio/wav",
		audioOpts,
		4500000, // x EMU
		1500000, // y EMU
		1800000, // cx EMU
		900000,  // cy EMU
	); err != nil {
		return fmt.Errorf("add audio: %w", err)
	}
	log.Println("Inserted WAV audio from bytes.")

	// --- 2. Video from file (optional fixture) ---
	videoPath := filepath.Join(assetsDir, "sample.mp4")
	posterPath := filepath.Join(assetsDir, "poster.png")
	if fileExists(videoPath) && fileExists(posterPath) {
		videoOpts := editor.NewAutoPlayVideoPlaybackOptions().
			WithLoop(false).
			WithMuted(true).
			WithVolume(0).
			WithAltText("Sample video clip")
		if _, err := ed.AddVideoFromFileWithPlaybackOptions(
			0, videoPath, posterPath, "video/mp4",
			videoOpts,
			600000, 1400000, 3600000, 2100000,
		); err != nil {
			return fmt.Errorf("add video: %w", err)
		}
		log.Printf("Inserted video from fixture: %s\n", videoPath)
	} else {
		log.Printf("Video fixture not found (%s + %s); skipping.\n", videoPath, posterPath)
	}

	// --- 3. OLE object from bytes (optional fixture) ---
	olePath := filepath.Join(assetsDir, "sample_ole.bin")
	if fileExists(olePath) && fileExists(posterPath) {
		oleData, readErr := os.ReadFile(olePath)
		if readErr != nil {
			return fmt.Errorf("read OLE fixture: %w", readErr)
		}
		if _, err := ed.AddOLEObject(
			0,
			oleData,
			tinyPNG(),
			"Excel.Sheet.12",
			4500000, 2800000, 2100000, 1400000,
		); err != nil {
			return fmt.Errorf("add OLE object: %w", err)
		}
		log.Printf("Inserted OLE object from fixture: %s\n", olePath)
	} else {
		log.Printf("OLE fixture not found (%s); skipping.\n", olePath)
	}

	outPath := filepath.Join(outputDir, outputFile)
	if err := ed.Save(outPath); err != nil {
		return fmt.Errorf("save output: %w", err)
	}

	log.Printf("Generated %s\n", outPath)
	return nil
}

func tinyWAV() []byte {
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

func tinyPNG() []byte {
	const b64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgA3FoXwAAAAASUVORK5CYII="
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return data
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
