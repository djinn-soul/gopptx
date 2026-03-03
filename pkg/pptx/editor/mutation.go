package editor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"regexp"
	"sort"
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
	rels, used, hasSectionRel, err := collectNonSlideRelationships(nonSlide, len(slides))
	if err != nil {
		return "", err
	}
	if hasSections && !hasSectionRel {
		rels = append(rels, makeSectionListRelationship(rels, slides))
	}
	if hasVBA {
		hasVBARel := false
		for _, r := range rels {
			if r.Type == vbaProjectRelType {
				hasVBARel = true
				break
			}
		}
		if !hasVBARel {
			rels = append(rels, makeVBARelationship(rels, slides))
		}
	}
	rels = appendMissingSlideRelationships(rels, used, slides)
	return renderRelationshipsXML(rels), nil
}

func renderRelationshipsXML(rels []common.EditorRelationship) string {
	sortRelationshipsByID(rels)
	return relationshipsXMLDocument(rels)
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
	doc, err := parseContentTypesDocument(current)
	if err != nil {
		return "", err
	}
	ensureContentTypeDefaults(&doc, mediaPaths, hasVBA)

	overrides := filterDynamicOverrides(doc.Overrides, len(slides))
	overrides = appendSlideOverrides(overrides, slides)
	overrides = appendOptionalContentTypeOverride(overrides, hasSections, "/ppt/sectionList.xml",
		"application/vnd.microsoft.powerpoint.sectionList+xml")
	overrides = appendPathOverrides(overrides, chartPaths,
		"application/vnd.openxmlformats-officedocument.drawingml.chart+xml")
	overrides = appendPathOverrides(overrides, notesPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml")
	overrides = appendPathOverrides(overrides, themePaths,
		"application/vnd.openxmlformats-officedocument.theme+xml")
	overrides = appendPathOverrides(overrides, layoutPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml")
	overrides = appendPathOverrides(overrides, masterPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasNotesMaster, "/ppt/notesMasters/notesMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasHandoutMaster, "/ppt/handoutMasters/handoutMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasCommentAuthors, "/ppt/commentAuthors.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.commentAuthors+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasVBA, "/ppt/vbaProject.bin",
		"application/vnd.ms-office.vbaProject")
	overrides = appendPathOverrides(overrides, commentPaths, commentsPartType)
	overrides = appendPathOverrides(overrides, customXMLPropsPaths,
		"application/vnd.openxmlformats-officedocument.customXmlProperties+xml")

	sort.Slice(overrides, func(i, j int) bool { return overrides[i].PartName < overrides[j].PartName })
	doc.Overrides = overrides
	return renderContentTypesDocument(doc)
}

func collectNonSlideRelationships(
	nonSlide []common.EditorRelationship,
	slideCapacity int,
) ([]common.EditorRelationship, map[string]struct{}, bool, error) {
	rels := make([]common.EditorRelationship, 0, len(nonSlide)+slideCapacity+1)
	used := map[string]struct{}{}
	hasSectionRel := false
	for _, rel := range nonSlide {
		id := strings.TrimSpace(rel.ID)
		if id == "" {
			return nil, nil, false, errors.New("non-slide relationship has empty Id")
		}
		if _, exists := used[id]; exists {
			return nil, nil, false, fmt.Errorf("duplicate relationship Id %q", id)
		}
		used[id] = struct{}{}
		rels = append(rels, rel)
		if rel.Type == common.RelTypeSectionList {
			hasSectionRel = true
		}
	}
	return rels, used, hasSectionRel, nil
}

func makeSectionListRelationship(
	rels []common.EditorRelationship,
	slides []common.EditorSlideRef,
) common.EditorRelationship {
	maxNum := maxRelationshipNumber(rels, slides)
	return common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", maxNum+1),
		Type:   common.RelTypeSectionList,
		Target: "sectionList.xml",
	}
}

const vbaProjectRelType = "http://schemas.microsoft.com/office/2006/relationships/vbaProject"

func makeVBARelationship(
	rels []common.EditorRelationship,
	slides []common.EditorSlideRef,
) common.EditorRelationship {
	maxNum := maxRelationshipNumber(rels, slides)
	return common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", maxNum+1),
		Type:   vbaProjectRelType,
		Target: "vbaProject.bin",
	}
}

func maxRelationshipNumber(rels []common.EditorRelationship, slides []common.EditorSlideRef) int {
	maxNum := 0
	for _, r := range rels {
		if n, ok := parseRelationshipNumber(r.ID); ok && n > maxNum {
			maxNum = n
		}
	}
	for _, slide := range slides {
		if n, ok := parseRelationshipNumber(slide.RelID); ok && n > maxNum {
			maxNum = n
		}
	}
	return maxNum
}

func appendMissingSlideRelationships(
	rels []common.EditorRelationship,
	used map[string]struct{},
	slides []common.EditorSlideRef,
) []common.EditorRelationship {
	for _, slide := range slides {
		if _, exists := used[slide.RelID]; exists {
			continue
		}
		rels = append(rels, common.EditorRelationship{
			ID:     slide.RelID,
			Type:   common.RelTypeSlide,
			Target: slide.Target,
		})
	}
	return rels
}

func sortRelationshipsByID(rels []common.EditorRelationship) {
	sort.Slice(rels, func(i, j int) bool {
		a, aok := common.ParseRelationshipNumber(rels[i].ID)
		b, bok := common.ParseRelationshipNumber(rels[j].ID)
		if aok && bok && a != b {
			return a < b
		}
		return rels[i].ID < rels[j].ID
	})
}

func relationshipsXMLDocument(rels []common.EditorRelationship) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n")
	b.WriteString(`<Relationships xmlns="` + common.RelationshipsXMLNS + `">`)
	for _, rel := range rels {
		writeRelationshipXML(&b, rel)
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

func writeRelationshipXML(b *strings.Builder, rel common.EditorRelationship) {
	b.WriteString("\n<Relationship Id=\"")
	b.WriteString(common.XMLEscape(rel.ID))
	b.WriteString("\" Type=\"")
	b.WriteString(common.XMLEscape(rel.Type))
	b.WriteString("\" Target=\"")
	b.WriteString(common.XMLEscape(rel.Target))
	b.WriteString("\"")
	if strings.TrimSpace(rel.TargetMode) != "" {
		b.WriteString(` TargetMode="` + common.XMLEscape(rel.TargetMode) + `"`)
	}
	b.WriteString("/>")
}

func parseContentTypesDocument(current []byte) (contentTypesDocument, error) {
	if len(current) == 0 {
		return contentTypesDocument{}, errors.New("missing content types content")
	}

	var doc contentTypesDocument
	if err := xml.Unmarshal(current, &doc); err != nil {
		return contentTypesDocument{}, fmt.Errorf("parse content types: %w", err)
	}
	if strings.TrimSpace(doc.XMLNS) == "" {
		doc.XMLNS = common.ContentTypesXMLNS
	}
	return doc, nil
}

func ensureContentTypeDefaults(doc *contentTypesDocument, mediaPaths []string, hasVBA bool) {
	exts := make(map[string]struct{}, len(doc.Defaults))
	for _, d := range doc.Defaults {
		exts[strings.ToLower(d.Extension)] = struct{}{}
	}
	for _, mediaPath := range mediaPaths {
		ext := strings.TrimPrefix(strings.ToLower(path.Ext(strings.TrimSpace(mediaPath))), ".")
		if _, ok := exts[ext]; ok {
			continue
		}
		contentType := contentTypeForExtension(ext)
		if contentType == "" {
			continue
		}
		doc.Defaults = append(doc.Defaults, contentTypeDefault{
			Extension:   ext,
			ContentType: contentType,
		})
		exts[ext] = struct{}{}
	}

	if hasVBA {
		if _, ok := exts["bin"]; !ok {
			doc.Defaults = append(doc.Defaults, contentTypeDefault{
				Extension:   "bin",
				ContentType: "application/vnd.ms-office.vbaProject",
			})
			exts["bin"] = struct{}{}
		}
	}
}

func contentTypeForExtension(ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	switch ext {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "bmp":
		return "image/bmp"
	case "tif", "tiff":
		return "image/tiff"
	case "wav":
		return "audio/wav"
	case "mp3":
		return "audio/mpeg"
	case "m4a":
		return "audio/mp4"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return ""
	}
}

func filterDynamicOverrides(existing []contentTypeOverride, slideCapacity int) []contentTypeOverride {
	filtered := make([]contentTypeOverride, 0, len(existing)+slideCapacity+1)
	for _, override := range existing {
		part := common.CanonicalPartPath(override.PartName)
		if shouldSkipOverridePart(part) {
			continue
		}
		filtered = append(filtered, override)
	}
	return filtered
}

func shouldSkipOverridePart(part string) bool {
	if isSlidePartOverride(part) {
		return true
	}
	return part == "ppt/sectionList.xml" ||
		part == "ppt/commentAuthors.xml" ||
		strings.HasPrefix(part, "ppt/charts/chart") ||
		strings.HasPrefix(part, "ppt/notesSlides/notesSlide") ||
		strings.HasPrefix(part, "ppt/notesMasters/notesMaster") ||
		strings.HasPrefix(part, "ppt/theme/theme") ||
		strings.HasPrefix(part, "ppt/slideLayouts/slideLayout") ||
		strings.HasPrefix(part, "ppt/slideMasters/slideMaster") ||
		strings.HasPrefix(part, "ppt/comments/comment")
}

func appendSlideOverrides(
	overrides []contentTypeOverride,
	slides []common.EditorSlideRef,
) []contentTypeOverride {
	for _, slide := range slides {
		overrides = append(overrides, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(slide.Part),
			ContentType: common.SlideContentType,
		})
	}
	return overrides
}

func appendPathOverrides(
	overrides []contentTypeOverride,
	paths []string,
	contentType string,
) []contentTypeOverride {
	for _, p := range paths {
		overrides = append(overrides, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: contentType,
		})
	}
	return overrides
}

func appendOptionalContentTypeOverride(
	overrides []contentTypeOverride,
	include bool,
	partName,
	contentType string,
) []contentTypeOverride {
	if !include {
		return overrides
	}
	return append(overrides, contentTypeOverride{
		PartName:    partName,
		ContentType: contentType,
	})
}

func renderContentTypesDocument(doc contentTypesDocument) (string, error) {
	rendered, err := xml.MarshalIndent(doc, "", "")
	if err != nil {
		return "", err
	}
	return xml.Header + strings.TrimSpace(string(rendered)), nil
}

type contentTypesDocument struct {
	XMLName   xml.Name              `xml:"Types"`
	XMLNS     string                `xml:"xmlns,attr,omitempty"`
	Defaults  []contentTypeDefault  `xml:"Default"`
	Overrides []contentTypeOverride `xml:"Override"`
}

type contentTypeDefault struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

type contentTypeOverride struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

func isSlidePartOverride(partName string) bool {
	clean := common.CanonicalPartPath(strings.TrimPrefix(strings.TrimSpace(partName), "/"))
	return strings.HasPrefix(clean, "ppt/slides/slide") && strings.HasSuffix(clean, ".xml")
}

func buildSectionListXML(sections []Section) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n")
	b.WriteString(`<s:sectionLst xmlns:s="http://schemas.microsoft.com/office/powerpoint/2010/main">`)
	for _, s := range sections {
		b.WriteString("\n  <s:section name=\"")
		b.WriteString(common.XMLEscape(s.Name))
		b.WriteString("\" id=\"")
		b.WriteString(s.GUID)
		b.WriteString("\">")
		b.WriteString("\n    <s:sldIdLst>")
		for _, slideID := range s.SlideIDs {
			b.WriteString(fmt.Sprintf("\n      <s:sldId id=\"%d\"/>", slideID))
		}
		b.WriteString("\n    </s:sldIdLst>")
		b.WriteString("\n  </s:section>")
	}
	b.WriteString("\n</s:sectionLst>")
	return b.String()
}

var (
	chartExternalDataPattern = regexp.MustCompile(`<c:externalData[^>]*r:id="([^"]*)"[^>]*/>`)
	chartRelIDAttrPattern    = regexp.MustCompile(`r:id="[^"]*"`)
)

func rewriteChartExternalData(current []byte, newRelID string) []byte {
	source := string(current)
	if !chartExternalDataPattern.MatchString(source) {
		return current
	}
	// Simplified replace that preserves other attributes if any (though usually it's just r:id)
	res := chartExternalDataPattern.ReplaceAllStringFunc(source, func(match string) string {
		// Just replace the r:id attribute value
		return chartRelIDAttrPattern.ReplaceAllString(match, fmt.Sprintf(`r:id="%s"`, newRelID))
	})
	return []byte(res)
}

var extLstPattern = regexp.MustCompile(`(?s)<p:extLst>.*?</p:extLst>|<p:extLst\s*/>`)

func rewritePresentationSections(current []byte, sections []Section) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	// Build the p14:sectionLst XML
	sectionXML := buildPresentationSectionExtensionXML(sections)

	// Extension URI for Section List
	extURI := "{521415D9-36F7-43E2-AB2F-B90AF26B5E84}"

	// Helper to build the full <p:ext> block
	fullExtBlock := fmt.Sprintf(`<p:ext uri="%s">%s</p:ext>`, extURI, sectionXML)

	// Strategy:
	// 1. If <p:extLst> exists:
	//    a. If matches our URI, replace the content.
	//    b. Else, append our <p:ext> to the list.
	// 2. If <p:extLst> missing:
	//    a. Insert new <p:extLst> at end of presentation (before closing tag).

	if extLstPattern.MatchString(source) {
		rewritten := extLstPattern.ReplaceAllStringFunc(source, func(match string) string {
			return rewriteExtListMatch(match, extURI, fullExtBlock)
		})
		return rewritten, nil
	}

	// No extLst, insert it.
	// Valid insertion point is after <p:defaultTextStyle> or just before </p:presentation>
	newExtLst := fmt.Sprintf("<p:extLst>\n%s\n</p:extLst>", fullExtBlock)
	if idx := strings.LastIndex(source, "</p:presentation>"); idx >= 0 {
		return source[:idx] + newExtLst + source[idx:], nil
	}

	return "", errors.New("presentation XML malformed (missing </p:presentation>)")
}

func rewriteExtListMatch(match, extURI, fullExtBlock string) string {
	if strings.Contains(match, extURI) {
		return replaceSectionExtension(match, extURI, fullExtBlock)
	}
	return appendSectionExtension(match, fullExtBlock)
}

func replaceSectionExtension(match, extURI, fullExtBlock string) string {
	pExtPattern := regexp.MustCompile(fmt.Sprintf(`(?s)<p:ext uri="%s">.*?</p:ext>`, regexp.QuoteMeta(extURI)))
	if pExtPattern.MatchString(match) {
		return pExtPattern.ReplaceAllString(match, fullExtBlock)
	}
	return appendSectionExtension(match, fullExtBlock)
}

func appendSectionExtension(match, fullExtBlock string) string {
	if strings.Contains(match, "</p:extLst>") {
		return strings.Replace(match, "</p:extLst>", "\n"+fullExtBlock+"\n</p:extLst>", 1)
	}
	return strings.Replace(match, "/>", ">\n"+fullExtBlock+"\n</p:extLst>", 1)
}

func buildPresentationSectionExtensionXML(sections []Section) string {
	if len(sections) == 0 {
		return ""
	}
	var b strings.Builder
	// Note: using p14 prefix requires definition, but p:ext usually allows it if defined in parent or self.
	// Usually presentation.xml has xmlns:p14. If not we might need to add it or use local xmlns.
	// Safest is local xmlns.
	b.WriteString(`<p14:sectionLst xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main">`)
	for _, s := range sections {
		// Section entries
		b.WriteString(fmt.Sprintf("\n<p14:section name=\"%s\" id=\"%s\">", common.XMLEscape(s.Name), s.GUID))
		b.WriteString("\n<p14:sldIdLst>")
		for _, sid := range s.SlideIDs {
			b.WriteString(fmt.Sprintf("\n<p14:sldId id=\"%d\"/>", sid))
		}
		b.WriteString("\n</p14:sldIdLst>")
		b.WriteString("\n</p14:section>")
	}
	b.WriteString("\n</p14:sectionLst>")
	return b.String()
}

var embeddedFontLstPattern = regexp.MustCompile(`(?s)<p:embeddedFontLst>.*?</p:embeddedFontLst>`)

func rewritePresentationEmbeddedFonts(current []byte, fontLst string) (string, error) {
	if fontLst == "" {
		return string(current), nil
	}
	source := string(current)

	// If it already has an embeddedFontLst, replace it.
	if embeddedFontLstPattern.MatchString(source) {
		return embeddedFontLstPattern.ReplaceAllString(source, fontLst), nil
	}

	// Otherwise, insert it. A safe place is before <p:extLst> or </p:presentation>.
	if strings.Contains(source, "<p:extLst>") {
		return strings.Replace(source, "<p:extLst>", fontLst+"\n<p:extLst>", 1), nil
	}

	idx := strings.LastIndex(source, "</p:presentation>")
	if idx >= 0 {
		return source[:idx] + fontLst + "\n" + source[idx:], nil
	}

	return "", errors.New("presentation XML malformed (missing </p:presentation>)")
}
