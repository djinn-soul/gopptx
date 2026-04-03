package pptxxml

import (
	"strings"
	"testing"
)

func TestRenderChart_AllKinds(t *testing.T) {
	categories := []string{"A", "B"}
	values := []float64{10.5, 20.0}
	xValues := []float64{1.0, 2.0}
	bubbleSizes := []float64{5.0, 10.0}

	kinds := []string{
		ChartKindBar,
		ChartKindBarHorizontal,
		ChartKindBarStacked,
		ChartKindBarStacked100,
		ChartKindLine,
		ChartKindLineMarkers,
		ChartKindLineStacked,
		ChartKindScatter,
		ChartKindArea,
		ChartKindAreaStacked,
		ChartKindAreaStacked100,
		ChartKindPie,
		ChartKindDoughnut,
		ChartKindBubble,
		ChartKindRadar,
		ChartKindRadarFilled,
		ChartKindStockHLC,
		ChartKindStockOHLC,
		ChartKindCombo,
	}

	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			spec := &ChartSpec{
				Kind:        kind,
				Title:       "Test Chart",
				Categories:  categories,
				Values:      values,
				XValues:     xValues,
				BubbleSizes: bubbleSizes,
				HighValues:  values,
				LowValues:   values,
				CloseValues: values,
				OpenValues:  values,
				BarSeries: []ChartSeries{
					{Name: "Bar1", Values: values},
				},
				LineSeries: []ChartSeries{
					{Name: "Line1", Values: values},
				},
			}
			xml := string(RenderChart(spec))
			if !strings.Contains(xml, "Test Chart") {
				t.Errorf("expected title 'Test Chart' in XML for kind %s", kind)
			}
			if !strings.Contains(xml, "http://schemas.openxmlformats.org/drawingml/2006/chart") {
				t.Errorf("expected chart namespace in XML for kind %s", kind)
			}
		})
	}
}

func TestRenderChart_Options(t *testing.T) {
	minValue := 0.0
	maxValue := 100.0
	spec := &ChartSpec{
		Kind:                       ChartKindBar,
		Title:                      "Options Chart",
		TitleOverlay:               true,
		ShowLegend:                 true,
		LegendPosition:             "b",
		LegendOverlay:              true,
		ShowDataLabels:             true,
		ShowMajorGridlines:         true,
		ShowCategoryMajorGridlines: true,
		CategoryAxisTitle:          "Cats",
		ValueAxisTitle:             "Vals",
		CategoryTickLabelPosition:  "low",
		ValueTickLabelPosition:     "high",
		CategoryAxisCrosses:        "min",
		ValueAxisCrosses:           "max",
		ValueFormat:                "0.00",
		ValueAxisCrossBetween:      "midCat",
		MinValue:                   &minValue,
		MaxValue:                   &maxValue,
		Categories:                 []string{"A"},
		Values:                     []float64{50},
		Smooth:                     true,
		ShowMarkers:                true,
	}

	xml := string(RenderChart(spec))

	checks := []string{
		`<c:overlay val="1"/>`,
		`<c:legendPos val="b"/>`,
		`<c:showVal val="1"/>`,
		`<c:majorGridlines/>`,
		`<c:tickLblPos val="low"/>`,
		`<c:tickLblPos val="high"/>`,
		`<c:crosses val="min"/>`,
		`<c:crosses val="max"/>`,
		`<a:t>Cats</a:t>`,
		`<a:t>Vals</a:t>`,
		`formatCode="0.00"`,
		`<c:crossBetween val="midCat"/>`,
		`<c:min val="0.000000"/>`,
		`<c:max val="100.000000"/>`,
	}

	for _, check := range checks {
		if !strings.Contains(xml, check) {
			t.Errorf("expected XML to contain %s", check)
		}
	}

	// Test line chart specific options
	spec.Kind = ChartKindLine
	xml = string(RenderChart(spec))
	if !strings.Contains(xml, `<c:smooth val="1"/>`) {
		t.Error("expected smooth val=1 for line chart")
	}
	if !strings.Contains(xml, `<c:marker>`) {
		t.Error("expected marker for line chart with ShowMarkers=true")
	}
}

func TestChartFrameShape(t *testing.T) {
	frame := &ChartFrame{
		RelID:        "rId1",
		X:            100,
		Y:            200,
		CX:           300,
		CY:           400,
		AltText:      "Alt Text",
		IsDecorative: false,
	}
	xml := chartFrameShape(frame, 5)
	if !strings.Contains(xml, `r:id="rId1"`) {
		t.Error("expected rId1 in chart frame XML")
	}
	if !strings.Contains(xml, `id="5"`) {
		t.Error("expected id=5 in chart frame XML")
	}
	if !strings.Contains(xml, `descr="Alt Text"`) {
		t.Error("expected Alt Text in chart frame XML")
	}
}

func TestChartPartXML(t *testing.T) {
	spec := &ChartSpec{
		Kind:       ChartKindPie,
		Categories: []string{"A"},
		Values:     []float64{1.0},
	}
	xml := ChartPartXML(spec)
	if !strings.Contains(xml, "<c:pieChart>") {
		t.Error("expected pieChart in ChartPartXML")
	}
}

func TestNormalizedLegendPosition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"l", "l"},
		{"t", "t"},
		{"b", "b"},
		{"r", "r"},
		{"  L  ", "l"},
		{"unknown", "r"},
		{"", "r"},
	}
	for _, tt := range tests {
		if got := normalizedLegendPosition(tt.input); got != tt.expected {
			t.Errorf("normalizedLegendPosition(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestChartValueFormatXML(t *testing.T) {
	if got := chartValueFormatXML(""); !strings.Contains(got, "General") {
		t.Error("expected General for empty format")
	}
	if got := chartValueFormatXML("General"); !strings.Contains(got, `sourceLinked="1"`) {
		t.Error("expected sourceLinked=1 for General format")
	}
}

func TestRenderChart_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("RenderChart did not panic on unsupported kind")
		}
	}()
	RenderChart(&ChartSpec{Kind: "invalid"})
}

func TestRenderChart_RadarAxisTicksAreNotDuplicated(t *testing.T) {
	spec := &ChartSpec{
		Kind:       ChartKindRadar,
		Title:      "Radar",
		Categories: []string{"A", "B"},
		Values:     []float64{1, 2},
	}
	xml := string(RenderChart(spec))
	if strings.Contains(xml, `<c:crosses val="autoZero"/><c:majorTickMark`) {
		t.Fatalf("unexpected duplicate majorTickMark after crosses in radar valAx: %s", xml)
	}
}
