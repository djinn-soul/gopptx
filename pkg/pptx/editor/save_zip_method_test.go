package editor

import (
	"archive/zip"
	"testing"
)

func TestSaveZipMethodStoresNotesAndMedia(t *testing.T) {
	cases := []struct {
		path string
		want uint16
	}{
		{path: "ppt/notesSlides/notesSlide1.xml", want: zip.Store},
		{path: "ppt/notesMasters/notesMaster1.xml", want: zip.Store},
		{path: "ppt/media/image1.png", want: zip.Store},
		{path: "ppt/media/audio1.mp3", want: zip.Store},
		{path: "ppt/slides/slide1.xml", want: zip.Deflate},
	}

	for _, tc := range cases {
		if got := saveZipMethod(tc.path); got != tc.want {
			t.Fatalf("saveZipMethod(%q) = %d, want %d", tc.path, got, tc.want)
		}
	}
}
