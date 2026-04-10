package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestEditorParagraphToSpec_RoundTripsSpaceBeforePts(t *testing.T) {
	spaceBefore := 1200
	p := &common.Paragraph{
		SpaceBeforePts: &spaceBefore,
	}
	xml, err := renderParagraphPropsXML(p)
	if err != nil {
		t.Fatalf("renderParagraphPropsXML failed: %v", err)
	}
	if !strings.Contains(xml, `<a:spcBef><a:spcPts val="1200"/></a:spcBef>`) {
		t.Fatalf("expected spcBef val=1200 round-trip, got: %s", xml)
	}
}

func TestEditorParagraphToSpec_PreservesFractionalSpacingUnits(t *testing.T) {
	lineSpacingPct := 115500
	spaceBefore := 250
	spaceAfter := 350
	p := &common.Paragraph{
		LineSpacingPct: &lineSpacingPct,
		SpaceBeforePts: &spaceBefore,
		SpaceAfterPts:  &spaceAfter,
	}
	xml, err := renderParagraphPropsXML(p)
	if err != nil {
		t.Fatalf("renderParagraphPropsXML failed: %v", err)
	}
	if !strings.Contains(xml, `<a:lnSpc><a:spcPct val="115500"/></a:lnSpc>`) {
		t.Fatalf("expected fractional line spacing pct round-trip, got: %s", xml)
	}
	if !strings.Contains(xml, `<a:spcBef><a:spcPts val="250"/></a:spcBef>`) {
		t.Fatalf("expected fractional space before round-trip, got: %s", xml)
	}
	if !strings.Contains(xml, `<a:spcAft><a:spcPts val="350"/></a:spcAft>`) {
		t.Fatalf("expected fractional space after round-trip, got: %s", xml)
	}
}

func TestEditorParagraphToSpec_PreservesFractionalLineSpacingPoints(t *testing.T) {
	lineSpacingPts := 1850
	p := &common.Paragraph{
		LineSpacingPts: &lineSpacingPts,
	}
	xml, err := renderParagraphPropsXML(p)
	if err != nil {
		t.Fatalf("renderParagraphPropsXML failed: %v", err)
	}
	if !strings.Contains(xml, `<a:lnSpc><a:spcPts val="1850"/></a:lnSpc>`) {
		t.Fatalf("expected fractional line spacing points round-trip, got: %s", xml)
	}
}
