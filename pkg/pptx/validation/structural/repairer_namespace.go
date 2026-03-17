package structural

import (
	"fmt"
	"strings"
)

const xmlDeclEndLength = 2

func (r *Repairer) repairMissingNamespace(issue Issue) error {
	data, ok := r.modifier.Get(issue.Path)
	if !ok {
		return nil
	}
	content := string(data)

	ns := issue.Context["ns"]
	if ns == "" {
		return nil
	}

	namespaceDeclaration := namespaceDecl(ns)
	if namespaceDeclaration == "" {
		return nil
	}

	xmlDeclIdx := strings.Index(content, "?>")
	searchStart := 0
	if xmlDeclIdx != -1 {
		searchStart = xmlDeclIdx + xmlDeclEndLength
	}

	firstBracket := strings.Index(content[searchStart:], "<")
	if firstBracket == -1 {
		return nil
	}
	rootStart := searchStart + firstBracket

	firstClose := strings.Index(content[rootStart:], ">")
	if firstClose == -1 {
		return nil
	}

	tagContent := content[rootStart : rootStart+firstClose]
	if strings.Contains(tagContent, namespaceDeclaration) {
		return nil
	}

	repaired := content[:rootStart+firstClose] + " " + namespaceDeclaration + content[rootStart+firstClose:]
	r.modifier.Set(issue.Path, []byte(repaired))
	return nil
}

func namespaceDecl(ns string) string {
	switch ns {
	case "p":
		return `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`
	case "a":
		return `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
	case "r":
		return `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
	default:
		return ""
	}
}

func (r *Repairer) repairEmptyRequiredElement(issue Issue) error {
	data, ok := r.modifier.Get(issue.Path)
	if !ok {
		return nil
	}
	content := string(data)

	element := issue.Context["element"]
	if element == "" {
		return nil
	}

	emptyPattern := fmt.Sprintf("<%s/>", element)
	filledPattern := fmt.Sprintf("<%s></%s>", element, element)

	repaired := strings.ReplaceAll(content, emptyPattern, filledPattern)
	r.modifier.Set(issue.Path, []byte(repaired))
	return nil
}
