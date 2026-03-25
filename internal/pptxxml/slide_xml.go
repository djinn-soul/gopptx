package pptxxml

import (
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
	slideWithLayoutGrowCap     = 3072
)

const startImageRID = 2

// PlaceholderOverrideSpec describes content override for a placeholder.
type PlaceholderOverrideSpec struct {
	Index int
	Type  string
	Text  string
	Image *ImageRef
	Table *TableSpec
	Chart *ChartFrame

	// Extension: Layout/Style Overrides
	X, Y, CX, CY      *int64
	TextStyle         *PlaceholderTextStyleSpec
	GeometryXML       string
	ForceRectGeometry *bool
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

// Precomputed header bodies for the two most common slide dimensions (zero allocs).
//
//nolint:gochecknoglobals // read-only precomputed constants, never mutated
var (
	slideHeaderEndBodyXML169  = slideHeaderEndBodyXMLFor("9144000", "6858000")
	slideHeaderEndBodyXML1610 = slideHeaderEndBodyXMLFor("12192000", "6858000")
)

func slideHeaderEndBodyXMLFor(wStr, hStr string) string {
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

func slideHeaderEndBodyXML(width, height int64) string {
	// Fast paths for the two most common slide dimensions — zero allocs.
	if width == 9144000 && height == 6858000 {
		return slideHeaderEndBodyXML169
	}
	if width == 12192000 && height == 6858000 {
		return slideHeaderEndBodyXML1610
	}
	wStr := strconv.FormatInt(width, 10)
	hStr := strconv.FormatInt(height, 10)
	return slideHeaderEndBodyXMLFor(wStr, hStr)
}

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
	// Pre-allocate ~3KB: covers header (~500) + title (~600) + 3-bullet content (~700) + footer (~100).
	// Avoids 2-3 reallocs for the common case; oversized slides will still realloc once.
	b.Grow(slideWithLayoutGrowCap)
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
