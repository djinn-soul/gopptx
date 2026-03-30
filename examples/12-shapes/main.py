"""Demonstrate shape creation with various shape types, fills, and text.

This example demonstrates:
- Adding rectangles, ellipses, rounded rectangles, and triangles
- Shapes with solid fill colors and text labels
- Using ShapeType constants for type-safe shape creation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_basic_shapes_slide(prs: Presentation) -> None:
    """Add a slide with basic rectangle and ellipse shapes."""
    prs.add_slide("Basic Shapes", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "RECTANGLE",
        bounds=(Inches(1), Inches(1), Inches(3), Inches(1.5)),
        text="Rectangle",
        properties={"fill": {"solid": "4472C4"}},
    )
    prs.add_shape(
        idx,
        "ELLIPSE",
        bounds=(Inches(5), Inches(1), Inches(3), Inches(1.5)),
        text="Ellipse",
        properties={"fill": {"solid": "ED7D31"}},
    )


def _add_rounded_and_triangle_slide(prs: Presentation) -> None:
    """Add a slide with rounded rectangle and triangle shapes."""
    prs.add_slide("Rounded Rectangle and Triangle", layout=SlideLayoutType.BLANK)
    idx = prs.slide_count - 1

    prs.add_shape(
        idx,
        "ROUNDED_RECTANGLE",
        bounds=(Inches(1), Inches(1), Inches(3), Inches(1.5)),
        text="Rounded Rect",
        properties={"fill": {"solid": "A9D18E"}},
    )
    prs.add_shape(
        idx,
        "TRIANGLE",
        bounds=(Inches(5), Inches(1), Inches(3), Inches(1.5)),
        text="Triangle",
        properties={"fill": {"solid": "FF0000"}},
    )


def _add_feature_overview_slide(prs: Presentation) -> None:
    """Add a slide describing the shape features."""
    prs.add_bullet_slide(
        "gopptx Shape Support",
        [
            "Rectangle, ellipse, rounded rectangle, triangle",
            "Solid fill colors (hex RGB)",
            "Text labels on shapes",
            "Precise positioning with Inches()",
        ],
    )


def main() -> None:
    """Create a presentation demonstrating shape creation and styling."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Shapes Demo") as prs:
        _add_feature_overview_slide(prs)
        _add_basic_shapes_slide(prs)
        _add_rounded_and_triangle_slide(prs)

        output_path = output_dir / "12-shapes.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Created 3 slides demonstrating shape features:")
    print("  - Overview of shape support")
    print("  - Rectangle and ellipse with solid fills")
    print("  - Rounded rectangle and triangle")


if __name__ == "__main__":
    main()
