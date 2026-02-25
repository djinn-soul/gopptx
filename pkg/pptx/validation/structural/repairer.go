package structural

import (
	"encoding/xml"
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
	case CodeInvalidXml:
		return r.repairInvalidXml(issue.Path)
	case CodeBrokenRelationship:
		return r.repairBrokenRelationship(issue.Path, issue.Description)
	case CodeOrphanSlide:
		return r.repairOrphanSlide(issue.Path)
	case CodeInvalidContentType:
		return r.repairInvalidContentType(issue.Path)
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

func (r *Repairer) repairInvalidXml(p string) error {
	data, ok := r.modifier.Get(p)
	if !ok {
		return fmt.Errorf("part not found: %s", p)
	}

	// Simple heuristic repair: ensure XML declaration and try to escape bare ampersands
	content := string(data)
	if !strings.HasPrefix(strings.TrimSpace(content), "<?xml") {
		content = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + content
	}

	// Escape bare ampersands - Go's regexp doesn't support negative lookahead,
	// so we use a multi-pass approach:
	// 1. Replace all existing valid entities with placeholders
	// 2. Escape remaining bare ampersands
	// 3. Restore the placeholders back to entities
	entityPattern := regexp.MustCompile(`&(amp|lt|gt|quot|apos|#\d+|#x[0-9a-fA-F]+);`)
	placeholder := "\x00ENTITY\x00"
	// Find all entities and replace with placeholders
	entities := entityPattern.FindAllString(content, -1)
	temp := entityPattern.ReplaceAllString(content, placeholder)
	// Replace bare ampersands
	temp = strings.ReplaceAll(temp, "&", "&amp;")
	// Restore entities
	for _, entity := range entities {
		temp = strings.Replace(temp, placeholder, entity, 1)
	}
	repaired := temp

	// Validate before setting
	decoder := xml.NewDecoder(strings.NewReader(repaired))
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML repair failed to produce valid XML: %v", err)
		}
	}

	r.modifier.Set(p, []byte(repaired))
	return nil
}

func (r *Repairer) repairBrokenRelationship(p, desc string) error {
	data, ok := r.modifier.Get(p)
	if !ok {
		return nil
	}

	// Extract target from description
	parts := strings.Split(desc, " -> ")
	if len(parts) < 2 {
		return nil
	}
	targetPart := strings.Split(parts[1], " (ID:")[0]

	// Proper fix would be to use encoding/xml to decode, filter, and re-encode.
	// For now, we use a more robust regex-based removal that isn't dependent on line breaks.
	content := string(data)
	relPattern := regexp.MustCompile(`(?s)<Relationship\s+[^>]*?Target="` + regexp.QuoteMeta(targetPart) + `"[^>]*?/>`)
	repaired := relPattern.ReplaceAllString(content, "")

	r.modifier.Set(p, []byte(repaired))
	return nil
}

func (r *Repairer) repairOrphanSlide(p string) error {
	// Either remove the orphaned slide file (pruning) or add it back to presentation.xml
	// For simplicity, we'll prune it if it's truly unreferenced.
	r.modifier.Delete(p)
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

	content := strings.Replace(string(data), "</Types>", newEntry+"\n</Types>", 1)
	r.modifier.Set(ctPath, []byte(content))
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
	sb.WriteString(`  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`)
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
