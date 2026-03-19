package slide

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	sldIDLstPattern          = regexp.MustCompile(`(?s)<p:sldIdLst>.*?</p:sldIdLst>|<p:sldIdLst\s*/>`)
	notesMasterIDListPattern = regexp.MustCompile(
		`(?s)<p:notesMasterIdLst>.*?</p:notesMasterIdLst>|<p:notesMasterIdLst\s*/>`,
	)
	slideMasterIDListPattern = regexp.MustCompile(`(?s)<p:sldMasterIdLst>.*?</p:sldMasterIdLst>|<p:sldMasterIdLst\s*/>`)
	slideMasterIDPattern     = regexp.MustCompile(`id="(\d+)"`)
)

const (
	initialSlideMasterID int64 = 2147483648
	patternSubmatchLen   int   = 2
)

func RewritePresentationSlideList(current []byte, slides []common.EditorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	replacement := BuildPresentationSlideListXML(slides)
	if !sldIDLstPattern.MatchString(source) {
		return "", errors.New("presentation XML does not contain <p:sldIdLst>")
	}

	found := false
	result := sldIDLstPattern.ReplaceAllStringFunc(source, func(match string) string {
		if found {
			return match
		}
		found = true
		return replacement
	})
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
		b.WriteString("\"/>")
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
		if notesMasterIDListPattern.MatchString(source) {
			return notesMasterIDListPattern.ReplaceAllString(source, ""), nil
		}
		return source, nil
	}
	if strings.TrimSpace(relID) == "" {
		return "", errors.New("notes master relationship id is required")
	}

	replacement := "<p:notesMasterIdLst>\n<p:notesMasterId r:id=\"" + common.XMLEscape(
		relID,
	) + "\"/>\n</p:notesMasterIdLst>"
	if notesMasterIDListPattern.MatchString(source) {
		return notesMasterIDListPattern.ReplaceAllString(source, replacement), nil
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

	matches := slideMasterIDPattern.FindAllStringSubmatch(source, -1)
	nextMasterID := initialSlideMasterID
	for _, m := range matches {
		if len(m) != patternSubmatchLen {
			continue
		}
		val, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			continue
		}
		if val >= nextMasterID {
			nextMasterID = val + 1
		}
	}

	newEntry := fmt.Sprintf(`<p:sldMasterId id="%d" r:id="%s"/>`, nextMasterID, common.XMLEscape(relID))
	if slideMasterIDListPattern.MatchString(source) {
		return slideMasterIDListPattern.ReplaceAllStringFunc(source, func(match string) string {
			if strings.Contains(match, "</p:sldMasterIdLst>") {
				if strings.Contains(match, "<p:sldMasterId") {
					return strings.Replace(match, "</p:sldMasterIdLst>", "\n"+newEntry+"\n</p:sldMasterIdLst>", 1)
				}
				return "<p:sldMasterIdLst>\n" + newEntry + "\n</p:sldMasterIdLst>"
			}
			return "<p:sldMasterIdLst>\n" + newEntry + "\n</p:sldMasterIdLst>"
		}), nil
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
