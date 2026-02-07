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

const slideHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:cSld>
<p:bg>
<p:bgRef idx="1001">
<a:schemeClr val="bg1"/>
</p:bgRef>
</p:bg>
<p:spTree>
<p:nvGrpSpPr>
<p:cNvPr id="1" name=""/>
<p:cNvGrpSpPr/>
<p:nvPr/>
</p:nvGrpSpPr>
<p:grpSpPr>
<a:xfrm>
<a:off x="0" y="0"/>
<a:ext cx="9144000" cy="6858000"/>
<a:chOff x="0" y="0"/>
<a:chExt cx="9144000" cy="6858000"/>
</a:xfrm>
</p:grpSpPr>`

const slideFooter = `
</p:spTree>
</p:cSld>
<p:clrMapOvr>
<a:masterClrMapping/>
</p:clrMapOvr>
</p:sld>`

// SlideWithContent renders a title+bullets slide with optional table, chart, and images.
func SlideWithContent(
	title string,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
) string {
	return SlideWithLayout(
		slideLayoutTitleAndContent,
		title,
		bullets,
		bulletStyles,
		bulletRuns,
		table,
		chart,
		images,
		nil,
		nil,
	)
}

// SlideWithLayout renders a slide using an explicit layout mode.
func SlideWithLayout(
	layout string,
	title string,
	bullets []string,
	bulletStyles []BulletParagraphSpec,
	bulletRuns [][]TextRunSpec,
	table *TableSpec,
	chart *ChartFrame,
	images []ImageRef,
	shapes []ShapeSpec,
	connectors []ConnectorSpec,
) string {
	var b strings.Builder
	layoutMode := normalizeSlideLayoutMode(layout)
	b.WriteString(slideHeader)

	nextID := 2
	if layoutMode != slideLayoutBlank {
		if layoutMode == slideLayoutCenteredTitle {
			b.WriteString(centeredTitleShape(title))
		} else {
			b.WriteString(titleShape(title))
		}
		nextID = 3
	}

	if table != nil {
		b.WriteString(tableShape(table, nextID))
		nextID++
	} else if len(bullets) > 0 {
		switch layoutMode {
		case slideLayoutTitleAndContent:
			b.WriteString(contentShape(bullets, bulletStyles, bulletRuns, nextID))
			nextID++
		case slideLayoutTitleBigContent:
			b.WriteString(bigContentShape(bullets, bulletStyles, bulletRuns, nextID))
			nextID++
		case slideLayoutTwoColumn:
			leftBullets, rightBullets := splitBulletsForTwoColumns(bullets)
			leftStyles, rightStyles := splitBulletStylesForTwoColumns(bulletStyles, len(leftBullets))
			leftRuns, rightRuns := splitBulletRunsForTwoColumns(bulletRuns, len(leftBullets))
			b.WriteString(leftTwoColumnShape(leftBullets, leftStyles, leftRuns, nextID))
			nextID++
			if len(rightBullets) > 0 {
				b.WriteString(rightTwoColumnShape(rightBullets, rightStyles, rightRuns, nextID))
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
	b.WriteString(slideFooter)
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

func SlideRelationshipsWithLayoutAndNotes(layoutTarget string, imageTargets []string, chartRel *ChartRel, notesTarget string) string {
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
