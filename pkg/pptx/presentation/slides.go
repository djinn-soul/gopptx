package presentation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const layoutsPerMaster = 6

var slideShowAttrPattern = regexp.MustCompile(`\s+show\s*=\s*("[^"]*"|'[^']*')`)

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
		slideXML, err = rewriteSlideHiddenAttribute(slideXML, slide.Hidden)
		if err != nil {
			return err
		}

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

func rewriteSlideHiddenAttribute(slideXML string, hidden bool) (string, error) {
	start := strings.Index(slideXML, "<p:sld")
	if start < 0 {
		return "", errors.New("slide XML does not contain <p:sld> root")
	}
	nameEnd := start + len("<p:sld")
	if nameEnd >= len(slideXML) {
		return "", errors.New("slide XML has malformed <p:sld> root")
	}
	next := slideXML[nameEnd]
	if next != ' ' && next != '\n' && next != '\r' && next != '\t' && next != '>' {
		return "", errors.New("slide XML root element is not <p:sld>")
	}
	endRel := strings.IndexByte(slideXML[start:], '>')
	if endRel < 0 {
		return "", errors.New("slide XML has unterminated <p:sld> root")
	}
	end := start + endRel
	tag := slideXML[start : end+1]
	tag = slideShowAttrPattern.ReplaceAllString(tag, "")
	if hidden {
		if prefix, ok := strings.CutSuffix(tag, "/>"); ok {
			tag = prefix + ` show="0"/>`
		} else {
			prefix, _ := strings.CutSuffix(tag, ">")
			tag = prefix + ` show="0">`
		}
	}
	return slideXML[:start] + tag + slideXML[end+1:], nil
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
