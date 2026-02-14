# Editing Existing Presentations

`gopptx` provides a powerful `PresentationEditor` API for reading, modifying, and saving existing PPTX files.

## Opening a Presentation

You can open an existing PPTX from a file path:

```go
editor, err := pptx.OpenPresentationEditor("my_deck.pptx")
if err != nil {
    log.Fatal(err)
}
```

## Inspecting Content

The editor provides access to presentation and slide-level metadata:

```go
// Get presentation metadata
meta := editor.Metadata()
fmt.Printf("Title: %s, Slides: %d
", meta.Title, meta.SlideCount)

// Get individual slide details
for _, slide := range editor.Slides() {
    fmt.Printf("[%d] ID: %d, Title: %s
", slide.Index, slide.SlideID, slide.Title)
}
```

## Modifying Slides

You can add, update, or remove slides within the presentation.

### Adding a Slide
```go
newSlide := pptx.NewSlide("New Summary").AddBullet("Item 1")
index, err := editor.AddSlide(newSlide)
```

### Updating a Slide
```go
updatedSlide := pptx.NewSlide("Updated Title").AddBullet("New Content")
err := editor.UpdateSlide(0, updatedSlide) // Replaces the first slide
```

## Editing Images (Reusing `shapes.Image`)

`PresentationEditor` uses the same `shapes.Image` model used by regular slide generation APIs, so image placement and options are shared between create/edit paths.

```go
updated := pptx.NewSlide("Product Shot").AddImage(
    shapes.NewImage("assets/product.png", 914400, 914400, 3657600, 2057400).
        WithAltText("Product image"),
)
if err := editor.UpdateSlide(0, updated); err != nil {
    log.Fatal(err)
}
```

For in-memory image bytes, use `shapes.NewImageFromBytes(...)`. If you need package-level reuse/dedup across slides, register once and reuse through editor flows:

```go
partPath, err := editor.RegisterImage(imageBytes, "png")
if err != nil {
    log.Fatal(err)
}
fmt.Println("registered media part:", partPath)
```

### Removing a Slide
```go
err := editor.RemoveSlide(1) // Removes the second slide
```

## Updating Theme and Slide Size

You can now apply a new theme and change presentation dimensions on an existing file:

```go
err := editor.ApplyTheme(styling.ThemeTech)
if err != nil {
    log.Fatal(err)
}

err = editor.SetSlideSize(pptx.SlideSize16x9)
if err != nil {
    log.Fatal(err)
}
```

## Merging Presentations

You can merge slides from another PPTX file into the current one:

```go
err := editor.MergeFromFile("another_deck.pptx")
```

## Managing Charts in Existing Slides

`PresentationEditor` supports chart discovery and data updates on existing slides, plus adding new charts.

```go
// Discover chart references on slide 0.
chartsOnSlide, err := editor.ListSlideCharts(0)
if err != nil {
    log.Fatal(err)
}

// Update the first chart by index.
idx := 0
err = editor.UpdateChartData(0, common.ChartSelector{Index: &idx}, common.ChartDataUpdate{
    Categories: []string{"Q1", "Q2", "Q3"},
    Series: []common.ChartSeriesData{
        {Values: []float64{12, 18, 21}},
    },
})
if err != nil {
    log.Fatal(err)
}
```

To add a chart to an existing slide:

```go
chartDef := charts.NewBarChart(
    []string{"A", "B"},
    []float64{10, 20},
).WithTitle("Quarterly")
if err := editor.AddChart(0, chartDef); err != nil {
    log.Fatal(err)
}
```

## Global Text Find/Replace (Shape Text Frames)

```go
count, err := editor.FindAndReplaceInShapes("Old Brand", "New Brand")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("replacements: %d\n", count)
```

## Swapping Existing Slide Images

```go
// Replace the first image relationship on slide 0.
err := editor.SwapImageByIndex(0, 0, newImageBytes, "png")
if err != nil {
    log.Fatal(err)
}
```

## Batch Shape Search

```go
matches, err := editor.SearchShapes(common.ShapeSearchQuery{
    TextContains:  "TODO",
    CaseSensitive: false,
})
if err != nil {
    log.Fatal(err)
}
for _, m := range matches {
    fmt.Printf("slide=%d shape=%d name=%q\n", m.SlideIndex, m.Shape.ID, m.Shape.Name)
}
```

## Saving Changes

Changes are committed only when you save the editor to a file or serialize it to bytes:

```go
// Save to a new file
err := editor.Save("edited_deck.pptx")

// Or get the bytes
data, err := editor.Bytes()
```

## Current Limitations

The `PresentationEditor` is under active development. The following constraints apply to **Add** and **Update** operations:

1.  **Layout Relationship Safety:** Updating a slide with image content requires an existing `slideLayout` relationship on that slide; malformed decks missing this relationship are rejected.
2.  **Complex Merges:** Merging slides that contain external assets (images/charts) may still require careful validation for all edge-case source packages.
3.  **Theme/Master Merge:** The editor can replace `ppt/theme/theme1.xml` in the active file, but does not yet merge multiple source themes or replace the global slide master.
