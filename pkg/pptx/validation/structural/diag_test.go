package structural

import (
	"strings"
	"testing"
)

type mockPartStore struct {
	parts map[string][]byte
}

func (m *mockPartStore) Has(path string) bool {
	_, ok := m.parts[path]
	return ok
}

func (m *mockPartStore) Get(path string) ([]byte, bool) {
	data, ok := m.parts[path]
	return data, ok
}

func (m *mockPartStore) Keys() []string {
	keys := make([]string, 0, len(m.parts))
	for k := range m.parts {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockPartStore) Set(path string, data []byte) {
	m.parts[path] = data
}

func (m *mockPartStore) Delete(path string) {
	delete(m.parts, path)
}

func TestValidator_MissingParts(t *testing.T) {
	m := &mockPartStore{parts: make(map[string][]byte)}
	v := NewValidator(m)
	issues := v.Validate()

	missingCount := 0
	for _, issue := range issues {
		if issue.Code == CodeMissingPart {
			missingCount++
		}
	}
	parts := requiredParts()
	if missingCount != len(parts) {
		t.Errorf("expected %d missing part issues, got %d", len(parts), missingCount)
	}
}

func TestValidator_InvalidXml(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"[Content_Types].xml": []byte("<invalid"),
	}}
	v := NewValidator(m)
	issues := v.Validate()

	found := false
	for _, issue := range issues {
		if issue.Code == CodeInvalidXML && issue.Path == "[Content_Types].xml" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected invalid XML issue for [Content_Types].xml")
	}
}

func TestRepairer_BasicRepair(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/slides/slide1.xml": []byte("Hello & World"), // Invalid XML (bare &)
	}}
	v := NewValidator(m)
	issues := v.Validate()

	// Initial validation should find missing parts and invalid XML
	r := NewRepairer(m)
	result := r.Repair(issues)

	if len(result.IssuesRepaired) == 0 {
		t.Error("expected issues to be repaired")
	}

	// Re-validate
	issues2 := v.Validate()
	for _, issue := range issues2 {
		if issue.Severity == SeverityError {
			t.Errorf("found unrepaired error: %v", issue)
		}
	}

	// Check if bare ampersand was fixed
	data, _ := m.Get("ppt/slides/slide1.xml")
	if !strings.Contains(string(data), "&amp;") {
		t.Errorf("bare ampersand was not repaired: %s", string(data))
	}
}

func TestValidator_BrokenRelationships(t *testing.T) {
	rels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://.../slide" Target="slides/slide1.xml"/>
  <Relationship Id="rId2" Type="http://.../slide" Target="slides/missing.xml"/>
</Relationships>`

	m := &mockPartStore{parts: map[string][]byte{
		"ppt/_rels/presentation.xml.rels": []byte(rels),
		"ppt/slides/slide1.xml":           []byte("<root/>"),
	}}
	v := NewValidator(m)
	issues := v.Validate()

	found := false
	for _, issue := range issues {
		if issue.Code == CodeBrokenRelationship && strings.Contains(issue.Description, "missing.xml") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected broken relationship issue for missing.xml")
	}
}

func TestRepairer_BrokenRelationshipPreservesOOXMLRelationshipsRoot(t *testing.T) {
	rels := `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://.../slide" Target="slides/slide1.xml"/>
  <Relationship Id="rId2" Type="http://.../slide" Target="slides/missing.xml"/>
</Relationships>`

	m := &mockPartStore{parts: map[string][]byte{
		"ppt/_rels/presentation.xml.rels": []byte(rels),
		"ppt/slides/slide1.xml":           []byte("<root/>"),
	}}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:        CodeBrokenRelationship,
		Path:        "ppt/_rels/presentation.xml.rels",
		Description: "Broken relationship: ppt/_rels/presentation.xml.rels -> slides/missing.xml (ID: rId2)",
		Repairable:  true,
	}})
	if len(result.IssuesUnrepaired) > 0 {
		t.Fatalf("expected issue to be repaired, got unrepaired: %+v", result.IssuesUnrepaired)
	}

	data, ok := m.Get("ppt/_rels/presentation.xml.rels")
	if !ok {
		t.Fatal("expected repaired relationships part to exist")
	}
	content := string(data)
	if strings.Contains(content, "<relationshipsXML>") {
		t.Fatalf("unexpected Go struct root found in relationships XML: %s", content)
	}
	if !strings.Contains(content, "<Relationships") {
		t.Fatalf("expected OOXML Relationships root, got: %s", content)
	}
	if !strings.Contains(content, `xmlns="`+packageRelationshipsXMLNS+`"`) {
		t.Fatalf("expected OOXML relationships namespace, got: %s", content)
	}
}

func TestRepairer_MissingSlideRef(t *testing.T) {
	presRels := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
</Relationships>`

	presXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation>
  <p:sldMasterIdLst><p:sldMasterId id="2147483648" r:id="rId1"/></p:sldMasterIdLst>
  <p:sldIdLst/>
</p:presentation>`

	m := &mockPartStore{parts: map[string][]byte{
		"ppt/presentation.xml":            []byte(presXml),
		"ppt/_rels/presentation.xml.rels": []byte(presRels),
		"ppt/slides/slide1.xml":           []byte("<root/>"), // Orphan/Missing ref
	}}

	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:        CodeMissingSlideRef,
		Path:        "ppt/slides/slide1.xml",
		Description: "Missing slide ref",
		Repairable:  true,
	}})

	if len(result.IssuesUnrepaired) > 0 {
		t.Fatalf("expected issue to be repaired, got unrepaired")
	}

	relsData, _ := m.Get("ppt/_rels/presentation.xml.rels")
	if !strings.Contains(string(relsData), `Target="slides/slide1.xml"`) {
		t.Errorf("expected slides/slide1.xml added to presentation.xml.rels, got %s", relsData)
	}
	if !strings.Contains(string(relsData), `Id="rId2"`) {
		t.Errorf("expected rId2 added to presentation.xml.rels, got %s", relsData)
	}

	presData, _ := m.Get("ppt/presentation.xml")
	if !strings.Contains(string(presData), `<p:sldId id="256" r:id="rId2"/>`) {
		t.Errorf("expected sldId added to presentation.xml, got %s", presData)
	}
}

func TestRepairer_MissingNamespace(t *testing.T) {
	slideXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sld xmlns:a="http://drawingml">
  <p:cSld/>
</p:sld>`

	m := &mockPartStore{parts: map[string][]byte{
		"ppt/slides/slide1.xml": []byte(slideXml),
	}}

	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:        CodeMissingNamespace,
		Path:        "ppt/slides/slide1.xml",
		Description: "Missing namespace p",
		Repairable:  true,
		Context:     map[string]string{"ns": "p"},
	}})

	if len(result.IssuesUnrepaired) > 0 {
		t.Fatalf("expected issue to be repaired")
	}

	data, _ := m.Get("ppt/slides/slide1.xml")
	content := string(data)
	if !strings.Contains(content, `xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"`) {
		t.Errorf("expected p namespace added to root element, got: %s", content)
	}
}

func TestRepairer_EmptyRequiredElement(t *testing.T) {
	presXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation>
  <p:sldIdLst/>
</p:presentation>`

	m := &mockPartStore{parts: map[string][]byte{
		"ppt/presentation.xml": []byte(presXml),
	}}

	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:        CodeEmptyRequiredElement,
		Path:        "ppt/presentation.xml",
		Description: "Empty sldIdLst",
		Repairable:  true,
		Context:     map[string]string{"element": "p:sldIdLst"},
	}})

	if len(result.IssuesUnrepaired) > 0 {
		t.Fatalf("expected issue to be repaired")
	}

	data, _ := m.Get("ppt/presentation.xml")
	content := string(data)
	if !strings.Contains(content, `<p:sldIdLst></p:sldIdLst>`) {
		t.Errorf("expected empty element expanded, got: %s", content)
	}
}
