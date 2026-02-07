package pptx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
)

const relationshipsXMLNS = "http://schemas.openxmlformats.org/package/2006/relationships"
const contentTypesXMLNS = "http://schemas.openxmlformats.org/package/2006/content-types"
const slideContentType = "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"

func rewritePresentationSlideList(current []byte, slides []editorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", fmt.Errorf("missing %s content", presentationXMLPath)
	}
	source := string(current)
	start := strings.Index(source, "<p:sldIdLst>")
	end := strings.Index(source, "</p:sldIdLst>")
	if start < 0 || end < 0 || end < start {
		return "", fmt.Errorf("%s does not contain <p:sldIdLst>", presentationXMLPath)
	}
	end += len("</p:sldIdLst>")

	replacement := buildPresentationSlideListXML(slides)
	return source[:start] + replacement + source[end:], nil
}

func buildPresentationSlideListXML(slides []editorSlideRef) string {
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

func renderPresentationRelsXML(nonSlide []editorRelationship, slides []editorSlideRef) (string, error) {
	rels := make([]editorRelationship, 0, len(nonSlide)+len(slides))
	used := map[string]struct{}{}

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
	}
	for _, slide := range slides {
		if _, exists := used[slide.RelID]; exists {
			return "", fmt.Errorf("duplicate relationship Id %q", slide.RelID)
		}
		used[slide.RelID] = struct{}{}
		rels = append(rels, editorRelationship{
			ID:     slide.RelID,
			Type:   relTypeSlide,
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
	b.WriteString(`<Relationships xmlns="` + relationshipsXMLNS + `">`)
	for _, rel := range rels {
		b.WriteString("\n  <Relationship Id=\"")
		b.WriteString(xmlEscape(rel.ID))
		b.WriteString("\" Type=\"")
		b.WriteString(xmlEscape(rel.Type))
		b.WriteString("\" Target=\"")
		b.WriteString(xmlEscape(rel.Target))
		b.WriteString("\"")
		if strings.TrimSpace(rel.TargetMode) != "" {
			b.WriteString(` TargetMode="` + xmlEscape(rel.TargetMode) + `"`)
		}
		b.WriteString("/>")
	}
	b.WriteString("\n</Relationships>")
	return b.String(), nil
}

func rewriteContentTypesSlideOverrides(current []byte, slides []editorSlideRef) (string, error) {
	if len(current) == 0 {
		return "", fmt.Errorf("missing %s content", contentTypesPath)
	}

	var doc contentTypesDocument
	if err := xml.Unmarshal(current, &doc); err != nil {
		return "", fmt.Errorf("parse %s: %w", contentTypesPath, err)
	}
	if strings.TrimSpace(doc.XMLNS) == "" {
		doc.XMLNS = contentTypesXMLNS
	}

	filtered := make([]contentTypeOverride, 0, len(doc.Overrides)+len(slides))
	for _, override := range doc.Overrides {
		if isSlidePartOverride(override.PartName) {
			continue
		}
		filtered = append(filtered, override)
	}
	for _, slide := range slides {
		filtered = append(filtered, contentTypeOverride{
			PartName:    "/" + canonicalPartPath(slide.Part),
			ContentType: slideContentType,
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
	clean := canonicalPartPath(strings.TrimPrefix(strings.TrimSpace(partName), "/"))
	return strings.HasPrefix(clean, "ppt/slides/slide") && strings.HasSuffix(clean, ".xml")
}

func xmlEscape(value string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}
