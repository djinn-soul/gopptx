package pptxxml

import (
	"strings"
	"testing"
)

func TestSignatureOriginXML(t *testing.T) {
	xml := SignatureOrigin()

	if !strings.Contains(
		xml,
		`<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`,
	) {
		t.Fatalf("unexpected signature origin xml: %s", xml)
	}
	if strings.Contains(xml, "<vnd.openxmlformats-package.digital-signature-origin") {
		t.Fatalf("signature origin must use SignatureOrigin element: %s", xml)
	}
}

func TestContentTypesCustomPropertiesOverride(t *testing.T) {
	withCustom := ContentTypes(1, nil, 0, 0, nil, false, 0, 1, 0, false, nil, true, false, false, false, false)
	if !strings.Contains(
		withCustom,
		`<Override PartName="/docProps/custom.xml" ContentType="application/vnd.openxmlformats-officedocument.custom-properties+xml"/>`,
	) {
		t.Fatalf("missing custom properties override in content types: %s", withCustom)
	}

	withoutCustom := ContentTypes(1, nil, 0, 0, nil, false, 0, 1, 0, false, nil, false, false, false, false, false)
	if strings.Contains(withoutCustom, `/docProps/custom.xml`) {
		t.Fatalf("unexpected custom properties override without custom props enabled: %s", withoutCustom)
	}
}
