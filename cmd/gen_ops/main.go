package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

type opSpec struct {
	PyName string
	Value  string
}

const (
	minCLIArgs     = 4
	initialOpsCap  = 64
	opNamePrefix   = "Op"
	opPyNamePrefix = "OP_"
)

func main() {
	if len(os.Args) < minCLIArgs {
		fmt.Fprintln(os.Stderr, "Usage: gen_ops <input_go_file> <output_py_file> <output_pyi_file>")
		os.Exit(1)
	}

	input := os.Args[1]
	outputPy := os.Args[2]
	outputPyi := os.Args[3]

	ops, err := parseOpsFromGo(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ops: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].PyName < ops[j].PyName
	})

	pyFile, err := os.Create(outputPy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}

	pyiFile, err := os.Create(outputPyi)
	if err != nil {
		if closeErr := pyFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", closeErr)
		}
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}

	if err := writeOpsPy(pyFile, ops); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing py file: %v\n", err)
		os.Exit(1)
	}
	if err := writeOpsPyi(pyiFile, ops); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing pyi file: %v\n", err)
		os.Exit(1)
	}
	if err := pyFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing file: %v\n", err)
		os.Exit(1)
	}
	if err := pyiFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing file: %v\n", err)
		os.Exit(1)
	}
}

func parseOpsFromGo(input string) ([]opSpec, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, input, nil, 0)
	if err != nil {
		return nil, err
	}

	ops := make([]opSpec, 0, initialOpsCap)
	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			op, ok, err := parseOpSpec(vspec)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			ops = append(ops, op)
		}
	}
	return ops, nil
}

func parseOpSpec(vspec *ast.ValueSpec) (opSpec, bool, error) {
	if len(vspec.Names) == 0 || len(vspec.Values) == 0 {
		return opSpec{}, false, nil
	}
	name := vspec.Names[0].Name
	if !strings.HasPrefix(name, opNamePrefix) {
		return opSpec{}, false, nil
	}
	lit, ok := vspec.Values[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return opSpec{}, false, fmt.Errorf("op %q must be a string literal", name)
	}
	value, unquoteErr := unquote(lit.Value)
	if unquoteErr != nil {
		return opSpec{}, false, fmt.Errorf("unquote op %q: %w", name, unquoteErr)
	}
	return opSpec{
		PyName: opPyNamePrefix + toSnakeCase(strings.TrimPrefix(name, opNamePrefix)),
		Value:  value,
	}, true, nil
}

func writeOpsPy(f *os.File, ops []opSpec) error {
	var out strings.Builder
	out.WriteString(`"""Operation constants shared by gopptx Python runtime."""` + "\n\n")
	out.WriteString("from __future__ import annotations\n\n")
	for _, op := range ops {
		out.WriteString(fmt.Sprintf("%s = %q\n", op.PyName, op.Value))
	}
	out.WriteString("\nSUPPORTED_OPS = (\n")
	for _, op := range ops {
		out.WriteString(fmt.Sprintf("    %s,\n", op.PyName))
	}
	out.WriteString(")\n\nSUPPORTED_OPS_SET = frozenset(SUPPORTED_OPS)\n")
	_, err := f.WriteString(out.String())
	return err
}

func writeOpsPyi(f *os.File, ops []opSpec) error {
	var out strings.Builder
	out.WriteString("from __future__ import annotations\n\n")
	for _, op := range ops {
		out.WriteString(fmt.Sprintf("%s: str\n", op.PyName))
	}
	out.WriteString("SUPPORTED_OPS: tuple[str, ...]\n")
	out.WriteString("SUPPORTED_OPS_SET: frozenset[str]\n")
	_, err := f.WriteString(out.String())
	return err
}

func unquote(s string) (string, error) { return strconv.Unquote(s) }

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if needsUnderscore(s, i, r) {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }

// needsUnderscore reports whether a '_' separator should be inserted before the
// character r at position i in s. Rules:
//   - lowercase → uppercase transition: always insert (e.g. shapeS → shape_S).
//   - uppercase run ending before lowercase: insert only when run is ≥2 chars long,
//     so "IDs" stays as "IDS" rather than splitting to "I_DS".
func needsUnderscore(s string, i int, r rune) bool {
	if i == 0 || !isUpper(r) {
		return false
	}
	prev := rune(s[i-1])
	if !isUpper(prev) {
		return true // lowercase → uppercase
	}
	// prev is uppercase: only insert at end of a ≥2-char run before a lowercase
	if i+1 >= len(s) || !isLower(rune(s[i+1])) {
		return false
	}
	return i >= 2 && isUpper(rune(s[i-2]))
}
