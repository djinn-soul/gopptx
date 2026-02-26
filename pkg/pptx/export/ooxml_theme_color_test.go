package export

import "testing"

func TestResolveOOXMLColorToken_Hex(t *testing.T) {
	r, g, b, ok := resolveOOXMLColorToken("#4F81BD")
	if !ok {
		t.Fatal("expected hex color to resolve")
	}
	if r != 0x4F || g != 0x81 || b != 0xBD {
		t.Fatalf("unexpected rgb: %d,%d,%d", r, g, b)
	}
}

func TestResolveOOXMLColorToken_SchemeAlias(t *testing.T) {
	r, g, b, ok := resolveOOXMLColorToken("tx1")
	if !ok {
		t.Fatal("expected tx1 alias to resolve")
	}
	if r != 0x00 || g != 0x00 || b != 0x00 {
		t.Fatalf("expected black for tx1 default, got %d,%d,%d", r, g, b)
	}
}

func TestResolveOOXMLColorToken_WithTransforms(t *testing.T) {
	baseR, baseG, baseB, ok := resolveOOXMLColorToken("accent1")
	if !ok {
		t.Fatal("expected accent1 to resolve")
	}
	r, g, b, ok := resolveOOXMLColorToken("accent1|tint=40000|lumMod=65000|lumOff=35000")
	if !ok {
		t.Fatal("expected transformed accent1 to resolve")
	}
	if r <= baseR || g <= baseG || b <= baseB {
		t.Fatalf("expected transformed color to lighten base (%d,%d,%d -> %d,%d,%d)", baseR, baseG, baseB, r, g, b)
	}
}

func TestResolveOOXMLColorToken_Invalid(t *testing.T) {
	_, _, _, ok := resolveOOXMLColorToken("not-a-color")
	if ok {
		t.Fatal("expected invalid color token to fail resolution")
	}
}
