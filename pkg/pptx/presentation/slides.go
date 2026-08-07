package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const layoutsPerMaster = pptxxml.LayoutsPerMaster

type slideParts struct {
	title        pptxxml.TitleSpec
	contentStyle pptxxml.ContentStyleSpec
	// tables holds the slide's tables in order: SlideContent.Table first, then
	// the SlideContent.Tables overflow the reader fills for slides that carry
	// more than one.
	tables               []pptxxml.TableSpec
	imageRefs            []pptxxml.ImageRef
	backgroundRID        string
	transitionXML        string
	placeholders         []pptxxml.PlaceholderOverrideSpec
	chartFrames          []pptxxml.ChartFrame
	chartRels            []pptxxml.ChartRel
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

		// Code blocks are shapes with a highlighted text body, appended after
		// the caller's own shapes so they paint on top of a background shape.
		slideShapes := append(append([]shapes.Shape(nil), slide.Shapes...), codeBlockShapes(slide)...)

		slideXML := pptxxml.SlideWithLayout(
			elements.SlideLayoutXMLMode(slide.Layout),
			parts.title,
			slide.Bullets,
			elements.ToXMLBulletParagraphStyles(slide.BulletStyles),
			elements.ToXMLTextRunRows(slide.BulletRuns, hyperlinkRIDs),
			parts.contentStyle,
			parts.tables,
			parts.chartFrames,
			parts.imageRefs,
			shapes.ToXMLShapeSpecs(slideShapes, hyperlinkRIDs),
			shapes.ToXMLConnectorSpecs(slide.Connectors, slideShapes, hyperlinkRIDs),
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
		slideXMLBytes, err := editorslide.RewriteSlideHidden([]byte(slideXML), slide.Hidden)
		if err != nil {
			return err
		}
		slideXML = string(slideXMLBytes)

		layoutTarget := elements.SlideLayoutTarget(slide.Layout)
		if masterCount > 1 {
			masterNum := (i % masterCount) + 1
			layoutTarget = layoutTargetForMaster(layoutTarget, masterNum)
		}

		commentTarget := ""
		if len(commentsBySlide[i]) > 0 {
			pw.AddPart(fmt.Sprintf("ppt/comments/comment%d.xml", num), pptxxml.CommentsXML(commentsBySlide[i]))
			commentTarget = fmt.Sprintf("../comments/comment%d.xml", num)
		}

		relsXML := pptxxml.SlideRelationshipsWithAll(
			layoutTarget,
			builder.targets,
			parts.chartRels,

			parts.placeholderChartRels,
			parts.smartArtRels,
			notesTargets[num],
			hyperlinks,
			commentTarget,
		)

		if err := writeSlideWithAttachments(pw, slides, i, slideXML, relsXML); err != nil {
			return err
		}
	}
	return nil
}

// writeSlideWithAttachments adds the ink and media a slide carries, then writes
// the slide and its relationships. Both attachments rewrite the shape tree they
// are given, so they run in sequence over the same markup.
func writeSlideWithAttachments(
	pw *pptxxml.PackageWriter,
	slides []elements.SlideContent,
	index int,
	slideXML, relsXML string,
) error {
	slide := slides[index]
	num := index + 1

	inkParts := attachInk(
		slideXML,
		relsXML,
		slideInkAnnotations(slide),
		inkPartStartIndex(slides, index),
	)
	for path, content := range inkParts.Parts {
		pw.AddPart(path, content)
	}

	mediaParts, err := attachMedia(
		inkParts.SlideXML,
		inkParts.RelsXML,
		slideMedia(slide),
		mediaPartStartIndex(slides, index),
		index,
		len(slides),
	)
	if err != nil {
		return err
	}
	for path, content := range mediaParts.Parts {
		pw.AddBinaryPart(path, content)
	}

	pw.AddPart(fmt.Sprintf("ppt/slides/slide%d.xml", num), mediaParts.SlideXML)
	pw.AddPart(fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", num), mediaParts.RelsXML)
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
	newNum := (masterNum-1)*layoutsPerMaster + num
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
