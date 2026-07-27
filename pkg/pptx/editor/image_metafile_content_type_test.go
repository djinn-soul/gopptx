package editor

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

// tinyEMF returns a minimal EMR_HEADER carrying the " EMF" signature at offset 40.
func tinyEMF() []byte {
	data := make([]byte, 88)
	copy(data, []byte{0x01, 0x00, 0x00, 0x00})
	copy(data[40:], []byte{0x20, 0x45, 0x4D, 0x46})
	return data
}

func TestImageContentTypeDetectsMetafiles(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		format string
		want   string
	}{
		{name: "emf bytes", data: tinyEMF(), format: "", want: "image/x-emf"},
		{name: "emf bytes with hint", data: tinyEMF(), format: "emf", want: "image/x-emf"},
		{
			name:   "wmf placeable",
			data:   []byte{0xD7, 0xCD, 0xC6, 0x9A, 0x00, 0x00},
			format: "",
			want:   "image/x-wmf",
		},
		{
			name:   "wmf disk header",
			data:   []byte{0x01, 0x00, 0x09, 0x00, 0x00, 0x00},
			format: "",
			want:   "image/x-wmf",
		},
		{
			name:   "wmf memory header",
			data:   []byte{0x02, 0x00, 0x09, 0x00, 0x00, 0x00},
			format: "",
			want:   "image/x-wmf",
		},
		{
			name:   "emf hint without magic",
			data:   []byte{0x00, 0x11, 0x22, 0x33},
			format: "emf",
			want:   "image/x-emf",
		},
		{
			name:   "wmf hint without magic",
			data:   []byte{0x00, 0x11, 0x22, 0x33},
			format: ".WMF",
			want:   "image/x-wmf",
		},
		{
			name:   "wdp hint",
			data:   []byte{0x00, 0x11, 0x22, 0x33},
			format: "wdp",
			want:   "image/vnd.ms-photo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := imageContentType(tc.data, tc.format); got != tc.want {
				t.Fatalf("imageContentType() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAddEMFImageEmitsContentTypeDefault(t *testing.T) {
	pptxPath := createPictureFixturePPTX(t, testutil.TinyPNG())

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	emf := tinyEMF()
	shapeID, err := ed.AddImageFromBytes(0, emf, "emf", 0, 0, 1000000, 1000000, nil)
	if err != nil {
		t.Fatalf("add emf image: %v", err)
	}

	meta, err := ed.GetImageMetadata(0, shapeID)
	if err != nil {
		t.Fatalf("get emf metadata: %v", err)
	}
	if meta.Format != "emf" {
		t.Fatalf("emf format = %q, want %q", meta.Format, "emf")
	}
	if meta.ContentType != "image/x-emf" {
		t.Fatalf("emf content type = %q, want image/x-emf", meta.ContentType)
	}

	saved, err := ed.SaveToBytes()
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(saved), int64(len(saved)))
	if err != nil {
		t.Fatalf("zip read: %v", err)
	}
	contentTypes := testutil.ReadZipFile(t, zr, "[Content_Types].xml")
	if !strings.Contains(contentTypes, `Extension="emf" ContentType="image/x-emf"`) {
		t.Fatalf("missing emf default in content types: %s", contentTypes)
	}
}

func TestImageContentTypeKeepsRasterDetection(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	if got := imageContentType(png, "png"); got != mimePNG {
		t.Fatalf("png content type = %q, want %q", got, mimePNG)
	}
	if got := imageContentType([]byte("GIF89a"), ""); got != mimeGIF {
		t.Fatalf("gif content type = %q, want %q", got, mimeGIF)
	}
}
