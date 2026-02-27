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
//
//nolint:funlen // Slide relationship XML must enumerate media/chart/notes/smartart links with stable ordering.
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
//
//nolint:funlen,gocognit // Large switch for XML generation is standard in this project.
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
		b.WriteString("\n<Relationship Id=\"")
		b.WriteString(FastEscapeRID(saRel.RID))
		b.WriteString("\" Type=\"")
		b.WriteString(Escape(saRel.Type))
		b.WriteString("\" Target=\"")
		b.WriteString(Escape(saRel.Target))
		b.WriteString("\"/>")
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
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(maxRID + 1))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide\" Target=\"")
		b.WriteString(Escape(notesTarget))
		b.WriteString("\"/>")
		maxRID++
	}
	if strings.TrimSpace(commentsTarget) != "" {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(maxRID + 1))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments\" Target=\"")
		b.WriteString(Escape(commentsTarget))
		b.WriteString("\"/>")
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}
