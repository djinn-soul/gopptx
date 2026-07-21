package structural

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var slideIDPattern = regexp.MustCompile(`<p:sldId\s+id="(\d+)"`)

func (r *Repairer) repairBrokenRelationship(issue Issue) error {
	data, ok := r.modifier.Get(issue.Path)
	if !ok {
		return nil
	}

	targetPart := issue.Context["target"]
	if targetPart == "" {
		return nil
	}

	var rels relationshipsXML
	if !tryUnmarshalRelationships(data, &rels) {
		content := string(data)
		relPattern := regexp.MustCompile(
			`(?s)<Relationship\s+[^>]*?Target="` + regexp.QuoteMeta(targetPart) + `"[^>]*?/>`,
		)
		repaired := relPattern.ReplaceAllString(content, "")
		r.modifier.Set(issue.Path, []byte(repaired))
		return nil
	}

	filtered := make([]relationshipXML, 0, len(rels.Relationships))
	for _, rel := range rels.Relationships {
		if rel.Target != targetPart {
			filtered = append(filtered, rel)
		}
	}
	rels.Relationships = filtered
	rels.XMLNS = packageRelationshipsXMLNS

	repaired, err := xml.Marshal(&rels)
	if err != nil {
		return fmt.Errorf("failed to re-encode relationships: %w", err)
	}

	repairedStr := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(repaired)
	r.modifier.Set(issue.Path, []byte(repairedStr))
	return nil
}

func (r *Repairer) repairOrphanSlide(p string) error {
	r.modifier.Delete(p)
	return nil
}

func (r *Repairer) repairMissingSlideRef(issue Issue) error {
	slidePath := issue.Path
	if !strings.HasPrefix(slidePath, "ppt/slides/") {
		if candidate := strings.TrimSpace(issue.Context["slide_part"]); strings.HasPrefix(candidate, "ppt/slides/") {
			slidePath = candidate
		} else {
			return fmt.Errorf("missing slide ref requires a slide part path, got %q", issue.Path)
		}
	}

	relsPath := presentationRelsPath
	relsData, ok := r.modifier.Get(relsPath)
	if !ok {
		return errors.New("cannot repair missing slide ref without presentation.xml.rels")
	}

	var rels relationshipsXML
	if !tryUnmarshalRelationships(relsData, &rels) {
		return errors.New("invalid presentation.xml.rels")
	}

	newRelID := nextRelationshipID(rels.Relationships)
	targetPath := strings.TrimPrefix(slidePath, "ppt/")
	rels.Relationships = append(rels.Relationships, relationshipXML{
		ID:     newRelID,
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide",
		Target: targetPath,
	})
	rels.XMLNS = packageRelationshipsXMLNS

	repairedRels, err := xml.Marshal(&rels)
	if err != nil {
		return fmt.Errorf("failed to encode updated presentation relationships: %w", err)
	}
	repairedRelsStr := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n" + string(repairedRels)
	r.modifier.Set(relsPath, []byte(repairedRelsStr))

	presPath := presentationPartPath
	presData, ok := r.modifier.Get(presPath)
	if !ok {
		return errors.New("cannot repair slide ref without presentation.xml")
	}
	presContent := string(presData)

	newSlideID := nextSlideID(presContent)
	newSlideElement := fmt.Sprintf(`<p:sldId id="%d" r:id="%s"/>`, newSlideID, newRelID)
	repairedContent := addSlideIDEntry(presContent, newSlideElement)
	r.modifier.Set(presPath, []byte(repairedContent))
	return nil
}

func nextRelationshipID(rels []relationshipXML) string {
	maxID := 0
	for _, rel := range rels {
		var id int
		if _, err := fmt.Sscanf(rel.ID, "rId%d", &id); err == nil && id > maxID {
			maxID = id
		}
	}
	return fmt.Sprintf("rId%d", maxID+1)
}

func nextSlideID(presentationXML string) uint32 {
	var maxSlideID uint32 = 255
	for _, match := range slideIDPattern.FindAllStringSubmatch(presentationXML, -1) {
		var id uint32
		if _, err := fmt.Sscanf(match[1], "%d", &id); err == nil && id > maxSlideID {
			maxSlideID = id
		}
	}
	return maxSlideID + 1
}

func addSlideIDEntry(presentationXML, slideElement string) string {
	if strings.Contains(presentationXML, "</p:sldIdLst>") {
		return strings.Replace(presentationXML, "</p:sldIdLst>", slideElement+"\n</p:sldIdLst>", 1)
	}
	if strings.Contains(presentationXML, "<p:sldIdLst/>") {
		return strings.Replace(
			presentationXML,
			"<p:sldIdLst/>",
			fmt.Sprintf("<p:sldIdLst>%s</p:sldIdLst>", slideElement),
			1,
		)
	}
	return strings.Replace(
		presentationXML,
		"</p:sldMasterIdLst>",
		"</p:sldMasterIdLst>\n<p:sldIdLst>"+slideElement+"</p:sldIdLst>",
		1,
	)
}
