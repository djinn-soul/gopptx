# Go Text, Paragraph, and Run Reference

This page documents the rich-text building blocks: `Run`, `Paragraph`, and `ParagraphStyle`.
These types are used wherever `AddBulletRuns`, `AddBulletRunsWithStyle`, and `WithRichNotes` accept structured text.

Primary source files:

- `pkg/pptx/text/text_run.go`
- `pkg/pptx/text/text_paragraph.go`
- `pkg/pptx/text/alignment.go`
- `pkg/pptx/text/bullet_style.go`

## Run

A `Run` is a contiguous span of text with uniform formatting.

### Constructor

- `NewRun(text string) Run`

### Fluent methods

#### Basic formatting

- `WithBold(bold bool) Run`
- `WithItalic(italic bool) Run`
- `WithUnderline(underline bool) Run`
- `WithUnderlineStyle(style string) Run`
- `WithStrikethrough(strikethrough bool) Run`
- `WithStrikethroughStyle(style string) Run`
- `WithAllCaps(allCaps bool) Run`
- `WithSmallCaps(smallCaps bool) Run`
- `WithSubscript(subscript bool) Run`
- `WithSuperscript(superscript bool) Run`

#### Color and highlight

- `WithColor(color string) Run` — hex color, e.g. `"FF0000"`
- `WithHighlight(color string) Run`

#### Font and size

- `WithFont(font string) Run`
- `WithSizePt(size int) Run`

#### Miscellaneous

- `WithCode(code bool) Run`
- `WithLang(lang string) Run` — BCP 47 language tag, e.g. `"en-US"`
- `WithOutline(color string, widthPt ...float64) Run`

#### Hyperlinks

- `WithHyperlink(link action.Hyperlink) Run` — click action
- `WithHoverAction(link action.Hyperlink) Run` — hover/mouse-over action

See [Go Hyperlinks Reference](go-hyperlinks.md) for how to build `Hyperlink` values.

### Underline style constants

- `UnderlineStyleNone`
- `UnderlineStyleSingle`
- `UnderlineStyleDouble`
- `UnderlineStyleDotted`

### Strikethrough style constants

- `StrikethroughStyleNone`
- `StrikethroughStyleSingle`
- `StrikethroughStyleDouble`

---

## Paragraph

A `Paragraph` is a list of `Run`s with an optional `ParagraphStyle`.

### Constructor

- `NewParagraph() Paragraph`

### Fluent methods

- `AddRun(run Run) Paragraph`
- `WithStyle(style ParagraphStyle) Paragraph`

---

## ParagraphStyle

Controls alignment, bullet/list type, indentation, spacing, and more.

### Constructors

- `NewParagraphStyle() ParagraphStyle`
- `DefaultParagraphStyle() ParagraphStyle`

### Alignment

- `WithAlign(align string) ParagraphStyle`
- `WithAlignLeft() ParagraphStyle`
- `WithAlignCenter() ParagraphStyle`
- `WithAlignRight() ParagraphStyle`
- `WithAlignJustify() ParagraphStyle`

### Bullet / list type

- `WithNumbered() ParagraphStyle`
- `WithLetteredLower() ParagraphStyle`
- `WithRomanUpper() ParagraphStyle`
- `WithBulletStyle(style string) ParagraphStyle`
- `WithCustomBullet(char string) ParagraphStyle`
- `WithNoBullet() ParagraphStyle`
- `WithBulletChar(char string) ParagraphStyle`
- `WithBulletColor(color string) ParagraphStyle`
- `WithBulletSize(size int) ParagraphStyle`

### Nesting level

- `WithLevel(level int) ParagraphStyle` — 0-based indent level

### Spacing

- `WithSpaceBeforePt(pt int) ParagraphStyle`
- `WithSpaceAfterPt(pt int) ParagraphStyle`
- `WithLineSpacingPct(pct int) ParagraphStyle` — e.g. `150` for 1.5×
- `WithLineSpacingPts(pt int) ParagraphStyle`

### Indentation

- `WithLeftIndent(emu styling.Length) ParagraphStyle`
- `WithRightIndent(emu styling.Length) ParagraphStyle`
- `WithHangingIndent(emu styling.Length) ParagraphStyle`
- `WithTabStops(stops ...styling.Length) ParagraphStyle`

### Direction

- `WithRTL(rtl bool) ParagraphStyle`

### Bullet style constants

- `BulletStyleBullet`
- `BulletStyleNumber`
- `BulletStyleLetterLower`
- `BulletStyleLetterUpper`
- `BulletStyleRomanLower`
- `BulletStyleRomanUpper`
- `BulletStyleCustom`
- `BulletStyleNone`

### Text alignment constants

- `TextAlignLeft`
- `TextAlignCenter`
- `TextAlignRight`
- `TextAlignJustify`

---

## Typical usage

```go
import "github.com/djinn-soul/gopptx/pkg/pptx"

runs := []pptx.Run{
    pptx.NewRun("Hello ").WithBold(true),
    pptx.NewRun("world").WithColor("FF0000").WithItalic(true),
}

style := pptx.NewParagraphStyle().
    WithAlignCenter().
    WithSpaceBeforePt(6)

slide := pptx.NewSlide("Rich Text").
    AddBulletRunsWithStyle(runs, style)
```

## See also

- [Go Slides Reference](go-slides.md)
- [Go Hyperlinks Reference](go-hyperlinks.md)
- [Go API Reference](go-api.md)
