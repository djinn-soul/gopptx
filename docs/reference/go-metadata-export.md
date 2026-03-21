# Go Metadata, Protection, and Export Reference

This page documents the metadata helpers, protection flags, and export-related package APIs in `pkg/pptx`.

Primary source files:

- `pkg/pptx/presentation.go`
- `pkg/pptx/common/metadata.go`
- `pkg/pptx/presentation_builder.go`

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

See also:

- [Go API Reference](go-api.md)
