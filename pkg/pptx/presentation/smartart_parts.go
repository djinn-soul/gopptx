package presentation

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type SmartArtPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.SmartArtSpec
	imageRels  []smartArtImageRel
}

// smartArtImageRel is one picture a diagram's data part points at. Node pictures
// are related from the data part rather than from the slide, which is why a
// diagram that uses them needs a rels file of its own.
type smartArtImageRel struct {
	relID     string
	mediaName string
}

const (
	smartArtRenderedFilesPerPart = 5 // data, layout, colors, quickStyle, drawing (no diagram rels file)
	smartArtOrderData            = 0
	smartArtOrderLayout          = 1
	smartArtOrderColors          = 2
	smartArtOrderQuickStyle      = 3
	smartArtOrderDrawing         = 4
)

func SmartArtPartCount(parts []SmartArtPart) int {
	return len(parts)
}

func BuildSmartArtParts(slides []elements.SlideContent, catalog *media.Catalog) []SmartArtPart {
	out := make([]SmartArtPart, 0)
	for i, slide := range slides {
		for _, sa := range slide.SmartArtDiagrams {
			spec := sa.ToSpec()
			rels := resolveSmartArtNodeImages(spec.Nodes, catalog)
			out = append(out, SmartArtPart{
				slideIndex: i,
				partNumber: len(out) + 1,
				spec:       spec,
				imageRels:  rels,
			})
		}
	}
	return out
}

// resolveSmartArtNodeImages registers each node picture with the media catalog
// and stamps the relationship it will be referenced by. Pictures that cannot be
// read are skipped, leaving the node's placeholder empty rather than failing the
// whole deck.
func resolveSmartArtNodeImages(nodes []pptxxml.SmartArtNodeSpec, catalog *media.Catalog) []smartArtImageRel {
	rels := make([]smartArtImageRel, 0)
	byMediaName := map[string]string{}

	var walk func(items []pptxxml.SmartArtNodeSpec)
	walk = func(items []pptxxml.SmartArtNodeSpec) {
		for i := range items {
			node := &items[i]
			if node.ImagePath != "" && catalog != nil {
				if mediaName, err := catalog.RegisterImage(shapes.Image{Path: node.ImagePath}); err == nil {
					relID, seen := byMediaName[mediaName]
					if !seen {
						relID = fmt.Sprintf("rId%d", len(rels)+1)
						byMediaName[mediaName] = relID
						rels = append(rels, smartArtImageRel{relID: relID, mediaName: mediaName})
					}
					node.ImageRelID = relID
				}
			}
			walk(node.Children)
		}
	}
	walk(nodes)
	return rels
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
	for _, part := range parts {
		if len(part.imageRels) == 0 {
			continue
		}
		pw.AddPart(
			fmt.Sprintf("ppt/diagrams/_rels/data%d.xml.rels", part.partNumber),
			smartArtDataRelsXML(part.imageRels),
		)
	}
	return nil
}

func smartArtDataRelsXML(rels []smartArtImageRel) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for _, rel := range rels {
		b.WriteString(`<Relationship Id="` + rel.relID +
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"` +
			` Target="../media/` + rel.mediaName + `"/>`)
	}
	b.WriteString(`</Relationships>`)
	return b.String()
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
	}
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
