"""Demonstrate adding shapes to slides.

This example demonstrates:
- Adding rectangles, ellipses, rounded rectangles, and triangles via add_shape()
- Adding textboxes with add_textbox()
- Adding connectors (straight, elbow, curved) with add_connector()
- Shape properties: fill color, text content
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.constants import ConnectorType, ShapeType
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches


def _add_basic_shapes_slide(prs: Presentation, slide_idx: int) -> None:
    """Add basic geometric shapes to the given slide."""
    prs.add_shape(
        slide_idx,
        ShapeType.RECTANGLE,
        (Inches(0.5), Inches(1.5), Inches(2), Inches(1.2)),
        text="Rectangle",
        properties={"fill_color": "4472C4"},
    )
    prs.add_shape(
        slide_idx,
        ShapeType.ELLIPSE,
        (Inches(3.0), Inches(1.5), Inches(2), Inches(1.2)),
        text="Ellipse",
        properties={"fill_color": "C0504D"},
    )
    prs.add_shape(
        slide_idx,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(5.5), Inches(1.5), Inches(2.5), Inches(1.2)),
        text="Rounded",
        properties={"fill_color": "F79646"},
    )
    prs.add_shape(
        slide_idx,
        ShapeType.TRIANGLE,
        (Inches(0.5), Inches(3.2), Inches(2), Inches(1.5)),
        text="Triangle",
        properties={"fill_color": "9BBB59"},
    )


def _add_textboxes_slide(prs: Presentation, slide_idx: int) -> None:
    """Add textboxes demonstrating text positioning."""
    prs.add_textbox(
        slide_idx,
        Inches(0.5),
        Inches(1.5),
        Inches(8),
        Inches(0.8),
        text="This is a textbox - positioned with Inches()",
    )
    prs.add_textbox(
        slide_idx,
        Inches(0.5),
        Inches(2.6),
        Inches(5),
        Inches(0.8),
        text="Textboxes have no background fill by default.",
    )
    prs.add_textbox(
        slide_idx,
        Inches(0.5),
        Inches(3.6),
        Inches(6),
        Inches(1.0),
        text="Use add_textbox(slide_idx, x, y, w, h, text=...) to place text anywhere.",
    )


def _add_connectors_slide(prs: Presentation, slide_idx: int) -> None:
    """Add connector shapes between points."""
    prs.add_connector(
        slide_idx,
        ConnectorType.STRAIGHT,
        Inches(0.5),
        Inches(2.0),
        Inches(4.0),
        Inches(2.0),
    )
    prs.add_connector(
        slide_idx,
        ConnectorType.ELBOW,
        Inches(0.5),
        Inches(3.0),
        Inches(4.0),
        Inches(4.5),
    )
    prs.add_connector(
        slide_idx,
        ConnectorType.CURVED,
        Inches(5.0),
        Inches(2.0),
        Inches(9.0),
        Inches(4.5),
    )


def main() -> None:
    """Create presentation demonstrating the shapes API."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    output_path = output_dir / "65_shapes_api.pptx"

    with Presentation.new("Shapes API Demo") as prs:
        # Slide 0: basic geometric shapes
        idx = prs.add_slide("Basic Shapes", layout=SlideLayoutType.TITLE_ONLY)
        _add_basic_shapes_slide(prs, idx)

        # Slide 1: textboxes
        idx = prs.add_slide("Textboxes", layout=SlideLayoutType.TITLE_ONLY)
        _add_textboxes_slide(prs, idx)

        # Slide 2: connectors
        idx = prs.add_slide("Connectors", layout=SlideLayoutType.TITLE_ONLY)
        prs.add_textbox(
            idx,
            Inches(0.5),
            Inches(1.2),
            Inches(8),
            Inches(0.6),
            text="Straight / Elbow / Curved connectors drawn below",
        )
        _add_connectors_slide(prs, idx)

        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Demonstrated: add_shape (RECTANGLE, ELLIPSE, ROUNDED_RECTANGLE, TRIANGLE),")
    print("  add_textbox, add_connector (STRAIGHT, ELBOW, CURVED)")


if __name__ == "__main__":
    main()
