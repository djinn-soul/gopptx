package export

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	placeholderTitle       = "title"
	placeholderCtrTitle    = "ctrtitle"
	placeholderBody        = "body"
	placeholderSubtitle    = "subtitle"
	placeholderObject      = "obj"
	placeholderContentName = "content"
	minTableBorderWidthPt  = 0.25
)

// SlidesFromPPTX reads an existing PPTX file and extracts slide content
// (title, bullets, shapes, embedded images) for the native PDF/HTML export pipeline.
func SlidesFromPPTX(pptxPath string) (string, []elements.SlideContent, error) {
	ed, err := editor.OpenPresentationEditor(pptxPath)
	if err != nil {
		return "", nil, fmt.Errorf("open PPTX: %w", err)
	}
	defer ed.Close()

	meta := ed.Metadata()
	presTitle := ""
	if meta != nil {
		presTitle = meta.Title
	}

	// Extract embedded images per slide via direct PPTX zip parsing.
	slideImages, err := extractSlideImages(pptxPath)
	if err != nil {
		// Best-effort image extraction; continue without images when parsing fails.
		slideImages = nil
	}

	slideMeta := ed.Slides()
	slideContents := make([]elements.SlideContent, 0, len(slideMeta))

	for _, sm := range slideMeta {
		slideContents = append(slideContents, extractSlideContent(ed, sm, slideImages))
	}

	if presTitle == "" && len(slideContents) > 0 {
		presTitle = slideContents[0].Title
	}

	return presTitle, slideContents, nil
}

//nolint:gocognit // Slide extraction maps multiple editor shape categories and image attachments in one pass.
func extractSlideContent(
	ed *editor.PresentationEditor,
	sm editorcommon.SlideMetadata,
	slideImages [][]SlideImage,
) elements.SlideContent {
	editorShapes, err := ed.GetShapes(sm.Index)
	if err != nil {
		editorShapes = nil
	}

	sc := elements.SlideContent{
		Title: sm.Title,
	}

	for _, es := range editorShapes {
		lowerType := strings.ToLower(es.Type)
		lowerName := strings.ToLower(strings.TrimSpace(es.Name))

		switch lowerType {
		case "graphicframe":
			if tbl := extractTableContent(ed, sm.Index, es); tbl != nil {
				sc.Table = tbl
				continue
			}
			sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
		case placeholderTitle, placeholderCtrTitle:
			if sc.Title == "" && es.Text != "" {
				sc.Title = es.Text
			}
		case placeholderBody, placeholderSubtitle, placeholderObject:
			sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
		default:
			switch {
			case isTitlePlaceholder(lowerType, lowerName):
				if sc.Title == "" && es.Text != "" {
					sc.Title = es.Text
				}
			case isBodyPlaceholder(lowerType, lowerName):
				sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
			default:
				sc.Shapes = append(sc.Shapes, editorShapeToShape(es))
			}
		}
	}

	// Attach images for this slide.
	if sm.Index < len(slideImages) {
		for _, img := range slideImages[sm.Index] {
			sc.Images = append(sc.Images, shapes.Image{
				Data:   img.Bytes,
				Format: img.Format,
				X:      styling.Emu(img.X),
				Y:      styling.Emu(img.Y),
				CX:     styling.Emu(img.CX),
				CY:     styling.Emu(img.CY),
			})
		}
	}
	return sc
}

func extractTableContent(
	ed *editor.PresentationEditor,
	slideIndex int,
	shape editorcommon.Shape,
) *tables.Table {
	if shape.ID <= 0 {
		return nil
	}
	info, err := ed.GetTable(slideIndex, shape.ID)
	if err != nil {
		return nil
	}
	tableAny, ok := info["table"]
	if !ok {
		return nil
	}
	tableMeta, ok := tableAny.(map[string]any)
	if !ok {
		return nil
	}
	rowsMeta := toMapSlice(tableMeta["rows"])
	if len(rowsMeta) == 0 {
		return nil
	}
	firstRow, _ := tableMeta["first_row"].(bool)
	bandRow, _ := tableMeta["band_row"].(bool)
	colWidths := toLengthSlice(tableMeta["column_widths"])
	rowHeights := toLengthSlice(tableMeta["row_heights"])
	table := tables.NewTable(colWidths)
	table = table.Position(styling.Emu(int64(shape.X)), styling.Emu(int64(shape.Y)))
	table = table.Size(styling.Emu(int64(shape.W)), styling.Emu(int64(shape.H)))
	if len(rowHeights) > 0 {
		table = table.WithRowHeights(rowHeights)
	}
	for rowIndex, row := range rowsMeta {
		cellsMeta := toMapSlice(row["cells"])
		if len(cellsMeta) == 0 {
			continue
		}
		values := make([]string, 0, len(cellsMeta))
		styled := make([]tables.TableCell, 0, len(cellsMeta))
		for _, cell := range cellsMeta {
			value, _ := cell["text"].(string)
			values = append(values, value)
			styledCell := tables.NewTableCell(value)
			applyTableCellLayout(&styledCell, cell)
			applyTableCellBorders(&styledCell, cell)
			if firstRow && rowIndex == 0 {
				styledCell = styledCell.WithBold(true).WithBackgroundColor("4472C4")
				styledCell.Color = "FFFFFF"
			} else if bandRow && rowIndex%2 == 1 {
				styledCell = styledCell.WithBackgroundColor("D0D8E8")
			}
			styled = append(styled, styledCell)
		}
		if firstRow || bandRow {
			table = table.AddStyledRow(styled)
			continue
		}
		table = table.AddRow(values)
	}
	return &table
}

func applyTableCellLayout(cell *tables.TableCell, cellMeta map[string]any) {
	if rawAlign, _ := cellMeta["v_align"].(string); rawAlign != "" {
		switch strings.ToLower(strings.TrimSpace(rawAlign)) {
		case "top":
			*cell = cell.WithVAlignTop()
		case "bottom":
			*cell = cell.WithVAlignBottom()
		case "middle", "center", "ctr":
			*cell = cell.WithVAlignMiddle()
		}
	}
	applyMargin := func(
		margin any,
		setter func(float64) tables.TableCell,
	) {
		emu, ok := parseNumericInt64(margin)
		if !ok || emu <= 0 {
			return
		}
		pts := math.Max(emuToPt(emu), 0)
		*cell = setter(pts)
	}
	applyMargin(cellMeta["margin_left"], cell.WithMarginLeftPt)
	applyMargin(cellMeta["margin_right"], cell.WithMarginRightPt)
	applyMargin(cellMeta["margin_top"], cell.WithMarginTopPt)
	applyMargin(cellMeta["margin_bottom"], cell.WithMarginBottomPt)
}

func applyTableCellBorders(cell *tables.TableCell, cellMeta map[string]any) {
	applyBorder := func(
		key string,
		setter func(float64, string) tables.TableCell,
	) {
		borderMeta, ok := cellMeta[key].(map[string]any)
		if !ok {
			return
		}
		widthEmu, ok := parseNumericInt64(borderMeta["width"])
		if !ok || widthEmu <= 0 {
			return
		}
		color := "B4B4B4"
		if c, ok := borderMeta["color"].(string); ok && strings.TrimSpace(c) != "" {
			color = c
		}
		*cell = setter(math.Max(emuToPt(widthEmu), minTableBorderWidthPt), color)
	}
	applyBorder("border_left", cell.WithLeftBorder)
	applyBorder("border_right", cell.WithRightBorder)
	applyBorder("border_top", cell.WithTopBorder)
	applyBorder("border_bottom", cell.WithBottomBorder)
}

func parseNumericInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func toMapSlice(value any) []map[string]any {
	switch items := value.(type) {
	case []map[string]any:
		return items
	case []any:
		out := make([]map[string]any, 0, len(items))
		for _, item := range items {
			m, ok := item.(map[string]any)
			if ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return nil
	}
}

func toLengthSlice(value any) []styling.Length {
	switch values := value.(type) {
	case []int64:
		out := make([]styling.Length, 0, len(values))
		for _, v := range values {
			out = append(out, styling.Emu(v))
		}
		return out
	case []any:
		out := make([]styling.Length, 0, len(values))
		for _, raw := range values {
			switch v := raw.(type) {
			case int64:
				out = append(out, styling.Emu(v))
			case int:
				out = append(out, styling.Emu(int64(v)))
			case float64:
				out = append(out, styling.Emu(int64(v)))
			}
		}
		return out
	default:
		return nil
	}
}

// editorShapeToShape maps an editor common.Shape to an export shapes.Shape.
// X, Y, W, H from the editor are in EMU (int).
func editorShapeToShape(es editorcommon.Shape) shapes.Shape {
	shapeFill := editorFillToExportFill(es.Fill)
	shapeLine := editorLineToExportLine(es.Line)
	adjustments := editorAdjustmentsToExport(es.Adjustments)
	return shapes.Shape{
		// Map OOXML preset geometry name directly — pdf_native.go uses these strings.
		Type:        editorTypeToPreset(es.Type),
		X:           styling.Emu(int64(es.X)),
		Y:           styling.Emu(int64(es.Y)),
		CX:          styling.Emu(int64(es.W)),
		CY:          styling.Emu(int64(es.H)),
		Text:        es.Text,
		Name:        es.Name,
		Fill:        shapeFill,
		Line:        shapeLine,
		Adjustments: adjustments,
	}
}

func editorFillToExportFill(fill *editorcommon.ShapeFill) *shapes.ShapeFill {
	if fill == nil || fill.Solid == nil || *fill.Solid == "" {
		return nil
	}
	return &shapes.ShapeFill{Color: *fill.Solid}
}

func editorLineToExportLine(line *editorcommon.ShapeLine) *shapes.ShapeLine {
	if line == nil || line.Color == nil || *line.Color == "" {
		return nil
	}
	width := styling.Emu(0)
	if line.WidthEmu != nil && *line.WidthEmu > 0 {
		width = styling.Emu(int64(*line.WidthEmu))
	}
	return &shapes.ShapeLine{
		Color: *line.Color,
		Width: width,
	}
}

func editorAdjustmentsToExport(
	adjustments []editorcommon.ShapeAdjustment,
) []shapes.ShapeAdjustment {
	if len(adjustments) == 0 {
		return nil
	}
	out := make([]shapes.ShapeAdjustment, 0, len(adjustments))
	for _, adjustment := range adjustments {
		if adjustment.Name == "" || adjustment.Formula == "" {
			continue
		}
		out = append(out, shapes.ShapeAdjustment{
			Name:    adjustment.Name,
			Formula: adjustment.Formula,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func isTitlePlaceholder(shapeType, shapeName string) bool {
	return shapeType == placeholderTitle ||
		shapeType == placeholderCtrTitle ||
		shapeName == placeholderTitle ||
		strings.Contains(shapeName, "title placeholder")
}

func isBodyPlaceholder(shapeType, shapeName string) bool {
	if shapeType == placeholderBody || shapeType == placeholderSubtitle || shapeType == placeholderObject {
		return true
	}
	return shapeName == "content" ||
		shapeName == placeholderBody ||
		strings.Contains(shapeName, "content placeholder") ||
		strings.Contains(shapeName, "body placeholder")
}

// editorTypeToPreset normalizes the editor shape type string to the OOXML
// preset geometry name used by our shape renderers.
func editorTypeToPreset(t string) string {
	switch strings.ToLower(t) {
	case "rect", "rectangle":
		return "rect"
	case "roundrect", "roundedrectangle":
		return "roundRect"
	case "ellipse", "oval", "circle":
		return "ellipse"
	case "triangle", "rt_triangle":
		return "triangle"
	case "rightarrow":
		return "rightArrow"
	case "leftarrow":
		return "leftArrow"
	default:
		// Pass through unknown presets as-is; renderers will fall back to rect.
		return t
	}
}
