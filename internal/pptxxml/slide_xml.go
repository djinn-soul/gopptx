package pptxxml

import (
	"path/filepath"
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

const (
	startImageRID = 2
)

// PlaceholderOverrideSpec describes content override for a placeholder.
type PlaceholderOverrideSpec struct {
	Index int
	Type  string
	Text  string
	Image *ImageRef
	Table *TableSpec
	Chart *ChartFrame

	// Extension: Layout/Style Overrides
	X, Y, CX, CY *int64
	TextStyle    *PlaceholderTextStyleSpec
}

// PlaceholderTextStyleSpec describes text formatting overrides for a placeholder.
type PlaceholderTextStyleSpec struct {
	SizePt    *int
	Color     *string
	Bold      *bool
	Italic    *bool
	Underline *string
	Align     *string
	Font      *string
}

const slideHeaderStart = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
	`xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" ` +
	`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
	`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
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
			xml += `
<a:solidFill>
<a:srgbClr val="` + strings.TrimPrefix(bg.SolidFill.Color, "#") + `"/>
</a:solidFill>`
		}
	case "gradient":
		if bg.GradientFill != nil {
			xml += shapeGradientFillXML(*bg.GradientFill)
		}
	case "picture":
		if bg.PictureFill != nil {
			xml += `
<a:blipFill>
<a:blip r:embed="` + FastEscapeRID(bg.PictureFill.RelID) + `"/>
<a:stretch>
<a:fillRect/>
</a:stretch>
</a:blipFill>`
		}
	}

	xml += `
<a:effectLst/>
</p:bgPr>
</p:bg>`
	return xml
}

func slideHeaderEndBodyXML(width, height int64) string {
	wStr := strconv.FormatInt(width, 10)
	hStr := strconv.FormatInt(height, 10)
	return `
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="` + wStr + `" cy="` + hStr + `"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="` + wStr + `" cy="` + hStr + `"/>
</a:xfrm>
</p:grpSpPr>`
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
	smartArtFrames []SmartArtFrame,
	background *SlideBackgroundSpec,

	transitionXML string,
	animationsXML string,
	showSlideNumber bool,
	footerText string,
	showDateTime bool,
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
		smartArtFrames,
		background,

		transitionXML,
		animationsXML,
		showSlideNumber,
		footerText,
		showDateTime,
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
	smartArtFrames []SmartArtFrame,
	background *SlideBackgroundSpec,

	transitionXML string,
	animationsXML string,
	showSlideNumber bool,
	footerText string,
	showDateTime bool,
	width, height int64,
) string {
	var b strings.Builder
	layoutMode := normalizeSlideLayoutMode(layout)

	b.WriteString(slideHeaderStart)
	b.WriteString(backgroundXML(background))
	b.WriteString(slideHeaderEndBodyXML(width, height))

	nextID := slideRenderBaseElements(
		&b, layoutMode, title, table, bullets, bulletStyles, bulletRuns, contentStyle, width, height,
	)

	if chart != nil {
		b.WriteString(chartFrameShape(chart, nextID))
		nextID++
	}

	nextID = slideRenderImages(&b, images, nextID)
	shapeIDs, nextID := slideRenderShapes(&b, shapes, nextID)
	nextID = slideRenderConnectors(&b, connectors, shapeIDs, nextID)

	for _, sa := range smartArtFrames {
		b.WriteString(smartArtFrameShape(&sa, nextID))
		nextID++
	}

	nextID = slideRenderPlaceholders(&b, placeholders, nextID)
	nextID = slideRenderOverlayShapes(
		&b,
		showSlideNumber,
		footerText,
		showDateTime,
		width,
		height,
		nextID,
	)

	b.WriteString(slideContentFooter)
	b.WriteString(slideFooterClrMap)
	slideRenderTimelineFeatures(&b, transitionXML, animationsXML)

	b.WriteString(slideFooterEnd)
	return b.String()
}

func slideRenderBaseElements(
	b *strings.Builder,
	layoutMode string,
	title TitleSpec,
	table *TableSpec,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	contentStyle ContentStyleSpec,
	width, height int64,
) int {
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
		b.WriteString(RenderTable(table, nextID))
		nextID++
	} else if len(bullets) > 0 {
		nextID = slideRenderBullets(
			b, layoutMode, bullets, bulletStyles,
			bulletRuns, contentStyle, nextID, width, height,
		)
	}
	return nextID
}

func slideRenderBullets(
	b *strings.Builder,
	layoutMode string,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	contentStyle ContentStyleSpec,
	nextID int,
	width, height int64,
) int {
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
			b.WriteString(rightTwoColumnShape(
				rightBullets, rightStyles, rightRuns, contentStyle, nextID, width, height,
			))
			nextID++
		}
	}
	return nextID
}

func slideRenderImages(b *strings.Builder, images []ImageRef, nextID int) int {
	for i, image := range images {
		b.WriteString(imageShape(image, nextID+i))
	}
	return nextID + len(images)
}

func slideRenderShapes(b *strings.Builder, shapes []ShapeSpec, nextID int) ([]int, int) {
	shapeIDs, shapeXMLParts := renderCustomShapeXMLConcurrently(shapes, nextID)
	for _, part := range shapeXMLParts {
		b.WriteString(part)
	}
	return shapeIDs, nextID + len(shapes)
}

func slideRenderConnectors(b *strings.Builder, connectors []ConnectorSpec, shapeIDs []int, nextID int) int {
	currentID := nextID
	for _, connector := range connectors {
		startShapeID := shapeAnchorID(shapeIDs, connector.StartShapeIndex)
		endShapeID := shapeAnchorID(shapeIDs, connector.EndShapeIndex)
		b.WriteString(connectorXML(connector, currentID, startShapeID, endShapeID))
		currentID++
		if labelXML := connectorLabelShape(connector, currentID); labelXML != "" {
			b.WriteString(labelXML)
			currentID++
		}
	}
	return currentID
}

func slideRenderPlaceholders(b *strings.Builder, placeholders []PlaceholderOverrideSpec, nextID int) int {
	for i, ph := range placeholders {
		b.WriteString(placeholderShape(ph, nextID+i))
	}
	return nextID + len(placeholders)
}

func slideRenderTimelineFeatures(
	b *strings.Builder,
	transitionXML, animationsXML string,
) {
	if tx := strings.TrimSpace(transitionXML); tx != "" {
		b.WriteString("\n")
		b.WriteString(tx)
	}
	if ax := strings.TrimSpace(animationsXML); ax != "" {
		b.WriteString("\n")
		b.WriteString(ax)
	}
}

func slideRenderOverlayShapes(
	b *strings.Builder,
	showSlideNumber bool,
	footerText string,
	showDateTime bool,
	width, height int64,
	nextID int,
) int {
	if showSlideNumber {
		b.WriteString(slideNumberShape(width, height, nextID))
		nextID++
	}
	if footerText != "" {
		b.WriteString(footerShape(footerText, width, height, nextID))
		nextID++
	}
	if showDateTime {
		b.WriteString(dateTimeShape(width, height, nextID))
		nextID++
	}
	return nextID
}

// ChartRel describes one chart relationship entry for slide relationships XML.
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
	Type     string
}

func SlideRelationshipsWithLayoutAndNotes(
	layoutTarget string,
	imageTargets []string,
	chartRel *ChartRel,
	notesTarget string,
) string {
	return SlideRelationshipsWithHyperlinks(layoutTarget, imageTargets, chartRel, notesTarget, nil)
}

// SlideRelationshipsWithHyperlinks extends slide relationships to include hyperlinks.
func SlideRelationshipsWithHyperlinks(
	layoutTarget string,
	imageTargets []string,
	chartRel *ChartRel,
	notesTarget string,
	hyperlinks []HyperlinkRel,
) string {
	return SlideRelationshipsWithMultiCharts(layoutTarget, imageTargets, chartRel, nil, nil, notesTarget, hyperlinks, "")
}

// SlideRelationshipsWithMultiCharts extends slide relationships to include multiple charts, SmartArt, and comments.
func SlideRelationshipsWithMultiCharts(
	layoutTarget string,
	imageTargets []string,
	chartRel *ChartRel,
	placeholderCharts []ChartRel,
	smartArtRels []SmartArtRel,
	notesTarget string,
	hyperlinks []HyperlinkRel,
	commentsTarget string,
) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" ` +
		`Target="` + Escape(layoutTarget) + `"/>`)
	maxRID := 1
	for i, target := range imageTargets {
		rid := i + startImageRID
		relType := slideMediaRelationshipType(target)
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(rid))
		b.WriteString("\" Type=\"")
		b.WriteString(relType)
		b.WriteString("\" Target=\"")
		b.WriteString(Escape(target))
		b.WriteString("\"/>")
		if rid > maxRID {
			maxRID = rid
		}
	}
	if chartRel != nil {
		b.WriteString("\n<Relationship Id=\"")
		b.WriteString(FastEscapeRID(chartRel.RID))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart\" Target=\"")
		b.WriteString(Escape(chartRel.Target))
		b.WriteString("\"/>")

		if rid := ridNumber(chartRel.RID); rid > maxRID {
			maxRID = rid
		}
	}
	for _, phChart := range placeholderCharts {
		b.WriteString("\n<Relationship Id=\"")
		b.WriteString(FastEscapeRID(phChart.RID))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart\" Target=\"")
		b.WriteString(Escape(phChart.Target))
		b.WriteString("\"/>")

		if rid := ridNumber(phChart.RID); rid > maxRID {
			maxRID = rid
		}
	}
	for _, saRel := range smartArtRels {
		b.WriteString(`
<Relationship Id="`)
		b.WriteString(FastEscapeRID(saRel.RID))
		b.WriteString(`" Type="`)
		b.WriteString(Escape(saRel.Type))
		b.WriteString(`" Target="`)
		b.WriteString(Escape(saRel.Target))
		b.WriteString(`"/>`)
		if rid := ridNumber(saRel.RID); rid > maxRID {
			maxRID = rid
		}
	}
	for _, hl := range hyperlinks {
		b.WriteString(HyperlinkRelationshipXML(hl.RID, hl.Target, hl.External, hl.Type))
		if rid := ridNumber(hl.RID); rid > maxRID {
			maxRID = rid
		}
	}
	if strings.TrimSpace(notesTarget) != "" {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(maxRID + 1))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide" Target="`)
		b.WriteString(Escape(notesTarget))
		b.WriteString(`"/>`)
		maxRID++
	}
	if strings.TrimSpace(commentsTarget) != "" {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(maxRID + 1))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments" Target="`)
		b.WriteString(Escape(commentsTarget))
		b.WriteString(`"/>`)
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

func slideMediaRelationshipType(target string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(target)))
	switch ext {
	case ".wav", ".mp3", ".m4a":
		return "http://schemas.openxmlformats.org/officeDocument/2006/relationships/audio"
	default:
		return "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	}
}
