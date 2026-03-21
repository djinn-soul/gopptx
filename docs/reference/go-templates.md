# Go Templates Reference

This page documents the high-level template builders in `pkg/pptx/templates`.

Primary source file:

- `pkg/pptx/templates/templates.go`

## Template interface

```go
type Template interface {
    Build() ([]elements.SlideContent, error)
}
```

## Common pattern

```go
slides, err := pptx.StatusTemplate{
    Project:   "Q4 Launch",
    OKRs:      []string{"Revenue on target", "Docs complete"},
    Risks:     []string{"Hiring delay"},
    NextSteps: []string{"Ship v1", "Review metrics"},
}.Build()
if err != nil {
    panic(err)
}
```

## `SimpleTemplate`

- `Title string`
- `Content string`
- `Branding BrandingSpec`
- `LayoutOverrides LayoutOverrides`

Builds a title slide and a content slide.

## `ProposalTemplate`

- `Title string`
- `Subtitle string`
- `Context string`
- `Solution string`
- `Pricing []PricingTier`
- `Timeline []Milestone`
- `Branding BrandingSpec`
- `LayoutOverrides LayoutOverrides`

Builds title, context, solution, pricing, and timeline slides.

## `TrainingTemplate`

- `Title string`
- `Agenda []string`
- `Concepts []string`
- `Summary string`
- `Branding BrandingSpec`
- `LayoutOverrides LayoutOverrides`

Builds a training deck with agenda, concept, and summary slides.

## `StatusTemplate`

- `Project string`
- `OKRs []string`
- `Risks []string`
- `NextSteps []string`
- `Branding BrandingSpec`
- `LayoutOverrides LayoutOverrides`

Builds a four-slide project status deck.

## `TechnicalTemplate`

- `Title string`
- `Architecture string`
- `DeepDive string`
- `Benchmarks string`
- `Branding BrandingSpec`
- `LayoutOverrides LayoutOverrides`

Builds a technical deep-dive deck.

## Supporting types

- `BrandingSpec`
- `LayoutOverrides`
- `PricingTier`
- `Milestone`

See also:

- [Go API Reference](go-api.md)
