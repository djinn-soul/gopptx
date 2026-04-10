package shape

import "testing"

func TestParseShapePropertiesAcceptsStandardAndLegacyCapsAttributes(t *testing.T) {
	tests := []struct {
		name            string
		runAttr         string
		expectAllCaps   bool
		expectSmallCaps bool
	}{
		{
			name:            "standard all caps",
			runAttr:         `cap="all"`,
			expectAllCaps:   true,
			expectSmallCaps: false,
		},
		{
			name:            "standard small caps",
			runAttr:         `cap="small"`,
			expectAllCaps:   false,
			expectSmallCaps: true,
		},
		{
			name:            "legacy all caps",
			runAttr:         `caps="all"`,
			expectAllCaps:   true,
			expectSmallCaps: false,
		},
		{
			name:            "legacy small caps",
			runAttr:         `smCaps="1"`,
			expectAllCaps:   false,
			expectSmallCaps: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shapeXML := []byte(
				`<p:sp xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">` +
					`<p:nvSpPr><p:cNvPr id="22" name="Caps Parse"/></p:nvSpPr>` +
					`<p:spPr><a:xfrm><a:off x="10" y="20"/><a:ext cx="300" cy="200"/></a:xfrm></p:spPr>` +
					`<p:txBody><a:bodyPr/><a:lstStyle/><a:p>` +
					`<a:r><a:rPr ` + tc.runAttr + ` lang="en-US"/><a:t>Text</a:t></a:r>` +
					`</a:p></p:txBody></p:sp>`,
			)

			shape, err := ParseShapeProperties(shapeXML)
			if err != nil {
				t.Fatalf("parseShapeProperties failed: %v", err)
			}
			if len(shape.Runs) != 1 {
				t.Fatalf("expected one parsed run, got %#v", shape.Runs)
			}
			if got := shape.Runs[0].AllCaps != nil && *shape.Runs[0].AllCaps; got != tc.expectAllCaps {
				t.Fatalf("all caps = %v, want %v", got, tc.expectAllCaps)
			}
			if got := shape.Runs[0].SmallCaps != nil && *shape.Runs[0].SmallCaps; got != tc.expectSmallCaps {
				t.Fatalf("small caps = %v, want %v", got, tc.expectSmallCaps)
			}
		})
	}
}
