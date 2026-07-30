package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestCacheOnlyPatchPreservesExternalFormulas(t *testing.T) {
	source := strings.ReplaceAll(
		twoSeriesChartXML,
		"Sheet1!",
		"'C:\\Reports\\[Sales.xlsx]Data'!",
	)
	req := common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Values: []float64{10, 20}},
			{Values: []float64{30, 40}},
		},
	}

	out, err := PatchChartDataCache([]byte(source), KindCategory, req, CachePatchOptions{KeepFormulas: true})
	if err != nil {
		t.Fatalf("patch cache: %v", err)
	}

	original := xmlFormulaPattern.FindAllString(source, -1)
	updated := xmlFormulaPattern.FindAllString(string(out), -1)
	if len(original) != len(updated) {
		t.Fatalf("formula count changed: %d -> %d", len(original), len(updated))
	}
	for i := range original {
		if updated[i] != original[i] {
			t.Fatalf("formula %d changed: %s -> %s", i, original[i], updated[i])
		}
	}
	if !strings.Contains(string(out), "<c:v>40</c:v>") {
		t.Fatalf("cache values were not refreshed: %s", out)
	}
}

func TestDataTablePartialPatchPreservesShapeAndRichTextProperties(t *testing.T) {
	fixture := strings.Replace(
		string(dataTableChartXML()),
		`<c:spPr><a:noFill/></c:spPr></c:plotArea>`,
		`<c:dTable><c:showHorzBorder val="1"/><c:showVertBorder val="1"/>`+
			`<c:showOutline val="1"/><c:showKeys val="1"/>`+
			`<c:spPr><a:solidFill><a:srgbClr val="ABCDEF"/></a:solidFill></c:spPr>`+
			`<c:txPr><a:bodyPr rot="60000"/><a:lstStyle/><a:p><a:pPr>`+
			`<a:defRPr sz="900" b="1"><a:latin typeface="Aptos"/></a:defRPr>`+
			`</a:pPr></a:p></c:txPr></c:dTable>`+
			`<c:spPr><a:noFill/></c:spPr></c:plotArea>`,
		1,
	)
	showKeys := false

	out, err := PatchChartFormatting([]byte(fixture), common.ChartFormatUpdate{
		DataTable: &common.ChartDataTable{Show: true, ShowKeys: &showKeys},
	})
	if err != nil {
		t.Fatalf("patch data table: %v", err)
	}
	block := reDataTableBlock.FindString(string(out))
	for _, want := range []string{
		`<a:srgbClr val="ABCDEF"/>`,
		`<a:bodyPr rot="60000"/>`,
		`<a:defRPr sz="900" b="1"><a:latin typeface="Aptos"/></a:defRPr>`,
	} {
		if !strings.Contains(block, want) {
			t.Fatalf("partial patch dropped %s: %s", want, block)
		}
	}
}

func TestDataTableFontPatchPreservesOtherTextFormatting(t *testing.T) {
	fontSize := 12
	existing := `<c:dTable><c:txPr><a:bodyPr rot="60000"/><a:lstStyle/>` +
		`<a:p><a:pPr><a:defRPr sz="900" b="1"><a:latin typeface="Aptos"/>` +
		`</a:defRPr></a:pPr></a:p></c:txPr></c:dTable>`

	got := buildDataTableBlock(existing, common.ChartDataTable{
		Show: true, FontSizePt: &fontSize,
	})

	for _, want := range []string{
		`<a:bodyPr rot="60000"/>`,
		`<a:defRPr sz="1200" b="1">`,
		`<a:latin typeface="Aptos"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("font patch dropped %s: %s", want, got)
		}
	}
}

func TestTitleCoordinatePatchPreservesManualLayoutDetails(t *testing.T) {
	fixture := strings.Replace(
		titledChartXML,
		`<c:overlay val="0"/>`,
		`<c:layout><c:manualLayout><c:layoutTarget val="inner"/>`+
			`<c:xMode val="factor"/><c:yMode val="edge"/>`+
			`<c:x val="0.1"/><c:y val="0.2"/><c:w val="0.7"/><c:h val="0.4"/>`+
			`</c:manualLayout></c:layout><c:overlay val="0"/>`,
		1,
	)

	out, err := PatchChartFormatting([]byte(fixture), common.ChartFormatUpdate{
		TitleX: floatPtrTL(0.3),
	})
	if err != nil {
		t.Fatalf("patch title: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		`<c:layoutTarget val="inner"/>`, `<c:xMode val="factor"/>`,
		`<c:yMode val="edge"/>`, `<c:x val="0.3"/>`,
		`<c:y val="0.2"/>`, `<c:w val="0.7"/>`, `<c:h val="0.4"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("title patch dropped %s: %s", want, got)
		}
	}
}

func TestDataLabelOffsetPreservesExistingLabelFormatting(t *testing.T) {
	fixture := strings.Replace(
		doughnutChartXML,
		`<c:cat>`,
		`<c:dLbls><c:dLbl><c:idx val="1"/><c:numFmt formatCode="0.0%" sourceLinked="0"/>`+
			`<c:tx><c:rich><a:bodyPr/><a:p><a:r><a:t>Custom</a:t></a:r></a:p></c:rich></c:tx>`+
			`<c:spPr><a:solidFill><a:srgbClr val="FF0000"/></a:solidFill></c:spPr>`+
			`<c:txPr><a:bodyPr/><a:p><a:pPr><a:defRPr b="1"/></a:pPr></a:p></c:txPr>`+
			`<c:separator>; </c:separator></c:dLbl></c:dLbls><c:cat>`,
		1,
	)

	out, err := PatchChartFormatting([]byte(fixture), common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 1, X: floatPtrDL(0.2), Y: floatPtrDL(-0.1)},
		},
	})
	if err != nil {
		t.Fatalf("patch label: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		`formatCode="0.0%"`, `<a:t>Custom</a:t>`, `val="FF0000"`,
		`<a:defRPr b="1"/>`, `<c:separator>; </c:separator>`,
		`<c:x val="0.2"/>`, `<c:y val="-0.1"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("label patch dropped %s: %s", want, got)
		}
	}
}
