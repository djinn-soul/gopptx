// Command gen_shape_presets generates the Python shape-preset surface from the
// ShapeType* constants declared in pkg/pptx/shapes.
//
// Python exposed 40 of the 197 presets Go knows, and an unknown preset name is
// not an error: presetGeometry falls back to a plain rectangle, so a typo
// silently drew the wrong shape. Generating the list keeps the two sides equal
// and makes the whole catalogue discoverable and checkable.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	minCLIArgs      = 3
	shapeTypePrefix = "ShapeType"
	sourcePattern   = "shape_types*.go"
)

type preset struct {
	// PythonName is the SHAPE_* constant, MemberName the enum member, and
	// Token the DrawingML value both carry.
	PythonName string
	MemberName string
	Token      string
	GoName     string
}

func main() {
	if len(os.Args) < minCLIArgs {
		fmt.Fprintln(os.Stderr, "Usage: gen_shape_presets <shapes_dir> <output_py>")
		os.Exit(1)
	}
	shapesDir := os.Args[1]
	outputFile := os.Args[2]

	presets, err := loadPresets(shapesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load presets: %v\n", err)
		os.Exit(1)
	}
	if len(presets) == 0 {
		fmt.Fprintln(os.Stderr, "no shape presets found")
		os.Exit(1)
	}
	// The surface is emitted as two modules. One file carrying both a constant
	// and an enum member per preset ran past the repository's 400-line ceiling,
	// and splitting it is what keeps a generated file from having to be
	// exempted from the rule every other file follows.
	enumFile := filepath.Join(filepath.Dir(outputFile), enumModuleFile)
	if err := os.WriteFile(outputFile, []byte(renderConstants(presets)), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "write output: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(enumFile, []byte(renderEnum(presets)), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "write enum output: %v\n", err)
		os.Exit(1)
	}
}

const (
	// enumModuleFile is the sibling module the ShapeType enum is written to.
	enumModuleFile = "shape_type_enum.py"
	// enumModuleImport is how the constants module imports it back, so
	// gopptx.shape_types.ShapeType keeps resolving for callers.
	enumModuleImport = "gopptx.shape_type_enum"
)

// loadPresets reads every ShapeType* constant with a string value, keyed by the
// DrawingML token so aliases for the same preset collapse to one entry.
func loadPresets(shapesDir string) ([]preset, error) {
	matches, err := filepath.Glob(filepath.Join(shapesDir, sourcePattern))
	if err != nil {
		return nil, err
	}

	byToken := map[string]preset{}
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		filePresets, parseErr := parseFile(path)
		if parseErr != nil {
			return nil, fmt.Errorf("%s: %w", path, parseErr)
		}
		for _, p := range filePresets {
			existing, seen := byToken[p.Token]
			// Keep the shortest Go name, which is the canonical constant
			// rather than a longer alias for the same token.
			if !seen || len(p.GoName) < len(existing.GoName) {
				byToken[p.Token] = p
			}
		}
	}

	out := make([]preset, 0, len(byToken))
	for _, p := range byToken {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].MemberName < out[j].MemberName })
	return out, nil
}

func parseFile(path string) ([]preset, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var out []preset
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			out = append(out, presetsFromSpec(spec)...)
		}
	}
	return out, nil
}

func presetsFromSpec(spec ast.Spec) []preset {
	valueSpec, ok := spec.(*ast.ValueSpec)
	if !ok {
		return nil
	}
	var out []preset
	for i, name := range valueSpec.Names {
		if !strings.HasPrefix(name.Name, shapeTypePrefix) || i >= len(valueSpec.Values) {
			continue
		}
		literal, isLiteral := valueSpec.Values[i].(*ast.BasicLit)
		if !isLiteral || literal.Kind != token.STRING {
			continue
		}
		token := strings.Trim(literal.Value, `"`)
		if token == "" {
			continue
		}
		member := screamingSnake(strings.TrimPrefix(name.Name, shapeTypePrefix))
		out = append(out, preset{
			PythonName: "SHAPE_" + member,
			MemberName: member,
			Token:      token,
			GoName:     name.Name,
		})
	}
	return out
}

// compoundWords are spellings Go writes as two words and the DrawingML and
// Python names treat as one, so FlowChartProcess stays FLOWCHART_PROCESS
// rather than becoming FLOW_CHART_PROCESS.
//
//nolint:gochecknoglobals // build-time generator config; behaves as a const map
var compoundWords = map[string]string{
	"FlowChart": "Flowchart",
}

// screamingSnake turns a Go identifier into the Python constant spelling:
// RoundedRectangle -> ROUNDED_RECTANGLE, BentConnector5 -> BENT_CONNECTOR_5.
func screamingSnake(name string) string {
	for from, to := range compoundWords {
		name = strings.ReplaceAll(name, from, to)
	}
	return screamingSnakeWords(name)
}

func screamingSnakeWords(name string) string {
	var b strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if i > 0 && needsSeparator(runes, i) {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

func needsSeparator(runes []rune, i int) bool {
	prev, cur := runes[i-1], runes[i]
	switch {
	case unicode.IsDigit(cur) && !unicode.IsDigit(prev):
		return true
	case unicode.IsUpper(cur) && !unicode.IsUpper(prev):
		return true
	case unicode.IsUpper(cur) && i+1 < len(runes) && unicode.IsLower(runes[i+1]):
		// The last capital of a run starts the next word: HTMLBlock.
		return true
	default:
		return false
	}
}

// renderConstants emits the SHAPE_* module constants, and re-exports the enum
// from its own module so gopptx.shape_types keeps offering the whole surface.
func renderConstants(presets []preset) string {
	var b strings.Builder
	b.WriteString(`"""Shape preset constants for PowerPoint auto shapes.

GENERATED FILE - DO NOT EDIT.

Generated by cmd/gen_shape_presets from the ShapeType constants in
pkg/pptx/shapes. Change the Go constants and re-run "go generate ./..."
rather than editing this file.

The ShapeType enum lives in ` + enumModuleImport + ` and is re-exported here,
so both spellings resolve off this module. __all__ names only the re-exports:
listing the SHAPE_* constants as well would put one name per formatted line and
push this module back over the line ceiling the split exists to respect. They
stay reachable as module attributes, which is how callers already read them.
"""

from __future__ import annotations

from ` + enumModuleImport + ` import MSO_SHAPE, ShapeType

__all__ = ["MSO_SHAPE", "ShapeType"]

`)

	for _, p := range presets {
		fmt.Fprintf(&b, "%s = %q\n", p.PythonName, p.Token)
	}

	// No __all__ here: the formatter puts one name per line, and 200 presets of
	// it would push this module back over the line ceiling the split exists to
	// respect. The constants are plain module-level names, which is surface
	// enough.
	return b.String()
}

// renderEnum emits the ShapeType enum. Members carry the DrawingML token
// directly rather than pointing at the constants module, which keeps the two
// files independent and free of a star import.
func renderEnum(presets []preset) string {
	var b strings.Builder
	b.WriteString(`"""The ShapeType enum of DrawingML presets.

GENERATED FILE - DO NOT EDIT.

Generated by cmd/gen_shape_presets from the ShapeType constants in
pkg/pptx/shapes. Change the Go constants and re-run "go generate ./..."
rather than editing this file.
"""

from __future__ import annotations

import sys
from enum import Enum

if sys.version_info >= (3, 11):
    from enum import StrEnum
else:

    class StrEnum(str, Enum):
        """Python 3.10 compatibility shim for enum.StrEnum."""

        __str__ = str.__str__
        __format__ = str.__format__  # type: ignore[assignment]


class ShapeType(StrEnum):
    """Every DrawingML preset gopptx can draw.

    Members subclass str, so a member compares equal to and serializes as its
    DrawingML token. A bare string still works; the members exist because an
    unknown token is not an error on the Go side — it falls back to a plain
    rectangle — so a typo would otherwise draw the wrong shape silently.
    """

`)
	for _, p := range presets {
		fmt.Fprintf(&b, "    %s = %q\n", p.MemberName, p.Token)
	}

	b.WriteString("\n\nMSO_SHAPE = ShapeType\n\n")
	writeDunderAll(&b, []string{"MSO_SHAPE", "ShapeType"})
	return b.String()
}

// writeDunderAll declares the module's public surface. It is what tells a reader
// — and the dead-code check — that these names are exported on purpose rather
// than left behind unused.
// It is written without a trailing comma and on one line, so the Python
// formatter wraps it to as many names per line as fit rather than exploding it
// to one per line — which on a 200-preset catalogue would double the module's
// length and push it back over the repository's line ceiling.
func writeDunderAll(b *strings.Builder, names []string) {
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	quoted := make([]string, 0, len(sorted))
	for _, name := range sorted {
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}
	fmt.Fprintf(b, "__all__ = [%s]\n", strings.Join(quoted, ", "))
}
