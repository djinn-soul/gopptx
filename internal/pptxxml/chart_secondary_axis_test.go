package pptxxml

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var reAxisIDValue = regexp.MustCompile(`<c:axId val="(\d+)"/>`)

func comboSpec(secondary bool) *ChartSpec {
	return &ChartSpec{
		Kind:       ChartKindCombo,
		Title:      "Revenue vs Growth",
		Categories: []string{"Q1", "Q2"},
		BarSeries:  []ChartSeries{{Name: "Revenue", Values: []float64{820, 932}}},
		LineSeries: []ChartSeries{{Name: "Growth", Values: []float64{4.5, 13.7}}},

		SecondaryAxis:           secondary,
		SecondaryValueAxisTitle: "Growth %",
	}
}

func plotAxisIDs(t *testing.T, xml, plot string) []string {
	t.Helper()
	block := regexp.MustCompile(`(?s)<c:` + plot + `>.*?</c:` + plot + `>`).FindString(xml)
	if block == "" {
		t.Fatalf("no <c:%s> element in %s", plot, xml)
	}
	ids := make([]string, 0, 2)
	for _, match := range reAxisIDValue.FindAllStringSubmatch(block, -1) {
		ids = append(ids, match[1])
	}
	return ids
}

// Without a secondary axis the two plots share one axis pair, which is the
// behaviour every existing deck was written with.
func TestComboChartSharesOneAxisPairByDefault(t *testing.T) {
	xml := comboChartPartXML(comboSpec(false))

	barIDs := plotAxisIDs(t, xml, "barChart")
	lineIDs := plotAxisIDs(t, xml, "lineChart")
	if len(barIDs) != 2 || len(lineIDs) != 2 {
		t.Fatalf("each plot needs two axId refs, got bar=%v line=%v", barIDs, lineIDs)
	}
	if barIDs[0] != lineIDs[0] || barIDs[1] != lineIDs[1] {
		t.Fatalf("plots should share one pair, got bar=%v line=%v", barIDs, lineIDs)
	}
	if count := strings.Count(xml, "<c:valAx>"); count != 1 {
		t.Fatalf("expected one value axis, got %d", count)
	}
}

// With a secondary axis the line plot must reference its own pair, otherwise
// both series are still measured against the same scale.
func TestComboChartBindsLinePlotToSecondaryAxis(t *testing.T) {
	xml := comboChartPartXML(comboSpec(true))

	barIDs := plotAxisIDs(t, xml, "barChart")
	lineIDs := plotAxisIDs(t, xml, "lineChart")
	if barIDs[0] == lineIDs[0] || barIDs[1] == lineIDs[1] {
		t.Fatalf("line plot must use its own axis pair, got bar=%v line=%v", barIDs, lineIDs)
	}

	unique := map[string]bool{}
	for _, id := range append(append([]string{}, barIDs...), lineIDs...) {
		unique[id] = true
	}
	if len(unique) != 4 {
		t.Fatalf("expected four distinct axis ids, got %v", unique)
	}

	if got := strings.Count(xml, "<c:valAx>"); got != 2 {
		t.Fatalf("expected two value axes, got %d", got)
	}
	if got := strings.Count(xml, "<c:catAx>"); got != 2 {
		t.Fatalf("expected two category axes, got %d", got)
	}
}

// Every axis a plot points at must be declared, or PowerPoint repairs the file.
func TestComboChartDeclaresEveryReferencedAxis(t *testing.T) {
	xml := comboChartPartXML(comboSpec(true))

	plots := regexp.MustCompile(`(?s)<c:(bar|line)Chart>.*?</c:(bar|line)Chart>`)
	declared := map[string]bool{}
	for _, match := range reAxisIDValue.FindAllStringSubmatch(plots.ReplaceAllString(xml, ""), -1) {
		declared[match[1]] = true
	}
	for _, id := range append(plotAxisIDs(t, xml, "barChart"), plotAxisIDs(t, xml, "lineChart")...) {
		if !declared[id] {
			t.Fatalf("axis id %s is referenced by a plot but never declared", id)
		}
	}
	for _, want := range []int{primaryCatAxID, primaryValAxID, secondaryCatAxID, secondaryValAxID} {
		if !declared[strconv.Itoa(want)] {
			t.Fatalf("axis id %d was not declared", want)
		}
	}
}

// The secondary value axis is drawn opposite the primary one, and the category
// axis it crosses is hidden so the chart does not show two identical scales.
func TestComboChartSecondaryAxisGeometry(t *testing.T) {
	xml := comboChartPartXML(comboSpec(true))

	valueAxes := regexp.MustCompile(`(?s)<c:valAx>.*?</c:valAx>`).FindAllString(xml, -1)
	if len(valueAxes) != 2 {
		t.Fatalf("expected two value axes, got %d", len(valueAxes))
	}
	primary, secondary := valueAxes[0], valueAxes[1]
	if !strings.Contains(primary, `<c:axPos val="l"/>`) {
		t.Fatalf("primary value axis should stay on the left: %s", primary)
	}
	if !strings.Contains(secondary, `<c:axPos val="r"/>`) {
		t.Fatalf("secondary value axis should be on the right: %s", secondary)
	}
	if !strings.Contains(secondary, `<c:crosses val="max"/>`) {
		t.Fatalf("secondary value axis should cross at the far end: %s", secondary)
	}
	if !strings.Contains(secondary, "Growth %") {
		t.Fatalf("secondary value axis lost its title: %s", secondary)
	}
	if strings.Contains(primary, "Growth %") {
		t.Fatalf("the secondary title leaked onto the primary axis: %s", primary)
	}

	categoryAxes := regexp.MustCompile(`(?s)<c:catAx>.*?</c:catAx>`).FindAllString(xml, -1)
	if !strings.Contains(categoryAxes[0], `<c:delete val="0"/>`) {
		t.Fatalf("primary category axis must stay visible: %s", categoryAxes[0])
	}
	if !strings.Contains(categoryAxes[1], `<c:delete val="1"/>`) {
		t.Fatalf("duplicate category axis must be hidden: %s", categoryAxes[1])
	}
}

// Charts that never opt in must be byte-identical to what they were before the
// secondary axis existed.
func TestSinglePlotChartsKeepThePrimaryAxisPair(t *testing.T) {
	for _, tc := range []struct {
		name string
		xml  string
	}{
		{"bar", barChartPartXML(&ChartSpec{Kind: ChartKindBar, Categories: []string{"a"}, Values: []float64{1}})},
		{"line", lineChartPartXML(&ChartSpec{Kind: ChartKindLine, Categories: []string{"a"}, Values: []float64{1}})},
		{"area", areaChartPartXML(&ChartSpec{Kind: ChartKindArea, Categories: []string{"a"}, Values: []float64{1}})},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wantCat := `<c:axId val="` + strconv.Itoa(primaryCatAxID) + `"/>`
			wantVal := `<c:axId val="` + strconv.Itoa(primaryValAxID) + `"/>`
			if !strings.Contains(tc.xml, wantCat) || !strings.Contains(tc.xml, wantVal) {
				t.Fatalf("expected the primary axis pair, got %s", tc.xml)
			}
			if strings.Contains(tc.xml, strconv.Itoa(secondaryValAxID)) {
				t.Fatalf("a single-plot chart must not reference the secondary axis: %s", tc.xml)
			}
			if got := strings.Count(tc.xml, "<c:valAx>"); got != 1 {
				t.Fatalf("expected one value axis, got %d", got)
			}
		})
	}
}
