package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsBulletStyleVariants(t *testing.T) {
	slide := NewSlide("Styles").
		AddBulletWithStyle("Number", NewTextParagraphStyle().WithNumbered()).
		AddBulletWithStyle("Letter", NewTextParagraphStyle().WithLetteredLower()).
		AddBulletWithStyle("Roman", NewTextParagraphStyle().WithRomanUpper().WithLevel(1)).
		AddBulletWithStyle("Custom", NewTextParagraphStyle().WithCustomBullet("~").WithLevel(2)).
		AddBulletWithStyle("None", NewTextParagraphStyle().WithNoBullet())

	data, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:buAutoNum type="arabicPeriod"/>`,
		`<a:buAutoNum type="alphaLcPeriod"/>`,
		`<a:buAutoNum type="romanUcPeriod"/>`,
		`<a:pPr lvl="1" marL="1371600" indent="-914400">`,
		`<a:pPr lvl="2" marL="2286000" indent="-1371600">`,
		`<a:buChar char="~"/>`,
		`<a:buNone/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesBulletStyleHelpers(t *testing.T) {
	slide := NewSlide("Helpers").
		WithBulletStyle(BulletStyleRomanUpper).
		AddBullet("Top").
		AddSubBullet("Nested").
		AddNumbered("Step").
		AddLettered("Option")

	data, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:buAutoNum type="romanUcPeriod"/>`,
		`<a:pPr lvl="1" marL="1371600" indent="-914400">`,
		`<a:buAutoNum type="arabicPeriod"/>`,
		`<a:buAutoNum type="alphaLcPeriod"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidBulletStyles(t *testing.T) {
	cases := []struct {
		name    string
		style   TextParagraphStyle
		wantErr string
	}{
		{
			name:    "unknown-style",
			style:   NewTextParagraphStyle().WithBulletStyle("triangle"),
			wantErr: "bullet style must be one of",
		},
		{
			name:    "bad-level",
			style:   NewTextParagraphStyle().WithLevel(9),
			wantErr: "level must be between 0 and 8",
		},
		{
			name:    "custom-missing-char",
			style:   NewTextParagraphStyle().WithBulletStyle(BulletStyleCustom),
			wantErr: "custom bullet must be a single character",
		},
		{
			name:    "char-with-number-style",
			style:   NewTextParagraphStyle().WithNumbered().WithCustomBullet("x").WithNumbered(),
			wantErr: "bullet char is only allowed with custom bullet style",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slide := NewSlide("Invalid").AddBulletWithStyle("one", tc.style)
			_, err := CreateWithSlides("Deck", []SlideContent{slide})
			if err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
