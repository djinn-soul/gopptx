package chart

import (
	"strings"
	"testing"
)

func TestCachePatchReplaceHelpers(t *testing.T) {
	field := `<c:val><c:f>Sheet1!$B$2:$B$4</c:f><c:numCache><c:ptCount val="1"/></c:numCache></c:val>`
	withFormula := applyFieldFormula(field, "Sheet1!$C$2:$C$4")
	if !strings.Contains(withFormula, "<c:f>Sheet1!$C$2:$C$4</c:f>") {
		t.Fatalf("applyFieldFormula failed: %s", withFormula)
	}

	out, err := applyStringValues("c:cat", `<c:cat><c:strCache></c:strCache></c:cat>`, []string{"A", "B"})
	if err != nil || !strings.Contains(out, "<c:pt idx=\"1\"><c:v>B</c:v></c:pt>") {
		t.Fatalf("applyStringValues failed: out=%s err=%v", out, err)
	}

	out, err = applyMultiLevelValues(
		"c:cat",
		`<c:cat><c:multiLvlStrCache></c:multiLvlStrCache></c:cat>`,
		[][]string{{"L1A", "L1B"}, {"L2A", "L2B"}},
	)
	if err != nil || !strings.Contains(out, "<c:lvl>") {
		t.Fatalf("applyMultiLevelValues failed: out=%s err=%v", out, err)
	}

	out, err = applyNumericValues(
		"c:val",
		`<c:val><c:numCache><c:formatCode>General</c:formatCode><c:ptCount val="0"/></c:numCache></c:val>`,
		[]float64{1.5, 2.5},
	)
	if err != nil || !strings.Contains(out, "<c:v>2.5</c:v>") {
		t.Fatalf("applyNumericValues failed: out=%s err=%v", out, err)
	}

	nums, err := convertStringsToFloats([]string{"1.25", "2"}, "c:numVal")
	if err != nil || len(nums) != 2 || nums[0] != 1.25 {
		t.Fatalf("convertStringsToFloats failed: nums=%v err=%v", nums, err)
	}
	if _, err = convertStringsToFloats([]string{"x"}, "c:numVal"); err == nil {
		t.Fatal("expected convertStringsToFloats parse error")
	}

	if got := buildStringData("c:strCache", []string{"A"}); !strings.Contains(got, `<c:ptCount val="1"/>`) {
		t.Fatalf("buildStringData unexpected output: %s", got)
	}
	if got := buildNumberData("c:numCache", "General", []float64{3.14}); !strings.Contains(got, "<c:formatCode>General</c:formatCode>") {
		t.Fatalf("buildNumberData unexpected output: %s", got)
	}
	if got := buildMultiLevelData("c:multiLvlStrCache", [][]string{{"A"}}); !strings.Contains(got, "<c:lvl>") {
		t.Fatalf("buildMultiLevelData unexpected output: %s", got)
	}

	if fc := extractFormatCode(`<c:numCache><c:formatCode>0.00</c:formatCode></c:numCache>`); fc != "0.00" {
		t.Fatalf("extractFormatCode=%q, want 0.00", fc)
	}
	if fc := extractFormatCode(`<c:numCache></c:numCache>`); fc != "General" {
		t.Fatalf("extractFormatCode default=%q, want General", fc)
	}

	if got := sheetRange("B", 3); got != "Sheet1!$B$2:$B$4" {
		t.Fatalf("sheetRange=%q", got)
	}
	if got := sheetRangeAcrossColumns(1, 3, 2); got != "Sheet1!$A$2:$C$3" {
		t.Fatalf("sheetRangeAcrossColumns=%q", got)
	}
}
