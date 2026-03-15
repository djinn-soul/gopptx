package structural

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

func (r *Repairer) repairInvalidContentType(p string) error {
	ctPath := "[Content_Types].xml"
	data, ok := r.modifier.Get(ctPath)
	if !ok {
		return r.repairMissingPart(ctPath)
	}

	ext := path.Ext(p)
	contentType := r.inferContentType(p)

	newEntry := ""
	if ext != "" {
		ext = strings.TrimPrefix(ext, ".")
		if !strings.Contains(string(data), fmt.Sprintf(`Extension="%s"`, ext)) {
			newEntry = fmt.Sprintf("\n  <Default Extension=\"%s\" ContentType=\"%s\"/>", ext, contentType)
		}
	}
	if newEntry == "" {
		newEntry = fmt.Sprintf("\n  <Override PartName=\"/%s\" ContentType=\"%s\"/>", p, contentType)
	}

	content := string(data)
	closingTagIdx := strings.LastIndex(strings.ToLower(content), "</types>")
	if closingTagIdx == -1 {
		return errors.New("invalid [Content_Types].xml: missing closing tag")
	}

	repaired := content[:closingTagIdx] + newEntry + "\n" + content[closingTagIdx:]
	r.modifier.Set(ctPath, []byte(repaired))
	return nil
}

func (r *Repairer) inferContentType(p string) string {
	if strings.Contains(p, "slideLayout") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"
	}
	if strings.Contains(p, "slideMaster") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"
	}
	if strings.Contains(p, "slide") && strings.HasSuffix(p, ".xml") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"
	}
	if strings.Contains(p, "presentation.xml") {
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"
	}
	if strings.HasSuffix(p, ".xml") {
		return "application/xml"
	}
	return "application/octet-stream"
}

func (r *Repairer) generateContentTypes() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	sb.WriteString("\n")
	sb.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	sb.WriteString("\n")
	sb.WriteString(
		`  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`,
	)
	sb.WriteString("\n")
	sb.WriteString(`  <Default Extension="xml" ContentType="application/xml"/>`)
	sb.WriteString("\n")

	for _, p := range r.modifier.Keys() {
		if strings.HasSuffix(p, ".rels") || p == "[Content_Types].xml" {
			continue
		}
		contentType := r.inferContentType(p)
		sb.WriteString(fmt.Sprintf("  <Override PartName=\"/%s\" ContentType=\"%s\"/>\n", p, contentType))
	}

	sb.WriteString(`</Types>`)
	return sb.String()
}
