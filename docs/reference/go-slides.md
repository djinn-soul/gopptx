# Go Slides Reference

This page documents `pptx.NewSlide` and the `SlideContent` fluent API.

Primary source files:

- `pkg/pptx/slide.go`
- `pkg/pptx/elements/slide.go`

## Construction

### `NewSlide(title string) SlideContent`

Create a new slide with the default title-and-content layout.

## Common slide methods

- `AddBullet(text string) SlideContent`
- `AddBulletWithStyle(text string, style ParagraphStyle) SlideContent`
- `AddBulletRuns(runs []Run) SlideContent`
- `AddBulletRunsWithStyle(runs []Run, style ParagraphStyle) SlideContent`
- `AddShape(sd shapes.ShapeDefinition) SlideContent`
- `AddImage(img shapes.Image) SlideContent`
- `AddComment(authorName, text string) SlideContent`
- `WithNotes(notes string) SlideContent`
- `WithRichNotes(body []Paragraph) SlideContent`
- `AddNoteParagraph(p Paragraph) SlideContent`
- `AddNoteBullet(text string) SlideContent`
- `AddNoteNumbered(text string) SlideContent`
- `AddNoteSubBullet(level int, text string) SlideContent`
- `WithTable(t tables.Table) SlideContent`
- `WithSlideNumber(show bool) SlideContent`
- `WithTransition(t transitions.SlideTransition) SlideContent`
- `WithTransitionOptions(opt transitions.TransitionOptions) SlideContent`
- `WithMorphTransition() SlideContent`
- `WithMorphTransitionOptions(option transitions.MorphOption) SlideContent`
- `WithTransitionSound(path string) SlideContent`
- `WithLayout(layout string) SlideContent`
- `WithBackgroundColor(color string) SlideContent`
- `WithBackground(bg SlideBackground) SlideContent`
- `WithGradientBackground(gradient shapes.ShapeGradientFill) SlideContent`
- `WithPictureBackground(img shapes.Image) SlideContent`
- `WithTitleSize(size int) SlideContent`
- `WithTitleColor(color string) SlideContent`
- `WithTitleBold(bold bool) SlideContent`
- `WithTitleItalic(italic bool) SlideContent`
- `WithTitleUnderline(underline bool) SlideContent`
- `WithTitleAlign(align string) SlideContent`
- `WithTitleFont(font string) SlideContent`
- `WithContentSize(size int) SlideContent`
- `WithContentColor(color string) SlideContent`
- `WithContentBold(bold bool) SlideContent`
- `WithContentItalic(italic bool) SlideContent`
- `WithContentUnderline(underline bool) SlideContent`
- `WithContentVAlign(align string) SlideContent`
- `WithDefaultBulletStyle(style ParagraphStyle) SlideContent`
- `WithBulletStyle(style ParagraphStyle) SlideContent`
- `WithBulletStyleName(styleName string) SlideContent`
- `WithTitleOnlyLayout() SlideContent`
- `WithBlankLayout() SlideContent`
- `WithCenteredTitleLayout() SlideContent`
- `WithTitleAndBigContentLayout() SlideContent`
- `WithTwoColumnLayout() SlideContent`
- `WithTitleAndContentLayout() SlideContent`

## Layout constants

- `SlideLayoutTitleAndContent`
- `SlideLayoutTitleOnly`
- `SlideLayoutBlank`
- `SlideLayoutCenteredTitle`
- `SlideLayoutTitleAndBigContent`
- `SlideLayoutTwoColumn`

## Usage note

Use this page when you want to inspect the slide model directly.
For deck-level construction, see [Go API Reference](go-api.md).

## See also

- [Go Notes, Comments, and Sections Reference](go-notes-comments-sections.md)
- [Go API Reference](go-api.md)
