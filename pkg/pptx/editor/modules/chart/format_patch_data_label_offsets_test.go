package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const doughnutChartXML = `<?xml version="1.0"?><c:chartSpace xmlns:c="http://x" xmlns:a="http://y">` +
	`<c:chart><c:plotArea><c:doughnutChart>` +
	`<c:ser><c:idx val="0"/><c:order val="0"/>` +
	`<c:cat><c:strRef><c:f>Sheet1!$A$1:$A$3</c:f></c:strRef></c:cat>` +
	`<c:val><c:numRef><c:f>Sheet1!$B$1:$B$3</c:f></c:numRef></c:val>` +
	`</c:ser></c:doughnutChart></c:plotArea></c:chart></c:chartSpace>`

func floatPtrDL(v float64) *float64 { return &v }

// Issue #1025: a doughnut label is moved with a manual layout, because the
// chart type rejects most c:dLblPos values.
func TestPatchDataLabelOffsetsCreatesManualLayout(t *testing.T) {
	out, err := PatchChartFormatting([]byte(doughnutChartXML), common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 1, X: floatPtrDL(-0.08), Y: floatPtrDL(0.05)},
		},
	})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		"<c:dLbls>", "<c:dLbl>", `<c:idx val="1"/>`,
		"<c:manualLayout>", `<c:x val="-0.08"/>`, `<c:y val="0.05"/>`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %s in %s", want, got)
		}
	}
	if strings.Contains(got, "<c:dLblPos") {
		t.Fatalf("must not emit dLblPos for a doughnut: %s", got)
	}

	// CT_DLbl orders idx before layout.
	label := got[strings.Index(got, "<c:dLbl>"):strings.Index(got, "</c:dLbl>")]
	if strings.Index(label, "<c:idx") > strings.Index(label, "<c:layout>") {
		t.Fatalf("idx must precede layout: %s", label)
	}
}

func TestPatchDataLabelOffsetsUpdatesExistingLabel(t *testing.T) {
	first, err := PatchChartFormatting([]byte(doughnutChartXML), common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 0, X: floatPtrDL(0.1), Y: floatPtrDL(0.2)},
		},
	})
	if err != nil {
		t.Fatalf("first patch: %v", err)
	}
	second, err := PatchChartFormatting(first, common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 0, X: floatPtrDL(0.3)},
		},
	})
	if err != nil {
		t.Fatalf("second patch: %v", err)
	}

	got := string(second)
	if strings.Count(got, "<c:dLbl>") != 1 {
		t.Fatalf("label duplicated: %s", got)
	}
	if !strings.Contains(got, `<c:x val="0.3"/>`) {
		t.Fatalf("x not updated: %s", got)
	}
	if !strings.Contains(got, `<c:y val="0.2"/>`) {
		t.Fatalf("y should be preserved: %s", got)
	}
}

func TestPatchDataLabelOffsetsOrdersPointsAscending(t *testing.T) {
	out, err := PatchChartFormatting([]byte(doughnutChartXML), common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 2, X: floatPtrDL(0.3)},
			{PointIndex: 0, X: floatPtrDL(0.1)},
		},
	})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	got := string(out)
	if strings.Index(got, `<c:idx val="0"/><c:layout>`) >
		strings.Index(got, `<c:idx val="2"/><c:layout>`) {
		t.Fatalf("labels should be ascending by point index: %s", got)
	}
}

// A <c:dLbl> does not inherit the series flags, so a moved label must repeat
// them or it loses its category name and percentage.
func TestPatchDataLabelOffsetsInheritsSeriesDisplayFlags(t *testing.T) {
	show := true
	withLabels, err := PatchChartFormatting([]byte(doughnutChartXML), common.ChartFormatUpdate{
		ShowDataLabels:        &show,
		DataLabelShowCategory: &show,
		DataLabelShowPercent:  &show,
	})
	if err != nil {
		t.Fatalf("enable labels: %v", err)
	}

	moved, err := PatchChartFormatting(withLabels, common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{PointIndex: 1, X: floatPtrDL(-0.1)},
		},
	})
	if err != nil {
		t.Fatalf("move label: %v", err)
	}

	got := string(moved)
	label := got[strings.Index(got, "<c:dLbl>"):strings.Index(got, "</c:dLbl>")]
	for _, want := range []string{`<c:showCatName val="1"/>`, `<c:showPercent val="1"/>`} {
		if !strings.Contains(label, want) {
			t.Fatalf("moved label dropped %s: %s", want, label)
		}
	}
}

func TestPatchDataLabelOffsetsIgnoresUnknownSeries(t *testing.T) {
	out, err := PatchChartFormatting([]byte(doughnutChartXML), common.ChartFormatUpdate{
		DataLabelOffsets: []common.DataLabelOffset{
			{SeriesIndex: 3, PointIndex: 0, X: floatPtrDL(0.1)},
		},
	})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if strings.Contains(string(out), "<c:dLbl>") {
		t.Fatalf("no label expected for a missing series: %s", out)
	}
}
