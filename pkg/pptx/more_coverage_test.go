package pptx

import (
	"testing"
)

func TestConnector_Validation(t *testing.T) {
	// Cover normalizers
	if got := normalizeConnectorType(" curved "); got != ConnectorTypeCurved {
		t.Errorf("expected ConnectorTypeCurved, got %q", got)
	}
	if got := normalizeConnectorType("c"); got != ConnectorTypeCurved {
		t.Errorf("expected ConnectorTypeCurved for short code 'c', got %q", got)
	}
	if got := normalizeConnectorType("unknown"); got != "unknown" {
		t.Errorf("expected 'unknown' for unknown input, got %q", got)
	}

	if got := normalizeArrowType(" TRIANGLE "); got != ArrowTypeTriangle {
		t.Errorf("expected ArrowTypeTriangle, got %q", got)
	}
	if got := normalizeArrowType("t"); got != ArrowTypeTriangle {
		t.Errorf("expected ArrowTypeTriangle for short code 't', got %q", got)
	}
	if got := normalizeArrowType("n"); got != ArrowTypeNone {
		t.Errorf("expected ArrowTypeNone for short code 'n', got %q", got)
	}

	if got := normalizeArrowSize(" Large "); got != ArrowSizeLarge {
		t.Errorf("expected ArrowSizeLarge, got %q", got)
	}
	if got := normalizeArrowSize("lg"); got != ArrowSizeLarge {
		t.Errorf("expected ArrowSizeLarge for short code 'lg', got %q", got)
	}

	if got := normalizeConnectionSite(" Top "); got != ConnectionSiteTop {
		t.Errorf("expected ConnectionSiteTop, got %q", got)
	}
	if got := normalizeConnectionSite("t"); got != ConnectionSiteTop {
		t.Errorf("expected ConnectionSiteTop for short code 't', got %q", got)
	}
	if got := normalizeConnectionSite("ctr"); got != ConnectionSiteCenter {
		t.Errorf("expected ConnectionSiteCenter for short code 'ctr', got %q", got)
	}
}
