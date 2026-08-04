package editor

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

// Namespaces and content type for the chart drawing part that holds user shapes.
const (
	chartDrawingNS = "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing"
	// RelTypeChartUserShapes is the relationship from a chart part to its
	// drawing part.
	RelTypeChartUserShapes = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chartUserShapes"
	chartUserShapesCT      = "application/vnd.openxmlformats-officedocument.drawingml.chartshapes+xml"
)

// shapeIDBase is the first id a user shape may take: id 1 belongs to the
// drawing's own group shape.
const shapeIDBase = 2

// userShapesRefPattern matches an existing c:userShapes reference in a chart part.
var userShapesRefPattern = regexp.MustCompile(`<c:userShapes\b[^>]*/>`)

// ChartUserShape is one shape drawn on top of a chart, anchored in fractions of
// the chart area so it moves and scales with the chart.
type ChartUserShape struct {
	// Text is the shape's caption.
	Text string
	// FromX, FromY, ToX and ToY are fractions of the chart area in [0,1].
	FromX, FromY float64
	ToX, ToY     float64
	// FontSizePt sets the caption size; zero leaves it to the chart defaults.
	FontSizePt float64
	// Bold draws the caption bold.
	Bold bool
	// Name is the shape name; a default is generated when empty.
	Name string
}

// AddChartUserShapes attaches shapes to a chart through its c:userShapes
// drawing part, which is what makes a caption part of the chart instead of a
// separate slide shape that has to be moved and copied alongside it
// (upstream issue #351).
//
// The call replaces any user shapes the chart already had.
func (e *PresentationEditor) AddChartUserShapes(
	slideIndex int,
	selector common.ChartSelector,
	userShapes []ChartUserShape,
) (string, error) {
	if len(userShapes) == 0 {
		return "", fmt.Errorf("no user shapes given for the chart on slide %d", slideIndex)
	}
	for i, shape := range userShapes {
		if err := validateChartUserShape(i, shape); err != nil {
			return "", err
		}
	}

	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return "", err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return "", err
	}
	chartPart := chartRef.ChartPart
	chartXML, ok := e.parts.Get(chartPart)
	if !ok {
		return "", fmt.Errorf("chart part %s not found", chartPart)
	}

	drawingPart := e.nextChartDrawingPart()
	e.parts.Set(drawingPart, []byte(renderChartUserShapesXML(userShapes)))
	e.addContentTypeOverride(drawingPart, chartUserShapesCT)

	relID, err := e.allocChartRelID(chartPart)
	if err != nil {
		return "", fmt.Errorf("allocate chart rel id: %w", err)
	}
	if addErr := e.addRelationship(
		chartPart,
		relID,
		RelTypeChartUserShapes,
		"../drawings/"+path.Base(drawingPart),
	); addErr != nil {
		return "", fmt.Errorf("add chart user shapes rel: %w", addErr)
	}

	updated, err := insertUserShapesRef(string(chartXML), relID)
	if err != nil {
		return "", err
	}
	e.parts.Set(chartPart, []byte(updated))
	return drawingPart, nil
}

func validateChartUserShape(index int, shape ChartUserShape) error {
	for name, value := range map[string]float64{
		"from_x": shape.FromX, "from_y": shape.FromY,
		"to_x": shape.ToX, "to_y": shape.ToY,
	} {
		if value < 0 || value > 1 {
			return fmt.Errorf("user shape %d: %s must be a fraction in [0,1], got %v", index, name, value)
		}
	}
	if shape.ToX <= shape.FromX || shape.ToY <= shape.FromY {
		return fmt.Errorf("user shape %d: the 'to' anchor must be past the 'from' anchor", index)
	}
	return nil
}

// nextChartDrawingPart picks a free ppt/drawings/drawingN.xml path.
func (e *PresentationEditor) nextChartDrawingPart() string {
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("ppt/drawings/drawing%d.xml", i)
		if !e.parts.Has(candidate) {
			return candidate
		}
	}
}

// insertUserShapesRef puts the c:userShapes reference in the chart part.
// CT_ChartSpace places it near the end, so writing it just before the closing
// tag keeps the element order valid.
func insertUserShapesRef(chartXML, relID string) (string, error) {
	ref := `<c:userShapes r:id="` + common.XMLEscape(relID) + `"/>`
	if userShapesRefPattern.MatchString(chartXML) {
		return userShapesRefPattern.ReplaceAllLiteralString(chartXML, ref), nil
	}
	const closeTag = "</c:chartSpace>"
	idx := strings.LastIndex(chartXML, closeTag)
	if idx < 0 {
		return "", fmt.Errorf("chart part has no %s end tag", closeTag)
	}
	return chartXML[:idx] + ref + chartXML[idx:], nil
}

// renderChartUserShapesXML renders the chart drawing part. Each shape uses a
// relative-size anchor, so it keeps its place when the chart is resized.
func renderChartUserShapesXML(userShapes []ChartUserShape) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<c:userShapes xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" `)
	b.WriteString(`xmlns:cdr="` + chartDrawingNS + `" `)
	b.WriteString(`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">`)

	for i, shape := range userShapes {
		name := shape.Name
		if strings.TrimSpace(name) == "" {
			name = fmt.Sprintf("Chart Shape %d", i+1)
		}
		fmt.Fprintf(&b,
			`<cdr:relSizeAnchor>`+
				`<cdr:from><cdr:x>%g</cdr:x><cdr:y>%g</cdr:y></cdr:from>`+
				`<cdr:to><cdr:x>%g</cdr:x><cdr:y>%g</cdr:y></cdr:to>`+
				`<cdr:sp macro="" textlink="">`+
				`<cdr:nvSpPr><cdr:cNvPr id="%d" name="%s"/>`+
				`<cdr:cNvSpPr txBox="1"/></cdr:nvSpPr>`+
				`<cdr:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/></a:xfrm>`+
				`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:noFill/></cdr:spPr>`+
				`<cdr:txBody><a:bodyPr wrap="square" rtlCol="0"/><a:lstStyle/>`+
				`<a:p><a:r>%s<a:t>%s</a:t></a:r></a:p>`+
				`</cdr:txBody></cdr:sp></cdr:relSizeAnchor>`,
			shape.FromX, shape.FromY, shape.ToX, shape.ToY,
			i+shapeIDBase, // shape ids start after the group id
			common.XMLEscape(name),
			chartUserShapeRunPropsXML(shape),
			common.XMLEscape(shape.Text),
		)
	}

	b.WriteString(`</c:userShapes>`)
	return b.String()
}

func chartUserShapeRunPropsXML(shape ChartUserShape) string {
	var b strings.Builder
	b.WriteString(`<a:rPr lang="en-US"`)
	if shape.FontSizePt > 0 {
		fmt.Fprintf(&b, ` sz="%d"`, int(shape.FontSizePt*100)) //nolint:mnd // sz is hundredths of a point
	}
	if shape.Bold {
		b.WriteString(` b="1"`)
	}
	b.WriteString(`/>`)
	return b.String()
}
