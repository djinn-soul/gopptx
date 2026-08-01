package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// twoValueAxisChart is the shape of a combo chart with a secondary axis: two
// c:valAx elements, the first primary and the second secondary.
const twoValueAxisChart = `<c:chartSpace><c:chart><c:plotArea>` +
	`<c:catAx><c:axId val="1"/><c:delete val="0"/><c:axPos val="b"/>` +
	`<c:crossAx val="2"/></c:catAx>` +
	`<c:valAx><c:axId val="2"/><c:delete val="0"/><c:axPos val="l"/>` +
	`<c:tickLblPos val="nextTo"/><c:crossAx val="1"/></c:valAx>` +
	`<c:valAx><c:axId val="4"/><c:delete val="0"/><c:axPos val="r"/>` +
	`<c:title><c:tx><c:rich><a:p><a:r><a:t>SECONDARY</a:t></a:r></a:p></c:rich></c:tx></c:title>` +
	`<c:tickLblPos val="nextTo"/><c:crossAx val="3"/></c:valAx>` +
	`</c:plotArea></c:chart></c:chartSpace>`

func valueAxisBlocks(t *testing.T, xml string) []string {
	t.Helper()
	blocks := make([]string, 0, 2)
	remaining := xml
	for {
		start := strings.Index(remaining, "<c:valAx>")
		if start < 0 {
			break
		}
		end := strings.Index(remaining[start:], "</c:valAx>")
		if end < 0 {
			break
		}
		end += start + len("</c:valAx>")
		blocks = append(blocks, remaining[start:end])
		remaining = remaining[end:]
	}
	if len(blocks) != 2 {
		t.Fatalf("expected two value axes, got %d", len(blocks))
	}
	return blocks
}

func secondaryAxisStrPtr(v string) *string { return &v }

// The ValueAxis* request fields describe the primary axis. Writing them onto
// every c:valAx would retitle and hide the secondary axis as a side effect.
func TestValueAxisTitleDoesNotLeakToTheSecondaryAxis(t *testing.T) {
	out := patchAxisDetails(twoValueAxisChart, common.ChartFormatUpdate{
		ValueAxisHasTitle: boolPtr(true),
		ValueAxisTitle:    secondaryAxisStrPtr("PRIMARY"),
	})

	blocks := valueAxisBlocks(t, out)
	if !strings.Contains(blocks[0], "PRIMARY") {
		t.Fatalf("primary axis did not get the title: %s", blocks[0])
	}
	if strings.Contains(blocks[1], "PRIMARY") {
		t.Fatalf("the primary title leaked onto the secondary axis: %s", blocks[1])
	}
	if !strings.Contains(blocks[1], "SECONDARY") {
		t.Fatalf("the secondary axis lost its own title: %s", blocks[1])
	}
}

func TestValueAxisVisibilityDoesNotHideTheSecondaryAxis(t *testing.T) {
	out := PatchAxisVisibility(twoValueAxisChart, common.ChartFormatUpdate{
		ValueAxisVisible: boolPtr(false),
	})

	blocks := valueAxisBlocks(t, out)
	if !strings.Contains(blocks[0], `<c:delete val="1"/>`) {
		t.Fatalf("primary axis should be hidden: %s", blocks[0])
	}
	if !strings.Contains(blocks[1], `<c:delete val="0"/>`) {
		t.Fatalf("secondary axis should stay visible: %s", blocks[1])
	}
}

func TestValueAxisScaleAppliesToThePrimaryAxisOnly(t *testing.T) {
	minimum, maximum := 0.0, 100.0
	out := patchAxisDetails(twoValueAxisChart, common.ChartFormatUpdate{
		ValueAxisMinimumScale: &minimum,
		ValueAxisMaximumScale: &maximum,
	})

	blocks := valueAxisBlocks(t, out)
	if !strings.Contains(blocks[0], `<c:max val="100"/>`) {
		t.Fatalf("primary axis missing the new scale: %s", blocks[0])
	}
	if strings.Contains(blocks[1], "<c:max") {
		t.Fatalf("the scale leaked onto the secondary axis: %s", blocks[1])
	}
}

// A chart with a single value axis must behave exactly as it did before.
func TestSingleValueAxisStillPatched(t *testing.T) {
	single := `<c:chartSpace><c:chart><c:plotArea>` +
		`<c:valAx><c:axId val="2"/><c:delete val="0"/><c:axPos val="l"/>` +
		`<c:tickLblPos val="nextTo"/><c:crossAx val="1"/></c:valAx>` +
		`</c:plotArea></c:chart></c:chartSpace>`

	out := patchAxisDetails(single, common.ChartFormatUpdate{
		ValueAxisHasTitle: boolPtr(true),
		ValueAxisTitle:    secondaryAxisStrPtr("ONLY"),
	})
	if !strings.Contains(out, "ONLY") {
		t.Fatalf("the only value axis should have been patched: %s", out)
	}
}
