package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	receiverType = "SlideContent"
	initialCap   = 128
	genSuffix    = "_gen"
)

// method is one generated delegating method.
type method struct {
	Name   string
	Params string // rendered parameter list, e.g. "text string, style ParagraphStyle"
	Args   string // matching argument list, e.g. "text, style"
	Doc    string // first line of the source doc comment, minus the method name
	Quals  []string
}

// collect returns the chainable SlideContent methods in dir, plus the import
// paths their parameter types need, keyed by package qualifier.
func collect(dir string) ([]method, map[string]string, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, nil, fmt.Errorf("glob %s: %w", dir, err)
	}

	fset := token.NewFileSet()
	methods := make([]method, 0, initialCap)
	available := map[string]string{}

	for _, path := range entries {
		if strings.HasSuffix(path, "_test.go") || strings.Contains(path, genSuffix+".go") {
			continue
		}
		file, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("parse %s: %w", path, parseErr)
		}
		collectImports(file, available)
		methods = append(methods, collectFileMethods(fset, file)...)
	}

	sort.Slice(methods, func(i, j int) bool { return methods[i].Name < methods[j].Name })
	return methods, available, nil
}

func collectImports(file *ast.File, into map[string]string) {
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		name := path[strings.LastIndex(path, "/")+1:]
		if spec.Name != nil {
			name = spec.Name.Name
		}
		into[name] = path
	}
}

func collectFileMethods(fset *token.FileSet, file *ast.File) []method {
	out := make([]method, 0, initialCap)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || !isChainable(fn) {
			continue
		}
		quals := map[string]bool{}
		params, args := renderParams(fset, fn.Type.Params, quals)
		out = append(out, method{
			Name:   fn.Name.Name,
			Params: params,
			Args:   args,
			Doc:    firstDocLine(fn),
			Quals:  sortedKeys(quals),
		})
	}
	return out
}

// isChainable reports whether fn is an exported SlideContent method whose only
// result is a SlideContent, i.e. one of the builder methods to wrap.
func isChainable(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || len(fn.Recv.List) != 1 || !fn.Name.IsExported() {
		return false
	}
	recv, ok := fn.Recv.List[0].Type.(*ast.Ident)
	if !ok || recv.Name != receiverType {
		return false
	}
	results := fn.Type.Results
	if results == nil || len(results.List) != 1 || len(results.List[0].Names) > 0 {
		return false
	}
	result, ok := results.List[0].Type.(*ast.Ident)
	return ok && result.Name == receiverType
}

// renderParams renders a parameter list and the matching argument list, naming
// any parameter the source left unnamed and forwarding variadics correctly.
func renderParams(fset *token.FileSet, fields *ast.FieldList, quals map[string]bool) (string, string) {
	if fields == nil || len(fields.List) == 0 {
		return "", ""
	}

	var params, args []string
	index := 0
	for _, field := range fields.List {
		typeText := renderExpr(fset, field.Type)
		recordQualifiers(field.Type, quals)

		names := make([]string, 0, len(field.Names))
		if len(field.Names) == 0 {
			names = append(names, fmt.Sprintf("arg%d", index))
			index++
		}
		for _, name := range field.Names {
			names = append(names, name.Name)
			index++
		}

		params = append(params, strings.Join(names, ", ")+" "+typeText)
		for _, name := range names {
			if strings.HasPrefix(typeText, "...") {
				args = append(args, name+"...")
				continue
			}
			args = append(args, name)
		}
	}
	return strings.Join(params, ", "), strings.Join(args, ", ")
}

func renderExpr(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return buf.String()
}

// recordQualifiers walks a type expression and notes every package qualifier,
// so each generated file imports exactly what its own signatures mention.
func recordQualifiers(expr ast.Expr, quals map[string]bool) {
	ast.Inspect(expr, func(node ast.Node) bool {
		sel, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if ident, ok := sel.X.(*ast.Ident); ok {
			quals[ident.Name] = true
		}
		return true
	})
}

func firstDocLine(fn *ast.FuncDecl) string {
	if fn.Doc == nil || len(fn.Doc.List) == 0 {
		return ""
	}
	line := strings.TrimSpace(strings.TrimPrefix(fn.Doc.List[0].Text, "//"))
	return strings.TrimPrefix(line, fn.Name.Name+" ")
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
