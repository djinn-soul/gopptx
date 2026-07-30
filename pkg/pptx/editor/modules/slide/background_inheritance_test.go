package slide

import (
	"strings"
	"testing"
)

func TestRewriteFollowMasterBackgroundRoundTrip(t *testing.T) {
	source := []byte(
		`<p:sld><p:cSld><p:spTree/></p:cSld><p:clrMapOvr/></p:sld>`,
	)
	if !ParseFollowMasterBackground(source) {
		t.Fatal("slide without p:bg must follow its master background")
	}

	independent, err := RewriteFollowMasterBackground(source, false)
	if err != nil {
		t.Fatalf("interrupt background inheritance: %v", err)
	}
	expected := `<p:cSld><p:bg><p:bgPr><a:noFill/><a:effectLst/></p:bgPr></p:bg><p:spTree/>`
	if !strings.Contains(string(independent), expected) {
		t.Fatalf("expected schema-ordered no-fill background: %s", independent)
	}
	if ParseFollowMasterBackground(independent) {
		t.Fatal("slide with p:bg must not follow its master background")
	}

	restored, err := RewriteFollowMasterBackground(independent, true)
	if err != nil {
		t.Fatalf("restore background inheritance: %v", err)
	}
	if strings.Contains(string(restored), "<p:bg") {
		t.Fatalf("restoring inheritance must remove p:bg: %s", restored)
	}
}

func TestRewriteFollowMasterBackgroundPreservesCustomBackground(t *testing.T) {
	source := []byte(
		`<p:sld><p:cSld><p:bg><p:bgPr><a:solidFill/></p:bgPr></p:bg>` +
			`<p:spTree/></p:cSld></p:sld>`,
	)
	unchanged, err := RewriteFollowMasterBackground(source, false)
	if err != nil {
		t.Fatalf("retain custom background: %v", err)
	}
	if string(unchanged) != string(source) {
		t.Fatalf("follow=false must preserve a custom background: %s", unchanged)
	}
}
