package editor_test

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func TestCommentsIntegration(t *testing.T) {
	srcPath := buildBaseCommentsPresentation(t)
	ed := mustOpenEditor(t, srcPath)

	a1 := mustAddAuthor(t, ed, "Alice", "AL")
	a2 := mustAddAuthor(t, ed, "Bob", "BB")
	mustAddComment(t, ed, 0, a1, "Alice's comment", 100, 100)
	mustAddComment(t, ed, 0, a1, "Alice's second comment", 200, 200)
	mustAddComment(t, ed, 0, a2, "Bob's comment", 300, 300)

	outPath := filepath.Join(t.TempDir(), "comments.pptx")
	if err := ed.Save(outPath); err != nil {
		t.Fatalf("save edited: %v", err)
	}
	if err := ed.Close(); err != nil {
		t.Errorf("close editor: %v", err)
	}

	ed2 := mustOpenEditor(t, outPath)
	defer func() {
		if err := ed2.Close(); err != nil {
			t.Errorf("close reopened editor: %v", err)
		}
	}()

	assertAuthorsCount(t, ed2, 2)
	assertInitialComments(t, ed2, a1, a2)

	if err := ed2.RemoveComment(0, a1, 1); err != nil {
		t.Fatalf("remove comment: %v", err)
	}
	assertCommentsAfterRemoval(t, ed2)
}

func buildBaseCommentsPresentation(t *testing.T) string {
	t.Helper()
	p := pptx.NewPresentationBuilder("Test Title")
	p.AddTitleSlide("Slide 1")
	srcPath := filepath.Join(t.TempDir(), "base.pptx")
	if err := p.WriteToFile(srcPath); err != nil {
		t.Fatalf("save base: %v", err)
	}
	return srcPath
}

func mustOpenEditor(t *testing.T, filePath string) *editor.PresentationEditor {
	t.Helper()
	ed, err := editor.OpenPresentationEditor(filePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	return ed
}

func mustAddAuthor(
	t *testing.T,
	ed *editor.PresentationEditor,
	name, initials string,
) int64 {
	t.Helper()
	author, err := ed.AddAuthor(name, initials)
	if err != nil {
		t.Fatalf("add author %s: %v", name, err)
	}
	return author.ID
}

func mustAddComment(
	t *testing.T,
	ed *editor.PresentationEditor,
	slideIndex int,
	authorID int64,
	text string,
	x, y int64,
) {
	t.Helper()
	if err := ed.AddComment(slideIndex, authorID, text, x, y); err != nil {
		t.Fatalf("add comment %q: %v", text, err)
	}
}

func assertAuthorsCount(t *testing.T, ed *editor.PresentationEditor, want int) {
	t.Helper()
	authors, err := ed.GetAuthors()
	if err != nil {
		t.Fatalf("get authors: %v", err)
	}
	if len(authors) != want {
		t.Errorf("expected %d authors, got %d", want, len(authors))
	}
}

func assertInitialComments(t *testing.T, ed *editor.PresentationEditor, aliceID, bobID int64) {
	t.Helper()
	comments, err := ed.GetComments(0)
	if err != nil {
		t.Fatalf("get comments: %v", err)
	}
	if len(comments) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(comments))
	}

	assertCommentField(t, comments[0].Text, "Alice's comment", "first comment text")
	assertCommentField(t, comments[0].AuthorID, aliceID, "first comment author")
	assertCommentField(t, comments[0].Index, 1, "first comment index")
	assertCommentField(t, comments[1].Text, "Alice's second comment", "second comment text")
	assertCommentField(t, comments[1].Index, 2, "second comment index")
	assertCommentField(t, comments[2].AuthorID, bobID, "third comment author")
}

func assertCommentsAfterRemoval(t *testing.T, ed *editor.PresentationEditor) {
	t.Helper()
	commentsAfter, err := ed.GetComments(0)
	if err != nil {
		t.Fatalf("get comments after remove: %v", err)
	}
	if len(commentsAfter) != 2 {
		t.Errorf("expected 2 comments after removal, got %d", len(commentsAfter))
	}
	assertCommentField(
		t,
		commentsAfter[0].Text,
		"Alice's second comment",
		"remaining first comment text",
	)
}

func assertCommentField[T comparable](t *testing.T, got, want T, field string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: expected %v, got %v", field, want, got)
	}
}
