package slide

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	initialSlideMasterID int64 = 2147483648
)

func RewritePresentationSlideList(current []byte, slides []common.EditorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	replacement := BuildPresentationSlideListXML(slides)
	result, found := replaceFirstXMLTagBlock(source, "p:sldIdLst", replacement)
	if !found {
		return "", errors.New("presentation XML does not contain <p:sldIdLst>")
	}
	return result, nil
}

func BuildPresentationSlideListXML(slides []common.EditorSlideRef) string {
	var b strings.Builder
	// Each entry: "\n<p:sldId id="NNN" r:id="rIdNNN"/>" ≈ 50 bytes.
	b.Grow(len("<p:sldIdLst></p:sldIdLst>") + len(slides)*50)
	b.WriteString("<p:sldIdLst>")
	for _, slide := range slides {
		b.WriteString("\n<p:sldId id=\"")
		b.WriteString(strconv.FormatInt(slide.SlideID, 10))
		b.WriteString("\" r:id=\"")
		b.WriteString(slide.RelID)
		b.WriteString("\"")
		b.WriteString("/>")
	}
	if len(slides) > 0 {
		b.WriteString("\n")
	}
	b.WriteString("</p:sldIdLst>")
	return b.String()
}

func RewritePresentationNotesMasterList(current []byte, relID string, enable bool) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	if !enable {
		if result, found := replaceAllXMLTagBlocks(source, "p:notesMasterIdLst", ""); found {
			return result, nil
		}
		return source, nil
	}
	if strings.TrimSpace(relID) == "" {
		return "", errors.New("notes master relationship id is required")
	}

	replacement := "<p:notesMasterIdLst>\n<p:notesMasterId r:id=\"" + common.XMLEscape(
		relID,
	) + "\"/>\n</p:notesMasterIdLst>"
	if result, found := replaceAllXMLTagBlocks(source, "p:notesMasterIdLst", replacement); found {
		return result, nil
	}

	if idx := strings.Index(source, "</p:sldMasterIdLst>"); idx >= 0 {
		insertPos := idx + len("</p:sldMasterIdLst>")
		return source[:insertPos] + "\n" + replacement + source[insertPos:], nil
	}
	if idx := strings.Index(source, "<p:sldIdLst"); idx >= 0 {
		return source[:idx] + replacement + "\n" + source[idx:], nil
	}
	return "", errors.New("presentation XML does not contain insertion point for notesMasterIdLst")
}

func RewritePresentationSlideMasterList(current []byte, relID string) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	if strings.TrimSpace(relID) == "" {
		return "", errors.New("slide master relationship id is required")
	}
	source := string(current)

	nextMasterID := scanNextIDAttribute(source, initialSlideMasterID)

	newEntry := fmt.Sprintf(`<p:sldMasterId id="%d" r:id="%s"/>`, nextMasterID, common.XMLEscape(relID))
	if result, found := replaceAllXMLTagBlocksFunc(source, "p:sldMasterIdLst", func(match string) string {
		return rewriteSlideMasterListBlock(match, newEntry)
	}); found {
		return result, nil
	}

	insertAfter := "</p:notesMasterIdLst>"
	if idx := strings.Index(source, insertAfter); idx >= 0 {
		pos := idx + len(insertAfter)
		replacement := "\n<p:sldMasterIdLst>\n" + newEntry + "\n</p:sldMasterIdLst>"
		return source[:pos] + replacement + source[pos:], nil
	}
	if idx := strings.Index(source, "<p:sldIdLst"); idx >= 0 {
		replacement := "<p:sldMasterIdLst>\n" + newEntry + "\n</p:sldMasterIdLst>\n"
		return source[:idx] + replacement + source[idx:], nil
	}
	return "", errors.New("presentation XML does not contain insertion point for sldMasterIdLst")
}

func scanNextIDAttribute(source string, initial int64) int64 {
	nextID := initial
	remaining := source

	for {
		idx := strings.Index(remaining, `id="`)
		if idx < 0 {
			return nextID
		}
		start := idx + len(`id="`)
		end := start
		for end < len(remaining) && remaining[end] >= '0' && remaining[end] <= '9' {
			end++
		}
		if end > start && end < len(remaining) && remaining[end] == '"' {
			if val, err := strconv.ParseInt(remaining[start:end], 10, 64); err == nil && val >= nextID {
				nextID = val + 1
			}
		}
		remaining = remaining[end:]
	}
}

func rewriteSlideMasterListBlock(match, newEntry string) string {
	const closeTag = "</p:sldMasterIdLst>"
	if strings.Contains(match, closeTag) && strings.Contains(match, "<p:sldMasterId") {
		return strings.Replace(match, closeTag, "\n"+newEntry+"\n"+closeTag, 1)
	}
	return "<p:sldMasterIdLst>\n" + newEntry + "\n</p:sldMasterIdLst>"
}
