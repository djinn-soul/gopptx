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

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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
) (string, error) {
	rels := make([]common.EditorRelationship, 0, len(nonSlide)+len(slides)+1)
	used := map[string]struct{}{}

	hasSectionRel := false
	for _, rel := range nonSlide {
		id := strings.TrimSpace(rel.ID)
		if id == "" {
			return "", errors.New("non-slide relationship has empty Id")
		}
		if _, exists := used[id]; exists {
			return "", fmt.Errorf("duplicate relationship Id %q", id)
		}
		used[id] = struct{}{}
		rels = append(rels, rel)
		if rel.Type == common.RelTypeSectionList {
			hasSectionRel = true
		}
	}

	if hasSections && !hasSectionRel {
		// Allocate next rId for sectionList
		maxNum := 0
		for _, r := range rels {
			if n, ok := parseRelationshipNumber(r.ID); ok && n > maxNum {
				maxNum = n
			}
		}
		for _, s := range slides {
			if n, ok := parseRelationshipNumber(s.RelID); ok && n > maxNum {
				maxNum = n
			}
		}
		sectionRelID := fmt.Sprintf("rId%d", maxNum+1)
		rels = append(rels, common.EditorRelationship{
			ID:     sectionRelID,
			Type:   common.RelTypeSectionList,
			Target: "sectionList.xml",
		})
	}

	for _, slide := range slides {
		if _, exists := used[slide.RelID]; exists {
			continue // Already in rels? (Shouldn't happen with slides usually)
		}
		rels = append(rels, common.EditorRelationship{
			ID:     slide.RelID,
			Type:   common.RelTypeSlide,
			Target: slide.Target,
		})
	}

	sort.Slice(rels, func(i, j int) bool {
		a, aok := parseRelationshipNumber(rels[i].ID)
		b, bok := parseRelationshipNumber(rels[j].ID)
		if aok && bok && a != b {
			return a < b
		}
		return rels[i].ID < rels[j].ID
	})

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n")
	b.WriteString(`<Relationships xmlns="` + common.RelationshipsXMLNS + `">`)
	for _, rel := range rels {
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
	b.WriteString("\n</Relationships>")
	return b.String(), nil
}

func renderRelationshipsXML(rels []common.EditorRelationship) string {
	sort.Slice(rels, func(i, j int) bool {
		a, aok := common.ParseRelationshipNumber(rels[i].ID)
		b, bok := common.ParseRelationshipNumber(rels[j].ID)
		if aok && bok && a != b {
			return a < b
		}
		return rels[i].ID < rels[j].ID
	})

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n")
	b.WriteString(`<Relationships xmlns="` + common.RelationshipsXMLNS + `">`)
	for _, rel := range rels {
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
	b.WriteString("\n</Relationships>")
	return b.String()
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
) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing content types content")
	}

	var doc contentTypesDocument
	if err := xml.Unmarshal(current, &doc); err != nil {
		return "", fmt.Errorf("parse content types: %w", err)
	}
	if strings.TrimSpace(doc.XMLNS) == "" {
		doc.XMLNS = common.ContentTypesXMLNS
	}

	// 1. Manage Defaults (Media Extensions)
	exts := make(map[string]struct{})
	for _, d := range doc.Defaults {
		exts[strings.ToLower(d.Extension)] = struct{}{}
	}
	for _, m := range mediaPaths {
		ext := strings.TrimPrefix(strings.ToLower(path.Ext(strings.TrimSpace(m))), ".")
		if _, ok := exts[ext]; !ok {
			contentType := ""
			switch ext {
			case "png":
				contentType = "image/png"
			case "jpg", "jpeg":
				contentType = "image/jpeg"
			case "gif":
				contentType = "image/gif"
			case "xlsx":
				contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			}
			if contentType != "" {
				doc.Defaults = append(doc.Defaults, contentTypeDefault{
					Extension:   ext,
					ContentType: contentType,
				})
				exts[ext] = struct{}{}
			}
		}
	}

	// 2. Manage Overrides (Slides, SectionList, and Comments)
	filtered := make([]contentTypeOverride, 0, len(doc.Overrides)+len(slides)+1)
	for _, override := range doc.Overrides {
		part := common.CanonicalPartPath(override.PartName)
		if isSlidePartOverride(part) {
			continue
		}
		if part == "ppt/sectionList.xml" ||
			part == "ppt/commentAuthors.xml" ||
			strings.HasPrefix(part, "ppt/charts/chart") ||
			strings.HasPrefix(part, "ppt/notesSlides/notesSlide") ||
			strings.HasPrefix(part, "ppt/notesMasters/notesMaster") ||
			strings.HasPrefix(part, "ppt/theme/theme") ||
			strings.HasPrefix(part, "ppt/slideLayouts/slideLayout") ||
			strings.HasPrefix(part, "ppt/slideMasters/slideMaster") ||
			strings.HasPrefix(part, "ppt/comments/comment") {
			continue
		}
		filtered = append(filtered, override)
	}

	for _, slide := range slides {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(slide.Part),
			ContentType: common.SlideContentType,
		})
	}
	if hasSections {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/ppt/sectionList.xml",
			ContentType: "application/vnd.microsoft.powerpoint.sectionList+xml",
		})
	}
	for _, p := range chartPaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.drawingml.chart+xml",
		})
	}
	for _, p := range notesPaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml",
		})
	}
	for _, p := range themePaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.theme+xml",
		})
	}
	for _, p := range layoutPaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml",
		})
	}
	for _, p := range masterPaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml",
		})
	}
	if hasNotesMaster {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/ppt/notesMasters/notesMaster1.xml",
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml",
		})
	}
	if hasCommentAuthors {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/ppt/commentAuthors.xml",
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.commentAuthors+xml",
		})
	}
	for _, p := range commentPaths {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + common.CanonicalPartPath(p),
			ContentType: "application/vnd.openxmlformats-officedocument.presentationml.comments+xml",
		})
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PartName < filtered[j].PartName
	})
	doc.Overrides = filtered

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

func buildSectionListXML(sections []EditorSection) string {
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

func rewritePresentationSections(current []byte, sections []EditorSection) (string, error) {
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
		return extLstPattern.ReplaceAllStringFunc(source, func(match string) string {
			// Check if our extension exists
			if strings.Contains(match, extURI) {
				// Regex to replace specific extension?
				// For simplicity/robustness, let's parse or strict regex.
				// Since we control usage, let's try a regex for the specific p:ext
				pExtPattern := regexp.MustCompile(
					fmt.Sprintf(`(?s)<p:ext uri="%s">.*?</p:ext>`, regexp.QuoteMeta(extURI)),
				)
				if pExtPattern.MatchString(match) {
					return pExtPattern.ReplaceAllString(match, fullExtBlock)
				}
				// URI not found, but list exists. Append.
				// Insert before closing </p:extLst>
				if strings.Contains(match, "</p:extLst>") {
					return strings.Replace(match, "</p:extLst>", "\n"+fullExtBlock+"\n</p:extLst>", 1)
				}
				// Self closing <p:extLst/>?
				return strings.Replace(match, "/>", ">"+fullExtBlock+"</p:extLst>", 1)
			}
			// Extension not present, append it.
			if strings.Contains(match, "</p:extLst>") {
				return strings.Replace(match, "</p:extLst>", "\n"+fullExtBlock+"\n</p:extLst>", 1)
			}
			return strings.Replace(match, "/>", ">\n"+fullExtBlock+"\n</p:extLst>", 1)
		}), nil
	}

	// No extLst, insert it.
	// Valid insertion point is after <p:defaultTextStyle> or just before </p:presentation>
	newExtLst := fmt.Sprintf("<p:extLst>\n%s\n</p:extLst>", fullExtBlock)
	if idx := strings.LastIndex(source, "</p:presentation>"); idx >= 0 {
		return source[:idx] + newExtLst + source[idx:], nil
	}

	return "", errors.New("presentation XML malformed (missing </p:presentation>)")
}

func buildPresentationSectionExtensionXML(sections []EditorSection) string {
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
