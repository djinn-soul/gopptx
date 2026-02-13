package editor_test

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func TestCommentsIntegration(t *testing.T) {
	// 1. Create a basic presentation
	p := pptx.NewPresentationBuilder("Test Title")
	p.AddTitleSlide("Slide 1") // Slide 0
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "base.pptx")
	if err := p.WriteToFile(srcPath); err != nil {
		t.Fatalf("save base: %v", err)
	}

	// 2. Open with Editor
	ed, err := editor.OpenPresentationEditor(srcPath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	// Do not defer close immediately because we close explicitly before save check? No, defer is fine if we close explicitly later.
	// But double close is safe? editor.Close is safe.

	// 3. Add Author
	a1, err := ed.AddAuthor("Alice", "AL")
	if err != nil {
		t.Fatalf("add author 1: %v", err)
	}
	a2, err := ed.AddAuthor("Bob", "BB")
	if err != nil {
		t.Fatalf("add author 2: %v", err)
	}

	// 4. Add Comments
	// Slide 0
	err = ed.AddComment(0, a1.ID, "Alice's comment", 100, 100)
	if err != nil {
		t.Fatalf("add comment 1: %v", err)
	}
	// Alice adds another
	err = ed.AddComment(0, a1.ID, "Alice's second comment", 200, 200)
	if err != nil {
		t.Fatalf("add comment 2: %v", err)
	}
	// Bob adds one
	err = ed.AddComment(0, a2.ID, "Bob's comment", 300, 300)
	if err != nil {
		t.Fatalf("add comment 3: %v", err)
	}

	// 5. Save
	outPath := filepath.Join(tempDir, "comments.pptx")
	if err := ed.Save(outPath); err != nil {
		t.Fatalf("save edited: %v", err)
	}
	if err := ed.Close(); err != nil {
		t.Errorf("close editor: %v", err)
	}

	// 6. Verify persistence
	ed2, err := editor.OpenPresentationEditor(outPath)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() {
		if closeErr := ed2.Close(); closeErr != nil {
			t.Errorf("close reopened editor: %v", closeErr)
		}
	}()

	authors, err := ed2.GetAuthors()
	if err != nil {
		t.Fatalf("get authors: %v", err)
	}
	if len(authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(authors))
	}

	comments, err := ed2.GetComments(0)
	if err != nil {
		t.Fatalf("get comments: %v", err)
	}
	if len(comments) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(comments))
	}

	// Verify order / content
	if comments[0].Text != "Alice's comment" {
		t.Errorf("expected 'Alice's comment', got '%s'", comments[0].Text)
	}
	if comments[0].AuthorID != a1.ID {
		t.Errorf("expected author ID %d, got %d", a1.ID, comments[0].AuthorID)
	}
	if comments[0].Index != 1 {
		t.Errorf("expected index 1, got %d", comments[0].Index)
	}

	if comments[1].Text != "Alice's second comment" {
		t.Errorf("expected 'Alice's second comment', got '%s'", comments[1].Text)
	}
	if comments[1].Index != 2 {
		t.Errorf("expected index 2, got %d", comments[1].Index)
	}

	if comments[2].AuthorID != a2.ID {
		t.Errorf("expected author ID %d, got %d", a2.ID, comments[2].AuthorID)
	}

	// 7. Remove Comment (Alice's first)
	// Assuming RemoveComment takes (slideIndex, authorID, authorIndex)
	err = ed2.RemoveComment(0, a1.ID, 1)
	if err != nil {
		t.Fatalf("remove comment: %v", err)
	}

	commentsAfter, err := ed2.GetComments(0)
	if err != nil {
		t.Fatalf("get comments after remove: %v", err)
	}
	if len(commentsAfter) != 2 {
		t.Errorf("expected 2 comments after removal, got %d", len(commentsAfter))
	}
	if commentsAfter[0].Text != "Alice's second comment" {
		t.Errorf("expected remaining comment to be 'Alice's second comment', got '%s'", commentsAfter[0].Text)
	}
}
