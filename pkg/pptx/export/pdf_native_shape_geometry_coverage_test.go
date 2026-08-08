package export

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/signintech/gopdf"
)

// declaredShapePresets reads the ShapeType* constants straight out of the
// shapes package, so a preset added there is covered by this test the moment it
// lands rather than whenever someone remembers to update a list here.
func declaredShapePresets(t *testing.T) map[string]string {
	t.Helper()
	matches, err := filepath.Glob("../shapes/shape_types*.go")
	if err != nil || len(matches) == 0 {
		t.Skipf("shape type sources not readable from here: %v", err)
	}
	presets := map[string]string{}
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		file, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if parseErr != nil {
			t.Fatalf("parse %s: %v", path, parseErr)
		}
		collectPresetConstants(file, presets)
	}
	return presets
}

func collectPresetConstants(file *ast.File, into map[string]string) {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range value.Names {
				if !strings.HasPrefix(name.Name, "ShapeType") || i >= len(value.Values) {
					continue
				}
				literal, isLiteral := value.Values[i].(*ast.BasicLit)
				if !isLiteral || literal.Kind != token.STRING {
					continue
				}
				into[name.Name] = strings.Trim(literal.Value, `"`)
			}
		}
	}
}

// drawsThroughATable reports whether any of the geometry tables claims the
// preset. The tables are the whole of the non-rectangle geometry, so a preset
// no table claims is one that silently draws as a box.
func drawsThroughATable(pdf *gopdf.GoPdf, fl flipState, preset string) bool {
	const x, y, w, h = 0.0, 0.0, 120.0, 80.0
	return drawPDFExtendedGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFCornerGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFSolidGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFPolygonGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFArrowGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFCurveGeometry(pdf, fl, preset, x, y, w, h, "D") ||
		drawPDFCalloutGeometry(pdf, fl, preset, x, y, w, h, "D")
}

// presetsDrawnByTheMainSwitch are handled by name in drawPDFGeometry itself
// rather than by one of the tables, so they never reach the fallback.
//
//nolint:gochecknoglobals // Test data; a map literal reads better than a func here.
var presetsDrawnByTheMainSwitch = map[string]bool{
	"rect": true, "roundRect": true, "ellipse": true, "triangle": true,
	"rtTriangle": true, "diamond": true, "hexagon": true, "pentagon": true,
	"octagon": true, "parallelogram": true, "trapezoid": true,
	"rightArrow": true, "leftArrow": true, "upArrow": true, "downArrow": true,
	"leftRightArrow": true, "upDownArrow": true, "chevron": true,
	"star4": true, "star5": true, "star6": true, "star7": true, "star8": true,
	"star10": true, "star12": true, "star16": true, "star24": true, "star32": true,
	"heart": true, "wedgeRectCallout": true, "wedgeRRectCallout": true,
	"wedgeEllipseCallout": true, "cloudCallout": true, "cloud": true,
	"pie": true, "chord": true, "pieWedge": true,
}

func TestEveryDeclaredPresetHasRealGeometry(t *testing.T) {
	pdf := newTestPDF(t)
	fl := flipState{unflippedShape: true}

	var missing []string
	presets := declaredShapePresets(t)
	for name, value := range presets {
		if presetsDrawnByTheMainSwitch[value] || drawsThroughATable(pdf, fl, value) {
			continue
		}
		missing = append(missing, name+" ("+value+")")
	}
	if len(missing) > 0 {
		t.Errorf("%d presets still fall through to the rectangle fallback:\n  %s",
			len(missing), strings.Join(missing, "\n  "))
	}
	if len(presets) < 200 {
		t.Errorf("only %d presets were read from the shapes package; the glob has gone stale", len(presets))
	}
	t.Logf("checked %d declared presets", len(presets))
}

func TestUnknownPresetIsClaimedByNoTable(t *testing.T) {
	pdf := newTestPDF(t)
	if drawsThroughATable(pdf, flipState{unflippedShape: true}, "definitelyNotAPreset") {
		t.Error("a geometry table claimed a preset that does not exist")
	}
}
