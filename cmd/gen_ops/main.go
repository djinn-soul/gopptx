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

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: gen_ops <input_go_file> <output_py_file> <output_pyi_file>")
		os.Exit(1)
	}

	input := os.Args[1]
	outputPy := os.Args[2]
	outputPyi := os.Args[3]

	ops, err := parseOpsFromGo(input)
	if err != nil {
		fmt.Printf("Error parsing ops: %v\n", err)
		os.Exit(1)
	}
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].PyName < ops[j].PyName
	})

	pyFile, err := os.Create(outputPy)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer pyFile.Close()

	pyiFile, err := os.Create(outputPyi)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer pyiFile.Close()

	writeOpsPy(pyFile, ops)
	writeOpsPyi(pyiFile, ops)
}

func parseOpsFromGo(input string) ([]opSpec, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, input, nil, 0)
	if err != nil {
		return nil, err
	}

	ops := make([]opSpec, 0, 64)
	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok || len(vspec.Names) == 0 || len(vspec.Values) == 0 {
				continue
			}
			name := vspec.Names[0].Name
			if !strings.HasPrefix(name, "Op") {
				continue
			}
			lit, ok := vspec.Values[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return nil, fmt.Errorf("op %q must be a string literal", name)
			}
			value, unquoteErr := unquote(lit.Value)
			if unquoteErr != nil {
				return nil, fmt.Errorf("unquote op %q: %w", name, unquoteErr)
			}

			ops = append(ops, opSpec{
				PyName: "OP_" + toSnakeCase(strings.TrimPrefix(name, "Op")),
				Value:  value,
			})
		}
	}
	return ops, nil
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
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToUpper(result.String())
}
