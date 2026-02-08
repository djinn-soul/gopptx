# Shape Types Comparison: gopptx vs. ppt-rs

This document compares the supported shape types between `gopptx` and the reference implementation `ppt-rs`.

## Summary

- **gopptx** focuses on a core set of 17 high-frequency shapes using a flexible string-based system.
- **ppt-rs** provides an exhaustive enum-based system supporting over 100 shape variants.

## Shared Core Shapes

Both libraries support the following essential shapes with identical OOXML preset names:

| gopptx Constant | ppt-rs Enum | OOXML Preset |
|-----------------|-------------|--------------|
| `ShapeTypeRectangle` | `Rectangle` | `rect` |
| `ShapeTypeRoundedRectangle` | `RoundedRectangle` | `roundRect` |
| `ShapeTypeEllipse` | `Ellipse` / `Circle` | `ellipse` |
| `ShapeTypeTriangle` | `Triangle` | `triangle` |
| `ShapeTypeRightTriangle` | `RightTriangle` | `rtTriangle` |
| `ShapeTypeDiamond` | `Diamond` | `diamond` |
| `ShapeTypePentagon` | `Pentagon` | `pentagon` |
| `ShapeTypeHexagon` | `Hexagon` | `hexagon` |
| `ShapeTypeParallelogram` | `Parallelogram` | `parallelogram` |
| `ShapeTypeCloud` | `Cloud` | `cloud` |
| `ShapeTypeRightArrow` | `RightArrow` | `rightArrow` |
| `ShapeTypeLeftArrow` | `LeftArrow` | `leftArrow` |
| `ShapeTypeUpArrow` | `UpArrow` | `upArrow` |
| `ShapeTypeDownArrow` | `DownArrow` | `downArrow` |
| `ShapeTypeFlowChartProcess` | `FlowChartProcess` | `flowChartProcess` |
| `ShapeTypeFlowChartDecision` | `FlowChartDecision` | `flowChartDecision` |
| `ShapeTypeFlowChartTerminator` | `FlowChartTerminator` | `flowChartTerminator` |

## Missing in gopptx (Gaps)

The following categories supported by `ppt-rs` are not yet explicitly defined as constants in `gopptx`:

### 1. Advanced Geometry
- **Octagon** (`octagon`)
- **Trapezoid** (`trapezoid` / `isoTrapezoid`)
- **Cube** (`cube`)
- **Cylinder** / **Can** (`can`)
- **Heart** (`heart`)
- **Lightning Bolt** (`lightningBolt`)

### 2. Stars and Banners
- **Stars (4 to 32 points)** (`star4` through `star32`)
- **Ribbons** (`ribbon`, `ribbon2`)
- **Wave** (`wave`)

### 3. Callouts
- **Wedge Rectangle Callout** (`wedgeRectCallout`)
- **Wedge Oval Callout** (`wedgeEllipseCallout`)
- **Cloud Callout** (`cloudCallout`)

### 4. Expanded Flow Chart
- **Document** (`flowChartDocument`)
- **Data / Input-Output** (`flowChartData`)
- **Internal Storage** (`flowChartInternalStorage`)
- **Magnetic Disk/Drum** (`flowChartMagneticDisk`)

### 5. Action Buttons
- **Home, Help, Back, Forward, Information, Sound, Movie**

## Implementation Differences

### Extension Model
- **gopptx** uses a `ShapeDefinition` interface, allowing you to easily add *any* OOXML shape by passing the raw preset string (e.g., `NewShape("star12", ...)`), even if it's not a predefined constant.
- **ppt-rs** requires the shape to be part of its `ShapeType` enum to be used safely.

### Alias Handling
- **gopptx** provides a `normalizeShapeType` function that maps common aliases like `"circle"` to `ellipse` automatically.
- **ppt-rs** uses Rust's pattern matching to handle aliases.
