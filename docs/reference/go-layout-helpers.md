# Go Layout Helpers Reference

This page documents the slide-layout utility functions in `pkg/pptx`.

Primary source file:

- `pkg/pptx/layout_helpers.go`

## Types

```go
type Point  = common.Point   // {X, Y styling.Length}
type Size   = common.Size    // {W, H styling.Length}
type Box    = common.Box     // {X, Y, W, H styling.Length}
```

## Constants

- `SlideWidth` — standard slide width (10 inches in EMU)
- `SlideHeight` — standard slide height (7.5 inches in EMU)
- `OrientationHorizontal` — `"horizontal"`
- `OrientationVertical` — `"vertical"`

## Functions

### `Center(cx, cy styling.Length) (styling.Length, styling.Length)`

Return the top-left position that centres a shape of the given size on a standard 10×7.5-inch slide.

### `CenterInSize(cx, cy, totalW, totalH styling.Length) (styling.Length, styling.Length)`

Like `Center` but against an arbitrary canvas size.

### `CenterInBox(cx, cy styling.Length, bounds common.Box) (styling.Length, styling.Length)`

Centre a shape inside an explicit bounding box.

### `Grid(rows, cols int, margin styling.Length) ([]common.Box, error)`

Divide the standard slide into an evenly spaced grid with the given margin between cells.
Returns one `Box` per cell in row-major order.

### `GridInBox(rows, cols int, margin styling.Length, bounds common.Box) ([]common.Box, error)`

Like `Grid` but within an explicit bounding box instead of the full slide.

### `Stack(orientation string, start common.Point, gap styling.Length, elements ...common.Size) ([]common.Point, error)`

Return the top-left position for each element when stacked horizontally or vertically
with a fixed gap between them. `orientation` is `OrientationHorizontal` or `OrientationVertical`.

### `DistributeUniform(orientation string, bounds common.Box, count int, elementSize styling.Length) ([]styling.Length, error)`

Evenly distribute `count` elements of `elementSize` along the given axis of `bounds`.
Returns one position (x or y) per element.

## Typical usage

```go
import "github.com/djinn-soul/gopptx/pkg/pptx"

// Centre a 4×3 inch shape on the slide
x, y := pptx.Center(pptx.Inches(4), pptx.Inches(3))

// Three-column grid with 0.1-inch gutters
cells, err := pptx.Grid(1, 3, pptx.Inches(0.1))
if err != nil {
    panic(err)
}
for i, cell := range cells {
    shape := pptx.NewRectangle(
        cell.X.Inches(), cell.Y.Inches(),
        cell.W.Inches(), cell.H.Inches(),
    )
    _ = i
    _ = shape
}
```

## See also

- [Go Shapes Reference](go-shapes.md)
- [Go API Reference](go-api.md)
