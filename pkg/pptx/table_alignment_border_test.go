package pptx

import (
	"archive/zip"
	"bytes"
	"math"
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
			msg:  "border width must be > 0 when",
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

func TestCreateWithSlidesEmbedsPerSideTableBorderStyles(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Edge").
				WithBorder(1.0, "112233").
				WithLeftBorderStyle(2.0, "AA0000", TableBorderDashDash).
				WithRightBorderStyle(1.5, "00AA00", TableBorderDashDot).
				WithTopBorderStyle(1.0, "0000AA", TableBorderDashLongDash).
				WithBottomBorderStyle(0.5, "ABCDEF", TableBorderDashDashDot),
		})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Border Sides").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`<a:lnL w="25400"><a:solidFill><a:srgbClr val="AA0000"/></a:solidFill><a:prstDash val="dash"/></a:lnL>`,
		`<a:lnR w="19050"><a:solidFill><a:srgbClr val="00AA00"/></a:solidFill><a:prstDash val="dot"/></a:lnR>`,
		`<a:lnT w="12700"><a:solidFill><a:srgbClr val="0000AA"/></a:solidFill><a:prstDash val="lgDash"/></a:lnT>`,
		`<a:lnB w="6350"><a:solidFill><a:srgbClr val="ABCDEF"/></a:solidFill><a:prstDash val="dashDot"/></a:lnB>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in table XML", needle)
		}
	}
}

func TestCreateWithSlidesEmbedsOnlyConfiguredBorderSides(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Edge").WithLeftBorder(1.0, "112233"),
		})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Border Left").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `<a:lnL w="12700">`) {
		t.Fatalf("expected left border in table XML")
	}
	if strings.Contains(slideXML, `<a:lnR w="`) {
		t.Fatalf("did not expect right border in table XML")
	}
	if strings.Contains(slideXML, `<a:lnT w="`) {
		t.Fatalf("did not expect top border in table XML")
	}
	if strings.Contains(slideXML, `<a:lnB w="`) {
		t.Fatalf("did not expect bottom border in table XML")
	}
}

func TestCreateWithSlidesRejectsStyledTableInvalidBorderDash(t *testing.T) {
	table := NewTable([]int64{2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("Header").WithLeftBorderStyle(1.0, "112233", "zigzag"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Styled Table").WithTable(table)})
	if err == nil {
		t.Fatalf("expected styled table border dash validation error")
	}
	if !strings.Contains(err.Error(), "left border dash style must be one of solid|dash|dot|dashDot|lgDash") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesEmbedsTableRowHeightsMarginsAndWrap(t *testing.T) {
	table := NewTable([]int64{2743400, 2743400}).
		AddStyledRow([]TableCell{
			NewTableCell("No wrap").
				WithMarginsPt(1.0).
				WithWrap(false),
			NewTableCell("Wrap enabled").
				WithMarginLeftPt(0.5).
				WithMarginRightPt(0.5).
				WithWrap(true),
		}).
		AddRow([]string{"R2C1", "R2C2"}).
		WithRowHeights([]int64{500000, 650000})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Advanced Table").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`<a:tr h="500000">`,
		`<a:tr h="650000">`,
		`<a:bodyPr wrap="none"/>`,
		`<a:bodyPr wrap="square"/>`,
		`<a:tcPr marL="12700" marR="12700" marT="12700" marB="12700">`,
		`<a:tcPr marL="6350" marR="6350">`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in table XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidTableRowHeightsAndMargins(t *testing.T) {
	tests := []struct {
		name  string
		table Table
		msg   string
	}{
		{
			name: "row heights count mismatch",
			table: NewTable([]int64{2000000}).
				AddRow([]string{"A"}).
				WithRowHeights([]int64{200000, 300000}),
			msg: "row heights count",
		},
		{
			name: "row height must be positive",
			table: NewTable([]int64{2000000}).
				AddRow([]string{"A"}).
				WithRowHeights([]int64{0}),
			msg: "height must be > 0",
		},
		{
			name: "negative margin",
			table: NewTable([]int64{2000000}).
				AddStyledRow([]TableCell{NewTableCell("A").WithMarginLeftPt(-1)}),
			msg: "left margin must be >= 0",
		},
		{
			name: "non-finite margin",
			table: NewTable([]int64{2000000}).
				AddStyledRow([]TableCell{NewTableCell("A").WithMarginTopPt(math.NaN())}),
			msg: "top margin must be finite",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Invalid Advanced Table").WithTable(tc.table)})
			if err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.msg) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
