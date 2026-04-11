# Go Slides Reference

This page documents `pptx.NewSlide` and the `SlideContent` fluent API.

Primary source files:

- `pkg/pptx/slide.go`
- `pkg/pptx/elements/slide.go`
- `pkg/pptx/elements/slide_style.go`
- `pkg/pptx/slide_animation.go`
- `pkg/pptx/smartart/smartart.go`

## Construction

### `NewSlide(title string) SlideContent`

Create a new slide with the default title-and-content layout.

## Common slide methods

### Content

- `AddBullet(text string) SlideContent`
- `AddBulletWithStyle(text string, style ParagraphStyle) SlideContent`
- `AddBulletRuns(runs []Run) SlideContent`
- `AddBulletRunsWithStyle(runs []Run, style ParagraphStyle) SlideContent`
- `AddNumbered(text string) SlideContent`
- `AddLettered(text string) SlideContent`
- `AddSubBullet(level int, text string) SlideContent`
- `AddShape(sd shapes.ShapeDefinition) SlideContent`
- `AddImage(img shapes.Image) SlideContent`
- `AddConnector(c shapes.Connector) SlideContent`
- `AutoRerouteConnectors() SlideContent`
- `AddComment(authorName, text string) SlideContent`
- `WithTable(t tables.Table) SlideContent`
- `AddAnimation(anim AnimationDefinition) SlideContent`
- `AddSmartArt(sa smartart.SmartArt) SlideContent`

### Notes

- `WithNotes(notes string) SlideContent`
- `WithRichNotes(body []Paragraph) SlideContent`
- `AddNoteParagraph(p Paragraph) SlideContent`
- `AddNoteBullet(text string) SlideContent`
- `AddNoteNumbered(text string) SlideContent`
- `AddNoteSubBullet(level int, text string) SlideContent`

### Display

- `WithSlideNumber(show bool) SlideContent`
- `WithLayout(layout string) SlideContent`

### Transitions

- `WithTransition(t transitions.SlideTransition) SlideContent`
- `WithTransitionOptions(opt transitions.TransitionOptions) SlideContent`
- `WithMorphTransition() SlideContent`
- `WithMorphTransitionOptions(option transitions.MorphOption) SlideContent`
- `WithTransitionSound(path string) SlideContent`

#### `MorphOption` constants

- `transitions.MorphOptionObject` — morphs at the shape level (default)
- `transitions.MorphOptionWord` — morphs at the word level
- `transitions.MorphOptionCharacter` — morphs at the character level

### Background

- `WithBackgroundColor(color string) SlideContent`
- `WithBackground(bg SlideBackground) SlideContent`
- `WithGradientBackground(gradient shapes.ShapeGradientFill) SlideContent`
- `WithPictureBackground(img shapes.Image) SlideContent`

### Title styling

- `WithTitleSize(size int) SlideContent`
- `WithTitleColor(color string) SlideContent`
- `WithTitleBold(bold bool) SlideContent`
- `WithTitleItalic(italic bool) SlideContent`
- `WithTitleUnderline(underline bool) SlideContent`
- `WithTitleAlign(align string) SlideContent`
- `WithTitleFont(font string) SlideContent`

### Content styling

- `WithContentSize(size int) SlideContent`
- `WithContentColor(color string) SlideContent`
- `WithContentBold(bold bool) SlideContent`
- `WithContentItalic(italic bool) SlideContent`
- `WithContentUnderline(underline bool) SlideContent`
- `WithContentVAlign(align string) SlideContent`
- `WithDefaultBulletStyle(style ParagraphStyle) SlideContent`
- `WithBulletStyle(style ParagraphStyle) SlideContent`
- `WithBulletStyleName(styleName string) SlideContent`

### Layout shortcuts

- `WithTitleOnlyLayout() SlideContent`
- `WithBlankLayout() SlideContent`
- `WithCenteredTitleLayout() SlideContent`
- `WithTitleAndBigContentLayout() SlideContent`
- `WithTwoColumnLayout() SlideContent`
- `WithTitleAndContentLayout() SlideContent`

## Background constructors

- `NewSolidBackground(color string) SlideBackground`
- `NewGradientBackground(gradient shapes.ShapeGradientFill) SlideBackground`
- `NewPictureBackground(img shapes.Image) SlideBackground`

## Master and notes-master constructors

- `NewMaster() *SlideMaster`
- `NewNotesMaster() *NotesMaster`

## Layout constants

### Slug-style (primary)

- `SlideLayoutTitleAndContent`
- `SlideLayoutTitleOnly`
- `SlideLayoutBlank`
- `SlideLayoutCenteredTitle`
- `SlideLayoutTitleAndBigContent`
- `SlideLayoutTwoColumn`

### Descriptive aliases (legacy / match by name)

- `SlideLayoutTitle` — `"Title Slide"`
- `SlideLayoutSectionHeader` — `"Section Header"`
- `SlideLayoutTwoContent` — `"Two Content"`
- `SlideLayoutComparison` — `"Comparison"`
- `SlideLayoutContentCaption` — `"Content with Caption"`
- `SlideLayoutPictureCaption` — `"Picture with Caption"`

## Animations

### Constructor

- `NewAnimation(shapeIndex int, effect AnimationEffect) Animation`

### Effect constants — Entrance

- `AnimationEntranceAppear`
- `AnimationEntranceFade`
- `AnimationEntranceFlyIn`
- `AnimationEntranceFloat`
- `AnimationEntranceSplit`
- `AnimationEntranceWipe`
- `AnimationEntranceShape`
- `AnimationEntranceWheel`
- `AnimationEntranceRandomBars`
- `AnimationEntranceGrowAndTurn`
- `AnimationEntranceZoom`
- `AnimationEntranceSwivel`
- `AnimationEntranceBounce`

### Effect constants — Exit

- `AnimationExitDisappear`
- `AnimationExitFadeOut`
- `AnimationExitFlyOut`
- `AnimationExitFloatOut`

### Effect constants — Emphasis

- `AnimationEmphasisPulse`
- `AnimationEmphasisColorPulse`
- `AnimationEmphasisTeeter`
- `AnimationEmphasisSpin`
- `AnimationEmphasisGrowShrink`
- `AnimationEmphasisDesaturate`
- `AnimationEmphasisDarken`
- `AnimationEmphasisLighten`
- `AnimationEmphasisTransparency`
- `AnimationEmphasisObjectColor`

### Effect constants — Motion paths

- `AnimationPathLines`
- `AnimationPathArcs`
- `AnimationPathTurns`
- `AnimationPathShapes`
- `AnimationPathLoops`
- `AnimationPathCustom`

### Trigger constants

- `AnimationOnClick`
- `AnimationWithPrevious`
- `AnimationAfterPrevious`

### Direction constants

- `AnimationDirIn` / `AnimationDirOut`
- `AnimationDirUp` / `AnimationDirDown`
- `AnimationDirLeft` / `AnimationDirRight`
- `AnimationDirUpLeft` / `AnimationDirUpRight`
- `AnimationDirDownLeft` / `AnimationDirDownRight`

## SmartArt

Source package: `pkg/pptx/smartart`

### Constructors

- `NewSmartArt(layout Layout) SmartArt` — create a diagram with a built-in layout
- `NewSmartArtWithLayout(layout LayoutProvider) SmartArt` — create with a custom layout URI

### `SmartArt` fluent methods

- `WithAltText(text string) SmartArt`
- `WithDecorative(enabled bool) SmartArt`
- `AddNode(node Node) SmartArt`
- `AddItems(items []string) SmartArt`
- `Position(x, y Length) SmartArt`
- `Size(cx, cy Length) SmartArt`
- `WithColorStyle(cs string) SmartArt`
- `WithQuickStyle(qs string) SmartArt`

### `Node` constructor and methods

- `NewNode(text string) Node`
- `(n Node) WithChild(child Node) Node` — append a child node (hierarchy layouts)
- `(n Node) WithColor(color string) Node`

### Custom layouts

- `CustomLayout(uri string) Layout` — create a `Layout` from an arbitrary URI (e.g. a custom `.glox`-based URI)

Implement the `LayoutProvider` interface to plug in layout strategies:

```go
type LayoutProvider interface {
    LayoutURI() string
}
```

Pass a `LayoutProvider` to `NewSmartArtWithLayout`.

### Built-in layout constants

Lists: `BasicBlockList`, `VerticalBlockList`, `HorizontalBulletLst`, `SquareAccentList`, `PictureAccentList`

Process: `BasicProcess`, `AccentProcess`, `AlternatingFlow`, `ContinuousBlockProcess`

Cycle: `BasicCycle`, `TextCycle`, `BlockCycle`

Hierarchy: `OrgChart`, `Hierarchy`, `HorizontalHierarchy`

Venn: `BasicVenn`, `LinearVenn`, `StackedVenn`

Radial: `BasicRadial`

Matrix: `BasicMatrix`, `TitledMatrix`

Pyramid: `BasicPyramid`, `InvertedPyramid`

Picture: `PictureStrips`, `PictureGrid`

### Adding SmartArt to a slide

```go
import "github.com/djinn-soul/gopptx/pkg/pptx/smartart"

sa := smartart.NewSmartArt(smartart.BasicBlockList).
    AddItems([]string{"Step 1", "Step 2", "Step 3"}).
    Position(pptx.Inches(1), pptx.Inches(2)).
    Size(pptx.Inches(8), pptx.Inches(4))

slide := pptx.NewSlide("My Slide").AddSmartArt(sa)
```

## Usage note

Use this page when you want to inspect the slide model directly.
For deck-level construction, see [Go API Reference](go-api.md).

## See also

- [Go Notes, Comments, and Sections Reference](go-notes-comments-sections.md)
- [Go API Reference](go-api.md)
