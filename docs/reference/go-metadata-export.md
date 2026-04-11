# Go Metadata, Protection, and Export Reference

This page documents the metadata helpers, protection flags, and export-related package APIs in `pkg/pptx`.

Primary source files:

- `pkg/pptx/presentation.go`
- `pkg/pptx/common/metadata.go`
- `pkg/pptx/presentation_builder.go`

## Opening existing files

Source file: `pkg/pptx/presentation_api.go`

- `Open(path string) (*Presentation, error)`
- `OpenFromReader(r io.Reader) (*Presentation, error)`
- `OpenFromBytes(data []byte) (*Presentation, error)`

Call `Close()` when done (or use `defer p.Close()`).

## Presentation runtime methods

These methods are on `*Presentation` returned by `Open*` or created at runtime.

### Saving

- `Save() error` — write back to the original path
- `SaveAs(path string) error` — write to a new path
- `SaveToBytes() ([]byte, error)` — return the file contents as a byte slice
- `SaveToWriter(w io.Writer) error` — stream the file to any `io.Writer`
- `Close() error` — release resources

### Inspection

- `SlideCount() int`
- `Validate() []structural.Issue`

### Core properties (fine-grained)

- `CoreProperties() common.CoreProperties` / `SetCoreProperties(props common.CoreProperties)`
- `Title() string` / `SetTitle(title string)`
- `Subject() string` / `SetSubject(subject string)`
- `Creator() string` / `SetCreator(creator string)`
- `Author() string` / `SetAuthor(author string)`
- `Keywords() string` / `SetKeywords(keywords string)`
- `Description() string` / `SetDescription(description string)`
- `LastModifiedBy() string` / `SetLastModifiedBy(lastModifiedBy string)`
- `Revision() string` / `SetRevision(revision string)`
- `Created() string` / `SetCreated(created string)`
- `Modified() string` / `SetModified(modified string)`
- `Category() string` / `SetCategory(category string)`
- `ContentStatus() string` / `SetContentStatus(status string)`

## Metadata types

- `Metadata`
- `MetadataFields`
- `SlideSize`
- `CustomXMLPart`
- `Section`

## Slide size helpers

- `SlideSize4x3() SlideSize`
- `SlideSize16x9() SlideSize`

## Package-level build helpers

- `Create(title string, slideCount int) ([]byte, error)`
- `CreateWithSlides(title string, slides []SlideContent) ([]byte, error)`
- `CreateWithMetadata(meta Metadata, slides []SlideContent) ([]byte, error)`
- `WriteFile(path string, title string, slides []SlideContent) error`
- `Validate(pptxData []byte) ([]structural.Issue, error)`
- `Repair(pptxData []byte) ([]byte, structural.RepairResult, error)`

## PresentationBuilder metadata and protection

- `WithMetadata(meta Metadata) *PresentationBuilder`
- `WithSlideSize(size SlideSize) *PresentationBuilder`
- `WithModifyPassword(password string) *PresentationBuilder`
- `WithMarkAsFinal(final bool) *PresentationBuilder`
- `WithSignaturesEnabled(enabled bool) *PresentationBuilder`
- `WithEncryptionPassword(password string) *PresentationBuilder`
- `WithFooter(text string) *PresentationBuilder`
- `WithDateTime(show bool) *PresentationBuilder`

## Core properties

The Go package exposes core properties through `Metadata` and the underlying `common.Metadata` fields:

- `Title`
- `Subject`
- `Creator`
- `Description`
- `FooterText`
- `ShowDateTime`
- `SlideSize`
- `CustomXML`
- `Protection`

## Typical usage

```go
builder := pptx.NewPresentationBuilder("Deck Title").
    WithSlideSize(pptx.SlideSize16x9()).
    WithMarkAsFinal(true).
    WithFooter("Confidential")

data, err := builder.Build()
if err != nil {
    panic(err)
}
_ = data
```

## Image constructors

Source file: `pkg/pptx/image.go`

- `NewImage(path string, x, y, cx, cy styling.Length) shapes.Image` — load from disk
- `NewImageFromBytes(data []byte, format string, x, y, cx, cy styling.Length) shapes.Image`
- `NewImageFromBase64(b64 string, format string, x, y, cx, cy styling.Length) (shapes.Image, error)`
- `NewImageFromURL(url string, x, y, cx, cy styling.Length) shapes.Image` — embed from URL at build time

`format` is the MIME sub-type, e.g. `"png"`, `"jpeg"`, `"gif"`.

## Markdown-to-slides

Source file: `pkg/pptx/markdown_compat.go`

- `SlidesFromMarkdown(content string) ([]SlideContent, error)`
- `SlidesFromMarkdownWithOptions(content string, options markdown.ParseOptions) ([]SlideContent, error)`
- `SlidesFromMarkdownFile(path string) ([]SlideContent, error)`

Parse Markdown into a `[]SlideContent` that can be passed to `CreateWithSlides` or `PresentationBuilder`.
`markdown.ParseOptions` is from `pkg/pptx/markdown`.

## Placeholder types

Source file: `pkg/pptx/placeholder.go`

- `PlaceholderType` — enum for placeholder kinds
- `Placeholder` — placeholder shape struct

### `PlaceholderType` constants

| Constant | Description |
|---|---|
| `PlaceholderTypeTitle` | Title |
| `PlaceholderTypeBody` | Body / content |
| `PlaceholderTypeCentrTitle` | Centered title |
| `PlaceholderTypeSubTitle` | Subtitle |
| `PlaceholderTypeDt` | Date/time |
| `PlaceholderTypeSldNum` | Slide number |
| `PlaceholderTypeFtr` | Footer |
| `PlaceholderTypeHdr` | Header |
| `PlaceholderTypeObj` | Object |
| `PlaceholderTypeChart` | Chart |
| `PlaceholderTypeTbl` | Table |
| `PlaceholderTypeClipArt` | Clip art |
| `PlaceholderTypeDgm` | Diagram |
| `PlaceholderTypeMedia` | Media |
| `PlaceholderTypeSldImg` | Slide image |
| `PlaceholderTypePic` | Picture |

See also:

- [Go API Reference](go-api.md)
