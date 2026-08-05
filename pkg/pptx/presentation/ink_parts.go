package presentation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/ink"
)

//nolint:gochecknoglobals // compiled once, immutable
var (
	relIDPattern      = regexp.MustCompile(`Id="rId(\d+)"`)
	shapeIDPattern    = regexp.MustCompile(`<p:cNvPr id="(\d+)"`)
	spTreeCloseMarker = "</p:spTree>"
)

// slideInkAnnotations returns the annotations on a slide that actually draw
// something, in slide order.
func slideInkAnnotations(slide elements.SlideContent) []*ink.Annotation {
	if len(slide.InkAnnotations) == 0 {
		return nil
	}
	out := make([]*ink.Annotation, 0, len(slide.InkAnnotations))
	for _, annotation := range slide.InkAnnotations {
		if annotation == nil || annotation.IsEmpty() {
			continue
		}
		out = append(out, annotation)
	}
	return out
}

// inkPartStartIndex returns the 1-based ink part number of the first
// annotation on slide slideIdx. Ink parts are numbered across the whole
// package in slide order, so the content-type pass and the slide-render pass
// agree on the part names.
func inkPartStartIndex(slides []elements.SlideContent, slideIdx int) int {
	index := 1
	for i := 0; i < slideIdx && i < len(slides); i++ {
		index += len(slideInkAnnotations(slides[i]))
	}
	return index
}

// totalInkParts counts the ink parts the package will contain.
func totalInkParts(slides []elements.SlideContent) int {
	total := 0
	for _, slide := range slides {
		total += len(slideInkAnnotations(slide))
	}
	return total
}

// addInkContentTypes registers one override per ink part.
func addInkContentTypes(contentTypes string, slides []elements.SlideContent) string {
	for i := 1; i <= totalInkParts(slides); i++ {
		contentTypes = pptxxml.WithContentTypeOverride(contentTypes, ink.PartPath(i), ink.PartContentType)
	}
	return contentTypes
}

// inkSlideParts is the result of attaching a slide's ink: the rewritten slide
// XML, the rewritten relationships, and the InkML parts to write.
type inkSlideParts struct {
	SlideXML string
	RelsXML  string
	Parts    map[string]string
}

// attachInk inserts the content-part markup into the slide shape tree, adds
// one relationship per annotation, and returns the InkML parts to write.
func attachInk(
	slideXML, relsXML string,
	annotations []*ink.Annotation,
	startPartIndex int,
) inkSlideParts {
	result := inkSlideParts{SlideXML: slideXML, RelsXML: relsXML, Parts: map[string]string{}}
	if len(annotations) == 0 {
		return result
	}

	nextRID := nextRelationshipID(relsXML)
	nextShapeID := nextShapeID(slideXML)

	var markup strings.Builder
	for i, annotation := range annotations {
		partIndex := startPartIndex + i
		relID := "rId" + strconv.Itoa(nextRID)
		nextRID++

		result.Parts[ink.PartPath(partIndex)] = annotation.InkML()
		result.RelsXML = pptxxml.WithRelationship(
			result.RelsXML,
			relID,
			ink.RelationshipType,
			ink.RelationshipTarget(partIndex),
		)
		markup.WriteString(annotation.ContentPartXML(relID, nextShapeID))
		nextShapeID++
	}

	result.SlideXML = insertBeforeShapeTreeClose(slideXML, markup.String())
	return result
}

// insertBeforeShapeTreeClose puts markup at the end of the shape tree, so ink
// paints above the shapes already on the slide.
func insertBeforeShapeTreeClose(slideXML, markup string) string {
	idx := strings.LastIndex(slideXML, spTreeCloseMarker)
	if idx < 0 {
		return slideXML
	}
	return slideXML[:idx] + markup + slideXML[idx:]
}

// nextRelationshipID returns one past the highest rId already in a rels part.
func nextRelationshipID(relsXML string) int {
	highest := 0
	for _, match := range relIDPattern.FindAllStringSubmatch(relsXML, -1) {
		if n, err := strconv.Atoi(match[1]); err == nil && n > highest {
			highest = n
		}
	}
	return highest + 1
}

// nextShapeID returns one past the highest shape id already in a slide, so ink
// content parts do not collide with existing shapes.
func nextShapeID(slideXML string) int {
	highest := 0
	for _, match := range shapeIDPattern.FindAllStringSubmatch(slideXML, -1) {
		if n, err := strconv.Atoi(match[1]); err == nil && n > highest {
			highest = n
		}
	}
	return highest + 1
}
