package editor

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestAddCommentPreservesExistingRelationshipTargetMode(t *testing.T) {
	ps := NewPartStore()
	ps.Set("ppt/slides/_rels/slide1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="https://example.com" TargetMode="External"/>
</Relationships>`))

	e := &PresentationEditor{
		parts:  ps,
		slides: []common.EditorSlideRef{{Part: "ppt/slides/slide1.xml"}},
	}

	author, err := e.AddAuthor("Alice", "AL")
	if err != nil {
		t.Fatalf("add author: %v", err)
	}
	if addErr := e.AddComment(0, author.ID, "hello", 100, 100); addErr != nil {
		t.Fatalf("add comment: %v", addErr)
	}

	rels, ok := e.parts.Get("ppt/slides/_rels/slide1.xml.rels")
	if !ok {
		t.Fatal("missing rewritten slide rels")
	}
	xml := string(rels)
	if !strings.Contains(xml, `TargetMode="External"`) {
		t.Fatalf("expected existing relationship TargetMode to be preserved, rels: %s", xml)
	}
}

func TestAddCommentParsesTimestampWithoutMilliseconds(t *testing.T) {
	ps := NewPartStore()
	ps.Set("ppt/slides/_rels/slide1.xml.rels", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments" Target="../comments/comment1.xml"/>
</Relationships>`))
	ps.Set("ppt/comments/comment1.xml", []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:cmLst xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cm authorId="1" dt="2026-02-13T10:00:00" idx="1">
    <p:pos x="10" y="20"/>
    <p:text>existing</p:text>
  </p:cm>
</p:cmLst>`))
	ps.Set("ppt/commentAuthors.xml", []byte(pptxxml.CommentAuthorsXML([]comments.Author{
		{ID: 1, Name: "Alice", Initials: "AL", LastIndex: 1, ColorIndex: 0},
	})))

	e := &PresentationEditor{
		parts:  ps,
		slides: []common.EditorSlideRef{{Part: "ppt/slides/slide1.xml"}},
	}

	if err := e.AddComment(0, 1, "new", 30, 40); err != nil {
		t.Fatalf("add comment with non-millisecond existing timestamp: %v", err)
	}

	updated, ok := e.parts.Get("ppt/comments/comment1.xml")
	if !ok {
		t.Fatal("missing updated comments part")
	}
	xml := string(updated)
	if strings.Contains(xml, `dt="0001-01-01T00:00:00.000"`) {
		t.Fatalf("unexpected zero timestamp rewrite: %s", xml)
	}
	if !strings.Contains(xml, "<p:text>existing</p:text>") {
		t.Fatalf("expected existing comment to persist after append: %s", xml)
	}
}
