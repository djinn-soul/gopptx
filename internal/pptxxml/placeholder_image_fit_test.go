package pptxxml

import (
	"strings"
	"testing"
)

// placeholderPicType is the p:ph type attribute for a picture placeholder.
const placeholderPicType = "pic"

// A placeholder picture must be able to carry an a:srcRect, otherwise a cover
// fit has no way to hide the overflowing axis and the bare a:fillRect distorts
// the image to the placeholder bounds.
func TestPlaceholderImageWritesSrcRectWhenCropped(t *testing.T) {
	spec := PlaceholderOverrideSpec{
		Index: 1,
		Type:  placeholderPicType,
		Image: &ImageRef{
			RelID: "rId2",
			Name:  "wide.png",
			X:     914400,
			Y:     914400,
			CX:    3657600,
			CY:    3657600,
			Crop:  &ImageCropRef{Left: 37500, Right: 37500},
		},
	}

	xml := PlaceholderShape(spec, 5)

	if !strings.Contains(xml, `<a:srcRect l="37500" r="37500"/>`) {
		t.Fatalf("expected the crop insets in the placeholder picture, got %s", xml)
	}
	if !strings.Contains(xml, "<a:stretch>") {
		t.Fatalf("a:srcRect must accompany a:stretch, not replace it: %s", xml)
	}
	if strings.Index(xml, "<a:srcRect") > strings.Index(xml, "<a:stretch>") {
		t.Fatalf("a:srcRect must precede a:stretch per CT_BlipFillProperties: %s", xml)
	}
	if !strings.Contains(xml, `<a:ext cx="3657600" cy="3657600"/>`) {
		t.Fatalf("cover keeps the box extent: %s", xml)
	}
}

func TestPlaceholderImageOmitsSrcRectWhenNotCropped(t *testing.T) {
	spec := PlaceholderOverrideSpec{
		Index: 1,
		Type:  placeholderPicType,
		Image: &ImageRef{RelID: "rId2", Name: "plain.png", CX: 100, CY: 100},
	}

	xml := PlaceholderShape(spec, 5)

	if strings.Contains(xml, "<a:srcRect") {
		t.Fatalf("an uncropped placeholder picture must not write a:srcRect: %s", xml)
	}
	if !strings.Contains(xml, "<a:stretch>") {
		t.Fatalf("expected the default stretch fill: %s", xml)
	}
}
