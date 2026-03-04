package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
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
	initialSlideMasterID   int64 = 2147483648
	patternSubmatchPairLen int   = 2
)

func rewritePresentationSlideList(current []byte, slides []common.EditorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	replacement := buildPresentationSlideListXML(slides)
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

func buildPresentationSlideListXML(slides []common.EditorSlideRef) string {
	var b strings.Builder
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

func rewritePresentationNotesMasterList(current []byte, relID string, enable bool) (string, error) {
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

func rewritePresentationSlideMasterList(current []byte, relID string) (string, error) {
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
		if len(m) != patternSubmatchPairLen {
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

func renderPresentationRelsXML(
	nonSlide []common.EditorRelationship,
	slides []common.EditorSlideRef,
	hasSections bool,
	hasVBA bool,
) (string, error) {
	return editorslide.RenderPresentationRelsXML(nonSlide, slides, hasSections, hasVBA)
}

func renderRelationshipsXML(rels []common.EditorRelationship) string {
	return editorslide.RenderRelationshipsXML(rels)
}

func rewriteContentTypes(
	current []byte,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	return editorslide.RewriteContentTypes(
		current,
		slides,
		mediaPaths,
		hasSections,
		chartPaths,
		notesPaths,
		themePaths,
		layoutPaths,
		masterPaths,
		hasNotesMaster,
		hasCommentAuthors,
		commentPaths,
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	)
}

func buildSectionListXML(sections []Section) string {
	moduleSections := toModuleSections(sections)
	return editorslide.BuildSectionListXML(moduleSections)
}

func toModuleSections(sections []Section) []editorslide.SectionData {
	moduleSections := make([]editorslide.SectionData, 0, len(sections))
	for _, s := range sections {
		moduleSections = append(moduleSections, editorslide.SectionData{
			Name:     s.Name,
			GUID:     s.GUID,
			SlideIDs: s.SlideIDs,
		})
	}
	return moduleSections
}

func rewriteChartExternalData(current []byte, newRelID string) []byte {
	return editorslide.RewriteChartExternalData(current, newRelID)
}

func rewritePresentationSections(current []byte, sections []Section) (string, error) {
	return editorslide.RewritePresentationSections(current, toModuleSections(sections))
}

func rewritePresentationEmbeddedFonts(current []byte, fontLst string) (string, error) {
	return editorslide.RewritePresentationEmbeddedFonts(current, fontLst)
}
