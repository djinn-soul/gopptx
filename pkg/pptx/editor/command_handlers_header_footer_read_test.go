package editor

import "testing"

func TestParseHeaderFooterXML_ConcatenatesAllTextRuns(t *testing.T) {
	slideXML := `<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><p:cSld/><p:hf sldNum="1" dt="1" ftr="1"><p:ftr><a:r><a:t>Team </a:t></a:r><a:r><a:t>Update</a:t></a:r></p:ftr><p:dt><a:r><a:t>2026-</a:t></a:r><a:r><a:t>04-04</a:t></a:r></p:dt></p:hf></p:sld>`

	hf := parseHeaderFooterXML(slideXML)

	if !hf.ShowSlideNum || !hf.ShowDateTime || !hf.ShowFooter {
		t.Fatalf("expected all visibility flags true, got %+v", hf)
	}
	if hf.Footer != "Team Update" {
		t.Fatalf("expected concatenated footer text, got %q", hf.Footer)
	}
	if hf.DateTimeText != "2026-04-04" {
		t.Fatalf("expected concatenated date text, got %q", hf.DateTimeText)
	}
}
