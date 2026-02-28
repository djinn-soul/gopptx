package presentation

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

func TestMetadata_Basic(t *testing.T) {
	meta := Metadata{}
	meta.Title = "Test Presentation"
	meta.Subject = "Testing"
	meta.Creator = "GoPPTX"
	meta.Description = "A test presentation"

	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}

func TestMetadata_Sections(t *testing.T) {
	meta := Metadata{
		Sections: []Section{
			{Name: "Section 1", SlideIndices: []int{0}},
			{Name: "Section 2", SlideIndices: []int{1}},
		},
	}

	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1"),
		elements.NewSlide("Slide 2"),
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := WritePresentationPackage(zw, meta, slides, len(slides))
	if err != nil {
		t.Fatalf("WritePresentationPackage failed: %v", err)
	}
	zw.Close()

	validatePPTX(t, buf.Bytes())
}

func TestConvertSections_Errors(t *testing.T) {
	sections := []Section{
		{Name: "Invalid", SlideIndices: []int{5}},
	}
	_, err := convertSections(sections, 2)
	if err == nil {
		t.Error("expected error for out of bounds slide index")
	}
}

func TestGenerateGUID(t *testing.T) {
	guid1, err := generateGUID()
	if err != nil {
		t.Fatalf("generateGUID failed: %v", err)
	}
	guid2, err := generateGUID()
	if err != nil {
		t.Fatalf("generateGUID failed: %v", err)
	}
	if guid1 == guid2 {
		t.Error("expected unique GUIDs")
	}
	if len(guid1) != 38 { // {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
		t.Errorf("expected GUID length 38, got %d", len(guid1))
	}
}

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

func validatePPTX(t *testing.T, blob []byte) {
	r, err := zip.NewReader(bytes.NewReader(blob), int64(len(blob)))
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}

	m := &mockPartStore{parts: make(map[string][]byte)}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("failed to open file %s: %v", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("failed to read file %s: %v", f.Name, err)
		}
		m.parts[f.Name] = data
	}

	v := structural.NewValidator(m)
	issues := v.Validate()
	for _, issue := range issues {
		if issue.Severity == structural.SeverityError {
			t.Errorf("validation error: %v", issue)
		}
	}
}
