package chart

import (
	"math"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestPatchChartFormatting_AxisDetails(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:chart><c:plotArea><c:barChart><c:axId val="1"/><c:axId val="2"/></c:barChart>
<c:catAx><c:axId val="1"/><c:crosses val="autoZero"/></c:catAx>
<c:valAx><c:axId val="2"/><c:crosses val="autoZero"/></c:valAx>
</c:plotArea></c:chart></c:chartSpace>`)
	categoryTitle := "Quarter"
	valueTitle := "Revenue"
	minimum := 0.0
	maximum := 200.0
	majorUnit := 25.0
	minorUnit := 5.0
	numberFormat := "$#,##0.00"
	linked := false

	got, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		CategoryAxisTitle:     &categoryTitle,
		ValueAxisTitle:        &valueTitle,
		ValueAxisMinimumScale: &minimum,
		ValueAxisMaximumScale: &maximum,
		ValueAxisMajorUnit:    &majorUnit,
		ValueAxisMinorUnit:    &minorUnit,
		ValueAxisNumberFormat: &numberFormat,
		ValueAxisFormatLinked: &linked,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	for _, want := range []string{
		`<a:t>Quarter</a:t>`,
		`<a:t>Revenue</a:t>`,
		`<c:min val="0"/>`,
		`<c:max val="200"/>`,
		`<c:majorUnit val="25"/>`,
		`<c:minorUnit val="5"/>`,
		`<c:numFmt formatCode="$#,##0.00" sourceLinked="0"/>`,
	} {
		if !strings.Contains(updated, want) {
			t.Fatalf("updated XML missing %q: %s", want, updated)
		}
	}
	valueAxisStart := strings.Index(updated, "<c:valAx>")
	valueAxisEnd := strings.Index(updated, "</c:valAx>")
	if valueAxisStart < 0 || valueAxisEnd <= valueAxisStart {
		t.Fatalf("value axis missing from updated XML: %s", updated)
	}
	valueAxis := updated[valueAxisStart:valueAxisEnd]
	assertXMLOrder(
		t,
		valueAxis,
		"<c:axId",
		"<c:scaling",
		"<c:title",
		"<c:numFmt",
		"<c:crosses",
		"<c:majorUnit",
		"<c:minorUnit",
	)

	state := ExtractChartState(got)
	if state.CategoryAx.Title != categoryTitle || state.ValueAx.Title != valueTitle {
		t.Fatalf("axis titles not preserved in state: %#v", state)
	}
	if state.ValueAx.MinimumScale == nil || *state.ValueAx.MinimumScale != minimum ||
		state.ValueAx.MaximumScale == nil || *state.ValueAx.MaximumScale != maximum {
		t.Fatalf("axis scales not preserved in state: %#v", state.ValueAx)
	}
	if state.ValueAx.MajorUnit == nil || *state.ValueAx.MajorUnit != majorUnit ||
		state.ValueAx.MinorUnit == nil || *state.ValueAx.MinorUnit != minorUnit {
		t.Fatalf("axis units not preserved in state: %#v", state.ValueAx)
	}
	if state.ValueAx.NumberFormat != numberFormat || state.ValueAx.FormatLinked == nil || *state.ValueAx.FormatLinked {
		t.Fatalf("axis number format not preserved in state: %#v", state.ValueAx)
	}
}

func TestPatchChartFormatting_AxisTitleVisibility(t *testing.T) {
	xml := []byte(`<c:chartSpace xmlns:c="x" xmlns:a="y"><c:chart><c:plotArea>` +
		`<c:valAx><c:axId val="2"/><c:title><c:tx><c:rich><a:p><a:r>` +
		`<a:t>Revenue</a:t></a:r></a:p></c:rich></c:tx></c:title>` +
		`<c:crosses val="autoZero"/></c:valAx></c:plotArea></c:chart></c:chartSpace>`)
	hide := false
	got, err := PatchChartFormatting(
		xml,
		common.ChartFormatUpdate{ValueAxisHasTitle: &hide},
	)
	if err != nil {
		t.Fatalf("hide axis title: %v", err)
	}
	if strings.Contains(string(got), "<c:title>") {
		t.Fatalf("axis title was not removed: %s", got)
	}
	if ExtractChartState(got).ValueAx.HasTitle {
		t.Fatal("axis state still reports a title after removal")
	}

	show := true
	got, err = PatchChartFormatting(
		got,
		common.ChartFormatUpdate{ValueAxisHasTitle: &show},
	)
	if err != nil {
		t.Fatalf("show empty axis title: %v", err)
	}
	if !ExtractChartState(got).ValueAx.HasTitle {
		t.Fatal("axis state does not report the restored empty title")
	}
}

func TestPatchChartFormattingRejectsScaleAgainstRetainedBound(t *testing.T) {
	xml := []byte(`<c:chartSpace xmlns:c="x"><c:chart><c:plotArea>` +
		`<c:valAx><c:axId val="2"/><c:scaling><c:min val="0"/><c:max val="100"/></c:scaling></c:valAx>` +
		`</c:plotArea></c:chart></c:chartSpace>`)
	minimum := 100.0
	if _, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		ValueAxisMinimumScale: &minimum,
	}); err == nil {
		t.Fatal("expected minimum equal to retained maximum to be rejected")
	}
}

func TestPatchChartFormattingRejectsScaleAgainstEveryRetainedAxis(t *testing.T) {
	xml := []byte(`<c:chartSpace xmlns:c="x"><c:chart><c:plotArea>` +
		`<c:valAx><c:axId val="2"/><c:scaling><c:max val="100"/></c:scaling></c:valAx>` +
		`<c:valAx><c:axId val="3"/><c:scaling><c:max val="10"/></c:scaling></c:valAx>` +
		`</c:plotArea></c:chart></c:chartSpace>`)
	minimum := 50.0
	if _, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		ValueAxisMinimumScale: &minimum,
	}); err == nil {
		t.Fatal("expected minimum conflicting with the second retained maximum to be rejected")
	}
}

func TestValidateChartFormatUpdateRejectsInvalidAxisDetails(t *testing.T) {
	minimum, maximum := 10.0, 5.0
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		ValueAxisMinimumScale: &minimum,
		ValueAxisMaximumScale: &maximum,
	}); err == nil {
		t.Fatal("expected invalid scale range error")
	}

	zero := 0.0
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{ValueAxisMajorUnit: &zero}); err == nil {
		t.Fatal("expected invalid axis unit error")
	}

	nan := math.NaN()
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{ValueAxisMinimumScale: &nan}); err == nil {
		t.Fatal("expected non-finite axis scale error")
	}
}

func TestPatchChartFormatting_CategoryAxisTickDetail(t *testing.T) {
	xml := []byte(`<c:chartSpace xmlns:c="x"><c:chart><c:plotArea>` +
		`<c:catAx><c:axId val="1"/><c:crossAx val="2"/><c:crosses val="autoZero"/>` +
		`<c:auto val="1"/><c:lblOffset val="100"/><c:noMultiLvlLbl val="0"/></c:catAx>` +
		`<c:valAx><c:axId val="2"/><c:crossAx val="1"/></c:valAx>` +
		`</c:plotArea></c:chart></c:chartSpace>`)
	skip := 3
	alignment := "l"

	got, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		CategoryAxisTickMarkSkip:   &skip,
		CategoryAxisLabelAlignment: &alignment,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	updated := string(got)
	for _, want := range []string{`<c:tickMarkSkip val="3"/>`, `<c:lblAlgn val="l"/>`} {
		if !strings.Contains(updated, want) {
			t.Fatalf("updated XML missing %q: %s", want, updated)
		}
	}
	if strings.Contains(categoryAxisBlock(t, updated, "<c:valAx>", "</c:valAx>"), "<c:tickMarkSkip") ||
		strings.Contains(categoryAxisBlock(t, updated, "<c:valAx>", "</c:valAx>"), "<c:lblAlgn") {
		t.Fatalf("category-only nodes leaked onto the value axis: %s", updated)
	}
	assertXMLOrder(
		t,
		categoryAxisBlock(t, updated, "<c:catAx>", "</c:catAx>"),
		"<c:crosses",
		"<c:auto",
		"<c:lblAlgn",
		"<c:lblOffset",
		"<c:tickMarkSkip",
		"<c:noMultiLvlLbl",
	)

	state := ExtractChartState(got)
	if state.CategoryAx.TickMarkSkip == nil || *state.CategoryAx.TickMarkSkip != skip {
		t.Fatalf("tick mark skip not preserved in state: %#v", state.CategoryAx)
	}
	if state.CategoryAx.LabelAlignment != alignment {
		t.Fatalf("label alignment not preserved in state: %#v", state.CategoryAx)
	}

	// Re-patching the same values must be idempotent, not append duplicates.
	again, err := PatchChartFormatting(got, common.ChartFormatUpdate{
		CategoryAxisTickMarkSkip:   &skip,
		CategoryAxisLabelAlignment: &alignment,
	})
	if err != nil {
		t.Fatalf("re-patch error: %v", err)
	}
	if string(again) != updated {
		t.Fatalf("re-patch was not idempotent:\n%s\n%s", updated, string(again))
	}
}

func TestPatchChartFormattingInsertsCategoryAxisTickDetailWithoutSiblings(t *testing.T) {
	xml := []byte(`<c:chartSpace xmlns:c="x"><c:chart><c:plotArea>` +
		`<c:catAx><c:axId val="1"/><c:crossAx val="2"/></c:catAx>` +
		`</c:plotArea></c:chart></c:chartSpace>`)
	skip := 2
	alignment := "ctr"
	got, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		CategoryAxisTickMarkSkip:   &skip,
		CategoryAxisLabelAlignment: &alignment,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}
	block := categoryAxisBlock(t, string(got), "<c:catAx>", "</c:catAx>")
	assertXMLOrder(t, block, "<c:crossAx", "<c:lblAlgn", "<c:tickMarkSkip")
}

func TestValidateChartFormatUpdateRejectsInvalidCategoryAxisTickDetail(t *testing.T) {
	zero := 0
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		CategoryAxisTickMarkSkip: &zero,
	}); err == nil {
		t.Fatal("expected tick mark skip below 1 to be rejected")
	}
	bad := "center"
	if err := ValidateChartFormatUpdate(common.ChartFormatUpdate{
		CategoryAxisLabelAlignment: &bad,
	}); err == nil {
		t.Fatal("expected invalid label alignment to be rejected")
	}
}

func categoryAxisBlock(t *testing.T, xml string, startTag string, endTag string) string {
	t.Helper()
	start := strings.Index(xml, startTag)
	end := strings.Index(xml, endTag)
	if start < 0 || end <= start {
		t.Fatalf("axis block %s missing from XML: %s", startTag, xml)
	}
	return xml[start:end]
}
