package shapes

import (
	"testing"
)

func TestNormalizeConnectorType(t *testing.T) {
	if got := NormalizeConnectorType(" curved "); got != ConnectorTypeCurved {
		t.Errorf("expected ConnectorTypeCurved, got %q", got)
	}
	if got := NormalizeConnectorType("c"); got != ConnectorTypeCurved {
		t.Errorf("expected ConnectorTypeCurved for short code 'c', got %q", got)
	}
	if got := NormalizeConnectorType("unknown"); got != "unknown" {
		t.Errorf("expected 'unknown' for unknown input, got %q", got)
	}
}

func TestNormalizeArrowType(t *testing.T) {
	if got := NormalizeArrowType(" TRIANGLE "); got != ArrowTypeTriangle {
		t.Errorf("expected ArrowTypeTriangle, got %q", got)
	}
	if got := NormalizeArrowType("t"); got != ArrowTypeTriangle {
		t.Errorf("expected ArrowTypeTriangle for short code 't', got %q", got)
	}
	if got := NormalizeArrowType("n"); got != ArrowTypeNone {
		t.Errorf("expected ArrowTypeNone for short code 'n', got %q", got)
	}
}

func TestNormalizeArrowSize(t *testing.T) {
	if got := NormalizeArrowSize(" Large "); got != ArrowSizeLarge {
		t.Errorf("expected ArrowSizeLarge, got %q", got)
	}
	if got := NormalizeArrowSize("lg"); got != ArrowSizeLarge {
		t.Errorf("expected ArrowSizeLarge for short code 'lg', got %q", got)
	}
}

func TestNormalizeConnectionSite(t *testing.T) {
	if got := NormalizeConnectionSite(" Top "); got != ConnectionSiteTop {
		t.Errorf("expected ConnectionSiteTop, got %q", got)
	}
	if got := NormalizeConnectionSite("t"); got != ConnectionSiteTop {
		t.Errorf("expected ConnectionSiteTop for short code 't', got %q", got)
	}
	if got := NormalizeConnectionSite("ctr"); got != ConnectionSiteCenter {
		t.Errorf("expected ConnectionSiteCenter for short code 'ctr', got %q", got)
	}
}
