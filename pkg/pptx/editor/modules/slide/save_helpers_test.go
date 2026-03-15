package slide

import (
	"archive/zip"
	"errors"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

func TestGenerateCustomXMLItem(t *testing.T) {
	xml, err := GenerateCustomXMLItem(common.CustomXMLPart{
		RootElement: "root",
		Namespace:   "urn:test",
		Properties: []common.CustomXMLKV{
			{Key: "Name", Value: `A&B`},
			{Key: "Flag", Value: `1<2`},
		},
	})
	if err != nil {
		t.Fatalf("GenerateCustomXMLItem failed: %v", err)
	}
	if !strings.Contains(xml, `<root xmlns="urn:test">`) {
		t.Fatalf("expected namespace in custom xml: %s", xml)
	}
	if !strings.Contains(xml, `<Name>A&amp;B</Name>`) || !strings.Contains(xml, `<Flag>1&lt;2</Flag>`) {
		t.Fatalf("expected escaped property values in custom xml: %s", xml)
	}

	xml, err = GenerateCustomXMLItem(common.CustomXMLPart{
		RootElement: "root",
		Content:     "<x>kept-raw</x>",
	})
	if err != nil {
		t.Fatalf("GenerateCustomXMLItem with direct content failed: %v", err)
	}
	if !strings.Contains(xml, "<x>kept-raw</x>") {
		t.Fatalf("expected direct content passthrough: %s", xml)
	}

	_, err = GenerateCustomXMLItem(common.CustomXMLPart{RootElement: "1invalid"})
	if err == nil {
		t.Fatal("expected invalid root element error")
	}
	_, err = GenerateCustomXMLItem(common.CustomXMLPart{
		RootElement: "root",
		Properties:  []common.CustomXMLKV{{Key: "bad-key!", Value: "x"}},
	})
	if err == nil {
		t.Fatal("expected invalid property element name error")
	}
}

func TestVbaProjectAndMapHelpers(t *testing.T) {
	if proj, ok := VbaProjectFromMetadata(nil); ok || proj != nil {
		t.Fatalf("expected nil project metadata to be rejected: %+v ok=%v", proj, ok)
	}
	if proj, ok := VbaProjectFromMetadata(&vba.VBAProject{}); ok || proj == nil {
		t.Fatalf("expected macro-disabled project pointer with ok=false: %+v ok=%v", proj, ok)
	}
	enabled := &vba.VBAProject{Data: []byte{1, 2, 3}}
	if proj, ok := VbaProjectFromMetadata(enabled); !ok || proj != enabled {
		t.Fatalf("expected macro-enabled project to be accepted: %+v ok=%v", proj, ok)
	}

	values := MapValues(map[string]string{"a": "x", "b": "y"})
	if len(values) != 2 {
		t.Fatalf("unexpected map value count: %v", values)
	}

	filtered := FilterXMLPartPaths([]string{"ppt/slides/slide1.xml", "ppt/media/image1.png", "docProps/core.xml"})
	if len(filtered) != 2 {
		t.Fatalf("unexpected filtered xml paths: %v", filtered)
	}

	if got := SaveZipMethod("ppt/notesSlides/notesSlide1.xml"); got != zip.Store {
		t.Fatalf("expected zip.Store for notes parts, got %d", got)
	}
	if got := SaveZipMethod("ppt/media/image1.png"); got != zip.Store {
		t.Fatalf("expected zip.Store for image parts, got %d", got)
	}
	if got := SaveZipMethod("ppt/slides/slide1.xml"); got != zip.Deflate {
		t.Fatalf("expected zip.Deflate for xml slide parts, got %d", got)
	}
}

func TestCommentAuthorAndRelationshipHelpers(t *testing.T) {
	data, ok, err := SerializeCommentAuthorsIfPopulated(nil, func() ([]comments.Author, error) {
		return nil, errors.New("should not run")
	})
	if err != nil || ok || data != nil {
		t.Fatalf("nil cache should skip serialization: data=%v ok=%v err=%v", data, ok, err)
	}

	_, _, err = SerializeCommentAuthorsIfPopulated(map[int64]comments.Author{}, func() ([]comments.Author, error) {
		return nil, errors.New("boom")
	})
	if err == nil {
		t.Fatal("expected getAuthors error to be wrapped")
	}

	xmlData, ok, err := SerializeCommentAuthorsIfPopulated(
		map[int64]comments.Author{},
		func() ([]comments.Author, error) {
			return []comments.Author{
				{ID: 3, Name: "C"},
				{ID: 1, Name: "A"},
			}, nil
		},
	)
	if err != nil || !ok || len(xmlData) == 0 {
		t.Fatalf("expected serialized comment author XML: len=%d ok=%v err=%v", len(xmlData), ok, err)
	}
	xml := string(xmlData)
	if strings.Index(xml, `id="1"`) > strings.Index(xml, `id="3"`) {
		t.Fatalf("expected comment authors sorted by id: %s", xml)
	}

	rels := []common.EditorRelationship{{ID: "rId1", Type: common.RelTypeSlide}}
	updated, next := EnsureCommentAuthorsRelationship(
		rels,
		2,
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors",
		"commentAuthors.xml",
	)
	if len(updated) != 2 || next != 3 {
		t.Fatalf("expected new commentAuthors relationship: rels=%+v next=%d", updated, next)
	}
	unchanged, next2 := EnsureCommentAuthorsRelationship(
		updated,
		next,
		"http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors",
		"commentAuthors.xml",
	)
	if len(unchanged) != len(updated) || next2 != next {
		t.Fatalf("expected existing relationship to be reused: rels=%+v next=%d", unchanged, next2)
	}
}

func TestResolveNotesMasterRelIDAndEscape(t *testing.T) {
	relID, err := ResolveNotesMasterRelID(nil, false, common.RelTypeNotesMaster)
	if err != nil || relID != "" {
		t.Fatalf("expected no-op when notes master disabled: id=%q err=%v", relID, err)
	}

	rels := []common.EditorRelationship{
		{ID: "rId9", Type: common.RelTypeNotesMaster, Target: "notesMasters/notesMaster1.xml"},
	}
	relID, err = ResolveNotesMasterRelID(rels, true, common.RelTypeNotesMaster)
	if err != nil || relID != "rId9" {
		t.Fatalf("expected notes master rel id, got id=%q err=%v", relID, err)
	}
	_, err = ResolveNotesMasterRelID(nil, true, common.RelTypeNotesMaster)
	if err == nil {
		t.Fatal("expected error when notes master exists without relationship")
	}

	escaped := EscapeCustomXML(`a&b<c>"d"'e'`)
	want := "a&amp;b&lt;c&gt;&quot;d&quot;&apos;e&apos;"
	if escaped != want {
		t.Fatalf("EscapeCustomXML=%q, want %q", escaped, want)
	}
}
