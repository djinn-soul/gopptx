// Command gen_chart_types generates the Python ChartType enum from the
// canonical XLChartType* constants declared in pkg/pptx/enums/shape_chart.go.
//
// Those constants are defined in terms of the ChartKind* string constants in
// internal/pptxxml, so the generator resolves the two files together: it reads
// the ChartKind* literals first, then maps each XLChartType* constant onto the
// literal it refers to.
//
// Go is the single source of truth. Editing the generated Python file directly
// will be overwritten by the next `go generate ./...`.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

const (
	minCLIArgs      = 4
	chartTypePrefix = "XLChartType"
	chartKindPrefix = "ChartKind"
	initialCap      = 32
)

// Python member names referenced from more than one place.
const (
	pyNameBar       = "BAR"
	pyNameColumn    = "COLUMN"
	pyNameStockOHLC = "STOCK_OHLC"
)

// chartType is one generated Python enum member.
type chartType struct {
	// PyName is the Python member name, e.g. "BAR_STACKED_100".
	PyName string
	// Value is the wire value the Go engine expects, e.g. "barStacked100".
	Value string
	// GoName is the Go constant it came from, used in the generated comment.
	GoName string
}

// precedingAliases lists extra Python member names to emit immediately before a
// given member, sharing its value. Python's Enum turns a repeated value into an
// alias of the first name declared, so the alias listed here becomes canonical
// and the Go-derived name becomes the alias.
//
// COLUMN exists for python-pptx familiarity: that library calls a vertical bar
// chart a column chart. It has no Go counterpart, so it is declared here rather
// than in the Go source.
//
//nolint:gochecknoglobals // build-time generator config; behaves as a const map
var precedingAliases = map[string][]string{
	pyNameBar: {pyNameColumn},
}

// aliasDoc documents each alias in the generated output.
//
//nolint:gochecknoglobals // build-time generator config; behaves as a const map
var aliasDoc = map[string]string{
	pyNameColumn: "Column/vertical bar chart. python-pptx spelling of BAR.",
}

// memberDoc gives each Go-derived member its Python docstring.
//
//nolint:gochecknoglobals // build-time generator config; behaves as a const map
var memberDoc = map[string]string{
	pyNameBar:          "Vertical bar chart.",
	"BAR_HORIZONTAL":   "Horizontal bar chart.",
	"BAR_STACKED":      "Stacked bar chart.",
	"BAR_STACKED_100":  "100% stacked bar chart.",
	"LINE":             "Line chart.",
	"LINE_MARKERS":     "Line chart with markers.",
	"LINE_STACKED":     "Stacked line chart.",
	"SCATTER":          "Scatter/XY chart. Takes x and y values, not categories.",
	"AREA":             "Area chart.",
	"AREA_STACKED":     "Stacked area chart.",
	"AREA_STACKED_100": "100% stacked area chart.",
	"PIE":              "Pie chart.",
	"DOUGHNUT":         "Doughnut chart.",
	"BUBBLE":           "Bubble chart. Takes x, y and size values, not categories.",
	"RADAR":            "Radar chart.",
	"RADAR_FILLED":     "Filled radar chart.",
	"STOCK_HLC":        "Stock high-low-close chart.",
	pyNameStockOHLC:    "Stock open-high-low-close chart.",
	"COMBO":            "Combo chart: bar and line series on one category axis.",
}

func main() {
	if len(os.Args) < minCLIArgs {
		fmt.Fprintln(os.Stderr,
			"Usage: gen_chart_types <enums_go_file> <chartkind_go_file> <output_py_file>")
		os.Exit(1)
	}

	enumsFile := os.Args[1]
	kindFile := os.Args[2]
	outputPy := os.Args[3]

	kinds, err := parseChartKinds(kindFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chart kinds: %v\n", err)
		os.Exit(1)
	}

	types, err := parseChartTypes(enumsFile, kinds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing chart types: %v\n", err)
		os.Exit(1)
	}
	if len(types) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no %s* constants found in %s\n", chartTypePrefix, enumsFile)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPy, []byte(renderPython(types)), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPy, err)
		os.Exit(1)
	}
}

// parseChartKinds collects the ChartKind* string constants, keyed by their name
// with the prefix stripped, e.g. "Bar" -> "bar".
func parseChartKinds(path string) (map[string]string, error) {
	const prefix = chartKindPrefix

	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	out := make(map[string]string, initialCap)
	forEachConst(file, func(vspec *ast.ValueSpec) {
		for i, name := range vspec.Names {
			if !strings.HasPrefix(name.Name, prefix) || i >= len(vspec.Values) {
				continue
			}
			lit, ok := vspec.Values[i].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			value, unquoteErr := strconv.Unquote(lit.Value)
			if unquoteErr != nil {
				continue
			}
			out[strings.TrimPrefix(name.Name, prefix)] = value
		}
	})
	return out, nil
}

// parseChartTypes collects XLChartType* constants and resolves each one to the
// string literal its ChartKind* right-hand side stands for.
func parseChartTypes(path string, kinds map[string]string) ([]chartType, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	types := make([]chartType, 0, initialCap)
	var resolveErr error
	forEachConst(file, func(vspec *ast.ValueSpec) {
		for i, name := range vspec.Names {
			if !strings.HasPrefix(name.Name, chartTypePrefix) || i >= len(vspec.Values) {
				continue
			}
			value, err := resolveChartKind(vspec.Values[i], kinds)
			if err != nil {
				resolveErr = fmt.Errorf("%s: %w", name.Name, err)
				return
			}
			types = append(types, chartType{
				PyName: toScreamingSnake(strings.TrimPrefix(name.Name, chartTypePrefix)),
				Value:  value,
				GoName: name.Name,
			})
		}
	})
	if resolveErr != nil {
		return nil, resolveErr
	}
	return types, nil
}

// resolveChartKind turns a `pptxxml.ChartKindFoo` selector, a bare `ChartKindFoo`
// identifier, or a plain string literal into its string value.
func resolveChartKind(expr ast.Expr, kinds map[string]string) (string, error) {
	var ident string
	switch node := expr.(type) {
	case *ast.SelectorExpr:
		ident = node.Sel.Name
	case *ast.Ident:
		ident = node.Name
	case *ast.BasicLit:
		if node.Kind != token.STRING {
			return "", fmt.Errorf("unsupported literal kind %s", node.Kind)
		}
		return strconv.Unquote(node.Value)
	default:
		return "", fmt.Errorf("unsupported value expression %T", expr)
	}

	value, ok := kinds[strings.TrimPrefix(ident, chartKindPrefix)]
	if !ok {
		return "", fmt.Errorf("no %s* constant named %q", chartKindPrefix, ident)
	}
	return value, nil
}

func forEachConst(file *ast.File, visit func(*ast.ValueSpec)) {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			if vspec, ok := spec.(*ast.ValueSpec); ok {
				visit(vspec)
			}
		}
	}
}

// toScreamingSnake converts a Go camel-case identifier to a Python constant
// name: BarStacked100 -> BAR_STACKED_100, StockOHLC -> STOCK_OHLC.
func toScreamingSnake(s string) string {
	var out strings.Builder
	for i, r := range s {
		if needsUnderscore(s, i, r) {
			out.WriteRune('_')
		}
		out.WriteRune(r)
	}
	return strings.ToUpper(out.String())
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// needsUnderscore reports whether a '_' belongs before s[i]. Rules:
//   - letter followed by the first digit of a run: split (Stacked100 -> STACKED_100).
//   - lowercase to uppercase: always split (BarStacked -> BAR_STACKED).
//   - an uppercase run splits only before its final letter when that letter
//     begins a lowercase word, so acronyms stay whole (StockOHLC -> STOCK_OHLC).
func needsUnderscore(s string, i int, r rune) bool {
	if i == 0 {
		return false
	}
	prev := rune(s[i-1])
	if isDigit(r) {
		return !isDigit(prev)
	}
	if !isUpper(r) {
		return false
	}
	if !isUpper(prev) {
		return true
	}
	if i+1 >= len(s) || !isLower(rune(s[i+1])) {
		return false
	}
	return i >= 2 && isUpper(rune(s[i-2]))
}

// renderPython emits the generated chart_types.py module: a fixed header, one
// member per chart type, then the fixed classmethods.
func renderPython(types []chartType) string {
	var out strings.Builder
	out.WriteString(pythonHeader)

	for _, t := range types {
		aliases := precedingAliases[t.PyName]
		for _, alias := range aliases {
			writeMember(&out, alias, t.Value, aliasDoc[alias], "", false)
		}
		// A member preceded by an alias repeats that alias's value on purpose,
		// which is how Python's Enum declares an alias. Ruff's PIE796 looks for
		// accidental duplicates, so silence it just on these lines.
		writeMember(&out, t.PyName, t.Value, memberDoc[t.PyName], t.GoName, len(aliases) > 0)
	}

	out.WriteString(pythonMethods)
	return out.String()
}

const pythonHeader = `"""Chart type constants for PowerPoint chart operations.

GENERATED FILE - DO NOT EDIT.

Generated by cmd/gen_chart_types from the XLChartType constants in
pkg/pptx/enums/shape_chart.go. Change the Go constants and re-run
"go generate ./..." rather than editing this file.
"""

from __future__ import annotations

from enum import Enum

__all__ = ["ChartType"]


class ChartType(str, Enum):
    """Chart types accepted by add_chart and friends.

    Members subclass str, so a member compares equal to and serializes as its
    wire value: ChartType.PIE == "pie" is True. Passing the bare string works
    too; the members exist so the set of valid values is discoverable, checkable
    by a type checker, and typo-proof.

    Examples:
        slide.add_chart(
            ChartType.PIE,
            ["A", "B", "C"],
            [25.0, 35.0, 40.0],
            title="Sales Mix",
            bounds=(Inches(1), Inches(1), Inches(4), Inches(3)),
        )
    """

`

const pythonMethods = `    # Report the wire value from str(), matching enum.StrEnum on Python 3.11+.
    # gopptx supports Python 3.10, where StrEnum does not exist.
    __str__ = str.__str__

    @classmethod
    def get_all(cls) -> dict[str, str]:
        """Map every member name, aliases included, to its wire value.

        Returns:
            Dictionary of constant name to chart type value.
        """
        return {name: member.value for name, member in cls.__members__.items()}

    @classmethod
    def values(cls) -> frozenset[str]:
        """Return the distinct wire values, without alias duplicates."""
        return frozenset(member.value for member in cls)

    @classmethod
    def validate(cls, chart_type: str | None) -> str:
        """Check a chart type and return it as a plain string.

        Accepts a ChartType member or its wire value (they are equal). Member
        *names* such as "COLUMN" are not values and are rejected.

        Args:
            chart_type: A ChartType member or wire value such as "pie".

        Returns:
            The wire value.

        Raises:
            ValueError: If chart_type is empty or not a known value.

        Examples:
            ChartType.validate(ChartType.COLUMN)   # -> "bar"
            ChartType.validate("pie")              # -> "pie"
            ChartType.validate("COLUMN")           # -> ValueError
        """
        if not chart_type:
            raise ValueError("chart_type cannot be empty")
        if chart_type in cls.values():
            return str(chart_type)
        valid_values = ", ".join(sorted(cls.values()))
        raise ValueError(" ".join((
            f"Invalid chart_type {chart_type!r}.",
            "Use a ChartType member such as ChartType.COLUMN or ChartType.PIE.",
            f"Valid values: {valid_values}",
        )))

    @classmethod
    def get_by_name(cls, name: str) -> str | None:
        """Look up a wire value by member name, case-sensitively.

        Args:
            name: Constant name such as "COLUMN", "PIE" or "LINE_MARKERS".

        Returns:
            The wire value, or None when no member has that name.
        """
        member = cls.__members__.get(name)
        return None if member is None else member.value
`

func writeMember(
	out *strings.Builder,
	name string,
	value string,
	doc string,
	goName string,
	isAliasOf bool,
) {
	if isAliasOf {
		// Deliberate Enum alias: same value as the member emitted just above.
		fmt.Fprintf(out, "    %s = %q  # noqa: PIE796\n", name, value)
	} else {
		fmt.Fprintf(out, "    %s = %q\n", name, value)
	}
	if doc == "" {
		doc = name + " chart."
	}
	if goName != "" {
		fmt.Fprintf(out, "    \"\"\"%s Go: enums.%s.\"\"\"\n\n", doc, goName)
		return
	}
	fmt.Fprintf(out, "    \"\"\"%s\"\"\"\n\n", doc)
}
