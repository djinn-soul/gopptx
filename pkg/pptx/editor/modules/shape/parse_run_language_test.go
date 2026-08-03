package shape

import (
	"encoding/xml"
	"testing"
)

// Without the language on the read side, a read-modify-write re-emits the
// editor's en-US default and drops the ea/cs typeface slot, undoing the font a
// non-Latin run was given.
func TestApplyRunStyleReadsLanguage(t *testing.T) {
	const runXML = `<r xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
		`<rPr lang="ja-JP" sz="1800"><latin typeface="Meiryo"/></rPr>` +
		`<t>日本語</t></r>`

	var parsed struct {
		RPr *runPropsXML `xml:"rPr"`
		T   string       `xml:"t"`
	}
	if err := xml.Unmarshal([]byte(runXML), &parsed); err != nil {
		t.Fatalf("unmarshal run: %v", err)
	}

	run, ok := parseTextRun(parsed)
	if !ok {
		t.Fatal("expected the run to parse")
	}
	if run.Language == nil {
		t.Fatal("run language was not read back")
	}
	if *run.Language != "ja-JP" {
		t.Fatalf("expected ja-JP, got %q", *run.Language)
	}
	if run.Font == nil || *run.Font != "Meiryo" {
		t.Fatalf("expected the Meiryo typeface, got %+v", run.Font)
	}
}

// A run with no lang attribute leaves the field unset, so the writer's own
// default still applies rather than an empty tag being emitted.
func TestApplyRunStyleLeavesLanguageUnsetWhenAbsent(t *testing.T) {
	var parsed struct {
		RPr *runPropsXML `xml:"rPr"`
		T   string       `xml:"t"`
	}
	const runXML = `<r><rPr sz="1200"/><t>plain</t></r>`
	if err := xml.Unmarshal([]byte(runXML), &parsed); err != nil {
		t.Fatalf("unmarshal run: %v", err)
	}

	run, ok := parseTextRun(parsed)
	if !ok {
		t.Fatal("expected the run to parse")
	}
	if run.Language != nil {
		t.Fatalf("expected no language, got %q", *run.Language)
	}
}
