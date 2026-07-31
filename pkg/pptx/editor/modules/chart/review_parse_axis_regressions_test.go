package chart

import (
	"regexp"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestTrendlineBoolValueAcceptsOOXMLBooleanLiterals(t *testing.T) {
	pattern := regexp.MustCompile(`<c:flag val="([^"]+)"`)
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{value: "1", want: true},
		{value: "true", want: true},
		{value: "TRUE", want: true},
		{value: "0", want: false},
		{value: "false", want: false},
	} {
		got := trendlineBoolValue(`<c:flag val="`+tc.value+`"/>`, pattern)
		if got == nil || *got != tc.want {
			t.Fatalf("value %q parsed as %v, want %v", tc.value, got, tc.want)
		}
	}
	if got := trendlineBoolValue(`<c:flag val="invalid"/>`, pattern); got != nil {
		t.Fatalf("invalid boolean parsed as %v", *got)
	}
}

func TestPatchAxisRotationPreservesExpandedBodyProperties(t *testing.T) {
	fixture := strings.Replace(
		axisFixtureXML,
		`<c:crossAx val="2"/>`,
		`<c:txPr><a:bodyPr wrap="square"><a:normAutofit/></a:bodyPr>`+
			`<a:lstStyle/><a:p/></c:txPr><c:crossAx val="2"/>`,
		1,
	)
	patched := PatchAxisVisibility(fixture, common.ChartFormatUpdate{
		CategoryAxisTickLabelRotation: floatPtr(30),
	})

	for _, want := range []string{
		`<a:bodyPr wrap="square" rot="1800000">`,
		`<a:normAutofit/>`,
	} {
		if !strings.Contains(patched, want) {
			t.Fatalf("rotation patch dropped %s:\n%s", want, patched)
		}
	}
}
