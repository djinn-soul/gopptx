package presentation_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation"
)

func TestWritePresentationPackageWritesSignatureOriginPart(t *testing.T) {
	meta := presentation.Metadata{
		Metadata: common.Metadata{
			Protection: common.Protection{
				SignaturesEnabled: true,
			},
		},
	}
	slide := elements.NewSlide("Signed")

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := presentation.WritePresentationPackage(zw, meta, []elements.SlideContent{slide}, 1); err != nil {
		t.Fatalf("write package: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}

	foundOrigin := false
	var rootRels string
	for _, f := range zr.File {
		if f.Name == "_xmlsignatures/origin.sigs" {
			foundOrigin = true
		}
		if f.Name == "_rels/.rels" {
			rc, openErr := f.Open()
			if openErr != nil {
				t.Fatalf("open root rels: %v", openErr)
			}
			data := new(bytes.Buffer)
			if _, readErr := data.ReadFrom(rc); readErr != nil {
				_ = rc.Close()
				t.Fatalf("read root rels: %v", readErr)
			}
			if closeErr := rc.Close(); closeErr != nil {
				t.Fatalf("close root rels: %v", closeErr)
			}
			rootRels = data.String()
		}
	}

	if !foundOrigin {
		t.Fatalf("missing _xmlsignatures/origin.sigs part")
	}
	if !strings.Contains(rootRels, `Target="_xmlsignatures/origin.sigs"`) {
		t.Fatalf("missing root relationship to signature origin part")
	}
}
