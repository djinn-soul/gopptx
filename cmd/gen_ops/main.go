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
		_ = pyFile.Close()
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}

	writeOpsPy(pyFile, ops)
	writeOpsPyi(pyiFile, ops)
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

func writeOpsPy(f *os.File, ops []opSpec) {
	fmt.Fprintln(f, "from __future__ import annotations")
	fmt.Fprintln(f)
	for _, op := range ops {
		fmt.Fprintf(f, "%s = %q\n", op.PyName, op.Value)
	}
	fmt.Fprintln(f)
	fmt.Fprintln(f, "SUPPORTED_OPS = (")
	for _, op := range ops {
		fmt.Fprintf(f, "    %s,\n", op.PyName)
	}
	fmt.Fprintln(f, ")")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "SUPPORTED_OPS_SET = frozenset(SUPPORTED_OPS)")
}

func writeOpsPyi(f *os.File, ops []opSpec) {
	fmt.Fprintln(f, "from __future__ import annotations")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "from typing import FrozenSet, Tuple")
	fmt.Fprintln(f)
	for _, op := range ops {
		fmt.Fprintf(f, "%s: str\n", op.PyName)
	}
	fmt.Fprintln(f, "SUPPORTED_OPS: Tuple[str, ...]")
	fmt.Fprintln(f, "SUPPORTED_OPS_SET: FrozenSet[str]")
}

func unquote(s string) (string, error) { return strconv.Unquote(s) }

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			if prev < 'A' || prev > 'Z' {
				result.WriteRune('_')
			} else if i+1 < len(s) {
				next := rune(s[i+1])
				if next >= 'a' && next <= 'z' {
					result.WriteRune('_')
				}
			}
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}
