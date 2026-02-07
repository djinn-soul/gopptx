package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"

	"github.com/djinn09/goppt/internal/pptxxml"
)

// Create builds a valid PPTX with generated slide titles.
func Create(title string, slideCount int) ([]byte, error) {
	if title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if slideCount < 1 {
		return nil, fmt.Errorf("slide count must be at least 1")
	}

	slides := make([]SlideContent, 0, slideCount)
	for i := 1; i <= slideCount; i++ {
		slideTitle := title
		if i > 1 {
			slideTitle = fmt.Sprintf("Slide %d", i)
		}
		slides = append(slides, NewSlide(slideTitle))
	}

	return CreateWithSlides(title, slides)
}

// CreateWithSlides builds a PPTX from caller-provided slide content.
func CreateWithSlides(title string, slides []SlideContent) ([]byte, error) {
	if title == "" {
		return nil, fmt.Errorf("presentation title cannot be empty")
	}
	if len(slides) == 0 {
		return nil, fmt.Errorf("at least one slide is required")
	}
	for i, slide := range slides {
		if err := validateSlide(slide, i+1); err != nil {
			return nil, err
		}
	}

	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)
	count := len(slides)

	if err := writePackageFiles(zw, title, slides, count); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteFile is a convenience helper that writes the generated PPTX to disk.
func WriteFile(path string, title string, slides []SlideContent) error {
	data, err := CreateWithSlides(title, slides)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func writePackageFiles(zw *zip.Writer, title string, slides []SlideContent, slideCount int) error {
	mediaCatalog, err := buildMediaCatalog(slides)
	if err != nil {
		return err
	}
	chartParts := buildChartParts(slides)
	chartBySlide := chartPartBySlide(chartParts)
	notesParts := buildRenderedNotesParts(slides)
	notesTargets := notesTargetBySlide(notesParts)
	hasNotes := len(notesParts) > 0

	files := []struct {
		name    string
		content string
	}{
		{"[Content_Types].xml", pptxxml.ContentTypes(slideCount, mediaCatalog.imageExtensions(), len(chartParts), notesSlideNumbers(notesParts), hasNotes)},
		{"_rels/.rels", pptxxml.RootRelationships()},
		{"ppt/_rels/presentation.xml.rels", pptxxml.PresentationRelationships(slideCount, hasNotes)},
		{"ppt/presentation.xml", pptxxml.Presentation(title, slideCount)},
		{"ppt/slideLayouts/slideLayout1.xml", pptxxml.SlideLayoutTitleAndContent()},
		{"ppt/slideLayouts/_rels/slideLayout1.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout2.xml", pptxxml.SlideLayoutTitleOnly()},
		{"ppt/slideLayouts/_rels/slideLayout2.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout3.xml", pptxxml.SlideLayoutBlank()},
		{"ppt/slideLayouts/_rels/slideLayout3.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout4.xml", pptxxml.SlideLayoutCenteredTitle()},
		{"ppt/slideLayouts/_rels/slideLayout4.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout5.xml", pptxxml.SlideLayoutTitleAndBigContent()},
		{"ppt/slideLayouts/_rels/slideLayout5.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideLayouts/slideLayout6.xml", pptxxml.SlideLayoutTwoColumn()},
		{"ppt/slideLayouts/_rels/slideLayout6.xml.rels", pptxxml.SlideLayoutRelationships()},
		{"ppt/slideMasters/slideMaster1.xml", pptxxml.SlideMaster()},
		{"ppt/slideMasters/_rels/slideMaster1.xml.rels", pptxxml.SlideMasterRelationships()},
		{"ppt/theme/theme1.xml", pptxxml.Theme()},
		{"docProps/core.xml", pptxxml.CoreProperties(title)},
		{"docProps/app.xml", pptxxml.AppProperties(slideCount, len(notesParts))},
	}
	if hasNotes {
		files = append(files,
			struct {
				name    string
				content string
			}{"ppt/notesMasters/notesMaster1.xml", pptxxml.NotesMaster()},
			struct {
				name    string
				content string
			}{"ppt/notesMasters/_rels/notesMaster1.xml.rels", pptxxml.NotesMasterRelationships()},
		)
	}

	for _, item := range files {
		if err := writeFile(zw, item.name, item.content); err != nil {
			return err
		}
	}

	if err := writeMediaFiles(zw, mediaCatalog); err != nil {
		return err
	}
	if err := writeChartFiles(zw, chartParts); err != nil {
		return err
	}
	if err := writeNotesFiles(zw, notesParts); err != nil {
		return err
	}

	for i, slide := range slides {
		slideNumber := i + 1

		var tableSpec *pptxxml.TableSpec
		if slide.Table != nil {
			styledRows, err := tableRowsWithMerges(*slide.Table, slideNumber)
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(styledRows))
			styledSpecRows := make([][]pptxxml.TableCellSpec, 0, len(styledRows))
			for _, srcRow := range styledRows {
				row := make([]string, len(srcRow))
				specRow := make([]pptxxml.TableCellSpec, len(srcRow))
				for i, cell := range srcRow {
					borders := cell.bordersForRender()
					row[i] = cell.Text
					specRow[i] = pptxxml.TableCellSpec{
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
				styledSpecRows = append(styledSpecRows, specRow)
			}
			columnWidths := make([]int64, len(slide.Table.ColumnWidths))
			copy(columnWidths, slide.Table.ColumnWidths)
			rowHeights := make([]int64, len(slide.Table.RowHeights))
			copy(rowHeights, slide.Table.RowHeights)
			tableSpec = &pptxxml.TableSpec{
				X:            slide.Table.X,
				Y:            slide.Table.Y,
				CX:           slide.Table.CX,
				CY:           slide.Table.CY,
				ColumnWidths: columnWidths,
				RowHeights:   rowHeights,
				Rows:         rows,
				StyledRows:   styledSpecRows,
			}
		}

		imageRefs := make([]pptxxml.ImageRef, 0, len(slide.Images))
		imageTargets := make([]string, 0, len(slide.Images))
		for imageIndex, image := range slide.Images {
			mediaName, ok := mediaCatalog.mediaNameForPath(image.Path)
			if !ok {
				return fmt.Errorf("slide %d image %d was not registered", slideNumber, imageIndex+1)
			}
			relID := fmt.Sprintf("rId%d", imageIndex+2)
			imageRefs = append(imageRefs, pptxxml.ImageRef{
				RelID: relID,
				Name:  fmt.Sprintf("Picture %d", imageIndex+1),
				X:     image.X,
				Y:     image.Y,
				CX:    image.CX,
				CY:    image.CY,
			})
			imageTargets = append(imageTargets, fmt.Sprintf("../media/%s", mediaName))
		}

		var chartFrame *pptxxml.ChartFrame
		var chartRel *pptxxml.ChartRel
		if part, ok := chartBySlide[i]; ok {
			rid := fmt.Sprintf("rId%d", len(imageTargets)+2)
			chartFrame = &pptxxml.ChartFrame{
				RelID: rid,
				X:     part.spec.X,
				Y:     part.spec.Y,
				CX:    part.spec.CX,
				CY:    part.spec.CY,
			}
			chartRel = &pptxxml.ChartRel{
				RID:    rid,
				Target: fmt.Sprintf("../charts/chart%d.xml", part.partNumber),
			}
		}

		bulletStyles := toXMLBulletParagraphStyles(slide.BulletStyles)
		bulletRuns := toXMLTextRunRows(slide.BulletRuns)
		shapeSpecs := toXMLShapeSpecs(slide.Shapes)
		connectorSpecs := toXMLConnectorSpecs(slide.Connectors)
		slideXML := pptxxml.SlideWithLayout(
			slideLayoutXMLMode(slide.Layout),
			slide.Title,
			slide.Bullets,
			bulletStyles,
			bulletRuns,
			tableSpec,
			chartFrame,
			imageRefs,
			shapeSpecs,
			connectorSpecs,
		)
		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		if err := writeFile(zw, slidePath, slideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", slideNumber)
		if err := writeFile(
			zw,
			relsPath,
			pptxxml.SlideRelationshipsWithLayoutAndNotes(
				slideLayoutTarget(slide.Layout),
				imageTargets,
				chartRel,
				notesTargets[slideNumber],
			),
		); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(zw *zip.Writer, path string, content string) error {
	w, err := zw.Create(path)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(content))
	return err
}

func writeMediaFiles(zw *zip.Writer, catalog *mediaCatalog) error {
	for _, asset := range catalog.ordered {
		path := fmt.Sprintf("ppt/media/%s", asset.mediaName)
		w, err := zw.Create(path)
		if err != nil {
			return err
		}
		if _, err := w.Write(asset.data); err != nil {
			return err
		}
	}
	return nil
}
