package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

const (
	chartTypePrefix = "XLChartType"
	chartKindPrefix = "ChartKind"
	initialCap      = 32
)

// chartType is one generated Python enum member.
type chartType struct {
	PyName string
	Value  string
	GoName string
}

// parseChartKinds collects ChartKind* string constants by suffix.
func parseChartKinds(path string) (map[string]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	out := make(map[string]string, initialCap)
	forEachConst(file, func(vspec *ast.ValueSpec) {
		for i, name := range vspec.Names {
			if !strings.HasPrefix(name.Name, chartKindPrefix) || i >= len(vspec.Values) {
				continue
			}
			lit, ok := vspec.Values[i].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			value, unquoteErr := strconv.Unquote(lit.Value)
			if unquoteErr == nil {
				out[strings.TrimPrefix(name.Name, chartKindPrefix)] = value
			}
		}
	})
	return out, nil
}

// parseChartTypes resolves XLChartType* constants to their ChartKind values.
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
			value, valueErr := resolveChartKind(vspec.Values[i], kinds)
			if valueErr != nil {
				resolveErr = fmt.Errorf("%s: %w", name.Name, valueErr)
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
