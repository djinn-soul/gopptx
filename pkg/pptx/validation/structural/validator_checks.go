package structural

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
	"strings"
	"sync"
)

func (v *Validator) checkSlideReferences() {
	if !v.provider.Has(presentationRelsPath) {
		return
	}

	data, _ := v.provider.Get(presentationRelsPath)
	referencedSlides := make(map[string]bool)

	// Parse relationships using proper XML decoder
	var rels relationshipsXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&rels); err != nil {
		// If XML parsing fails, the XML validity check will report it
		return
	}

	for _, rel := range rels.Relationships {
		// Check if this is a slide relationship (but not slideMaster or slideLayout)
		// The type should end with "/slide" not contain "/slide" somewhere
		if strings.HasSuffix(rel.Type, "/slide") && rel.Target != "" {
			fullPath := v.resolvePath(presentationRelsPath, rel.Target)
			referencedSlides[fullPath] = true
		}
	}

	// Find actual slide files
	for _, p := range v.keys {
		if strings.HasPrefix(p, "ppt/slides/slide") && strings.HasSuffix(p, ".xml") && !strings.Contains(p, "_rels") {
			if !referencedSlides[p] {
				v.issues = append(v.issues, Issue{
					Code:        CodeOrphanSlide,
					Severity:    SeverityInfo,
					Path:        p,
					Description: fmt.Sprintf("Orphan slide: %s is not referenced in presentation.xml", p),
					Repairable:  true,
				})
			}
		}
	}

	// Check for missing slide files that are referenced
	for slidePath := range referencedSlides {
		if !v.provider.Has(slidePath) {
			v.issues = append(v.issues, Issue{
				Code:        CodeMissingSlideRef,
				Severity:    SeverityError,
				Path:        presentationRelsPath,
				Description: fmt.Sprintf("Referenced slide not found: %s", slidePath),
				Repairable:  true,
			})
		}
	}
}

func (v *Validator) checkContentTypes() {
	data, ok := v.provider.Get(contentTypesPath)
	if !ok {
		return
	}

	var ct contentTypesXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&ct); err != nil {
		return
	}

	overrides := make(map[string]bool)
	for _, o := range ct.Overrides {
		overrides[strings.TrimPrefix(o.PartName, "/")] = true
	}

	defaults := make(map[string]bool)
	for _, d := range ct.Defaults {
		defaults[strings.ToLower(d.Extension)] = true
	}

	for _, p := range v.keys {
		if p == contentTypesPath || strings.HasSuffix(p, ".rels") {
			continue
		}

		if overrides[p] {
			continue
		}

		ext := strings.TrimPrefix(path.Ext(p), ".")
		if ext != "" && defaults[strings.ToLower(ext)] {
			continue
		}

		v.issues = append(v.issues, Issue{
			Code:        CodeInvalidContentType,
			Severity:    SeverityError,
			Path:        p,
			Description: fmt.Sprintf("Part %s has no content type registration", p),
			Repairable:  true,
		})
	}
}

func (v *Validator) checkNamespaces() {
	paths := v.keys
	issuesChan := make(chan []Issue, len(paths))
	var wg sync.WaitGroup

	for _, p := range paths {
		if !strings.HasSuffix(p, ".xml") {
			continue
		}
		wg.Add(1)
		go func(partPath string) {
			defer wg.Done()
			issuesChan <- v.checkNamespacesForPart(partPath)
		}(p)
	}

	wg.Wait()
	close(issuesChan)

	for issues := range issuesChan {
		v.issues = append(v.issues, issues...)
	}
}

func (v *Validator) checkEmptyElements() {
	// A basic implementation to catch some known fragile elements
	// that shouldn't be empty self-closing (like p:sldIdLst)
	// unless they are explicitly handled correctly by a reader.
	paths := v.keys
	issuesChan := make(chan []Issue, len(paths))
	var wg sync.WaitGroup

	for _, p := range paths {
		if !strings.HasSuffix(p, ".xml") {
			continue
		}
		wg.Add(1)
		go func(partPath string) {
			defer wg.Done()
			issuesChan <- v.checkEmptyElementsForPart(partPath)
		}(p)
	}

	wg.Wait()
	close(issuesChan)

	for issues := range issuesChan {
		v.issues = append(v.issues, issues...)
	}
}

var (
	xmlDeclEnd  = []byte("?>")
	xmlTagOpen  = []byte("<")
	xmlTagClose = []byte(">")
)

func (v *Validator) checkNamespacesForPart(partPath string) []Issue {
	data, ok := v.provider.Get(partPath)
	if !ok {
		return nil
	}

	// Operate on []byte directly to avoid string(data) copy.
	startIdx := 0
	if idx := bytes.Index(data, xmlDeclEnd); idx != -1 {
		startIdx = idx + xmlDeclSuffixLength
	}

	firstTagStart := bytes.Index(data[startIdx:], xmlTagOpen)
	if firstTagStart == -1 {
		return nil
	}

	firstTagEnd := bytes.Index(data[startIdx+firstTagStart:], xmlTagClose)
	if firstTagEnd == -1 {
		return nil
	}

	openingTag := data[startIdx+firstTagStart : startIdx+firstTagStart+firstTagEnd+1]
	if !isPresentationNamespacePart(partPath) {
		return nil
	}

	issues := make([]Issue, 0, namespaceIssueCap)
	if !bytes.Contains(openingTag, []byte("xmlns:p=")) &&
		!bytes.Contains(openingTag, []byte(`xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"`)) {
		issues = append(issues, Issue{
			Code:        CodeMissingNamespace,
			Severity:    SeverityWarning,
			Path:        partPath,
			Description: "Missing presentationml namespace declaration",
			Repairable:  true,
			Context:     map[string]string{"ns": "p"},
		})
	}
	if !bytes.Contains(openingTag, []byte("xmlns:a=")) &&
		!bytes.Contains(openingTag, []byte(`xmlns="http://schemas.openxmlformats.org/drawingml/2006/main"`)) {
		issues = append(issues, Issue{
			Code:        CodeMissingNamespace,
			Severity:    SeverityWarning,
			Path:        partPath,
			Description: "Missing drawingml namespace declaration",
			Repairable:  true,
			Context:     map[string]string{"ns": "a"},
		})
	}
	if !bytes.Contains(openingTag, []byte("xmlns:r=")) &&
		!bytes.Contains(openingTag, []byte(`xmlns="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`)) {
		issues = append(issues, Issue{
			Code:        CodeMissingNamespace,
			Severity:    SeverityWarning,
			Path:        partPath,
			Description: "Missing relationships namespace declaration",
			Repairable:  true,
			Context:     map[string]string{"ns": "r"},
		})
	}
	return issues
}

var emptySldIdLst = []byte("<p:sldIdLst/>")

func (v *Validator) checkEmptyElementsForPart(partPath string) []Issue {
	data, ok := v.provider.Get(partPath)
	if !ok {
		return nil
	}
	if !bytes.Contains(data, emptySldIdLst) {
		return nil
	}
	return []Issue{{
		Code:        CodeEmptyRequiredElement,
		Severity:    SeverityInfo,
		Path:        partPath,
		Description: "Found empty self-closing <p:sldIdLst/> element",
		Repairable:  true,
		Context:     map[string]string{"element": "p:sldIdLst"},
	}}
}

func isPresentationNamespacePart(partPath string) bool {
	return strings.HasPrefix(partPath, "ppt/presentation.xml") ||
		strings.HasPrefix(partPath, "ppt/slides/") ||
		strings.HasPrefix(partPath, "ppt/slideMasters/") ||
		strings.HasPrefix(partPath, "ppt/slideLayouts/")
}
