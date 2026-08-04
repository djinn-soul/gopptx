package shape

import (
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestParseBodyPrAutoFit_AcceptsLegacyNormAutofitTag(t *testing.T) {
	bodyPr := &bodyPrXML{
		NormAutofit: &normAutofitXML{},
	}

	autoFitType, autoFit := parseBodyPrAutoFit(bodyPr)
	if autoFitType == nil || *autoFitType != "normal" {
		t.Fatalf("expected normal autofit type, got %+v", autoFitType)
	}
	if autoFit == nil || !*autoFit {
		t.Fatalf("expected legacy autofit bool true, got %+v", autoFit)
	}
}

// A deck PowerPoint has laid out carries the shrink amounts it computed. Losing
// them on read means a later edit rewrites the frame at full size and the text
// spills out of the shape again (upstream issue #969).
func TestApplyParsedShrinkAmounts_ReportsPercentages(t *testing.T) {
	fontScale := 62500
	lnSpcReduction := 20000
	bodyPr := &bodyPrXML{
		NormAutofit: &normAutofitXML{FontScale: &fontScale, LnSpcReduction: &lnSpcReduction},
	}
	frame := &common.TextFrame{}

	applyParsedShrinkAmounts(frame, bodyPr)

	if frame.FontScale == nil || *frame.FontScale != 62.5 {
		t.Fatalf("expected font scale 62.5%%, got %+v", frame.FontScale)
	}
	if frame.LineSpaceReduction == nil || *frame.LineSpaceReduction != 20 {
		t.Fatalf("expected line space reduction 20%%, got %+v", frame.LineSpaceReduction)
	}
}
