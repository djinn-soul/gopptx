package pptxxml_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func BenchmarkContentTypesLargeDeck(b *testing.B) {
	imageExts := []string{"png", "jpg", "gif", "bmp"}
	notesSlides := make([]int, 200)
	for i := range notesSlides {
		notesSlides[i] = i + 1
	}

	b.ReportAllocs()
	for b.Loop() {
		_ = pptxxml.ContentTypes(
			200,
			imageExts,
			80,
			50,
			notesSlides,
			true,
			10,
			4,
			5,
			false,
			nil,
			false,
			false,
			false,
			false,
			false,
		)
	}
}
