package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	slideLayoutTitleAndContent = "titleAndContent"
	slideLayoutTitleOnly       = "titleOnly"
	slideLayoutBlank           = "blank"
	slideLayoutCenteredTitle   = "centeredTitle"
	slideLayoutTitleBigContent = "titleAndBigContent"
	slideLayoutTwoColumn       = "twoColumn"
)

// PlaceholderOverrideSpec describes content override for a placeholder.
type PlaceholderOverrideSpec struct {
	Index int
	Type  string
	Text  string
	Image *ImageRef
	Table *TableSpec
	Chart *ChartFrame
}

const slideHeaderStart = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>`

const slideDefaultBackground = `
<p:bg>
<p:bgRef idx="1001">
<a:schemeClr val="bg1"/>
</p:bgRef>
</p:bg>`

func backgroundXML(bg *SlideBackgroundSpec) string {
	if bg == nil || bg.Type == "" {
		return ""
	}

	xml := `
<p:bg>
<p:bgPr>`

	switch bg.Type {
	case "solid":
		if bg.SolidFill != nil {
			xml += fmt.Sprintf(`
<a:solidFill>
<a:srgbClr val="%s"/>
</a:solidFill>`, strings.TrimPrefix(bg.SolidFill.Color, "#"))
		}
	case "gradient":
		if bg.GradientFill != nil {
			xml += shapeGradientFillXML(*bg.GradientFill)
		}
	case "picture":
		if bg.PictureFill != nil {
			xml += fmt.Sprintf(`
<a:blipFill>
<a:blip r:embed="%s"/>
<a:stretch>
<a:fillRect/>
</a:stretch>
</a:blipFill>`, Escape(bg.PictureFill.RelID))
		}
	}

	xml += `
<a:effectLst/>
</p:bgPr>
</p:bg>`
	return xml
}

func slideHeaderEndBodyXML(width, height int64) string {
	return fmt.Sprintf(`
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="%d" cy="%d"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="%d" cy="%d"/>
</a:xfrm>
</p:grpSpPr>`, width, height, width, height)
}

const slideContentFooter = `
</p:spTree>
</p:cSld>`

const slideFooterClrMap = `
<p:clrMapOvr>
<a:masterClrMapping/>
</p:clrMapOvr>`

const slideFooterEnd = `
</p:sld>`

// TitleSpec describes title formatting.
type TitleSpec struct {
	Text      string
	SizePt    int
	Color     string
	Bold      bool
	Italic    bool
	Underline bool
	Align     string
	Font      string
}

// ContentStyleSpec describes default content formatting.
type ContentStyleSpec struct {
	SizePt    int
	Color     string
	Bold      bool
	Italic    bool
	Underline bool
	VAlign    string
}

// SlideBackgroundSpec describes how the slide background should be filled.
type SlideBackgroundSpec struct {
	Type         string // "solid", "gradient", "picture"
	SolidFill    *ShapeFillSpec
	GradientFill *ShapeGradientFillSpec
	PictureFill  *ImageRef
}

// SlideWithContent renders a title+bullets slide with optional table, chart, and images.
func SlideWithContent(
	title TitleSpec,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	contentStyle ContentStyleSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
	background *SlideBackgroundSpec,
	transitionXML string,
	animationsXML string,
	showSlideNumber bool,
	width, height int64,
) string {
	return SlideWithLayout(
		slideLayoutTitleAndContent,
		title,
		bullets,
		bulletStyles,
		bulletRuns,
		contentStyle,
		table,
		chart,
		images,
		nil,
		nil,
		nil,
		background,
		transitionXML,
		animationsXML,
		showSlideNumber,
		width,
		height,
	)
}

// SlideWithLayout renders a slide using an explicit layout mode.
func SlideWithLayout(
	layout string,
	title TitleSpec,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	contentStyle ContentStyleSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
	shapes []ShapeSpec,
	connectors []ConnectorSpec,
	placeholders []PlaceholderOverrideSpec,
	background *SlideBackgroundSpec,
	transitionXML string,
	animationsXML string,
	showSlideNumber bool,
	width, height int64,
) string {
	var b strings.Builder
	layoutMode := normalizeSlideLayoutMode(layout)
	b.WriteString(slideHeaderStart)
	b.WriteString(backgroundXML(background))
	b.WriteString(slideHeaderEndBodyXML(width, height))

	nextID := 2
	if layoutMode != slideLayoutBlank {
		if layoutMode == slideLayoutCenteredTitle {
			b.WriteString(centeredTitleShape(title, width, height))
		} else {
			b.WriteString(titleShape(title, width, height))
		}
		nextID = 3
	}

	if table != nil {
		b.WriteString(tableShape(table, nextID))
		nextID++
	} else if len(bullets) > 0 {
		switch layoutMode {
		case slideLayoutTitleAndContent:
			b.WriteString(contentShape(bullets, bulletStyles, bulletRuns, contentStyle, nextID, width, height))
			nextID++
		case slideLayoutTitleBigContent:
			b.WriteString(bigContentShape(bullets, bulletStyles, bulletRuns, contentStyle, nextID, width, height))
			nextID++
		case slideLayoutTwoColumn:
			leftBullets, rightBullets := splitBulletsForTwoColumns(bullets)
			leftStyles, rightStyles := splitBulletStylesForTwoColumns(bulletStyles, len(leftBullets))
			leftRuns, rightRuns := splitBulletRunsForTwoColumns(bulletRuns, len(leftBullets))
			b.WriteString(leftTwoColumnShape(leftBullets, leftStyles, leftRuns, contentStyle, nextID, width, height))
			nextID++
			if len(rightBullets) > 0 {
				b.WriteString(rightTwoColumnShape(rightBullets, rightStyles, rightRuns, contentStyle, nextID, width, height))
				nextID++
			}
		}
	}

	if chart != nil {
		b.WriteString(chartFrameShape(chart, nextID))
		nextID++
	}

	for i, image := range images {
		b.WriteString(imageShape(image, nextID+i))
	}
	nextID += len(images)

	shapeIDs, shapeXMLParts := renderCustomShapeXMLConcurrently(shapes, nextID)
	for _, part := range shapeXMLParts {
		b.WriteString(part)
	}
	nextID += len(shapes)

	for i, connector := range connectors {
		startShapeID := shapeAnchorID(shapeIDs, connector.StartShapeIndex)
		endShapeID := shapeAnchorID(shapeIDs, connector.EndShapeIndex)
		b.WriteString(connectorXML(connector, nextID+i, startShapeID, endShapeID))
	}
	nextID += len(connectors)

	b.WriteString(slideContentFooter)
	b.WriteString(slideFooterClrMap)

	// Placeholders
	for i, ph := range placeholders {
		b.WriteString(placeholderShape(ph, nextID+i))
	}

	if tx := strings.TrimSpace(transitionXML); tx != "" {
		b.WriteString("\n")
		b.WriteString(tx)
	}
	if ax := strings.TrimSpace(animationsXML); ax != "" {
		b.WriteString("\n")
		b.WriteString(ax)
	}

	if showSlideNumber {
		b.WriteString(slideNumberShape(width, height, nextID))
	}

	b.WriteString(slideFooterEnd)
	return b.String()
}

// SlideRelationships renders ppt/slides/_rels/slideN.xml.rels.
type ChartRel struct {
	RID    string
	Target string
}

func SlideRelationships(imageTargets []string, chartRel *ChartRel) string {
	return SlideRelationshipsWithLayout("../slideLayouts/slideLayout1.xml", imageTargets, chartRel)
}

func SlideRelationshipsWithLayout(layoutTarget string, imageTargets []string, chartRel *ChartRel) string {
	return SlideRelationshipsWithLayoutAndNotes(layoutTarget, imageTargets, chartRel, "")
}

// HyperlinkRel describes a hyperlink relationship for slide rels.
type HyperlinkRel struct {
	RID      string
	Target   string
	External bool
}

func SlideRelationshipsWithLayoutAndNotes(layoutTarget string, imageTargets []string, chartRel *ChartRel, notesTarget string) string {
	return SlideRelationshipsWithHyperlinks(layoutTarget, imageTargets, chartRel, notesTarget, nil)
}

// SlideRelationshipsWithHyperlinks extends slide relationships to include hyperlinks.
func SlideRelationshipsWithHyperlinks(layoutTarget string, imageTargets []string, chartRel *ChartRel, notesTarget string, hyperlinks []HyperlinkRel) string {
	return SlideRelationshipsWithMultiCharts(layoutTarget, imageTargets, chartRel, nil, notesTarget, hyperlinks)
}

// SlideRelationshipsWithMultiCharts extends slide relationships to include multiple charts.
func SlideRelationshipsWithMultiCharts(
	layoutTarget string,
	imageTargets []string,
	chartRel *ChartRel,
	placeholderCharts []ChartRel,
	notesTarget string,
	hyperlinks []HyperlinkRel,
) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="` + Escape(layoutTarget) + `"/>`)
	maxRID := 1
	for i, target := range imageTargets {
		rid := i + 2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`, rid, Escape(target)))
		if rid > maxRID {
			maxRID = rid
		}
	}
	if chartRel != nil {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" Target="%s"/>`,
			Escape(chartRel.RID),
			Escape(chartRel.Target),
		))
		if rid := ridNumber(chartRel.RID); rid > maxRID {
			maxRID = rid
		}
	}
	for _, phChart := range placeholderCharts {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart" Target="%s"/>`,
			Escape(phChart.RID),
			Escape(phChart.Target),
		))
		if rid := ridNumber(phChart.RID); rid > maxRID {
			maxRID = rid
		}
	}
	for _, hl := range hyperlinks {
		b.WriteString(HyperlinkRelationshipXML(hl.RID, hl.Target, hl.External))
		if rid := ridNumber(hl.RID); rid > maxRID {
			maxRID = rid
		}
	}
	if strings.TrimSpace(notesTarget) != "" {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" Target="%s"/>`,
			maxRID+1,
			Escape(notesTarget),
		))
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}

func normalizeSlideLayoutMode(layout string) string {
	switch strings.ToLower(strings.TrimSpace(layout)) {
	case slideLayoutTitleOnly, "title_only", "title-only", "titleonly":
		return slideLayoutTitleOnly
	case slideLayoutBlank:
		return slideLayoutBlank
	case slideLayoutCenteredTitle, "centered_title", "centered-title", "centeredtitle":
		return slideLayoutCenteredTitle
	case slideLayoutTitleBigContent, "title_and_big_content", "title-and-big-content", "titleandbigcontent":
		return slideLayoutTitleBigContent
	case slideLayoutTwoColumn, "two_column", "two-column", "twocolumn":
		return slideLayoutTwoColumn
	default:
		return slideLayoutTitleAndContent
	}
}

func shapeAnchorID(shapeIDs []int, shapeIndex int) int {
	if shapeIndex <= 0 || shapeIndex > len(shapeIDs) {
		return 0
	}
	return shapeIDs[shapeIndex-1]
}

func ridNumber(rid string) int {
	if !strings.HasPrefix(rid, "rId") {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimPrefix(rid, "rId"))
	if err != nil {
		return 0
	}
	return n
}
