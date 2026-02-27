package itest

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

// The OOXML <a:rPr> maps Title/Content formatting. We verify they exist.

func getSlide1XML(t *testing.T, slides []pptx.SlideContent) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("test", slides)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("unzip: %v", err)
	}
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			rc, _ := f.Open()
			var c bytes.Buffer
			_, _ = c.ReadFrom(rc)
			_ = rc.Close()
			return c.String()
		}
	}
	t.Fatal("no slide1.xml found")
	return ""
}

func mustContainAll(t *testing.T, xml string, substrs ...string) {
	t.Helper()
	for _, sub := range substrs {
		if !strings.Contains(xml, sub) {
			t.Errorf("missing expected XML %q\nGot:\n%s", sub, xml)
		}
	}
}

func TestTextStyleParity(t *testing.T) {
	t.Run("TitleBold and ContentItalic", func(t *testing.T) {
		s := pptx.NewSlide("T").
			WithTitleBold(true).
			AddBullet("B").
			WithContentItalic(true)

		xml := getSlide1XML(t, []pptx.SlideContent{s})
		mustContainAll(t, xml, `b="1"`) // title bold
		mustContainAll(t, xml, `i="1"`) // content italic
	})

	t.Run("Underline", func(t *testing.T) {
		s := pptx.NewSlide("T").
			WithTitleUnderline(true).
			AddBullet("B").
			WithContentUnderline(true)

		xml := getSlide1XML(t, []pptx.SlideContent{s})
		// both title and content get underline
		if strings.Count(xml, `u="sng"`) < 2 {
			t.Errorf("expected 2+ underline bounds, got less. XML:\n%s", xml)
		}
	})

	t.Run("Colors and Size", func(t *testing.T) {
		s := pptx.NewSlide("T").
			WithTitleColor("FF0000").
			WithTitleSize(42).
			AddBullet("B").
			WithContentColor("00FF00").
			WithContentSize(12)

		xml := getSlide1XML(t, []pptx.SlideContent{s})
		// Title: 42pt = 4200 hundredths of pt; red = FF0000
		mustContainAll(t, xml, `sz="4200"`, `<a:srgbClr val="FF0000"/>`)
		// Content: 12pt = 1200 hundredths of pt; green = 00FF00
		mustContainAll(t, xml, `sz="1200"`, `<a:srgbClr val="00FF00"/>`)
	})
}
