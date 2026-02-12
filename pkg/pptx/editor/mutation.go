package editor

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var sldIdLstPattern = regexp.MustCompile(`(?s)<p:sldIdLst>.*?</p:sldIdLst>|<p:sldIdLst\s*/>`)

func rewritePresentationSlideList(current []byte, slides []common.EditorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", fmt.Errorf("missing presentation XML content")
	}
	source := string(current)

	replacement := buildPresentationSlideListXML(slides)
	if !sldIdLstPattern.MatchString(source) {
		return "", fmt.Errorf("presentation XML does not contain <p:sldIdLst>")
	}

	found := false
	result := sldIdLstPattern.ReplaceAllStringFunc(source, func(match string) string {
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
		b.WriteString(fmt.Sprintf("%d", slide.SlideID))
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

func renderPresentationRelsXML(nonSlide []common.EditorRelationship, slides []common.EditorSlideRef, hasSections bool) (string, error) {
	rels := make([]common.EditorRelationship, 0, len(nonSlide)+len(slides)+1)
	used := map[string]struct{}{}

	hasSectionRel := false
	for _, rel := range nonSlide {
		id := strings.TrimSpace(rel.ID)
		if id == "" {
			return "", fmt.Errorf("non-slide relationship has empty Id")
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
		b.WriteString("\n  <Relationship Id=\"")
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

func rewriteContentTypes(current []byte, slides []common.EditorSlideRef, mediaPaths []string, hasSections bool) (string, error) {
	if len(current) == 0 {
		return "", fmt.Errorf("missing content types content")
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
		ext := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(strings.ToLower(strings.TrimPrefix(m, "."))), "."))
		if idx := strings.LastIndex(m, "."); idx >= 0 {
			ext = strings.ToLower(m[idx+1:])
		}
		if _, ok := exts[ext]; !ok {
			contentType := ""
			switch ext {
			case "png":
				contentType = "image/png"
			case "jpg", "jpeg":
				contentType = "image/jpeg"
			case "gif":
				contentType = "image/gif"
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

	// 2. Manage Overrides (Slides and SectionList)
	filtered := make([]contentTypeOverride, 0, len(doc.Overrides)+len(slides)+1)
	for _, override := range doc.Overrides {
		part := common.CanonicalPartPath(override.PartName)
		if isSlidePartOverride(part) {
			continue
		}
		if part == "ppt/sectionList.xml" {
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
