package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func paragraphXML(runs ...string) []byte {
	var b strings.Builder
	b.WriteString(`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:p>`)
	for _, run := range runs {
		b.WriteString(`<a:r><a:rPr lang="en-US"/><a:t>`)
		b.WriteString(run)
		b.WriteString(`</a:t></a:r>`)
	}
	b.WriteString(`</a:p></p:sld>`)
	return []byte(b.String())
}

func TestReplaceTextRuns_MatchesAcrossRuns(t *testing.T) {
	cases := []struct {
		name     string
		runs     []string
		find     string
		replace  string
		want     int
		wantText string
	}{
		{"two runs", []string{"Hello ", "World"}, "Hello World", "Hello Mars", 1, "Hello Mars"},
		{"three runs", []string{"Hel", "lo Wo", "rld"}, "Hello World", "Bye", 1, "Bye"},
		{"single run still works", []string{"Hello World"}, "World", "Mars", 1, "Hello Mars"},
		{"replacement longer", []string{"a", "bc"}, "abc", "abcdef", 1, "abcdef"},
		{"replacement shorter", []string{"abc", "def"}, "abcdef", "x", 1, "x"},
		{"split template tag", []string{"Hello {{ na", "me }}"}, "{{ name }}", "Ada", 1, "Hello Ada"},
		{"no match", []string{"Hello ", "World"}, "Jupiter", "Mars", 0, "Hello World"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			updated, count := replaceTextRuns(paragraphXML(tc.runs...), tc.find, tc.replace)
			if count != tc.want {
				t.Errorf("count = %d, want %d", count, tc.want)
			}
			if got := joinRunText(string(updated)); got != tc.wantText {
				t.Errorf("text = %q, want %q", got, tc.wantText)
			}
		})
	}
}

// A phrase spanning a paragraph break must not match: the runs belong to
// different paragraphs and the joined text is never contiguous on screen.
func TestReplaceTextRuns_DoesNotMatchAcrossParagraphs(t *testing.T) {
	content := []byte(`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<a:p><a:r><a:t>Hello </a:t></a:r></a:p>` +
		`<a:p><a:r><a:t>World</a:t></a:r></a:p></p:sld>`)

	updated, count := replaceTextRuns(content, "Hello World", "Mars")
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
	if string(updated) != string(content) {
		t.Errorf("content changed:\n%s", updated)
	}
}

func TestReplaceTextRuns_PreservesRunPropertiesAndEscapes(t *testing.T) {
	content := []byte(`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:p>` +
		`<a:r><a:rPr b="1"/><a:t>R&amp;D bud</a:t></a:r>` +
		`<a:r><a:rPr i="1"/><a:t>get</a:t></a:r>` +
		`</a:p></p:sld>`)

	updated, count := replaceTextRuns(content, "R&D budget", "R&D <spend>")
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	out := string(updated)
	if !strings.Contains(out, `<a:rPr b="1"/>`) || !strings.Contains(out, `<a:rPr i="1"/>`) {
		t.Errorf("run properties were dropped:\n%s", out)
	}
	if strings.Contains(out, "<spend>") {
		t.Errorf("replacement text was not escaped:\n%s", out)
	}
	if got := joinRunText(out); got != "R&D <spend>" {
		t.Errorf("decoded text = %q, want %q", got, "R&D <spend>")
	}
}

// The raw-bytes prefilter compares against escaped, run-split XML, so it can
// only rule a slide out for needles that cannot be affected by either.
func TestPrefilterIsReliable(t *testing.T) {
	for needle, want := range map[string]bool{
		"budget":       true,
		"R&D":          false,
		"a<b":          false,
		"Hello World":  false,
		"quote\"here":  false,
		"tab\tbetween": false,
	} {
		if got := prefilterIsReliable(needle); got != want {
			t.Errorf("prefilterIsReliable(%q) = %v, want %v", needle, got, want)
		}
	}
}

func TestTextReplacementPartsScopes(t *testing.T) {
	e := newTableEditorFixture()
	e.parts.Set("ppt/notesSlides/notesSlide1.xml", []byte("<p:notes/>"))
	e.parts.Set("ppt/notesSlides/_rels/notesSlide1.xml.rels", []byte("<Relationships/>"))
	e.parts.Set("ppt/slideLayouts/slideLayout1.xml", []byte("<p:sldLayout/>"))
	e.parts.Set("ppt/slideMasters/slideMaster1.xml", []byte("<p:sldMaster/>"))

	slides, err := e.textReplacementParts(common.TextScopeSlides)
	if err != nil || len(slides) != 1 {
		t.Fatalf("slides scope = %v, err=%v", slides, err)
	}
	withNotes, err := e.textReplacementParts(common.TextScopeSlidesAndNotes)
	if err != nil || len(withNotes) != 2 {
		t.Fatalf("slides+notes scope = %v, err=%v", withNotes, err)
	}
	all, err := e.textReplacementParts(common.TextScopeAll)
	if err != nil {
		t.Fatalf("all scope err=%v", err)
	}
	if len(all) != 4 {
		t.Fatalf("all scope = %v, want 4 parts", all)
	}
	for _, part := range all {
		if strings.Contains(part, "/_rels/") {
			t.Errorf("rels part leaked into scope: %s", part)
		}
	}
	if _, err = e.textReplacementParts("bogus"); err == nil {
		t.Error("expected an error for an unknown scope")
	}
}

// joinRunText concatenates the <a:t> values of an XML fragment, undoing the
// escaping, which is what a caller sees as the shape's text.
func joinRunText(xml string) string {
	var out strings.Builder
	for _, match := range textRunPattern.FindAllStringSubmatch(xml, -1) {
		out.WriteString(match[2])
	}
	text := out.String()
	for _, pair := range [][2]string{
		{"&amp;", "&"}, {"&lt;", "<"}, {"&gt;", ">"},
		{"&#34;", `"`}, {"&#39;", "'"},
	} {
		text = strings.ReplaceAll(text, pair[0], pair[1])
	}
	return text
}
