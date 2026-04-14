package tplx

import (
	"fmt"
	"strings"
)

// collectSlideParts returns all ppt/slides/slideN.xml paths in order.
func collectSlideParts(parts map[string][]byte) []string {
	var slides []string
	for name := range parts {
		if strings.HasPrefix(name, "ppt/slides/slide") && strings.HasSuffix(name, ".xml") &&
			!strings.Contains(name, "_rels") {
			slides = append(slides, name)
		}
	}
	sortStrings(slides)
	return slides
}

// expandSlideLoops expands {{#each KEY}} slide-level loops.
func expandSlideLoops(
	slideParts []string,
	parts map[string][]byte,
	ctx Context,
) ([]string, map[string][]byte) {
	var newSlides []string
	nextNum := maxSlideNumber(slideParts) + 1

	for _, name := range slideParts {
		data := parts[name]
		cond := detectEachSlide(data)
		if cond == nil {
			newSlides = append(newSlides, name)
			continue
		}
		templateRels, hasTemplateRels := parts[slideRelsPath(name)]

		rowsAny, ok := ctx[cond.key]
		if !ok {
			delete(parts, name)
			deleteSlideRels(parts, name)
			continue
		}
		rows, ok := toRows(rowsAny)
		if !ok {
			delete(parts, name)
			deleteSlideRels(parts, name)
			continue
		}

		expanded := expandSlide(data, rows, ctx)
		delete(parts, name)
		deleteSlideRels(parts, name)

		for _, slideXML := range expanded {
			newName := fmt.Sprintf("ppt/slides/slide%d.xml", nextNum)
			parts[newName] = slideXML
			newSlides = append(newSlides, newName)
			if hasTemplateRels {
				parts[slideRelsPath(newName)] = templateRels
			}
			nextNum++
		}
	}
	return newSlides, parts
}

func maxSlideNumber(slides []string) int {
	maxSlideNum := 0
	for _, s := range slides {
		n := parseSlideNumber(s)
		if n > maxSlideNum {
			maxSlideNum = n
		}
	}
	return maxSlideNum
}

func parseSlideNumber(name string) int {
	base := trimPrefix(name, "ppt/slides/slide")
	base = trimSuffix(base, ".xml")
	n := 0
	for _, c := range base {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*numericBaseTen + int(c-'0')
	}
	return n
}

func deleteSlideRels(parts map[string][]byte, slideName string) {
	delete(parts, slideRelsPath(slideName))
}

func slideRelsPath(slideName string) string {
	return "ppt/slides/_rels/" + lastSegment(slideName) + ".rels"
}

func lastSegment(s string) string {
	idx := strings.LastIndex(s, "/")
	if idx < 0 {
		return s
	}
	return s[idx+1:]
}

// isNonSlideTextPart returns true for XML parts that may contain template tokens.
func isNonSlideTextPart(name string) bool {
	if strings.Contains(name, "_rels") {
		return false
	}
	return strings.HasSuffix(name, ".xml") &&
		(strings.HasPrefix(name, "ppt/notesSlides/") ||
			strings.HasPrefix(name, "ppt/slideLayouts/") ||
			name == "ppt/presentation.xml")
}
