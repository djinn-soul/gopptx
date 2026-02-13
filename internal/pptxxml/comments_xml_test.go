package pptxxml_test

import (
	"strings"
	"testing"
	"time"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
)

func TestCommentAuthorsXML(t *testing.T) {
	authors := []comments.Author{
		{ID: 1, Name: "User One", Initials: "UO", LastIndex: 5, ColorIndex: 0},
		{ID: 2, Name: "User Two", Initials: "UT", LastIndex: 10, ColorIndex: 1},
	}
	xml := pptxxml.CommentAuthorsXML(authors)

	if !strings.Contains(xml, `<p:cmAuthor id="1" name="User One" initials="UO" lastIdx="5" clrIdx="0"/>`) {
		t.Errorf("expected author 1 xml, got: %s", xml)
	}
	if !strings.Contains(xml, `<p:cmAuthor id="2" name="User Two" initials="UT" lastIdx="10" clrIdx="1"/>`) {
		t.Errorf("expected author 2 xml, got: %s", xml)
	}
}

func TestCommentsXML(t *testing.T) {
	date := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)
	slideComments := []comments.Comment{
		{AuthorID: 1, Text: "First comment", Date: date, X: 100, Y: 200, Index: 1},
		{AuthorID: 2, Text: "Second comment", Date: date.Add(time.Hour), X: 300, Y: 400, Index: 2},
	}
	xml := pptxxml.CommentsXML(slideComments)

	expectedDate := "2023-10-27T10:00:00.000"
	if !strings.Contains(xml, `dt="`+expectedDate+`"`) {
		t.Errorf("expected date %s, got: %s", expectedDate, xml)
	}
	if !strings.Contains(xml, `<p:cm authorId="1"`) {
		t.Errorf("expected authorId 1, got: %s", xml)
	}
	if !strings.Contains(xml, `<p:text>First comment</p:text>`) {
		t.Errorf("expected first comment text, got: %s", xml)
	}
	if !strings.Contains(xml, `<p:pos x="100" y="200"/>`) {
		t.Errorf("expected position 100,200, got: %s", xml)
	}
}
