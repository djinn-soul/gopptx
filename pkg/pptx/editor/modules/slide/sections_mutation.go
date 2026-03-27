package slide

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type SectionData struct {
	Name     string
	GUID     string
	SlideIDs []int64
}

const sectionExtURI = "{521415D9-36F7-43E2-AB2F-B90AF26B5E84}"

var (
	chartExternalDataPattern = regexp.MustCompile(`<c:externalData[^>]*r:id="([^"]*)"[^>]*/>`)
	chartRelIDAttrPattern    = regexp.MustCompile(`r:id="[^"]*"`)
	extLstPattern            = regexp.MustCompile(`(?s)<p:extLst>.*?</p:extLst>|<p:extLst\s*/>`)
	embeddedFontLstPattern   = regexp.MustCompile(`(?s)<p:embeddedFontLst>.*?</p:embeddedFontLst>`)
	sectionExtPattern        = regexp.MustCompile(`(?s)<p:ext uri="` + regexp.QuoteMeta(sectionExtURI) + `">.*?</p:ext>`)
)

func BuildSectionListXML(sections []SectionData) string {
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
			b.WriteString("\n      <s:sldId id=\"")
			b.WriteString(strconv.FormatInt(slideID, 10))
			b.WriteString("\"/>")
		}
		b.WriteString("\n    </s:sldIdLst>")
		b.WriteString("\n  </s:section>")
	}
	b.WriteString("\n</s:sectionLst>")
	return b.String()
}

func RewriteChartExternalData(current []byte, newRelID string) []byte {
	source := string(current)
	if !chartExternalDataPattern.MatchString(source) {
		return current
	}
	res := chartExternalDataPattern.ReplaceAllStringFunc(source, func(match string) string {
		return chartRelIDAttrPattern.ReplaceAllString(match, `r:id="`+newRelID+`"`)
	})
	return []byte(res)
}

func RewritePresentationSections(current []byte, sections []SectionData) (string, error) {
	if len(current) == 0 {
		return "", errors.New("missing presentation XML content")
	}
	source := string(current)

	sectionXML := buildPresentationSectionExtensionXML(sections)
	fullExtBlock := `<p:ext uri="` + sectionExtURI + `">` + sectionXML + `</p:ext>`

	if extLstPattern.MatchString(source) {
		rewritten := extLstPattern.ReplaceAllStringFunc(source, func(match string) string {
			return rewriteExtListMatch(match, fullExtBlock)
		})
		return rewritten, nil
	}

	newExtLst := "<p:extLst>\n" + fullExtBlock + "\n</p:extLst>"
	if idx := strings.LastIndex(source, "</p:presentation>"); idx >= 0 {
		return source[:idx] + newExtLst + source[idx:], nil
	}

	return "", errors.New("presentation XML malformed (missing </p:presentation>)")
}

func RewritePresentationEmbeddedFonts(current []byte, fontLst string) (string, error) {
	if fontLst == "" {
		return string(current), nil
	}
	source := string(current)

	if embeddedFontLstPattern.MatchString(source) {
		return embeddedFontLstPattern.ReplaceAllString(source, fontLst), nil
	}

	if strings.Contains(source, "<p:extLst>") {
		return strings.Replace(source, "<p:extLst>", fontLst+"\n<p:extLst>", 1), nil
	}

	idx := strings.LastIndex(source, "</p:presentation>")
	if idx >= 0 {
		return source[:idx] + fontLst + "\n" + source[idx:], nil
	}

	return "", errors.New("presentation XML malformed (missing </p:presentation>)")
}

func ExtractEmbeddedFontList(xmlData []byte) string {
	match := embeddedFontLstPattern.Find(xmlData)
	if match == nil {
		return ""
	}
	return string(match)
}

func rewriteExtListMatch(match, fullExtBlock string) string {
	if strings.Contains(match, sectionExtURI) {
		return replaceSectionExtension(match, fullExtBlock)
	}
	return appendSectionExtension(match, fullExtBlock)
}

func replaceSectionExtension(match, fullExtBlock string) string {
	if sectionExtPattern.MatchString(match) {
		return sectionExtPattern.ReplaceAllString(match, fullExtBlock)
	}
	return appendSectionExtension(match, fullExtBlock)
}

func appendSectionExtension(match, fullExtBlock string) string {
	if strings.Contains(match, "</p:extLst>") {
		return strings.Replace(match, "</p:extLst>", "\n"+fullExtBlock+"\n</p:extLst>", 1)
	}
	return strings.Replace(match, "/>", ">\n"+fullExtBlock+"\n</p:extLst>", 1)
}

func buildPresentationSectionExtensionXML(sections []SectionData) string {
	if len(sections) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<p14:sectionLst xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main">`)
	for _, s := range sections {
		b.WriteString("\n<p14:section name=\"")
		b.WriteString(common.XMLEscape(s.Name))
		b.WriteString("\" id=\"")
		b.WriteString(s.GUID)
		b.WriteString("\">")
		b.WriteString("\n<p14:sldIdLst>")
		for _, sid := range s.SlideIDs {
			b.WriteString("\n<p14:sldId id=\"")
			b.WriteString(strconv.FormatInt(sid, 10))
			b.WriteString("\"/>")
		}
		b.WriteString("\n</p14:sldIdLst>")
		b.WriteString("\n</p14:section>")
	}
	b.WriteString("\n</p14:sectionLst>")
	return b.String()
}
