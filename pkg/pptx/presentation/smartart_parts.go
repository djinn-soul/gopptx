package presentation

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type SmartArtPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.SmartArtSpec
}

func BuildSmartArtParts(slides []elements.SlideContent) []SmartArtPart {
	out := make([]SmartArtPart, 0)
	for i, slide := range slides {
		for _, sa := range slide.SmartArtDiagrams {
			out = append(out, SmartArtPart{
				slideIndex: i,
				partNumber: len(out) + 1,
				spec:       sa.ToSpec(),
			})
		}
	}
	return out
}

func smartArtPartBySlide(parts []SmartArtPart) map[int][]SmartArtPart {
	bySlide := make(map[int][]SmartArtPart, len(parts))
	for _, part := range parts {
		bySlide[part.slideIndex] = append(bySlide[part.slideIndex], part)
	}
	return bySlide
}

func writeSmartArtFiles(pw *pptxxml.PackageWriter, parts []SmartArtPart) error {
	for _, part := range parts {
		// Write 5 parts for each diagram
		num := part.partNumber

		// Data
		pw.AddPart(fmt.Sprintf("ppt/diagrams/data%d.xml", num),
			pptxxml.SmartArtDataXML(part.spec))

		// Layout (stub)
		category := categoryFromURI(part.spec.LayoutURI)
		pw.AddPart(fmt.Sprintf("ppt/diagrams/layout%d.xml", num),
			pptxxml.SmartArtLayoutXML(part.spec.LayoutURI, category))

		// Colors (stub)
		pw.AddPart(fmt.Sprintf("ppt/diagrams/colors%d.xml", num),
			pptxxml.SmartArtColorsXML(part.spec.ColorStyleID))

		// QuickStyle (stub)
		pw.AddPart(fmt.Sprintf("ppt/diagrams/quickStyle%d.xml", num),
			pptxxml.SmartArtStyleXML(part.spec.QuickStyleID))
	}
	return nil
}

func categoryFromURI(uri string) string {
	// Simple heuristic based on URI keywords.
	// URI format: urn:microsoft.com/office/officeart/2005/8/layout/<name>
	if strings.Contains(uri, "process") {
		return "process"
	} else if strings.Contains(uri, "cycle") {
		return "cycle"
	} else if strings.Contains(uri, "hierarchy") || strings.Contains(uri, "orgChart") {
		return "hierarchy"
	} else if strings.Contains(uri, "venn") || strings.Contains(uri, "radial") || strings.Contains(uri, "target") {
		return "relationship"
	} else if strings.Contains(uri, "matrix") {
		return "matrix"
	} else if strings.Contains(uri, "pyramid") {
		return "pyramid"
	} else if strings.Contains(uri, "picture") {
		return "picture"
	}
	return "list" // Default fallback
}
