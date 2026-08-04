package mathml

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestToOMMLSupportedConstructs(t *testing.T) {
	cases := []struct {
		name  string
		latex string
		want  []string
	}{
		{"literal", "x + 1", []string{"<m:t>x+1</m:t>"}},
		{"greek", `\alpha\beta`, []string{"<m:t>α</m:t>", "<m:t>β</m:t>"}},
		{"symbol", `a \leq b`, []string{"<m:t>≤</m:t>"}},
		{"superscript", "x^2", []string{"<m:sSup>", "<m:t>x</m:t>", "<m:sup>"}},
		{"subscript", "x_i", []string{"<m:sSub>", "<m:sub>"}},
		{"braced script", "x^{n+1}", []string{"<m:sSup>", "<m:t>n+1</m:t>"}},
		{"fraction", `\frac{a}{b}`, []string{"<m:f>", "<m:num>", "<m:den>"}},
		{"radical", `\sqrt{2}`, []string{"<m:rad>", `<m:degHide m:val="1"/>`}},
		{"nth root", `\sqrt[3]{x}`, []string{"<m:rad>", `<m:degHide m:val="0"/>`, "<m:t>3</m:t>"}},
		{"summation", `\sum_{i=1}^{n} i`, []string{"<m:nary>", `<m:chr m:val="∑"/>`, "<m:sub>", "<m:sup>"}},
		{"integral", `\int_0^1 x`, []string{`<m:chr m:val="∫"/>`}},
		{"nested fraction", `\frac{\alpha}{\sqrt{2}}`, []string{"<m:f>", "<m:rad>", "<m:t>α</m:t>"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToOMML(tc.latex)
			if err != nil {
				t.Fatalf("ToOMML(%q) returned error: %v", tc.latex, err)
			}
			if !strings.HasPrefix(got, "<m:oMath ") || !strings.HasSuffix(got, "</m:oMath>") {
				t.Fatalf("ToOMML(%q) is not wrapped in m:oMath: %s", tc.latex, got)
			}
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Errorf("ToOMML(%q) = %s, want it to contain %q", tc.latex, got, want)
				}
			}
		})
	}
}

// A script applied after an element that contains self-closing tags used to
// slice that element in half, producing XML no parser would accept.
func TestToOMMLScriptAfterNaryStaysWellFormed(t *testing.T) {
	cases := []string{
		`\int_0^\infty e^{-x} dx = 1`,
		`\sum_{i=1}^{n} i^2`,
		`\sqrt{2}^3`,
		`\frac{a}{b}^2`,
	}
	for _, latex := range cases {
		got, err := ToOMML(latex)
		if err != nil {
			t.Fatalf("ToOMML(%q) returned error: %v", latex, err)
		}
		if err := xml.Unmarshal([]byte(got), new(any)); err != nil {
			t.Errorf("ToOMML(%q) is not well-formed XML: %v\n%s", latex, err, got)
		}
	}
}

func TestToOMMLRejectsUnsupportedInput(t *testing.T) {
	cases := []string{
		`\unsupported{x}`,
		`{x`,
		`x}`,
		`\frac{a}`,
		`^2`,
		`\`,
	}
	for _, latex := range cases {
		if _, err := ToOMML(latex); err == nil {
			t.Errorf("ToOMML(%q) succeeded, want an error", latex)
		}
	}
}

func TestToOMMLEscapesMarkup(t *testing.T) {
	got, err := ToOMML("a < b & c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(got, "a < b") || !strings.Contains(got, "&lt;") || !strings.Contains(got, "&amp;") {
		t.Errorf("markup characters were not escaped: %s", got)
	}
}

func TestParagraphXMLCarriesChoiceAndFallback(t *testing.T) {
	got, err := ParagraphXML(`\frac{1}{2}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		"<mc:AlternateContent",
		`Requires="a14"`,
		"<a14:m>",
		"<m:oMathPara",
		"<mc:Fallback>",
		`\frac{1}{2}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("ParagraphXML output missing %q: %s", want, got)
		}
	}
	// The inner oMath must not redeclare the namespace the para already has.
	if strings.Count(got, `xmlns:m="`+MathNS+`"`) != 1 {
		t.Errorf("math namespace should be declared once: %s", got)
	}
}
