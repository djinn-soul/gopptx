package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsBulletParagraphStyle(t *testing.T) {
	slide := NewSlide("Paragraph Style").
		AddBulletWithStyle(
			"Styled plain bullet",
			NewTextParagraphStyle().
				WithAlignCenter().
				WithLineSpacingPct(125).
				WithSpaceBeforePt(4).
				WithSpaceAfterPt(6),
		).
		AddBulletRunsWithStyle(
			[]TextRun{
				NewTextRun("Styled run bullet").WithBold(true),
			},
			NewTextParagraphStyle().WithAlignRight().WithSpaceAfterPt(3),
		)

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:pPr lvl="0" marL="457200" indent="-457200" algn="ctr">`,
		`<a:lnSpc><a:spcPct val="125000"/></a:lnSpc>`,
		`<a:spcBef><a:spcPts val="400"/></a:spcBef>`,
		`<a:spcAft><a:spcPts val="600"/></a:spcAft>`,
		`<a:pPr lvl="0" marL="457200" indent="-457200" algn="r"><a:buChar char="•"/><a:spcAft><a:spcPts val="300"/></a:spcAft></a:pPr>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidTextParagraphAlign(t *testing.T) {
	slide := NewSlide("Bad").
		AddBulletWithStyle("one", NewTextParagraphStyle().WithAlign("diagonal"))

	_, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected paragraph align validation error")
	}
	if !strings.Contains(err.Error(), "align must be one of l|ctr|r|just") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsInvalidTextParagraphSpacing(t *testing.T) {
	tests := []struct {
		name  string
		style TextParagraphStyle
		msg   string
	}{
		{
			name:  "negative before",
			style: NewTextParagraphStyle().WithSpaceBeforePt(-1),
			msg:   "space-before must be >= 0",
		},
		{
			name:  "negative after",
			style: NewTextParagraphStyle().WithSpaceAfterPt(-1),
			msg:   "space-after must be >= 0",
		},
		{
			name:  "negative line spacing",
			style: NewTextParagraphStyle().WithLineSpacingPct(-1),
			msg:   "line-spacing must be >= 0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slide := NewSlide("Bad").AddBulletWithStyle("one", tc.style)
			_, err := CreateWithSlides("Demo", []SlideContent{slide})
			if err == nil {
				t.Fatalf("expected paragraph spacing validation error")
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
