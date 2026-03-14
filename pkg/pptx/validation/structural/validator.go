package structural

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
)

const packageRelationshipsXMLNS = "http://schemas.openxmlformats.org/package/2006/relationships"

// PartProvider defines the interface for accessing package parts.
type PartProvider interface {
	Has(path string) bool
	Get(path string) ([]byte, bool)
	Keys() []string
}

// Checker defines the interface for custom validation logic.
type Checker interface {
	Check(ps PartProvider) []Issue
}

// Validator provides methods for validating PPTX packages.
type Validator struct {
	provider PartProvider
	issues   []Issue
	checkers []Checker
}

// NewValidator creates a new validator using the given part provider.
func NewValidator(provider PartProvider) *Validator {
	return &Validator{
		provider: provider,
	}
}

// AddChecker registers a custom checker with the validator.
func (v *Validator) AddChecker(c Checker) {
	v.checkers = append(v.checkers, c)
}

// requiredParts identifies the minimum set of parts for a valid PPTX.
func requiredParts() map[string]string {
	return map[string]string{
		"[Content_Types].xml":             "Content types definition",
		"_rels/.rels":                     "Package relationships",
		"ppt/presentation.xml":            "Presentation document",
		"ppt/_rels/presentation.xml.rels": "Presentation relationships",
	}
}

// Validate performs a comprehensive validation check on the package.
func (v *Validator) Validate() []Issue {
	v.issues = nil
	v.checkRequiredParts()
	v.checkXMLValidity()
	v.checkRelationships()
	v.checkSlideReferences()
	v.checkContentTypes()
	v.checkNamespaces()
	v.checkEmptyElements()

	for _, c := range v.checkers {
		v.issues = append(v.issues, c.Check(v.provider)...)
	}

	return v.issues
}

func (v *Validator) checkRequiredParts() {
	for p, desc := range requiredParts() {
		if !v.provider.Has(p) {
			v.issues = append(v.issues, Issue{
				Code:        CodeMissingPart,
				Severity:    SeverityError,
				Path:        p,
				Description: fmt.Sprintf("Missing required part: %s", desc),
				Repairable:  true,
			})
		}
	}
}

func (v *Validator) checkXMLValidity() {
	paths := v.provider.Keys()
	issuesChan := make(chan Issue, len(paths))
	var wg sync.WaitGroup

	for _, p := range paths {
		if !strings.HasSuffix(p, ".xml") && !strings.HasSuffix(p, ".rels") {
			continue
		}
		wg.Add(1)
		go func(partPath string) {
			defer wg.Done()

			data, ok := v.provider.Get(partPath)
			if !ok {
				return
			}

			if err := v.validateXML(data); err != nil {
				issuesChan <- Issue{
					Code:        CodeInvalidXML,
					Severity:    SeverityError,
					Path:        partPath,
					Description: fmt.Sprintf("Invalid XML: %v", err),
					Repairable:  true,
				}
			}
		}(p)
	}

	wg.Wait()
	close(issuesChan)

	for issue := range issuesChan {
		v.issues = append(v.issues, issue)
	}
}

func (v *Validator) validateXML(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty part")
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		_, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) checkRelationships() {
	paths := v.provider.Keys()
	issuesChan := make(chan []Issue, len(paths))
	var wg sync.WaitGroup

	for _, p := range paths {
		if !strings.HasSuffix(p, ".rels") {
			continue
		}
		wg.Add(1)
		go func(relsPath string) {
			defer wg.Done()
			issuesChan <- v.checkRelsFile(relsPath)
		}(p)
	}

	wg.Wait()
	close(issuesChan)

	for issues := range issuesChan {
		v.issues = append(v.issues, issues...)
	}
}

// relationshipXML represents a Relationship element in .rels files.
type relationshipXML struct {
	ID     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

// relationshipsXML represents the root Relationships element.
type relationshipsXML struct {
	XMLName       xml.Name          `xml:"Relationships"`
	XMLNS         string            `xml:"xmlns,attr,omitempty"`
	Relationships []relationshipXML `xml:"Relationship"`
}

// contentTypeDefault represents a Default element in [Content_Types].xml.
type contentTypeDefault struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// contentTypeOverride represents an Override element in [Content_Types].xml.
type contentTypeOverride struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// contentTypesXML represents the root Types element in [Content_Types].xml.
type contentTypesXML struct {
	XMLName   xml.Name              `xml:"Types"`
	Defaults  []contentTypeDefault  `xml:"Default"`
	Overrides []contentTypeOverride `xml:"Override"`
}

func (v *Validator) checkRelsFile(relsPath string) []Issue {
	var issues []Issue
	data, ok := v.provider.Get(relsPath)
	if !ok {
		return issues
	}

	// Parse relationships using proper XML decoder
	var rels relationshipsXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&rels); err != nil {
		// If XML parsing fails, the XML validity check will report it
		return issues
	}

	for _, rel := range rels.Relationships {
		target := rel.Target
		if target == "" {
			continue
		}

		// Skip external targets
		if strings.HasPrefix(target, "http") || strings.HasPrefix(target, "mailto") {
			continue
		}

		fullPath := v.resolvePath(relsPath, target)
		if !v.provider.Has(fullPath) {
			issues = append(issues, Issue{
				Code:        CodeBrokenRelationship,
				Severity:    SeverityWarning,
				Path:        relsPath,
				Description: fmt.Sprintf("Broken relationship: %s -> %s (ID: %s)", relsPath, target, rel.ID),
				Repairable:  true,
				Context:     map[string]string{"target": target},
			})
		}
	}
	return issues
}

func (v *Validator) resolvePath(relsPath, target string) string {
	if trimmed, ok := strings.CutPrefix(target, "/"); ok {
		return trimmed
	}

	dir := path.Dir(relsPath)
	// OpenXML rels folder is usually '_rels'
	baseDir := strings.TrimSuffix(dir, "_rels")
	resolved := path.Join(baseDir, target)
	return strings.ReplaceAll(resolved, "\\", "/")
}

func (v *Validator) checkSlideReferences() {
	presentationRels := "ppt/_rels/presentation.xml.rels"
	if !v.provider.Has(presentationRels) {
		return
	}

	data, _ := v.provider.Get(presentationRels)
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
			fullPath := v.resolvePath(presentationRels, rel.Target)
			referencedSlides[fullPath] = true
		}
	}

	// Find actual slide files
	for _, p := range v.provider.Keys() {
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
				Path:        presentationRels,
				Description: fmt.Sprintf("Referenced slide not found: %s", slidePath),
				Repairable:  true,
			})
		}
	}
}

func (v *Validator) checkContentTypes() {
	ctPath := "[Content_Types].xml"
	data, ok := v.provider.Get(ctPath)
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

	for _, p := range v.provider.Keys() {
		if p == ctPath || strings.HasSuffix(p, ".rels") {
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

// ValidateZipIntegrity checks if the reader points to a valid zip archive.
func ValidateZipIntegrity(r io.ReaderAt, size int64) error {
	_, err := zip.NewReader(r, size)
	return err
}

func (v *Validator) checkNamespaces() {
	paths := v.provider.Keys()
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
	paths := v.provider.Keys()
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

func (v *Validator) checkNamespacesForPart(partPath string) []Issue {
	data, ok := v.provider.Get(partPath)
	if !ok {
		return nil
	}

	content := string(data)
	xmlDeclIdx := strings.Index(content, "?>")
	startIdx := xmlDeclIdx + 2
	if xmlDeclIdx == -1 {
		startIdx = 0
	}

	firstTagStart := strings.Index(content[startIdx:], "<")
	if firstTagStart == -1 {
		return nil
	}

	firstTagEnd := strings.Index(content[startIdx+firstTagStart:], ">")
	if firstTagEnd == -1 {
		return nil
	}

	openingTag := content[startIdx+firstTagStart : startIdx+firstTagStart+firstTagEnd+1]
	if !isPresentationNamespacePart(partPath) {
		return nil
	}

	issues := make([]Issue, 0, 3)
	if !strings.Contains(openingTag, "xmlns:p=") &&
		!strings.Contains(openingTag, `xmlns="http://schemas.openxmlformats.org/presentationml/2006/main"`) {
		issues = append(issues, Issue{
			Code:        CodeMissingNamespace,
			Severity:    SeverityWarning,
			Path:        partPath,
			Description: "Missing presentationml namespace declaration",
			Repairable:  true,
			Context:     map[string]string{"ns": "p"},
		})
	}
	if !strings.Contains(openingTag, "xmlns:a=") &&
		!strings.Contains(openingTag, `xmlns="http://schemas.openxmlformats.org/drawingml/2006/main"`) {
		issues = append(issues, Issue{
			Code:        CodeMissingNamespace,
			Severity:    SeverityWarning,
			Path:        partPath,
			Description: "Missing drawingml namespace declaration",
			Repairable:  true,
			Context:     map[string]string{"ns": "a"},
		})
	}
	if !strings.Contains(openingTag, "xmlns:r=") &&
		!strings.Contains(openingTag, `xmlns="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`) {
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

func (v *Validator) checkEmptyElementsForPart(partPath string) []Issue {
	data, ok := v.provider.Get(partPath)
	if !ok {
		return nil
	}
	content := string(data)
	if !strings.Contains(content, "<p:sldIdLst/>") {
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
