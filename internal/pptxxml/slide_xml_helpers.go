package pptxxml

import (
	"path/filepath"
	"strconv"
	"strings"
)

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
