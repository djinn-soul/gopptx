package pptx

import (
	"bytes"
	"encoding/xml"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var coreTitlePattern = regexp.MustCompile(`(?s)<dc:title[^>]*>(.*?)</dc:title>`)

func extractCoreTitle(content []byte) string {
	if len(content) == 0 {
		return ""
	}
	matches := coreTitlePattern.FindSubmatch(content)
	if len(matches) != 2 {
		return ""
	}
	return strings.TrimSpace(htmlUnescapeXML(string(matches[1])))
}

func htmlUnescapeXML(value string) string {
	replacer := strings.NewReplacer(
		"&lt;", "<",
		"&gt;", ">",
		"&amp;", "&",
		"&quot;", `"`,
		"&apos;", "'",
	)
	return replacer.Replace(value)
}

func nextSlideID(slides []editorSlideRef) int64 {
	var maxID int64 = 255
	for _, slide := range slides {
		if slide.SlideID > maxID {
			maxID = slide.SlideID
		}
	}
	return maxID + 1
}

func nextRelationshipNumber(rels []editorRelationship) int {
	maxNum := 0
	for _, rel := range rels {
		num, ok := parseRelationshipNumber(rel.ID)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func nextSlidePartNumber(slides []editorSlideRef) int {
	maxNum := 0
	for _, slide := range slides {
		num, ok := parseSlidePartNumber(slide.Part)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func parseRelationshipNumber(id string) (int, bool) {
	trimmed := strings.TrimSpace(id)
	if !strings.HasPrefix(trimmed, "rId") {
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimPrefix(trimmed, "rId"))
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

func parseSlidePartNumber(partPath string) (int, bool) {
	base := path.Base(strings.TrimSpace(partPath))
	if !strings.HasPrefix(base, "slide") || !strings.HasSuffix(base, ".xml") {
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(base, "slide"), ".xml"))
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

func (e *PresentationEditor) populateSlideTitlesConcurrently() {
	if e == nil || len(e.slides) == 0 {
		return
	}

	type result struct {
		index int
		title string
	}
	ch := make(chan result, len(e.slides))
	var wg sync.WaitGroup

	for idx := range e.slides {
		idx := idx
		wg.Add(1)
		go func() {
			defer wg.Done()
			title := extractFirstAText(e.parts[e.slides[idx].Part])
			ch <- result{index: idx, title: title}
		}()
	}
	wg.Wait()
	close(ch)

	results := make([]result, 0, len(e.slides))
	for item := range ch {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].index < results[j].index })
	for _, item := range results {
		e.slides[item.index].Title = item.title
	}
}

func extractFirstAText(content []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return ""
			}
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if err := decoder.DecodeElement(&value, &start); err != nil {
			return ""
		}
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
}
