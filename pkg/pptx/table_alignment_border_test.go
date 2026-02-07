package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsTableCellAlignmentAndBorder(t *testing.T) {
	table := NewTable([]int64{2743400, 2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").
				WithAlignCenter().
				WithVAlignMiddle().
				WithBorder(1.0, "112233"),
			NewTableCell("Value").WithAlignRight(),
		})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Styled Table").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`<a:pPr algn="ctr"/>`,
		`<a:pPr algn="r"/>`,
		`<a:tcPr anchor="ctr">`,
		`<a:lnL w="12700">`,
		`<a:lnR w="12700">`,
		`<a:lnT w="12700">`,
		`<a:lnB w="12700">`,
		`<a:srgbClr val="112233"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in table XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsStyledTableInvalidAlign(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").WithAlign("diagonal"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Styled Table").WithTable(table)})
	if err == nil {
		t.Fatalf("expected styled table align validation error")
	}
	if !strings.Contains(err.Error(), "align must be one of l|ctr|r|just") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsStyledTableInvalidVAlign(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").WithVAlign("middle"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Styled Table").WithTable(table)})
	if err == nil {
		t.Fatalf("expected styled table valign validation error")
	}
	if !strings.Contains(err.Error(), "valign must be one of t|ctr|b") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsStyledTableInvalidBorder(t *testing.T) {
	tests := []struct {
		name string
		cell TableCell
		msg  string
	}{
		{
			name: "invalid border color",
			cell: NewTableCell("Header").WithBorder(1.0, "NOTHEX"),
			msg:  "border color must be 6-digit RGB hex",
		},
		{
			name: "missing width",
			cell: NewTableCell("Header").WithBorder(0, "112233"),
			msg:  "border width must be > 0 when border color is set",
		},
		{
			name: "negative width",
			cell: NewTableCell("Header").WithBorder(-1, "112233"),
			msg:  "border width must be >= 0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := NewTable([]int64{2743400}).
				AddStyledRow([]TableCell{tc.cell})

			_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Styled Table").WithTable(table)})
			if err == nil {
				t.Fatalf("expected styled table border validation error")
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
