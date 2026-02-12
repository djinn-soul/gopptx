package editor

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// EditorSection describes a PowerPoint section entry.
type EditorSection struct {
	Name     string
	GUID     string
	SlideIDs []int64
}

// PresentationEditor provides read/modify/save operations for existing PPTX files.
type PresentationEditor struct {
	parts map[string][]byte

	slides       []common.EditorSlideRef
	nextSlideID  int64
	nextRelIDNum int
	nextSlideNum int

	metadata        common.PresentationMetadata
	nonSlideRels    []common.EditorRelationship
	presentationXML string

	// Media inventory for deduplication (SHA1 -> PartPath)
	mediaInventory map[string]string
	nextMediaNum   int

	// Section management
	sections []EditorSection
}

// Metadata returns presentation-level metadata parsed from the package.
func (e *PresentationEditor) Metadata() common.PresentationMetadata {
	return e.metadata
}

// SlideCount returns the number of slides currently tracked by the editor.
func (e *PresentationEditor) SlideCount() int {
	if e == nil {
		return 0
	}
	return len(e.slides)
}

// Slides returns ordered slide metadata snapshots (0-based indexes).
func (e *PresentationEditor) Slides() []common.SlideMetadata {
	if e == nil || len(e.slides) == 0 {
		return nil
	}
	out := make([]common.SlideMetadata, 0, len(e.slides))
	for idx, slide := range e.slides {
		out = append(out, common.SlideMetadata{
			Index:          idx,
			SlideID:        slide.SlideID,
			RelationshipID: slide.RelID,
			PartName:       slide.Part,
			Title:          slide.Title,
		})
	}
	return out
}

func cloneParts(parts map[string][]byte) map[string][]byte {
	out := make(map[string][]byte, len(parts))
	for path, content := range parts {
		clone := make([]byte, len(content))
		copy(clone, content)
		out[path] = clone
	}
	return out
}

func requirePart(parts map[string][]byte, path string) ([]byte, error) {
	content, ok := parts[path]
	if !ok {
		return nil, fmt.Errorf("missing required package part %q", path)
	}
	return content, nil
}

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

func nextSlideID(slides []common.EditorSlideRef) int64 {
	var maxID int64 = 255
	for _, slide := range slides {
		if slide.SlideID > maxID {
			maxID = slide.SlideID
		}
	}
	return maxID + 1
}

func nextRelationshipNumber(rels []common.EditorRelationship) int {
	maxNum := 0
	for _, rel := range rels {
		num, ok := parseRelationshipNumber(rel.ID)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func nextSlidePartNumber(slides []common.EditorSlideRef) int {
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
