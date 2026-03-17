package presentation

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type SmartArtPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.SmartArtSpec
}

const (
	smartArtRenderedFilesPerPart = 5
	smartArtOrderData            = 0
	smartArtOrderLayout          = 1
	smartArtOrderColors          = 2
	smartArtOrderQuickStyle      = 3
	smartArtOrderDrawing         = 4
	smartArtOrderDataRels        = 5
)

func SmartArtPartCount(parts []SmartArtPart) int {
	return len(parts)
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
	renderedParts, err := renderSmartArtPartsParallel(parts)
	if err != nil {
		return err
	}
	for _, rendered := range renderedParts {
		pw.AddPart(rendered.path, rendered.content)
	}
	return nil
}

type smartArtRenderedPart struct {
	partNumber int
	order      int
	path       string
	content    string
}

func renderSmartArtPartsParallel(parts []SmartArtPart) ([]smartArtRenderedPart, error) {
	for _, part := range parts {
		if part.partNumber <= 0 {
			return nil, fmt.Errorf("invalid SmartArt part number: %d", part.partNumber)
		}
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make([]smartArtRenderedPart, 0, len(parts)*smartArtRenderedFilesPerPart)
	)

	for _, part := range parts {
		wg.Go(func() {
			rendered := renderSmartArtPart(part)
			mu.Lock()
			results = append(results, rendered...)
			mu.Unlock()
		})
	}
	wg.Wait()
	sort.Slice(results, func(i, j int) bool {
		if results[i].partNumber != results[j].partNumber {
			return results[i].partNumber < results[j].partNumber
		}
		return results[i].order < results[j].order
	})
	return results, nil
}

func renderSmartArtPart(part SmartArtPart) []smartArtRenderedPart {
	num := part.partNumber
	category := categoryFromURI(part.spec.LayoutURI)
	return []smartArtRenderedPart{
		{
			partNumber: num,
			order:      smartArtOrderData,
			path:       fmt.Sprintf("ppt/diagrams/data%d.xml", num),
			content:    pptxxml.SmartArtDataXML(part.spec),
		},
		{
			partNumber: num,
			order:      smartArtOrderLayout,
			path:       fmt.Sprintf("ppt/diagrams/layout%d.xml", num),
			content:    pptxxml.SmartArtLayoutXML(part.spec.LayoutURI, category),
		},
		{
			partNumber: num,
			order:      smartArtOrderColors,
			path:       fmt.Sprintf("ppt/diagrams/colors%d.xml", num),
			content:    pptxxml.SmartArtColorsXML(part.spec.ColorStyleID),
		},
		{
			partNumber: num,
			order:      smartArtOrderQuickStyle,
			path:       fmt.Sprintf("ppt/diagrams/quickStyle%d.xml", num),
			content:    pptxxml.SmartArtStyleXML(part.spec.QuickStyleID),
		},
		{
			partNumber: num,
			order:      smartArtOrderDrawing,
			path:       fmt.Sprintf("ppt/diagrams/drawing%d.xml", num),
			content:    pptxxml.SmartArtDrawingXML(part.spec),
		},
		{
			partNumber: num,
			order:      smartArtOrderDataRels,
			path:       fmt.Sprintf("ppt/diagrams/_rels/data%d.xml.rels", num),
			content:    smartArtDataRelsXML(num),
		},
	}
}

func smartArtDataRelsXML(partNumber int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId6" Type="http://schemas.microsoft.com/office/2007/relationships/diagramDrawing" Target="drawing%d.xml"/>
</Relationships>`, partNumber)
}

func categoryFromURI(uri string) string {
	// Simple heuristic based on URI keywords.
	// URI format: urn:microsoft.com/office/officeart/2005/8/layout/<name>
	switch {
	case strings.Contains(uri, "process"):
		return "process"
	case strings.Contains(uri, "cycle"):
		return "cycle"
	case strings.Contains(uri, "hierarchy"), strings.Contains(uri, "orgChart"):
		return "hierarchy"
	case strings.Contains(uri, "venn"), strings.Contains(uri, "radial"), strings.Contains(uri, "target"):
		return "relationship"
	case strings.Contains(uri, "matrix"):
		return "matrix"
	case strings.Contains(uri, "pyramid"):
		return "pyramid"
	case strings.Contains(uri, "picture"):
		return "picture"
	}
	return "list" // Default fallback
}
