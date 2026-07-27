package slide

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const chartExContentType = "application/vnd.ms-office.chartex+xml"

const contentTypesWithChartEx = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>
<Override PartName="/ppt/slides/slide1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>
<Override PartName="/ppt/charts/chart1.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.chart+xml"/>
<Override PartName="/ppt/charts/chartEx1.xml" ContentType="` + chartExContentType + `"/>
</Types>`

// A deck holding a box & whisker, waterfall, funnel, treemap, sunburst or
// histogram carries a chart extension part. It shares the ppt/charts/chart
// prefix with a classic chart but has its own content type, and rewriting it as
// a DrawingML chart makes PowerPoint reject the package.
func TestRewriteContentTypesPreservesChartExOverride(t *testing.T) {
	out, err := RewriteContentTypes(
		[]byte(contentTypesWithChartEx),
		[]common.EditorSlideRef{{Part: "ppt/slides/slide1.xml"}},
		nil,   // mediaPaths
		false, // hasSections
		FilterDrawingMLChartPaths([]string{"ppt/charts/chart1.xml", "ppt/charts/chartEx1.xml"}),
		nil, nil, nil, nil,
		false, false, nil, false, false, nil,
	)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}

	if !strings.Contains(out, `PartName="/ppt/charts/chartEx1.xml" ContentType="`+chartExContentType+`"`) {
		t.Fatalf("chartEx override lost or rewritten: %s", out)
	}
	if strings.Count(out, `PartName="/ppt/charts/chartEx1.xml"`) != 1 {
		t.Fatalf("chartEx override duplicated: %s", out)
	}
	if !strings.Contains(
		out,
		`PartName="/ppt/charts/chart1.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.chart+xml"`,
	) {
		t.Fatalf("classic chart override missing: %s", out)
	}
}

func TestFilterDrawingMLChartPathsExcludesChartEx(t *testing.T) {
	got := FilterDrawingMLChartPaths([]string{
		"ppt/charts/chart1.xml",
		"ppt/charts/chartEx1.xml",
		"ppt/charts/chart10.xml",
		"ppt/charts/chartEx22.xml",
	})
	want := []string{"ppt/charts/chart1.xml", "ppt/charts/chart10.xml"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
