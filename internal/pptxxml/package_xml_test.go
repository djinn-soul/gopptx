package pptxxml

import (
	"strings"
	"testing"
)

func TestSignatureOriginXML(t *testing.T) {
	xml := SignatureOrigin()

	if !strings.Contains(xml, `<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`) {
		t.Fatalf("unexpected signature origin xml: %s", xml)
	}
	if strings.Contains(xml, "<vnd.openxmlformats-package.digital-signature-origin") {
		t.Fatalf("signature origin must use SignatureOrigin element: %s", xml)
	}
}
