package editorcommon

import (
	"testing"
)

func TestPathHelpers(t *testing.T) {
	t.Run("ResolveRelationshipTarget", func(t *testing.T) {
		res := ResolveRelationshipTarget("ppt/slides/slide1.xml", "../media/image1.png")
		if res != "ppt/media/image1.png" {
			t.Errorf("expected ppt/media/image1.png, got %q", res)
		}
	})

	t.Run("MakeRelativePath", func(t *testing.T) {
		res := MakeRelativePath("ppt/slides/slide1.xml", "ppt/media/image1.png")
		if res != "../media/image1.png" {
			t.Errorf("expected ../media/image1.png, got %q", res)
		}
	})

	t.Run("CanonicalPartPath", func(t *testing.T) {
		if CanonicalPartPath(`\ppt\slides\slide1.xml`) != "ppt/slides/slide1.xml" {
			t.Error("CanonicalPartPath failed")
		}
		if CanonicalPartPath("/ppt/slides/slide1.xml") != "ppt/slides/slide1.xml" {
			t.Error("CanonicalPartPath failed")
		}
	})

	t.Run("SlideRelsPartName", func(t *testing.T) {
		if SlideRelsPartName("ppt/slides/slide1.xml") != "ppt/slides/_rels/slide1.xml.rels" {
			t.Error("SlideRelsPartName failed")
		}
		if SlideRelsPartName("presentation.xml") != "_rels/presentation.xml.rels" {
			t.Error("SlideRelsPartName failed for root")
		}
	})
	
	t.Run("RelsPathFor", func(t *testing.T) {
		if RelsPathFor("ppt/slides/slide1.xml") != "ppt/slides/_rels/slide1.xml.rels" {
			t.Error("RelsPathFor failed")
		}
	})
}

func TestParsingHelpers(t *testing.T) {
	t.Run("ParseRelationshipNumber", func(t *testing.T) {
		n, ok := ParseRelationshipNumber("rId5")
		if !ok || n != 5 { t.Error("ParseRelationshipNumber failed") }
		_, ok = ParseRelationshipNumber("invalid")
		if ok { t.Error("ParseRelationshipNumber should fail") }
	})

	t.Run("ParseSlidePartNumber", func(t *testing.T) {
		n, ok := ParseSlidePartNumber("ppt/slides/slide12.xml")
		if !ok || n != 12 { t.Error("ParseSlidePartNumber failed") }
		_, ok = ParseSlidePartNumber("invalid.xml")
		if ok { t.Error("ParseSlidePartNumber should fail") }
	})
}

func TestXMLEscape(t *testing.T) {
	got := XMLEscape("<tag> & \"")
	expected := "&lt;tag&gt; &amp; &#34;"
	if got != expected {
		t.Errorf("XMLEscape failed: expected %q, got %q", expected, got)
	}
}
