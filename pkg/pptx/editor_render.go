package pptx

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func renderEditorSlideParts(slide SlideContent, slideNumber int, notesTarget string, width, height int64) (string, string, error) {
	tableSpec, err := renderEditorTableSpec(slide, slideNumber)
	if err != nil {
		return "", "", err
	}
	layoutMode := elements.SlideLayoutXMLMode(slide.Layout)
	hyperlinkRIDs, hyperlinks, _ := buildSlideHyperlinkRels(slide, 2)
	shapeIDs := elements.CalculateShapeIDs(slide)
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
		elements.SlideLayoutTarget(slide.Layout),
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

	return slide.Table.ToTableSpec(slideNumber)
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
