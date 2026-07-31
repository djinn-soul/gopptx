package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestToScreamingSnake(t *testing.T) {
	cases := map[string]string{
		"Bar":            "BAR",
		"BarHorizontal":  "BAR_HORIZONTAL",
		"BarStacked":     "BAR_STACKED",
		"BarStacked100":  "BAR_STACKED_100",
		"AreaStacked100": "AREA_STACKED_100",
		"LineMarkers":    "LINE_MARKERS",
		"RadarFilled":    "RADAR_FILLED",
		"StockHLC":       "STOCK_HLC",
		"StockOHLC":      "STOCK_OHLC",
		"Combo":          "COMBO",
	}
	for in, want := range cases {
		if got := toScreamingSnake(in); got != want {
			t.Errorf("toScreamingSnake(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseChartTypesUsesStablePythonNameOverride(t *testing.T) {
	dir := t.TempDir()
	kindPath := filepath.Join(dir, "kinds.go")
	enumPath := filepath.Join(dir, "enums.go")
	if err := os.WriteFile(kindPath, []byte(`package pptxxml
const ChartKindThreeDPie = "pie3D"
`), 0o600); err != nil {
		t.Fatalf("write kinds: %v", err)
	}
	if err := os.WriteFile(enumPath, []byte(`package enums
const XLChartTypeThreeDPie XLChartType = pptxxml.ChartKindThreeDPie
`), 0o600); err != nil {
		t.Fatalf("write enums: %v", err)
	}

	kinds, err := parseChartKinds(kindPath)
	if err != nil {
		t.Fatalf("parseChartKinds: %v", err)
	}
	types, err := parseChartTypes(enumPath, kinds)
	if err != nil {
		t.Fatalf("parseChartTypes: %v", err)
	}
	if len(types) != 1 || types[0].PyName != pyNameThreeDPie {
		t.Fatalf("generated names = %+v, want THREE_D_PIE", types)
	}
}

func TestParseStringConsts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kinds.go")
	src := `package pptxxml
const (
	ChartKindBar = "bar"
	ChartKindPie = "pie"
	Unrelated    = "nope"
)`
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	kinds, err := parseChartKinds(path)
	if err != nil {
		t.Fatalf("parseChartKinds: %v", err)
	}
	if len(kinds) != 2 {
		t.Fatalf("got %d kinds, want 2: %v", len(kinds), kinds)
	}
	if kinds["Bar"] != "bar" || kinds["Pie"] != "pie" {
		t.Errorf("unexpected kinds: %v", kinds)
	}
	if _, ok := kinds["Unrelated"]; ok {
		t.Error("constant without the prefix leaked into the result")
	}
}

// parseChartTypes must resolve the pptxxml.ChartKind* selector on the
// right-hand side, since XLChartType constants never spell the value directly.
func TestParseChartTypesResolvesSelectors(t *testing.T) {
	dir := t.TempDir()
	kindPath := filepath.Join(dir, "kinds.go")
	enumPath := filepath.Join(dir, "enums.go")

	if err := os.WriteFile(kindPath, []byte(`package pptxxml
const (
	ChartKindBar          = "bar"
	ChartKindStockOHLC    = "stockOHLC"
)`), 0o600); err != nil {
		t.Fatalf("write kinds: %v", err)
	}
	if err := os.WriteFile(enumPath, []byte(`package enums
type XLChartType string
const (
	XLChartTypeBar       XLChartType = pptxxml.ChartKindBar
	XLChartTypeStockOHLC XLChartType = pptxxml.ChartKindStockOHLC
)`), 0o600); err != nil {
		t.Fatalf("write enums: %v", err)
	}

	kinds, err := parseChartKinds(kindPath)
	if err != nil {
		t.Fatalf("parseChartKinds: %v", err)
	}
	types, err := parseChartTypes(enumPath, kinds)
	if err != nil {
		t.Fatalf("parseChartTypes: %v", err)
	}

	want := map[string]string{"BAR": "bar", "STOCK_OHLC": "stockOHLC"}
	if len(types) != len(want) {
		t.Fatalf("got %d types, want %d: %+v", len(types), len(want), types)
	}
	for _, ct := range types {
		if want[ct.PyName] != ct.Value {
			t.Errorf("%s = %q, want %q", ct.PyName, ct.Value, want[ct.PyName])
		}
	}
}

// An XLChartType pointing at a ChartKind that does not exist is a wiring
// mistake, and must fail loudly rather than emit an empty Python value.
func TestParseChartTypesRejectsUnknownKind(t *testing.T) {
	dir := t.TempDir()
	enumPath := filepath.Join(dir, "enums.go")
	if err := os.WriteFile(enumPath, []byte(`package enums
const (
	XLChartTypeGhost XLChartType = pptxxml.ChartKindGhost
)`), 0o600); err != nil {
		t.Fatalf("write enums: %v", err)
	}

	_, err := parseChartTypes(enumPath, map[string]string{"Bar": "bar"})
	if err == nil {
		t.Fatal("expected an error for an unresolvable ChartKind")
	}
	if !strings.Contains(err.Error(), "XLChartTypeGhost") {
		t.Errorf("error %q does not name the offending constant", err)
	}
}

func TestRenderPythonEmitsAliasBeforeMember(t *testing.T) {
	out := renderPython([]chartType{
		{PyName: "BAR", Value: "bar", GoName: "XLChartTypeBar"},
		{PyName: "PIE", Value: "pie", GoName: "XLChartTypePie"},
	})

	columnAt := strings.Index(out, "COLUMN = ")
	barAt := strings.Index(out, "BAR = ")
	if columnAt < 0 || barAt < 0 {
		t.Fatalf("COLUMN and BAR must both be emitted:\n%s", out)
	}
	// Python's Enum makes the first name with a value canonical, so COLUMN has
	// to precede BAR for BAR to be the alias.
	if columnAt > barAt {
		t.Error("COLUMN must be emitted before BAR so that BAR becomes the alias")
	}
	if !strings.Contains(out, "class ChartType(str, Enum):") {
		t.Error("generated class must subclass str and Enum")
	}
	if !strings.Contains(out, "GENERATED FILE - DO NOT EDIT") {
		t.Error("generated file must be marked as generated")
	}
}

// The committed Python enum must match what the current Go constants produce.
// Complements the repo-level `task check:generated`.
func TestGeneratedFileMatchesGoSource(t *testing.T) {
	root := repoRoot(t)
	kindPath := filepath.Join(root, "internal", "pptxxml", "chart_spec.go")
	enumPath := filepath.Join(root, "pkg", "pptx", "enums", "shape_chart.go")
	outPath := filepath.Join(
		root, "python", "gopptx", "presentation", "charts", "chart_types.py",
	)

	kinds, err := parseChartKinds(kindPath)
	if err != nil {
		t.Fatalf("parseChartKinds: %v", err)
	}
	types, err := parseChartTypes(enumPath, kinds)
	if err != nil {
		t.Fatalf("parseChartTypes: %v", err)
	}
	if len(types) == 0 {
		t.Fatal("no XLChartType constants parsed from the Go source")
	}

	committed, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if string(committed) != renderPython(types) {
		t.Error("chart_types.py is stale; run 'go generate ./...' and commit the result")
	}
}
