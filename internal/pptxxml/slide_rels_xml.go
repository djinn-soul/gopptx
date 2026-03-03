package pptxxml

import (
	"strconv"
	"strings"
)

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
	return SlideRelationshipsWithMultiCharts(
		layoutTarget,
		imageTargets,
		chartRel,
		nil,
		nil,
		notesTarget,
		hyperlinks,
		"",
	)
}

// SlideRelationshipsWithMultiCharts extends slide relationships to include multiple charts, SmartArt, 3D models, and comments.
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
	return SlideRelationshipsWithAll(
		layoutTarget,
		imageTargets,
		chartRel,
		placeholderCharts,
		smartArtRels,
		notesTarget,
		hyperlinks,
		commentsTarget,
	)
}

// SlideRelationshipsWithAll is the comprehensive relationship generator.
func SlideRelationshipsWithAll(
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
	writeSlideRelsHeader(&b, layoutTarget)
	maxRID := appendImageRelationships(&b, imageTargets)
	maxRID = appendPrimaryChartRelationship(&b, chartRel, maxRID)
	maxRID = appendPlaceholderChartRelationships(&b, placeholderCharts, maxRID)
	maxRID = appendSmartArtRelationships(&b, smartArtRels, maxRID)
	maxRID = appendHyperlinkRelationships(&b, hyperlinks, maxRID)
	maxRID = appendNotesRelationship(&b, notesTarget, maxRID)
	appendCommentsRelationship(&b, commentsTarget, maxRID)
	b.WriteString("\n</Relationships>")
	return b.String()
}

func writeSlideRelsHeader(b *strings.Builder, layoutTarget string) {
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" ` +
		`Target="` + Escape(layoutTarget) + `"/>`)
}

func appendImageRelationships(b *strings.Builder, imageTargets []string) int {
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
	return maxRID
}

func appendPrimaryChartRelationship(b *strings.Builder, chartRel *ChartRel, maxRID int) int {
	if chartRel == nil {
		return maxRID
	}
	b.WriteString("\n<Relationship Id=\"")
	b.WriteString(FastEscapeRID(chartRel.RID))
	b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart\" Target=\"")
	b.WriteString(Escape(chartRel.Target))
	b.WriteString("\"/>")
	return max(maxRID, ridNumber(chartRel.RID))
}

func appendPlaceholderChartRelationships(b *strings.Builder, placeholderCharts []ChartRel, maxRID int) int {
	for _, phChart := range placeholderCharts {
		b.WriteString("\n<Relationship Id=\"")
		b.WriteString(FastEscapeRID(phChart.RID))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart\" Target=\"")
		b.WriteString(Escape(phChart.Target))
		b.WriteString("\"/>")
		maxRID = max(maxRID, ridNumber(phChart.RID))
	}
	return maxRID
}

func appendSmartArtRelationships(b *strings.Builder, smartArtRels []SmartArtRel, maxRID int) int {
	for _, saRel := range smartArtRels {
		b.WriteString("\n<Relationship Id=\"")
		b.WriteString(FastEscapeRID(saRel.RID))
		b.WriteString("\" Type=\"")
		b.WriteString(Escape(saRel.Type))
		b.WriteString("\" Target=\"")
		b.WriteString(Escape(saRel.Target))
		b.WriteString("\"/>")
		maxRID = max(maxRID, ridNumber(saRel.RID))
	}
	return maxRID
}

func appendHyperlinkRelationships(b *strings.Builder, hyperlinks []HyperlinkRel, maxRID int) int {
	for _, hl := range hyperlinks {
		b.WriteString(HyperlinkRelationshipXML(hl.RID, hl.Target, hl.External, hl.Type))
		maxRID = max(maxRID, ridNumber(hl.RID))
	}
	return maxRID
}

func appendNotesRelationship(b *strings.Builder, notesTarget string, maxRID int) int {
	if strings.TrimSpace(notesTarget) == "" {
		return maxRID
	}
	b.WriteString("\n<Relationship Id=\"rId")
	b.WriteString(strconv.Itoa(maxRID + 1))
	b.WriteString(
		"\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide\" Target=\"",
	)
	b.WriteString(Escape(notesTarget))
	b.WriteString("\"/>")
	return maxRID + 1
}

func appendCommentsRelationship(b *strings.Builder, commentsTarget string, maxRID int) {
	if strings.TrimSpace(commentsTarget) == "" {
		return
	}
	b.WriteString("\n<Relationship Id=\"rId")
	b.WriteString(strconv.Itoa(maxRID + 1))
	b.WriteString(
		"\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments\" Target=\"",
	)
	b.WriteString(Escape(commentsTarget))
	b.WriteString("\"/>")
}
