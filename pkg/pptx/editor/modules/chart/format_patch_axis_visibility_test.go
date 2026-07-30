package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const axisFixtureXML = `<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" ` +
	`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><c:chart><c:plotArea>` +
	`<c:catAx><c:axId val="1"/><c:scaling><c:orientation val="minMax"/></c:scaling>` +
	`<c:axPos val="b"/><c:crossAx val="2"/></c:catAx>` +
	`<c:valAx><c:axId val="2"/><c:scaling><c:orientation val="minMax"/></c:scaling>` +
	`<c:axPos val="l"/><c:crossAx val="1"/><c:crosses val="autoZero"/></c:valAx>` +
	`</c:plotArea></c:chart></c:chartSpace>`

func boolPtr(value bool) *bool        { return &value }
func floatPtr(value float64) *float64 { return &value }
func stringPtr(value string) *string  { return &value }

// Upstream #473 and #852: hiding an axis was impossible, so a chart always drew
// both axes.
func TestPatchAxisVisibilityHidesAndShows(t *testing.T) {
	hidden := PatchAxisVisibility(axisFixtureXML, common.ChartFormatUpdate{
		ValueAxisVisible: boolPtr(false),
	})
	if !strings.Contains(hidden, `<c:delete val="1"/>`) {
		t.Fatalf("expected the value axis to be deleted:\n%s", hidden)
	}
	// The axis element itself must survive: the series still reference its id.
	if !strings.Contains(hidden, `<c:valAx>`) || !strings.Contains(hidden, `<c:axId val="2"/>`) {
		t.Fatalf("the axis element was removed instead of hidden:\n%s", hidden)
	}
	// <c:delete> belongs between scaling and axPos.
	if strings.Index(hidden, `<c:delete`) > strings.Index(hidden, `<c:axPos val="l"`) {
		t.Fatalf("c:delete is out of schema order:\n%s", hidden)
	}

	shown := PatchAxisVisibility(hidden, common.ChartFormatUpdate{
		ValueAxisVisible: boolPtr(true),
	})
	if strings.Contains(shown, `<c:delete val="1"/>`) {
		t.Fatalf("expected the axis to be shown again:\n%s", shown)
	}
	if strings.Count(shown, "<c:delete") != 1 {
		t.Fatalf("expected the existing c:delete to be rewritten, not duplicated:\n%s", shown)
	}
}

// Upstream #329: long category labels need an angle.
func TestPatchAxisTickLabelRotation(t *testing.T) {
	patched := PatchAxisVisibility(axisFixtureXML, common.ChartFormatUpdate{
		CategoryAxisTickLabelRotation: floatPtr(-45),
	})
	if !strings.Contains(patched, `rot="-2700000"`) {
		t.Fatalf("expected -45 degrees as 60000ths:\n%s", patched)
	}
	if strings.Index(patched, "<c:txPr>") > strings.Index(patched, `<c:crossAx val="2"`) {
		t.Fatalf("c:txPr is out of schema order:\n%s", patched)
	}

	// A second call rewrites the angle instead of stacking another txPr.
	again := PatchAxisVisibility(patched, common.ChartFormatUpdate{
		CategoryAxisTickLabelRotation: floatPtr(90),
	})
	if strings.Count(again, "<c:txPr>") != 1 {
		t.Fatalf("expected one c:txPr, got:\n%s", again)
	}
	if !strings.Contains(again, `rot="5400000"`) {
		t.Fatalf("expected the angle to be rewritten:\n%s", again)
	}
}

// Upstream #349: a line chart should be able to start on the axis rather than
// between tick marks.
func TestPatchAxisCrossBetween(t *testing.T) {
	patched := PatchAxisVisibility(axisFixtureXML, common.ChartFormatUpdate{
		ValueAxisCrossBetween: stringPtr("midCat"),
	})
	if !strings.Contains(patched, `<c:crossBetween val="midCat"/>`) {
		t.Fatalf("expected crossBetween on the value axis:\n%s", patched)
	}
	categoryBlock, _, found := strings.Cut(patched, "</c:catAx>")
	if !found {
		t.Fatalf("fixture lost its category axis:\n%s", patched)
	}
	if strings.Contains(categoryBlock, "<c:crossBetween") {
		t.Fatalf("crossBetween belongs to CT_ValAx only:\n%s", patched)
	}
	if strings.Index(patched, "<c:crossBetween") < strings.Index(patched, `<c:crosses val="autoZero"`) {
		t.Fatalf("crossBetween must follow c:crosses:\n%s", patched)
	}
}

func TestValidateAxisVisibilityRejectsBadValues(t *testing.T) {
	if err := ValidateAxisVisibility(common.ChartFormatUpdate{
		ValueAxisCrossBetween: stringPtr("sometimes"),
	}); err == nil {
		t.Fatalf("expected an unknown crossBetween value to be refused")
	}
	if err := ValidateAxisVisibility(common.ChartFormatUpdate{
		CategoryAxisTickLabelRotation: floatPtr(720),
	}); err == nil {
		t.Fatalf("expected an out-of-range rotation to be refused")
	}
}

// The state reader reports what was written, so a caller can read a hidden axis
// back.
func TestAxisStateReportsVisibilityRotationAndCrossBetween(t *testing.T) {
	patched := PatchAxisVisibility(axisFixtureXML, common.ChartFormatUpdate{
		ValueAxisVisible:              boolPtr(false),
		CategoryAxisTickLabelRotation: floatPtr(-45),
		ValueAxisCrossBetween:         stringPtr("midCat"),
	})

	state := ExtractChartState([]byte(patched))
	if state.ValueAx.Visible == nil || *state.ValueAx.Visible {
		t.Fatalf("expected the value axis to read back as hidden, got %+v", state.ValueAx.Visible)
	}
	if state.CategoryAx.TickLabelRotation == nil || *state.CategoryAx.TickLabelRotation != -45 {
		t.Fatalf("expected -45 degrees, got %+v", state.CategoryAx.TickLabelRotation)
	}
	if state.ValueAx.CrossBetween != "midCat" {
		t.Fatalf("expected midCat, got %q", state.ValueAx.CrossBetween)
	}
}
