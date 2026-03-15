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
	outputFile = "23_media_embed_editor.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-media-embed-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	basePath := filepath.Join(tmpDir, "media_embed_base.pptx")
	outPath := filepath.Join(outputDir, outputFile)

	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Media Embed (Editor)").
			AddBullet("Video, audio, and OLE media relationships inserted by PresentationEditor.").
			AddBullet("Provide real fixture files in examples/assets/23 to enable inserts."),
	}
	if err := pptx.WriteFile(basePath, "Media Embed Demo", baseSlides); err != nil {
		return fmt.Errorf("write base deck: %w", err)
	}

	ed, err := pptx.OpenPresentationEditor(basePath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	posterPath := filepath.Join("examples", "assets", "23", "poster.png")
	if _, err := ensureMediaFixture(posterPath); err != nil {
		log.Printf("Poster fixture missing (%s): %v", posterPath, err)
	}

	videoFixturePath := filepath.Join("examples", "assets", "23", "sample.mp4")
	if fileExists(videoFixturePath) && fileExists(posterPath) {
		videoOpts := editor.NewAutoPlayVideoPlaybackOptions().
			WithLoop(true).
			WithMuted(true).
			WithVolume(0).
			WithAltText("Demo launch reel")
		if _, err := ed.AddVideoFromFileWithPlaybackOptions(
			0,
			videoFixturePath,
			posterPath,
			"video/mp4",
			videoOpts,
			600000,
			1400000,
			3600000,
			2100000,
		); err != nil {
			return fmt.Errorf("add video from fixture: %w", err)
		}
		log.Printf("Added video fixture: %s\n", videoFixturePath)
	} else {
		log.Printf("Video fixtures missing (need sample.mp4 + poster.png in examples/assets/23); skipping video insertion.")
	}

	audioMP3FixturePath := filepath.Join("examples", "assets", "23", "sample.mp3")
	audioWAVFixturePath := filepath.Join("examples", "assets", "23", "sample.wav")
	switch {
	case fileExists(audioMP3FixturePath):
		audioOpts := editor.NewAutoPlayAudioPlaybackOptions().
			WithLoop(true).
			WithPlayAcrossSlides(true).
			WithVolume(85).
			WithAltText("Background soundtrack")
		if _, err := ed.AddAudioFromFileWithPlaybackOptions(
			0,
			audioMP3FixturePath,
			"audio/mpeg",
			audioOpts,
			4500000,
			1500000,
			1800000,
			900000,
		); err != nil {
			return fmt.Errorf("add audio from mp3 fixture: %w", err)
		}
		log.Printf("Added audio fixture: %s\n", audioMP3FixturePath)
	case fileExists(audioWAVFixturePath):
		audioOpts := editor.NewAudioPlaybackOptions().
			WithVolume(85).
			WithAltText("Background soundtrack")
		if _, err := ed.AddAudioFromFileWithPlaybackOptions(
			0,
			audioWAVFixturePath,
			"audio/wav",
			audioOpts,
			4500000,
			1500000,
			1800000,
			900000,
		); err != nil {
			return fmt.Errorf("add audio from wav fixture: %w", err)
		}
		log.Printf("Added audio fixture: %s\n", audioWAVFixturePath)
	default:
		if _, err := ed.AddAudioWithIcon(0, tinyWAV(), tinyPNG(), "audio/wav", 4500000, 1500000, 1800000, 900000); err != nil {
			return fmt.Errorf("add fallback audio with icon: %w", err)
		}
		log.Printf(
			"Audio fixtures missing (sample.mp3/sample.wav); inserted fallback tiny WAV with generated icon.",
		)
	}

	oleFixturePath := filepath.Join("examples", "assets", "23", "sample_ole.bin")
	if fileExists(oleFixturePath) && fileExists(posterPath) {
		oleData, readErr := os.ReadFile(oleFixturePath)
		if readErr != nil {
			return fmt.Errorf("read ole fixture: %w", readErr)
		}
		if _, err := ed.AddOLEObject(0, oleData, tinyPNG(), "Excel.Sheet.12", 4500000, 2800000, 2100000, 1400000); err != nil {
			return fmt.Errorf("add ole object: %w", err)
		}
		log.Printf("Added OLE fixture: %s\n", oleFixturePath)
	} else {
		log.Printf("OLE fixtures missing (need sample_ole.bin + poster.png in examples/assets/23); skipping OLE insertion.")
	}

	if err := ed.Save(outPath); err != nil {
		return fmt.Errorf("save output: %w", err)
	}

	log.Printf("Generated media embed example: %s\n", outPath)
	return nil
}

func tinyPNG() []byte {
	const b64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgA3FoXwAAAAASUVORK5CYII="
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	return data
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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func ensureMediaFixture(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}
