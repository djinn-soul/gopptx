package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

type slideParts struct {
	title                pptxxml.TitleSpec
	contentStyle         pptxxml.ContentStyleSpec
	table                *pptxxml.TableSpec
	imageRefs            []pptxxml.ImageRef
	backgroundRID        string
	transitionXML        string
	placeholders         []pptxxml.PlaceholderOverrideSpec
	chartFrame           *pptxxml.ChartFrame
	chartRel             *pptxxml.ChartRel
	placeholderChartRels []pptxxml.ChartRel
	smartArtFrames       []pptxxml.SmartArtFrame
	smartArtRels         []pptxxml.SmartArtRel
}

func renderSlides(
	pw *pptxxml.PackageWriter,
	meta Metadata,
	slides []elements.SlideContent,
	catalog *media.Catalog,
	chartBySlide map[int][]ChartPart,
	smartArtBySlide map[int][]SmartArtPart,
	notesTargets map[int]string,
	masterCount int,
	commentsBySlide map[int][]comments.Comment,
) error {
	for i, slide := range slides {
		num := i + 1
		builder := newSlidePartBuilder(num, catalog)
		parts, err := builder.build(i, slide, chartBySlide, smartArtBySlide)
		if err != nil {
			return err
		}

		hyperlinkRIDs, hyperlinks, _ := elements.BuildSlideHyperlinkRels(
			slide,
			builder.ridNext,
		)

		slideXML := pptxxml.SlideWithLayout(
			elements.SlideLayoutXMLMode(slide.Layout),
			parts.title,
			slide.Bullets,
			elements.ToXMLBulletParagraphStyles(slide.BulletStyles),
			elements.ToXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs),
			parts.contentStyle,
			parts.table,
			parts.chartFrame,
			parts.imageRefs,
			shapes.ToXMLShapeSpecs(slide.Shapes, hyperlinkRIDs),
			shapes.ToXMLConnectorSpecs(slide.Connectors, slide.Shapes),
			parts.placeholders,
			parts.smartArtFrames,
			elements.ToXMLBackgroundSpec(slide.Background, parts.backgroundRID),
			parts.transitionXML,
			elements.SlideAnimationsXML(slide, elements.CalculateShapeIDs(slide)),
			slide.ShowSlideNumber,
			func() string {
				if slide.FooterText != "" {
					return slide.FooterText
				}
				return meta.FooterText
			}(),
			meta.ShowDateTime,
			meta.SlideSize.Width,
			meta.SlideSize.Height,
		)

		layoutTarget := elements.SlideLayoutTarget(slide.Layout)
		if masterCount > 1 {
			masterNum := (i % masterCount) + 1
			layoutTarget = layoutTargetForMaster(layoutTarget, masterNum)
		}

		pw.AddPart(fmt.Sprintf("ppt/slides/slide%d.xml", num), slideXML)
		commentTarget := ""
		if len(commentsBySlide[i]) > 0 {
			pw.AddPart(fmt.Sprintf("ppt/comments/comment%d.xml", num), pptxxml.CommentsXML(commentsBySlide[i]))
			commentTarget = fmt.Sprintf("../comments/comment%d.xml", num)
		}

		pw.AddPart(fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num), pptxxml.SlideRelationshipsWithAll(
			layoutTarget,
			builder.targets,
			parts.chartRel,

			parts.placeholderChartRels,
			parts.smartArtRels,
			notesTargets[num],
			hyperlinks,
			commentTarget,
		))
	}
	return nil
}

func layoutTargetForMaster(target string, masterNum int) string {
	if masterNum <= 1 {
		return target
	}
	// Target is like "../slideLayouts/slideLayout1.xml"
	// For master 2, layouts are 7-12
	var num int
	if n, _ := fmt.Sscanf(target, "../slideLayouts/slideLayout%d.xml", &num); n != 1 {
		return target
	}
	newNum := (masterNum-1)*6 + num
	return fmt.Sprintf("../slideLayouts/slideLayout%d.xml", newNum)
}

func mapOptionalLength(l *styling.Length) *int64 {
	if l == nil {
		return nil
	}
	v := l.Emu()
	return &v
}

func mapPlaceholderTextStyle(ts *shapes.PlaceholderTextStyle) *pptxxml.PlaceholderTextStyleSpec {
	if ts == nil {
		return nil
	}
	return &pptxxml.PlaceholderTextStyleSpec{
		SizePt:    ts.SizePt,
		Color:     ts.Color,
		Bold:      ts.Bold,
		Italic:    ts.Italic,
		Underline: ts.Underline,
		Align:     ts.Align,
		Font:      ts.Font,
	}
}
