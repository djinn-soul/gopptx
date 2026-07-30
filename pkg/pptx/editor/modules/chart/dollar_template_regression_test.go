package chart

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Regexp.ReplaceAllString reads "$1", "$0" and "${name}" in the *replacement*
// as capture-group references. Currency format codes and "$"-prefixed bucket
// labels are ordinary chart inputs, so every patch site that embeds caller text
// must replace literally.

func TestPatchChartFormatting_DollarInAxisNumberFormatIsLiteral(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:chart><c:plotArea><c:barChart><c:axId val="1"/><c:axId val="2"/></c:barChart>
<c:catAx><c:axId val="1"/><c:crosses val="autoZero"/></c:catAx>
<c:valAx><c:axId val="2"/><c:numFmt formatCode="General" sourceLinked="1"/><c:crosses val="autoZero"/></c:valAx>
</c:plotArea></c:chart></c:chartSpace>`)
	numberFormat := "$0.00"

	got, err := PatchChartFormatting(xml, common.ChartFormatUpdate{
		ValueAxisNumberFormat: &numberFormat,
	})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}

	updated := string(got)
	if !strings.Contains(updated, `<c:numFmt formatCode="$0.00" sourceLinked="1"/>`) {
		t.Errorf("format code was expanded as a template:\n%s", updated)
	}
	if strings.Contains(updated, `formatCode="<c:numFmt`) {
		t.Errorf("matched node was spliced into the attribute value:\n%s", updated)
	}
}

func TestPatchChartFormatting_DollarInTitleIsLiteral(t *testing.T) {
	xml := []byte(`<?xml version="1.0"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<c:chart><c:title><c:tx><c:rich><a:p><a:r><a:t>Old</a:t></a:r></a:p></c:rich></c:tx></c:title>
<c:plotArea><c:barChart><c:axId val="1"/><c:axId val="2"/></c:barChart></c:plotArea></c:chart></c:chartSpace>`)
	title := "Revenue $1M"

	got, err := PatchChartFormatting(xml, common.ChartFormatUpdate{Title: &title})
	if err != nil {
		t.Fatalf("PatchChartFormatting error: %v", err)
	}

	if updated := string(got); !strings.Contains(updated, `<a:t>Revenue $1M</a:t>`) {
		t.Errorf("title lost its $-reference:\n%s", updated)
	}
}

func TestApplyStringValues_DollarCategoriesAreLiteral(t *testing.T) {
	field := `<c:cat><c:strCache><c:ptCount val="0"/></c:strCache></c:cat>`

	out, err := applyStringValues("c:cat", field, []string{"$0-10k", "$1M+"})
	if err != nil {
		t.Fatalf("applyStringValues error: %v", err)
	}

	for _, want := range []string{`<c:v>$0-10k</c:v>`, `<c:v>$1M+</c:v>`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "<c:strCache><c:strCache") {
		t.Errorf("matched node was spliced into a value:\n%s", out)
	}
}

func TestApplyStringValues_DollarFormatCodeSurvivesNumericFallback(t *testing.T) {
	field := `<c:val><c:numCache><c:formatCode>$#,##0.00</c:formatCode>` +
		`<c:ptCount val="0"/></c:numCache></c:val>`

	out, err := applyNumericValues("c:val", field, []float64{1, 2})
	if err != nil {
		t.Fatalf("applyNumericValues error: %v", err)
	}

	if !strings.Contains(out, `<c:formatCode>$#,##0.00</c:formatCode>`) {
		t.Errorf("existing $-bearing format code was expanded:\n%s", out)
	}
}
