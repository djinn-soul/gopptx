package structural

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
)

// PartModifier defines the interface for modifying package parts.
type PartModifier interface {
	PartProvider
	Set(path string, data []byte)
	Delete(path string)
}

// Repairer provides methods for repairing detected diagnostic issues.
type Repairer struct {
	modifier PartModifier
}

// NewRepairer creates a new repairer using the given part modifier.
func NewRepairer(modifier PartModifier) *Repairer {
	return &Repairer{
		modifier: modifier,
	}
}

// RepairResult summarizes the outcome of a repair operation.
type RepairResult struct {
	IssuesRepaired   []Issue
	IssuesUnrepaired []Issue
}

// Repair attempts to fix the provided list of issues.
func (r *Repairer) Repair(issues []Issue) RepairResult {
	result := RepairResult{}
	for _, issue := range issues {
		if !issue.Repairable {
			result.IssuesUnrepaired = append(result.IssuesUnrepaired, issue)
			continue
		}

		if err := r.repairIssue(issue); err != nil {
			result.IssuesUnrepaired = append(result.IssuesUnrepaired, issue)
		} else {
			result.IssuesRepaired = append(result.IssuesRepaired, issue)
		}
	}
	return result
}

func (r *Repairer) repairIssue(issue Issue) error {
	switch issue.Code {
	case CodeMissingPart:
		return r.repairMissingPart(issue.Path)
	case CodeInvalidXML:
		return r.repairInvalidXML(issue.Path)
	case CodeBrokenRelationship:
		return r.repairBrokenRelationship(issue)
	case CodeOrphanSlide:
		return r.repairOrphanSlide(issue.Path)
	case CodeMissingSlideRef:
		return r.repairMissingSlideRef(issue)
	case CodeInvalidContentType:
		return r.repairInvalidContentType(issue.Path)
	case CodeMissingNamespace:
		return r.repairMissingNamespace(issue)
	case CodeEmptyRequiredElement:
		return r.repairEmptyRequiredElement(issue)
	default:
		return fmt.Errorf("unsupported repair: %s", issue.Code)
	}
}

func (r *Repairer) repairMissingPart(p string) error {
	var content string
	switch p {
	case "[Content_Types].xml":
		content = r.generateContentTypes()
	case "_rels/.rels":
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`
	case "ppt/presentation.xml":
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>
  <p:sldIdLst/>
  <p:sldSz cx="9144000" cy="6858000"/>
  <p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`
	case "ppt/_rels/presentation.xml.rels":
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
</Relationships>`
	default:
		return fmt.Errorf("cannot auto-generate part: %s", p)
	}
	r.modifier.Set(p, []byte(content))
	return nil
}

func (r *Repairer) repairInvalidXML(p string) error {
	data, ok := r.modifier.Get(p)
	if !ok {
		return fmt.Errorf("part not found: %s", p)
	}

	content := string(data)
	if !strings.HasPrefix(strings.TrimSpace(content), "<?xml") {
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + content
	}

	repaired := escapeBareAmpersands(content)

	// Validate before setting
	decoder := xml.NewDecoder(strings.NewReader(repaired))
	for {
		_, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("XML repair failed to produce valid XML: %w", err)
		}
	}

	r.modifier.Set(p, []byte(repaired))
	return nil
}

var (
	entityPattern = regexp.MustCompile(`&(amp|lt|gt|quot|apos|#\d+|#x[0-9a-fA-F]+);`)
)

func escapeBareAmpersands(s string) string {
	var result strings.Builder
	last := 0
	for _, match := range entityPattern.FindAllStringIndex(s, -1) {
		// Escape ampersands in the text before this entity
		result.WriteString(strings.ReplaceAll(s[last:match[0]], "&", "&amp;"))
		// Write the entity itself
		result.WriteString(s[match[0]:match[1]])
		last = match[1]
	}
	result.WriteString(strings.ReplaceAll(s[last:], "&", "&amp;"))
	return result.String()
}

func (r *Repairer) repairBrokenRelationship(issue Issue) error {
	data, ok := r.modifier.Get(issue.Path)
	if !ok {
		return nil
	}

	targetPart := issue.Context["target"]
	if targetPart == "" {
		return nil
	}

	// Use XML parsing for robust handling of different attribute orderings and whitespace.
	var rels relationshipsXML
	if !tryUnmarshalRelationships(data, &rels) {
		// If XML parsing fails, fall back to regex-based removal.
		content := string(data)
		relPattern := regexp.MustCompile(
			`(?s)<Relationship\s+[^>]*?Target="` + regexp.QuoteMeta(targetPart) + `"[^>]*?/>`,
		)
		repaired := relPattern.ReplaceAllString(content, "")
		r.modifier.Set(issue.Path, []byte(repaired))
		return nil
	}

	// Filter out relationships with the broken target
	filtered := make([]relationshipXML, 0, len(rels.Relationships))
	for _, rel := range rels.Relationships {
		if rel.Target != targetPart {
			filtered = append(filtered, rel)
		}
	}
	rels.Relationships = filtered
	rels.XMLNS = packageRelationshipsXMLNS

	// Re-encode the relationships
	repaired, err := xml.Marshal(&rels)
	if err != nil {
		return fmt.Errorf("failed to re-encode relationships: %w", err)
	}

	// Add XML header and fix namespace
	repairedStr := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(repaired)
	r.modifier.Set(issue.Path, []byte(repairedStr))
	return nil
}

func tryUnmarshalRelationships(data []byte, rels *relationshipsXML) bool {
	return xml.Unmarshal(data, rels) == nil
}

func (r *Repairer) repairOrphanSlide(p string) error {
	// Either remove the orphaned slide file (pruning) or add it back to presentation.xml.
	// We'll safely remove unreferenced files to avoid corrupting the presentation structure
	r.modifier.Delete(p)
	return nil
}

func (r *Repairer) repairMissingSlideRef(issue Issue) error {
	// A slide file exists but is missing from presentation.xml / presentation.xml.rels
	slidePath := issue.Path
	if !strings.HasPrefix(slidePath, "ppt/slides/") {
		return nil
	}

	relsPath := "ppt/_rels/presentation.xml.rels"
	relsData, ok := r.modifier.Get(relsPath)
	if !ok {
		return errors.New("cannot repair missing slide ref without presentation.xml.rels")
	}

	var rels relationshipsXML
	if !tryUnmarshalRelationships(relsData, &rels) {
		return errors.New("invalid presentation.xml.rels")
	}

	// Generate a new unique rId
	maxID := 0
	for _, rel := range rels.Relationships {
		var id int
		if _, err := fmt.Sscanf(rel.ID, "rId%d", &id); err == nil && id > maxID {
			maxID = id
		}
	}
	newRId := fmt.Sprintf("rId%d", maxID+1)

	// Add to rels
	targetPath := strings.TrimPrefix(slidePath, "ppt/")
	rels.Relationships = append(rels.Relationships, relationshipXML{
		ID:     newRId,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide",
		Target: targetPath,
	})
	rels.XMLNS = packageRelationshipsXMLNS

	repairedRels, _ := xml.Marshal(&rels)
	repairedRelsStr := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(repairedRels)
	r.modifier.Set(relsPath, []byte(repairedRelsStr))

	// Now add it to presentation.xml
	presPath := "ppt/presentation.xml"
	presData, ok := r.modifier.Get(presPath)
	if !ok {
		return errors.New("cannot repair slide ref without presentation.xml")
	}
	presContent := string(presData)

	// Find the highest slide id
	var maxSldId uint32 = 255
	idRe := regexp.MustCompile(`<p:sldId\s+id="(\d+)"`)
	for _, match := range idRe.FindAllStringSubmatch(presContent, -1) {
		var id uint32
		if _, err := fmt.Sscanf(match[1], "%d", &id); err == nil && id > maxSldId {
			maxSldId = id
		}
	}
	newSldId := maxSldId + 1

	newSldElement := fmt.Sprintf(`<p:sldId id="%d" r:id="%s"/>`, newSldId, newRId)

	if strings.Contains(presContent, "</p:sldIdLst>") {
		presContent = strings.Replace(presContent, "</p:sldIdLst>", newSldElement+"\n</p:sldIdLst>", 1)
	} else if strings.Contains(presContent, "<p:sldIdLst/>") {
		presContent = strings.Replace(presContent, "<p:sldIdLst/>", fmt.Sprintf("<p:sldIdLst>%s</p:sldIdLst>", newSldElement), 1)
	} else {
		// Attempt to insert after sldMasterIdLst
		presContent = strings.Replace(presContent, "</p:sldMasterIdLst>", "</p:sldMasterIdLst>\n<p:sldIdLst>"+newSldElement+"</p:sldIdLst>", 1)
	}

	r.modifier.Set(presPath, []byte(presContent))
	return nil
}

func (r *Repairer) repairInvalidContentType(p string) error {
	ctPath := "[Content_Types].xml"
	data, ok := r.modifier.Get(ctPath)
	if !ok {
		return r.repairMissingPart(ctPath)
	}

	ext := path.Ext(p)
	ct := r.inferContentType(p)

	newEntry := ""
	if ext != "" {
		ext = strings.TrimPrefix(ext, ".")
		// Check if extension default is missing or if we should use Override
		if !strings.Contains(string(data), fmt.Sprintf(`Extension="%s"`, ext)) {
			newEntry = fmt.Sprintf("\n  <Default Extension=\"%s\" ContentType=\"%s\"/>", ext, ct)
		}
	}

	if newEntry == "" {
		newEntry = fmt.Sprintf("\n  <Override PartName=\"/%s\" ContentType=\"%s\"/>", p, ct)
	}

	// Use a slightly more robust replacement that handles potential whitespace
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
	if strings.Contains(p, "slide") && strings.HasSuffix(p, ".xml") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"
	}
	if strings.Contains(p, "slideLayout") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"
	}
	if strings.Contains(p, "slideMaster") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"
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

	// Add overrides for existing parts
	for _, p := range r.modifier.Keys() {
		if strings.HasSuffix(p, ".rels") || p == "[Content_Types].xml" {
			continue
		}
		ct := r.inferContentType(p)
		sb.WriteString(fmt.Sprintf("  <Override PartName=\"/%s\" ContentType=\"%s\"/>\n", p, ct))
	}

	sb.WriteString(`</Types>`)
	return sb.String()
}

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
	
	nsDecl := ""
	switch ns {
	case "p":
		nsDecl = `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`
	case "a":
		nsDecl = `xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"`
	case "r":
		nsDecl = `xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`
	default:
		return nil
	}
	
	// Find the first '>' which should be closing the root or XML decl
	// We want to skip the <?xml ... ?> if it exists.
	xmlDeclIdx := strings.Index(content, "?>")
	searchStart := 0
	if xmlDeclIdx != -1 {
		searchStart = xmlDeclIdx + 2
	}
	
	firstBracket := strings.Index(content[searchStart:], "<")
	if firstBracket == -1 {
		return nil
	}
	rootStart := searchStart + firstBracket
	
	// Find the end of this opening tag
	firstClose := strings.Index(content[rootStart:], ">")
	if firstClose == -1 {
		return nil
	}
	
	// Check if namespace already exists here just in case
	tagContent := content[rootStart : rootStart+firstClose]
	if !strings.Contains(tagContent, nsDecl) {
		// Insert right before the closing >
		repaired := content[:rootStart+firstClose] + " " + nsDecl + content[rootStart+firstClose:]
		r.modifier.Set(issue.Path, []byte(repaired))
	}
	
	return nil
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
