# Go Hyperlinks Reference

This page documents the hyperlink and action constructors in `pkg/pptx`.

Primary source file:

- `pkg/pptx/hyperlink.go`

## Types

- `Hyperlink` — a clickable or hover action attached to a shape or run
- `HyperlinkAction` — the target of the action (URL, slide, email, file, etc.)

## Constructors

### `NewHyperlink(action HyperlinkAction) Hyperlink`

Wrap an action into a `Hyperlink` that can be passed to `Run.WithHyperlink` or `Run.WithHoverAction`.

## Action constructors

### Navigation within the presentation

- `HyperlinkFirstSlide() HyperlinkAction`
- `HyperlinkLastSlide() HyperlinkAction`
- `HyperlinkNextSlide() HyperlinkAction`
- `HyperlinkPreviousSlide() HyperlinkAction`
- `HyperlinkSlide(slideNum uint32) HyperlinkAction` — 1-based slide number
- `HyperlinkEndShow() HyperlinkAction`

### External targets

- `HyperlinkURL(urlStr string) HyperlinkAction`
- `HyperlinkEmail(address string) HyperlinkAction`
- `HyperlinkEmailWithSubject(address, subject string) HyperlinkAction`
- `HyperlinkFile(path string) HyperlinkAction`
- `HyperlinkProgram(path string) HyperlinkAction`

## Attaching hyperlinks to text

```go
import "github.com/djinn-soul/gopptx/pkg/pptx"

run := pptx.NewRun("Visit our site").
    WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com")))

slide := pptx.NewSlide("Links").AddBulletRuns([]pptx.Run{run})
```

## Attaching hover actions to a run

```go
run := pptx.NewRun("Hover me").
    WithHoverAction(pptx.NewHyperlink(pptx.HyperlinkNextSlide()))
```

## See also

- [Go Text Reference](go-text.md)
- [Go Slides Reference](go-slides.md)
- [Go API Reference](go-api.md)
