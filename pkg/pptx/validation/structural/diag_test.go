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
	if missingCount != len(RequiredParts) {
		t.Errorf("expected %d missing part issues, got %d", len(RequiredParts), missingCount)
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
		if issue.Code == CodeInvalidXml && issue.Path == "[Content_Types].xml" {
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
