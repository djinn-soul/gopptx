package pptx

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func renderEditorSlideParts(slide SlideContent, slideNumber int, notesTarget string, width, height int64) (string, string, error) {
	tableSpec, err := renderEditorTableSpec(slide, slideNumber)
	if err != nil {
		return "", "", err
	}
	layoutMode := slideLayoutXMLMode(slide.Layout)
	hyperlinkRIDs, hyperlinks, _ := buildSlideHyperlinkRels(slide, 2)
	shapeIDs := calculateShapeIDs(slide)
	animationsXML := slideAnimationsXML(slide, shapeIDs)

	titleSpec := pptxxml.TitleSpec{
		Text:      slide.Title,
		SizePt:    slide.TitleSize,
		Color:     slide.TitleColor,
		Bold:      slide.TitleBold,
		Italic:    slide.TitleItalic,
		Underline: slide.TitleUnderline,
	}
	contentStyle := pptxxml.ContentStyleSpec{
		SizePt:    slide.ContentSize,
		Color:     slide.ContentColor,
		Bold:      slide.ContentBold,
		Italic:    slide.ContentItalic,
		Underline: slide.ContentUnderline,
	}

	slideXML := pptxxml.SlideWithLayout(
		layoutMode,
		titleSpec,
		slide.Bullets,
		toXMLBulletParagraphStyles(slide.BulletStyles),
		toXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs),
		contentStyle,
		tableSpec,
		nil,
		nil,
		toXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
		toXMLConnectorSpecs(slide.Connectors, slide.Shapes),
		nil,
		slideTransitionXML(slide),
		animationsXML,
		width,
		height,
	)
	relsXML := pptxxml.SlideRelationshipsWithHyperlinks(
		slideLayoutTarget(slide.Layout),
		nil,
		nil,
		notesTarget,
		hyperlinks,
	)
	return slideXML, relsXML, nil
}

func renderEditorTableSpec(slide SlideContent, slideNumber int) (*pptxxml.TableSpec, error) {
	if slide.Table == nil {
		return nil, nil
	}

	styledRows, err := tableRowsWithMerges(*slide.Table, slideNumber)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(styledRows))
	specRows := make([][]pptxxml.TableCellSpec, 0, len(styledRows))
	for _, srcRow := range styledRows {
		row := make([]string, len(srcRow))
		specRow := make([]pptxxml.TableCellSpec, len(srcRow))
		for idx, cell := range srcRow {
			borders := cell.bordersForRender()
			row[idx] = cell.Text
			specRow[idx] = pptxxml.TableCellSpec{
				Text:            cell.Text,
				Bold:            cell.Bold,
				BackgroundColor: cell.BackgroundColor,
				Align:           cell.Align,
				VAlign:          cell.VAlign,
				MarginLeft:      tableMarginEMU(cell.MarginLeftPt),
				MarginRight:     tableMarginEMU(cell.MarginRightPt),
				MarginTop:       tableMarginEMU(cell.MarginTopPt),
				MarginBottom:    tableMarginEMU(cell.MarginBottomPt),
				WrapText:        cloneBoolPointer(cell.WrapText),
				RowSpan:         cell.RowSpan,
				ColSpan:         cell.ColSpan,
				VMerge:          cell.VMerge,
				HMerge:          cell.HMerge,
				BorderColor:     cell.BorderColor,
				BorderWidth:     tableBorderWidthEMU(cell.BorderWidthPt),
				BorderLeft:      toXMLTableBorderSpec(borders.Left),
				BorderRight:     toXMLTableBorderSpec(borders.Right),
				BorderTop:       toXMLTableBorderSpec(borders.Top),
				BorderBottom:    toXMLTableBorderSpec(borders.Bottom),
			}
		}
		rows = append(rows, row)
		specRows = append(specRows, specRow)
	}

	columnWidths := make([]int64, len(slide.Table.ColumnWidths))
	copy(columnWidths, slide.Table.ColumnWidths)
	rowHeights := make([]int64, len(slide.Table.RowHeights))
	copy(rowHeights, slide.Table.RowHeights)
	return &pptxxml.TableSpec{
		X:            slide.Table.X,
		Y:            slide.Table.Y,
		CX:           slide.Table.CX,
		CY:           slide.Table.CY,
		ColumnWidths: columnWidths,
		RowHeights:   rowHeights,
		Rows:         rows,
		StyledRows:   specRows,
	}, nil
}

func editorEnsureSlideRelsExist(parts map[string][]byte, slidePart string) error {
	relsPath := slideRelsPartName(slidePart)
	if _, ok := parts[relsPath]; ok {
		return nil
	}
	return fmt.Errorf("missing slide relationships part %q", relsPath)
}

func buildSlideHyperlinkRels(slide SlideContent, firstRID int) (map[*Hyperlink]string, []pptxxml.HyperlinkRel, int) {
	hyperlinkRIDs := make(map[*Hyperlink]string)
	hyperlinks := make([]pptxxml.HyperlinkRel, 0)
	nextRID := firstRID

	addHyperlink := func(h *Hyperlink) {
		if h == nil {
			return
		}
		if _, exists := hyperlinkRIDs[h]; exists {
			return
		}

		rid := fmt.Sprintf("rId%d", nextRID)
		hyperlinkRIDs[h] = rid
		nextRID++

		hyperlinks = append(hyperlinks, pptxxml.HyperlinkRel{
			RID:      rid,
			Target:   h.Action.RelationshipTarget(),
			External: h.Action.IsExternal(),
		})
	}

	for _, shape := range slide.Shapes {
		addHyperlink(shape.Hyperlink)
	}
	for _, runRow := range slide.BulletRuns {
		for _, run := range runRow {
			addHyperlink(run.Hyperlink)
		}
	}

	return hyperlinkRIDs, hyperlinks, nextRID
}
