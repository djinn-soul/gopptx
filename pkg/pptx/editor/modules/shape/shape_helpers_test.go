package shape

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestBasicParsingAndNormalizationHelpers(t *testing.T) {
	if got := GetStr(nil); got != "" {
		t.Fatalf("GetStr(nil)=%q, want empty", got)
	}
	s := " value "
	if got := GetStr(&s); got != " value " {
		t.Fatalf("GetStr(pointer)=%q", got)
	}

	if got := ParseIntAttr(nil); got != 0 {
		t.Fatalf("ParseIntAttr(nil)=%d, want 0", got)
	}
	v := "42"
	if got := ParseIntAttr(&v); got != 42 {
		t.Fatalf("ParseIntAttr(42)=%d, want 42", got)
	}
	bad := "x"
	if got := ParseIntAttr(&bad); got != 0 {
		t.Fatalf("ParseIntAttr(invalid)=%d, want 0", got)
	}

	for _, raw := range []string{"1", "true", "on", "yes", " YES "} {
		c := raw
		if !ParseXMLBoolAttr(&c) {
			t.Fatalf("ParseXMLBoolAttr(%q)=false, want true", raw)
		}
	}
	off := "0"
	if ParseXMLBoolAttr(&off) {
		t.Fatal("ParseXMLBoolAttr(\"0\") should be false")
	}

	escaped := XMLEscape(`A&B<"C">`)
	if escaped != `A&amp;B&lt;&#34;C&#34;&gt;` {
		t.Fatalf("unexpected XMLEscape output: %q", escaped)
	}

	if color, err := NormalizeHexColor("#a1b2c3"); err != nil || color != "A1B2C3" {
		t.Fatalf("NormalizeHexColor valid failed: color=%q err=%v", color, err)
	}
	if _, err := NormalizeHexColor("ABC"); err == nil {
		t.Fatal("expected short hex color error")
	}
	if _, err := NormalizeHexColor("ZZZZZZ"); err == nil {
		t.Fatal("expected invalid hex color error")
	}

	if dash, err := NormalizeLineDashStyle("long-dash-dot"); err != nil || dash != "lgDashDot" {
		t.Fatalf("NormalizeLineDashStyle alias failed: dash=%q err=%v", dash, err)
	}
	if dash, err := NormalizeLineDashStyle("sysDashDot"); err != nil || dash != "sysDashDot" {
		t.Fatalf("NormalizeLineDashStyle canonical failed: dash=%q err=%v", dash, err)
	}
	if _, err := NormalizeLineDashStyle("unknown"); err == nil {
		t.Fatal("expected unsupported dash style error")
	}
}

func TestHyperlinkHelpers(t *testing.T) {
	if got := DeriveActionURL(nil); got != "" {
		t.Fatalf("DeriveActionURL(nil)=%q, want empty", got)
	}
	targetSlide := 2
	if got := DeriveActionURL(&common.Hyperlink{TargetSlide: &targetSlide}); got != "ppaction://hlinksldjump" {
		t.Fatalf("DeriveActionURL(target_slide)=%q", got)
	}
	jump := "nextslide"
	if got := DeriveActionURL(&common.Hyperlink{TargetJump: &jump}); got != "ppaction://hlinkshowjump?jump=nextslide" {
		t.Fatalf("DeriveActionURL(jump)=%q", got)
	}
	macro := "RunMe"
	if got := DeriveActionURL(&common.Hyperlink{Macro: &macro}); got != "ppaction://macro?name=RunMe" {
		t.Fatalf("DeriveActionURL(macro)=%q", got)
	}

	addr := "https://example.com"
	hl := &common.Hyperlink{Address: &addr, Macro: &macro}
	if err := ValidateHyperlinkAction(hl); err == nil {
		t.Fatal("expected mutually exclusive selector error")
	}
	badJump := "somewhere"
	if err := ValidateHyperlinkAction(&common.Hyperlink{TargetJump: &badJump}); err == nil {
		t.Fatal("expected invalid jump target error")
	}
	validJump := "lastslide"
	if err := ValidateHyperlinkAction(&common.Hyperlink{TargetJump: &validJump}); err != nil {
		t.Fatalf("expected valid jump target: %v", err)
	}
}

func TestTextRunHelpers(t *testing.T) {
	src := []common.TextRun{{Text: "a"}, {Text: "b"}}
	copyRuns := CopyTextRuns(src)
	if len(copyRuns) != 2 || copyRuns[1].Text != "b" {
		t.Fatalf("CopyTextRuns unexpected result: %+v", copyRuns)
	}
	copyRuns[0].Text = "mutated"
	if src[0].Text != "a" {
		t.Fatal("CopyTextRuns must detach from source slice")
	}

	updated, err := UpdateRunText(src, 1, "B")
	if err != nil || updated[1].Text != "B" || src[1].Text != "b" {
		t.Fatalf("UpdateRunText failed: updated=%+v src=%+v err=%v", updated, src, err)
	}
	if _, err = UpdateRunText(src, 5, "x"); err == nil {
		t.Fatal("expected out-of-range update error")
	}

	appended := AppendRun(src, common.TextRun{Text: "c"})
	if len(appended) != 3 || appended[2].Text != "c" {
		t.Fatalf("AppendRun unexpected output: %+v", appended)
	}
}

func TestShapeXMLMutationHelpers(t *testing.T) {
	xmlData := []byte(
		`<p:sp><p:spPr><a:xfrm><a:off x="1" y="2"/><a:ext cx="3" cy="4"/></a:xfrm><a:solidFill/></p:spPr></p:sp>`,
	)
	transformed := string(UpdateShapeTransforms(xmlData, 10, 20, 30, 40))
	if !strings.Contains(transformed, `<a:off x="10" y="20"/>`) ||
		!strings.Contains(transformed, `<a:ext cx="30" cy="40"/>`) {
		t.Fatalf("UpdateShapeTransforms output unexpected: %s", transformed)
	}

	base := []byte(
		`<p:sp><p:spPr bwMode="auto"><a:solidFill/><a:ln w="1"/><a:effectLst/><a:prstGeom prst="rect"/></p:spPr></p:sp>`,
	)
	updated := string(ReplaceStyleInSpPr(base, `<a:noFill/>`, true, false, true))
	if strings.Contains(updated, "<a:solidFill") || strings.Contains(updated, "<a:effectLst") {
		t.Fatalf("ReplaceStyleInSpPr should remove fill/effects blocks: %s", updated)
	}
	if !strings.Contains(updated, "<a:ln w=\"1\"/>") {
		t.Fatalf("ReplaceStyleInSpPr should preserve line when applyLine=false: %s", updated)
	}
	if !strings.Contains(updated, "<a:noFill/>") {
		t.Fatalf("ReplaceStyleInSpPr should inject styleXML: %s", updated)
	}

	inner := `<a:solidFill/><a:ln/><a:effectLst/>`
	stripped := stripSelectiveStyleBlocks(inner, false, true, true)
	if strings.Contains(stripped, "<a:ln") || strings.Contains(stripped, "<a:effectLst") {
		t.Fatalf("stripSelectiveStyleBlocks should remove line/effects: %s", stripped)
	}
	if !strings.Contains(removeFillBlocks(`<a:solidFill/><a:gradFill/><a:ln/>`), "<a:ln/>") {
		t.Fatal("removeFillBlocks should not remove line blocks")
	}
}

func TestPresetShapeRenderingHelpers(t *testing.T) {
	if got := presetGeometry("oval"); got != "ellipse" {
		t.Fatalf("presetGeometry(oval)=%q, want ellipse", got)
	}
	if got := presetGeometry("triangle"); got != "triangle" {
		t.Fatalf("presetGeometry(triangle)=%q", got)
	}
	if got := presetGeometry("unknown"); got != "rect" {
		t.Fatalf("presetGeometry(default)=%q", got)
	}

	xml := string(BuildPresetShapeXML(
		5,
		`Name "A"`,
		"oval",
		`<a:hlinkClick/>`,
		`<a:hlinkHover/>`,
		1, 2, 3, 4,
		`<a:noFill/>`,
		`<p:txBody/>`,
	))
	if !strings.Contains(xml, `name="Name &#34;A&#34;"`) {
		t.Fatalf("BuildPresetShapeXML should escape name: %s", xml)
	}
	if !strings.Contains(xml, `prst="ellipse"`) || !strings.Contains(xml, `<a:noFill/>`) ||
		!strings.Contains(xml, `<p:txBody/>`) {
		t.Fatalf("BuildPresetShapeXML missing expected blocks: %s", xml)
	}
}
