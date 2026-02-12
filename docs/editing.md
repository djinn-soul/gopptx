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

1.  **Unsupported Content:** Adding or updating slides with images, charts, or speaker notes is not yet supported.
2.  **Complex Merges:** Merging slides that contain external assets (images/charts) will currently fail to protect the integrity of the package relationships.
3.  **Theme/Master Merge:** The editor can replace `ppt/theme/theme1.xml` in the active file, but does not yet merge multiple source themes or replace the global slide master.
