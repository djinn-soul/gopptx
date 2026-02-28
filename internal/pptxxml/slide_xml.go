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
	slideBackgroundPicture     = "picture"
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
	_ = slideRenderOverlayShapes(
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

	if len(bullets) > 0 {
		nextID = slideRenderBullets(
			b, layoutMode, bullets, bulletStyles,
			bulletRuns, contentStyle, nextID, width, height,
		)
	}
	if table != nil {
		b.WriteString(RenderTable(table, nextID))
		nextID++
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
