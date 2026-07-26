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
