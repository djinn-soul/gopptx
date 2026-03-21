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

Attach a slide comment in the slide model.

### Presentation comment helpers

- `GetAuthors()`
- `AddAuthor(name, initials)`
- `GetComments(slideIndex)`
- `AddComment(...)`
- `RemoveComment(...)`

## Sections

### Presentation section helpers

- `GetSections()`
- `AddSection(name, slideIndices)`
- `RemoveSection(name)`
- `RenameSection(oldName, newName)`

## See also

- [Go Slides Reference](go-slides.md)
- [Go API Reference](go-api.md)
