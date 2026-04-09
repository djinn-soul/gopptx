package shape

import "testing"

func TestParseBodyPrAutoFit_AcceptsLegacyNormAutofitTag(t *testing.T) {
	bodyPr := &bodyPrXML{
		NormAutofit: &struct{}{},
	}

	autoFitType, autoFit := parseBodyPrAutoFit(bodyPr)
	if autoFitType == nil || *autoFitType != "normal" {
		t.Fatalf("expected normal autofit type, got %+v", autoFitType)
	}
	if autoFit == nil || !*autoFit {
		t.Fatalf("expected legacy autofit bool true, got %+v", autoFit)
	}
}
