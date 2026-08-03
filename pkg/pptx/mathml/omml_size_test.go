package mathml

import (
	"encoding/xml"
	"strings"
	"testing"
)

// The size has to land on the math runs. Run properties on a following
// paragraph do not flow backwards into an equation, so a size written only on
// endParaRPr left the formula at whatever the theme said.
func TestParagraphXMLSizedPutsSizeOnMathRuns(t *testing.T) {
	paragraph, err := ParagraphXMLSized(`x^2 + \frac{1}{2}`, 3200)
	if err != nil {
		t.Fatalf("render equation: %v", err)
	}

	runs := strings.Count(paragraph, "<m:r>")
	if runs == 0 {
		t.Fatal("expected the equation to contain math runs")
	}
	sized := strings.Count(paragraph, `<m:r><a:rPr lang="en-US" sz="3200" dirty="0"/>`)
	if sized != runs {
		t.Fatalf("expected all %d math runs sized, got %d:\n%s", runs, sized, paragraph)
	}
	if !strings.Contains(paragraph, `<mc:Fallback><a:r><a:rPr lang="en-US" sz="3200"`) {
		t.Errorf("the fallback run should carry the size too:\n%s", paragraph)
	}
}

// A zero size means "leave it to the theme", and must not write an sz at all.
func TestParagraphXMLWithoutSizeWritesNoSizeAttribute(t *testing.T) {
	paragraph, err := ParagraphXML(`a + b`)
	if err != nil {
		t.Fatalf("render equation: %v", err)
	}
	if strings.Contains(paragraph, "sz=") {
		t.Fatalf("expected no size attribute:\n%s", paragraph)
	}
}

// Injecting the run properties must not break the markup.
func TestParagraphXMLSizedStaysWellFormed(t *testing.T) {
	paragraph, err := ParagraphXMLSized(`\sum_{i=1}^{n} x_i^2`, 2400)
	if err != nil {
		t.Fatalf("render equation: %v", err)
	}
	wrapped := `<root xmlns:a="a" xmlns:m="m" xmlns:mc="mc" xmlns:a14="a14">` +
		paragraph + `</root>`
	var discard any
	if err := xml.Unmarshal([]byte(wrapped), &discard); err != nil {
		t.Fatalf("equation paragraph is not well formed: %v\n%s", err, paragraph)
	}
}
