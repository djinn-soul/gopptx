package pptxxml_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func TestSlideMasterTxStyles(t *testing.T) {
	spec := &pptxxml.SlideMasterSpec{
		TxStyles: &pptxxml.TxStylesSpec{
			TitleStyle: []pptxxml.TextLevelStyle{
				{Level: 0, Font: "Arial", SizePt: 44, Bold: true, Color: "003366"},
			},
			BodyStyle: []pptxxml.TextLevelStyle{
				{Level: 0, Font: "Calibri", SizePt: 28, Color: "333333"},
				{Level: 1, Font: "Calibri", SizePt: 24, Color: "666666", BulletChar: "–", IndentEMU: 457200},
				{Level: 2, Font: "Calibri", SizePt: 20, Italic: true},
			},
			OtherStyle: []pptxxml.TextLevelStyle{
				{Level: 0, SizePt: 18},
			},
		},
	}

	xml := pptxxml.SlideMaster(spec)

	// Title style checks
	if !strings.Contains(xml, "<p:titleStyle>") {
		t.Error("missing <p:titleStyle>")
	}
	if !strings.Contains(xml, `<a:latin typeface="Arial"/>`) {
		t.Error("missing Arial typeface in titleStyle")
	}
	if !strings.Contains(xml, `sz="4400"`) {
		t.Error("expected 44pt = sz 4400")
	}
	if !strings.Contains(xml, `b="1"`) {
		t.Error("expected bold flag")
	}
	if !strings.Contains(xml, `<a:srgbClr val="003366"/>`) {
		t.Error("expected title color")
	}

	// Body style checks
	if !strings.Contains(xml, "<p:bodyStyle>") {
		t.Error("missing <p:bodyStyle>")
	}
	if !strings.Contains(xml, `<a:lvl2pPr`) {
		t.Error("expected level 2 paragraph props")
	}
	if !strings.Contains(xml, `<a:buChar char="–"/>`) {
		t.Error("expected bullet char in lvl2")
	}
	if !strings.Contains(xml, `indent="457200"`) {
		t.Error("expected indent in lvl2")
	}
	if !strings.Contains(xml, `<a:lvl3pPr>`) {
		t.Error("expected level 3 paragraph props")
	}
	if !strings.Contains(xml, `i="1"`) {
		t.Error("expected italic in lvl3")
	}

	// Other style checks
	if !strings.Contains(xml, "<p:otherStyle>") {
		t.Error("missing <p:otherStyle>")
	}
	if !strings.Contains(xml, `sz="1800"`) {
		t.Error("expected 18pt = sz 1800 in otherStyle")
	}

	// Enclosing element
	if !strings.Contains(xml, "<p:txStyles>") {
		t.Error("missing <p:txStyles>")
	}
	if !strings.Contains(xml, "</p:txStyles>") {
		t.Error("missing </p:txStyles>")
	}
}

func TestSlideMasterWithoutTxStyles(t *testing.T) {
	xml := pptxxml.SlideMaster(nil)
	if strings.Contains(xml, "<p:txStyles>") {
		t.Error("nil spec should not emit txStyles")
	}
}

func TestPresentationMultiMaster(t *testing.T) {
	xml := pptxxml.Presentation("Test", 2, false, 12192000, 6858000, 3, nil, nil, false, nil, nil)

	// Should have 3 master IDs
	if count := strings.Count(xml, "sldMasterId"); count != 2+3 {
		// 2 for opening/closing sldMasterIdLst, 3 for individual entries + their closing
		// Actually <p:sldMasterId has the substring, let's count self-closing tags
		t.Errorf("expected 5 occurrences of 'sldMasterId' (2 tags + 3 entries), got %d", count)
	}
	if !strings.Contains(xml, `id="2147483648"`) {
		t.Error("missing first master ID")
	}
	if !strings.Contains(xml, `id="2147483655"`) {
		t.Error("missing second master ID")
	}
	if !strings.Contains(xml, `id="2147483662"`) {
		t.Error("missing third master ID")
	}

	// Slides are offset by masterCount + theme relationship.
	if !strings.Contains(xml, `r:id="rId5"`) {
		t.Error("first slide should reference rId5 with 3 masters")
	}
	if !strings.Contains(xml, `r:id="rId6"`) {
		t.Error("second slide should reference rId6 with 3 masters")
	}
}

func TestPresentationRelationshipsMultiMaster(t *testing.T) {
	xml := pptxxml.PresentationRelationships(2, false, 0, 2, false, false, false, false, 0)

	// Should have 2 master relationships (checked by target presence)
	if !strings.Contains(xml, `Target="slideMasters/slideMaster1.xml"`) {
		t.Error("missing slideMaster1 target")
	}
	if !strings.Contains(xml, `Target="slideMasters/slideMaster2.xml"`) {
		t.Error("missing slideMaster2 target")
	}
	if strings.Contains(xml, `Target="slideMasters/slideMaster3.xml"`) {
		t.Error("should not have slideMaster3 target")
	}

	// Slides should start at rId3 with 2 masters
	if !strings.Contains(xml, `Id="rId3"`) {
		t.Error("first slide should be rId3 with 2 masters")
	}
}

func TestPresentationModifyVerifierUsesPowerPointFields(t *testing.T) {
	xml := pptxxml.Presentation(
		"Protected",
		1,
		false,
		12192000,
		6858000,
		1,
		&pptxxml.ProtectionInfo{
			HashAlgSID: 14,
			HashData:   "HASH",
			SaltData:   "SALT",
			SpinCount:  100000,
		},
		nil,
		false,
		nil,
		nil,
	)

	if !strings.Contains(xml, `cryptProviderType="rsaAES"`) {
		t.Error("missing rsaAES provider")
	}
	if !strings.Contains(xml, `cryptAlgorithmSid="14"`) {
		t.Error("missing SHA-512 algorithm SID")
	}
	if !strings.Contains(xml, `saltData="SALT"`) {
		t.Error("missing saltData attribute")
	}
	if !strings.Contains(xml, `hashData="HASH"`) {
		t.Error("missing hashData attribute")
	}
	if strings.Contains(xml, `saltValue=`) {
		t.Error("legacy saltValue attribute should not be emitted")
	}
	if strings.Contains(xml, `hashValue=`) {
		t.Error("legacy hashValue attribute should not be emitted")
	}
}

func TestContentTypesMultiMaster(t *testing.T) {
	xml := pptxxml.ContentTypes(1, nil, 0, 0, nil, false, 0, 2, 0, false, nil, false, false, false, false, false)
	if !strings.Contains(xml, `/ppt/slideMasters/slideMaster1.xml`) {
		t.Error("missing slideMaster1 content-type override")
	}
	if !strings.Contains(xml, `/ppt/slideMasters/slideMaster2.xml`) {
		t.Error("missing slideMaster2 content-type override")
	}
	if strings.Contains(xml, `/ppt/slideMasters/slideMaster3.xml`) {
		t.Error("unexpected slideMaster3 content-type override")
	}
}
