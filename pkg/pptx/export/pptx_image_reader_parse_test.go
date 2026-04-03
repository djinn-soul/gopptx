package export

import (
	"encoding/xml"
	"testing"
)

// makeStringPtr returns a pointer to a string value — test helper.
func makeStringPtr(s string) *string { return &s }

// unmarshalPicReaderXML parses raw XML into picReaderXML for test cases.
func unmarshalPicReaderXML(t *testing.T, rawXML string) *picReaderXML {
	t.Helper()
	var src picReaderXML
	if err := xml.Unmarshal([]byte(rawXML), &src); err != nil {
		t.Fatalf("xml.Unmarshal: %v", err)
	}
	return &src
}

// --- picCNvPrIsDecorative ---

func TestPicCNvPrIsDecorative_NoExtLst(t *testing.T) {
	c := picCNvPrXML{}
	if picCNvPrIsDecorative(c) {
		t.Error("expected false when no extLst present")
	}
}

func TestPicCNvPrIsDecorative_WrongURI(t *testing.T) {
	const picXML = `<pic xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <nvPicPr>
    <cNvPr descr="">
      <a:extLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
        <a:ext uri="{WRONG-URI-0000-0000-0000-000000000000}">
          <adec:decorative xmlns:adec="http://schemas.microsoft.com/office/drawing/2016/decorative" val="1"/>
        </a:ext>
      </a:extLst>
    </cNvPr>
  </nvPicPr>
</pic>`
	src := unmarshalPicReaderXML(t, picXML)
	_, isDecorative := picAltText(src)
	if isDecorative {
		t.Error("expected IsDecorative=false for wrong extension URI")
	}
}

func TestPicCNvPrIsDecorative_ExplicitTrue(t *testing.T) {
	const picXML = `<pic>
  <nvPicPr>
    <cNvPr descr="">
      <a:extLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
        <a:ext uri="{C183D7F6-72BE-476a-BEBA-66C5E2CAE503}">
          <adec:decorative xmlns:adec="http://schemas.microsoft.com/office/drawing/2016/decorative" val="1"/>
        </a:ext>
      </a:extLst>
    </cNvPr>
  </nvPicPr>
</pic>`
	src := unmarshalPicReaderXML(t, picXML)
	_, isDecorative := picAltText(src)
	if !isDecorative {
		t.Error("expected IsDecorative=true for explicit decorative extension with val=1")
	}
}

// --- picAltText: P2 regression suite ---

func TestPicAltText_EmptyDescrIsNotDecorative(t *testing.T) {
	// P2 regression: blank descr attr MUST NOT infer IsDecorative=true.
	src := &picReaderXML{}
	src.NvPicPr.CNvPr.Descr = makeStringPtr("")

	altText, isDecorative := picAltText(src)
	if isDecorative {
		t.Error("expected IsDecorative=false for empty descr without explicit extension, got true")
	}
	if altText != "" {
		t.Errorf("expected empty altText, got %q", altText)
	}
}

func TestPicAltText_NilDescrIsNotDecorative(t *testing.T) {
	src := &picReaderXML{}
	_, isDecorative := picAltText(src)
	if isDecorative {
		t.Error("expected IsDecorative=false when descr attr is absent")
	}
}

func TestPicAltText_NonEmptyDescrReturnsAltText(t *testing.T) {
	src := &picReaderXML{}
	src.NvPicPr.CNvPr.Descr = makeStringPtr("company logo")

	altText, isDecorative := picAltText(src)
	if altText != "company logo" {
		t.Errorf("expected altText=%q, got %q", "company logo", altText)
	}
	if isDecorative {
		t.Error("expected IsDecorative=false for image with non-empty descr")
	}
}

func TestPicAltText_TitleFallbackWhenNoDescr(t *testing.T) {
	src := &picReaderXML{}
	src.NvPicPr.CNvPr.Title = makeStringPtr("  Chart image  ")

	altText, isDecorative := picAltText(src)
	if altText != "Chart image" {
		t.Errorf("expected title fallback %q, got %q", "Chart image", altText)
	}
	if isDecorative {
		t.Error("expected IsDecorative=false for image with title fallback")
	}
}

func TestPicAltText_NilSrc(t *testing.T) {
	altText, isDecorative := picAltText(nil)
	if altText != "" || isDecorative {
		t.Errorf("expected (\"\", false) for nil src, got (%q, %v)", altText, isDecorative)
	}
}
