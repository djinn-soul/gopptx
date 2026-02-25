package structural

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"
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

// RequiredParts identifies the minimum set of parts for a valid PPTX.
var RequiredParts = map[string]string{
	"[Content_Types].xml":             "Content types definition",
	"_rels/.rels":                     "Package relationships",
	"ppt/presentation.xml":            "Presentation document",
	"ppt/_rels/presentation.xml.rels": "Presentation relationships",
}

// Validate performs a comprehensive validation check on the package.
func (v *Validator) Validate() []Issue {
	v.issues = nil
	v.checkRequiredParts()
	v.checkXmlValidity()
	v.checkRelationships()
	v.checkSlideReferences()
	v.checkContentTypes()

	for _, c := range v.checkers {
		v.issues = append(v.issues, c.Check(v.provider)...)
	}

	return v.issues
}

func (v *Validator) checkRequiredParts() {
	for p, desc := range RequiredParts {
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

func (v *Validator) checkXmlValidity() {
	for _, p := range v.provider.Keys() {
		if !strings.HasSuffix(p, ".xml") && !strings.HasSuffix(p, ".rels") {
			continue
		}
		data, ok := v.provider.Get(p)
		if !ok {
			continue
		}

		if err := v.validateXml(data); err != nil {
			v.issues = append(v.issues, Issue{
				Code:        CodeInvalidXml,
				Severity:    SeverityError,
				Path:        p,
				Description: fmt.Sprintf("Invalid XML: %v", err),
				Repairable:  true,
			})
		}
	}
}

func (v *Validator) validateXml(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty part")
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) checkRelationships() {
	for _, p := range v.provider.Keys() {
		if !strings.HasSuffix(p, ".rels") {
			continue
		}
		v.checkRelsFile(p)
	}
}

// relationshipXML represents a Relationship element in .rels files
type relationshipXML struct {
	ID     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

// relationshipsXML represents the root Relationships element
type relationshipsXML struct {
	Relationships []relationshipXML `xml:"Relationship"`
}

func (v *Validator) checkRelsFile(relsPath string) {
	data, ok := v.provider.Get(relsPath)
	if !ok {
		return
	}

	// Parse relationships using proper XML decoder
	var rels relationshipsXML
	decoder := xml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&rels); err != nil {
		// If XML parsing fails, the XML validity check will report it
		return
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
			v.issues = append(v.issues, Issue{
				Code:        CodeBrokenRelationship,
				Severity:    SeverityWarning,
				Path:        relsPath,
				Description: fmt.Sprintf("Broken relationship: %s -> %s (ID: %s)", relsPath, target, rel.ID),
				Repairable:  true,
			})
		}
	}
}

func (v *Validator) resolvePath(relsPath, target string) string {
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(target, "/")
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

	ctStr := string(data)
	for _, p := range v.provider.Keys() {
		if p == ctPath || strings.HasSuffix(p, ".rels") {
			continue
		}

		partName := "/" + p
		// Search for exactly PartName="/path" with either single or double quotes
		hasOverride := strings.Contains(ctStr, `PartName="`+partName+`"`) ||
			strings.Contains(ctStr, `PartName='`+partName+`'`)

		if !hasOverride {
			// Check if extension has a default
			ext := path.Ext(p)
			if ext != "" {
				ext = strings.TrimPrefix(ext, ".")
				if strings.Contains(ctStr, `Extension="`+ext+`"`) {
					continue
				}
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
}

// ValidateZipIntegrity checks if the reader points to a valid zip archive.
func ValidateZipIntegrity(r io.ReaderAt, size int64) error {
	_, err := zip.NewReader(r, size)
	return err
}
