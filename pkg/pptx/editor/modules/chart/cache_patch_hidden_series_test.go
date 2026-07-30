package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const twoSeriesChartXML = `<?xml version="1.0"?><c:chartSpace xmlns:c="http://x" xmlns:a="http://y">` +
	`<c:chart><c:plotArea><c:barChart>` +
	`<c:ser><c:idx val="0"/><c:order val="0"/>` +
	`<c:cat><c:strRef><c:f>Sheet1!$A$2:$A$3</c:f><c:strCache><c:ptCount val="2"/></c:strCache></c:strRef></c:cat>` +
	`<c:val><c:numRef><c:f>Sheet1!$B$2:$B$3</c:f><c:numCache><c:ptCount val="2"/></c:numCache></c:numRef></c:val>` +
	`</c:ser>` +
	`<c:ser><c:idx val="1"/><c:order val="1"/>` +
	`<c:cat><c:strRef><c:f>Sheet1!$A$2:$A$3</c:f><c:strCache><c:ptCount val="2"/></c:strCache></c:strRef></c:cat>` +
	`<c:val><c:numRef><c:f>Sheet1!$C$2:$C$3</c:f><c:numCache><c:ptCount val="2"/></c:numCache></c:numRef></c:val>` +
	`</c:ser>` +
	`</c:barChart></c:plotArea></c:chart></c:chartSpace>`

func seriesName(name string) *string { return &name }

// Issue #1043: the workbook behind the chart may hold more columns than the
// chart plots.
func TestPatchChartDataCacheSkipsHiddenSeries(t *testing.T) {
	req := common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Name: seriesName("Plotted A"), Values: []float64{1, 2}},
			{Name: seriesName("Backing only"), Values: []float64{9, 9}, Hidden: true},
			{Name: seriesName("Plotted B"), Values: []float64{3, 4}},
		},
	}

	out, err := PatchChartDataCache([]byte(twoSeriesChartXML), KindCategory, req, CachePatchOptions{})
	if err != nil {
		t.Fatalf("patch: %v", err)
	}

	got := string(out)
	if strings.Count(got, "<c:ser>") != 2 {
		t.Fatalf("chart should still hold two series: %s", got)
	}
	// The second plotted series is payload index 2, so it reads column D:
	// B is the first series, C belongs to the hidden one.
	if !strings.Contains(got, "Sheet1!$B$2:$B$3") {
		t.Fatalf("first plotted series should read column B: %s", got)
	}
	if !strings.Contains(got, "Sheet1!$D$2:$D$3") {
		t.Fatalf("second plotted series should skip the hidden column and read D: %s", got)
	}
	if strings.Contains(got, "<c:v>9</c:v>") {
		t.Fatalf("hidden series values must not be plotted: %s", got)
	}
}

func TestPlottedSeriesKeepsWorkbookIndex(t *testing.T) {
	plotted := PlottedSeries([]common.ChartSeriesData{
		{Hidden: true},
		{Name: seriesName("visible")},
		{Hidden: true},
		{Name: seriesName("also visible")},
	})
	if len(plotted) != 2 {
		t.Fatalf("expected 2 plotted series, got %d", len(plotted))
	}
	if plotted[0].WorkbookIndex != 1 || plotted[1].WorkbookIndex != 3 {
		t.Fatalf("workbook indexes wrong: %+v", plotted)
	}
}

func TestPatchChartDataCacheCountsPlottedSeriesOnly(t *testing.T) {
	req := common.ChartDataUpdate{
		Categories: []string{"Q1", "Q2"},
		Series: []common.ChartSeriesData{
			{Name: seriesName("only one"), Values: []float64{1, 2}},
			{Name: seriesName("hidden"), Values: []float64{5, 6}, Hidden: true},
		},
	}
	_, err := PatchChartDataCache([]byte(twoSeriesChartXML), KindCategory, req, CachePatchOptions{})
	if err == nil {
		t.Fatal("expected a mismatch error: the chart draws two series, the payload plots one")
	}
	if !strings.Contains(err.Error(), "plotted") {
		t.Fatalf("error should name the plotted count: %v", err)
	}
}
