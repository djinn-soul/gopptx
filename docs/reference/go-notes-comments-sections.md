# Go Notes, Comments, and Sections Reference

This page documents the slide-level notes, comments, and section helpers in `pkg/pptx`.

Primary source files:

- `pkg/pptx/elements/slide.go`
- `pkg/pptx/presentation_builder.go`
- `pkg/pptx/presentation.go`

## Notes

### `WithNotes(notes string) SlideContent`

Attach plain-text speaker notes to a slide.

### `WithRichNotes(body []Paragraph) SlideContent`

Attach rich-text speaker notes.

### `AddNoteParagraph(p Paragraph) SlideContent`

Append one paragraph to the notes body.

### `AddNoteBullet(text string) SlideContent`

Append a bulleted notes paragraph.

### `AddNoteNumbered(text string) SlideContent`

Append a numbered notes paragraph.

### `AddNoteSubBullet(level int, text string) SlideContent`

Append an indented bulleted notes paragraph.

### `get_notes` style accessors

- `GetNotes` behavior is exposed through the slide model in the builder/runtime surfaces.

## Comments

### `AddComment(authorName, text string) SlideContent`

Attach a slide comment during deck construction (builder path).

### `PresentationEditor` comment helpers

The following methods exist on `*PresentationEditor` (the JSON-bridge / editor path), **not** on `*Presentation`:

- `GetAuthors() ([]comments.Author, error)`
- `GetComments(slideIndex int) ([]comments.Comment, error)`
- `AddComment(slideIndex int, authorID int64, text string, x, y int64) error`
- `RemoveComment(slideIndex int, authorID int64, authorIndex int) error`

To open a `PresentationEditor` use `editor.Open(path)` from `pkg/pptx/editor`.

## Sections

### `PresentationEditor` section helpers

The following methods exist on `*PresentationEditor`, **not** on `*Presentation`:

- `GetSections() ([]SectionData, error)`
- `AddSection(name string, slideIndices []int) error`
- `RemoveSection(name string) error`
- `RenameSection(oldName, newName string) error`

## See also

- [Go Slides Reference](go-slides.md)
- [Go API Reference](go-api.md)
