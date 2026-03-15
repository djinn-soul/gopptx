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

func TestRepairer_MissingSlideRef_NonSlidePathIsNotMarkedRepaired(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/presentation.xml":            []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><p:presentation><p:sldMasterIdLst/><p:sldIdLst/></p:presentation>`),
		"ppt/_rels/presentation.xml.rels": []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`),
	}}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:        CodeMissingSlideRef,
		Path:        "ppt/_rels/presentation.xml.rels",
		Description: "missing slide ref",
		Repairable:  true,
	}})
	if len(result.IssuesRepaired) != 0 {
		t.Fatalf("expected non-slide missing-ref issue to remain unrepaired")
	}
	if len(result.IssuesUnrepaired) != 1 {
		t.Fatalf("expected one unrepaired issue, got repaired=%d unrepaired=%d", len(result.IssuesRepaired), len(result.IssuesUnrepaired))
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

func TestIssue_String(t *testing.T) {
	s := SeverityInfo
	if s.String() != "INFO" {
		t.Errorf("expected INFO, got %s", s.String())
	}
	s = SeverityWarning
	if s.String() != "WARNING" {
		t.Errorf("expected WARNING, got %s", s.String())
	}
	s = SeverityError
	if s.String() != "ERROR" {
		t.Errorf("expected ERROR, got %s", s.String())
	}
	s = Severity(99)
	if s.String() != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %s", s.String())
	}

	issue := Issue{
		Code:        CodeMissingPart,
		Severity:    SeverityError,
		Path:        "test.xml",
		Description: "test description",
		Repairable:  true,
	}
	expected := "[ERROR] MISSING_PART (test.xml): test description (Repairable: true)"
	if issue.String() != expected {
		t.Errorf("expected %s, got %s", expected, issue.String())
	}
}

type mockChecker struct {
	called bool
}

func (m *mockChecker) Check(ps PartProvider) []Issue {
	m.called = true
	return []Issue{{Code: "MOCK_ISSUE"}}
}

func TestValidator_AddChecker(t *testing.T) {
	m := &mockPartStore{parts: make(map[string][]byte)}
	v := NewValidator(m)
	checker := &mockChecker{}
	v.AddChecker(checker)

	issues := v.Validate()
	if !checker.called {
		t.Error("expected custom checker to be called")
	}
	found := false
	for _, issue := range issues {
		if issue.Code == "MOCK_ISSUE" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected mock issue from custom checker")
	}
}

func TestValidateZipIntegrity(t *testing.T) {
	r := strings.NewReader("not a zip file")
	err := ValidateZipIntegrity(r, int64(r.Len()))
	if err == nil {
		t.Error("expected error for invalid zip data")
	}
}

func TestRepairer_OrphanSlide(t *testing.T) {
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/slides/slide1.xml": []byte("<root/>"),
	}}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:       CodeOrphanSlide,
		Path:       "ppt/slides/slide1.xml",
		Repairable: true,
	}})
	if len(result.IssuesUnrepaired) > 0 {
		t.Errorf("expected orphan slide issue to be repaired")
	}
	if m.Has("ppt/slides/slide1.xml") {
		t.Errorf("expected orphan slide to be deleted")
	}
}

func TestRepairer_InvalidContentType(t *testing.T) {
	ctXml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`
	m := &mockPartStore{parts: map[string][]byte{
		"[Content_Types].xml":               []byte(ctXml),
		"ppt/slides/slide1.xml":             []byte("<root/>"),
		"ppt/slideLayouts/slideLayout1.xml": []byte("<root/>"),
		"ppt/slideMasters/slideMaster1.xml": []byte("<root/>"),
		"ppt/presentation.xml":              []byte("<root/>"),
		"docProps/core.xml":                 []byte("<root/>"),
		"image.png":                         []byte("image data"),
	}}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:       CodeInvalidContentType,
		Path:       "ppt/slides/slide1.xml",
		Repairable: true,
	}, {
		Code:       CodeInvalidContentType,
		Path:       "ppt/slideLayouts/slideLayout1.xml",
		Repairable: true,
	}, {
		Code:       CodeInvalidContentType,
		Path:       "ppt/slideMasters/slideMaster1.xml",
		Repairable: true,
	}, {
		Code:       CodeInvalidContentType,
		Path:       "ppt/presentation.xml",
		Repairable: true,
	}, {
		Code:       CodeInvalidContentType,
		Path:       "docProps/core.xml",
		Repairable: true,
	}, {
		Code:       CodeInvalidContentType,
		Path:       "image.png",
		Repairable: true,
	}})

	if len(result.IssuesUnrepaired) > 0 {
		t.Errorf("expected issues to be repaired, got %d unrepaired", len(result.IssuesUnrepaired))
	}

	data, _ := m.Get("[Content_Types].xml")
	content := string(data)

	failed := false
	if !strings.Contains(content, "presentationml.slide+xml") {
		t.Errorf("missing slide content type")
		failed = true
	}
	if !strings.Contains(content, "presentationml.slideLayout+xml") {
		t.Errorf("missing layout content type")
		failed = true
	}
	if !strings.Contains(content, "presentationml.slideMaster+xml") {
		t.Errorf("missing master content type")
		failed = true
	}
	if !strings.Contains(content, "presentationml.presentation.main+xml") {
		t.Errorf("missing presentation content type")
		failed = true
	}
	if !strings.Contains(content, "application/xml") {
		t.Errorf("missing xml content type")
		failed = true
	}
	if !strings.Contains(content, "application/octet-stream") {
		t.Errorf("missing default content type")
		failed = true
	}

	if failed {
		t.Logf("Generated content:\n%s", content)
	}

	m.Delete("[Content_Types].xml")
	result = r.Repair([]Issue{{
		Code:       CodeInvalidContentType,
		Path:       "ppt/slides/slide1.xml",
		Repairable: true,
	}})
	if len(result.IssuesUnrepaired) > 0 {
		t.Errorf("expected issue to be repaired when [Content_Types].xml is missing")
	}
	if !m.Has("[Content_Types].xml") {
		t.Errorf("expected [Content_Types].xml to be generated")
	}
}

func TestRepairer_EscapeBareAmpersands(t *testing.T) {
	input := "Here is a & and an encoded &amp; and &gt; but wait & more."
	expected := "Here is a &amp; and an encoded &amp; and &gt; but wait &amp; more."
	escaped := escapeBareAmpersands(input)
	if escaped != expected {
		t.Errorf("expected %s, got %s", expected, escaped)
	}
}

func TestRepairer_BrokenRelationshipFallback(t *testing.T) {
	rels := `<?xml version="1.0" encoding="UTF-8"?><InvalidXML`
	m := &mockPartStore{parts: map[string][]byte{
		"ppt/_rels/presentation.xml.rels": []byte(rels),
	}}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:       CodeBrokenRelationship,
		Path:       "ppt/_rels/presentation.xml.rels",
		Context:    map[string]string{"target": "slides/missing.xml"},
		Repairable: true,
	}})
	if len(result.IssuesUnrepaired) > 0 {
		t.Errorf("expected broken relationship to be repaired using fallback")
	}
}

func TestRepairer_UnsupportedRepair(t *testing.T) {
	m := &mockPartStore{parts: make(map[string][]byte)}
	r := NewRepairer(m)
	result := r.Repair([]Issue{{
		Code:       CodeModelValidationError,
		Repairable: true,
	}, {
		Code:       CodeMissingPart,
		Repairable: false,
	}})

	if len(result.IssuesUnrepaired) != 2 {
		t.Errorf("expected 2 unrepaired issues")
	}
}
