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

const (
	presentationRelsPath = "ppt/_rels/presentation.xml.rels"
	contentTypesPath     = "[Content_Types].xml"
	xmlDeclSuffixLength  = 2
	namespaceIssueCap    = 3
)

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
	keys     []string // computed once per Validate() call, shared across check functions
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

// requiredParts is the minimum set of parts for a valid PPTX.
//
//nolint:gochecknoglobals // read-only lookup table, never mutated
var requiredParts = map[string]string{
	contentTypesPath:       "Content types definition",
	"_rels/.rels":          "Package relationships",
	"ppt/presentation.xml": "Presentation document",
	presentationRelsPath:   "Presentation relationships",
}

// Validate performs a comprehensive validation check on the package.
func (v *Validator) Validate() []Issue {
	v.issues = nil
	v.keys = v.provider.Keys() // computed once; shared by all check functions below
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
	for p, desc := range requiredParts {
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
	paths := v.keys
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
	paths := v.keys
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

// ValidateZipIntegrity checks if the reader points to a valid zip archive.
func ValidateZipIntegrity(r io.ReaderAt, size int64) error {
	_, err := zip.NewReader(r, size)
	return err
}
